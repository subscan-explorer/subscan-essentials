package script

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/itering/subscan/util/mq"
	"io"
	"os"

	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
)

func Install(conf string) {
	// create database
	// conf
	_ = fileCopy(fmt.Sprintf("%s/config.yaml.example", conf), fmt.Sprintf("%s/config.yaml", conf))
	func() {
		dbHost := util.GetEnv("MYSQL_HOST", "127.0.0.1")
		dbUser := util.GetEnv("MYSQL_USER", "root")
		dbPass := util.GetEnv("MYSQL_PASS", "")
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, "")
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = db.Close()
		}()
		_, err = db.Exec("CREATE DATABASE IF NOT EXISTS subscan DEFAULT CHARACTER SET = `utf8mb4`")
		if err != nil {
			panic(err)
		}
		fmt.Println("Create database success!!!")

	}()

}

func fileCopy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// CheckCompleteness Check blocks Completeness
func CheckCompleteness(startBlock uint, fastMode bool) {
	srv := service.New()
	defer srv.Close()
	c := context.TODO()

	// latest fill block num
	latest, err := srv.GetFinalizedBlock(c)
	if err != nil {
		panic(err)
	}
	const holdOnNum uint = 20

	util.Debug(fmt.Sprintf("Now: block height %d", latest))
	var (
		latestBlockNum uint
	)

	var fillBlock = func(num uint) {
		if fastMode {
			_ = mq.Instant.Publish("block", "block", map[string]interface{}{"block_num": num, "finalized": true, "force": true})
			return
		}
		// re-sync
		err = srv.FillBlockData(c, num, true)
		if err != nil {
			panic(fmt.Errorf("not found the block num %d %v", num, err))
		}
	}

	if startBlock > 0 {
		latestBlockNum = startBlock - 1
	} else {
		latestBlockNum = 1
	}

	for {
		if latestBlockNum >= uint(latest)-holdOnNum {
			break
		}
		endBlockNum := latestBlockNum + 3000
		if endBlockNum/model.SplitTableBlockNum != (latestBlockNum+1)/model.SplitTableBlockNum {
			endBlockNum = (endBlockNum/model.SplitTableBlockNum)*model.SplitTableBlockNum - 1
		}
		if endBlockNum > uint(latest)-holdOnNum {
			endBlockNum = uint(latest) - holdOnNum
		}

		util.Logger().Info(fmt.Sprintf("Start checkout block %d, end block %d", latestBlockNum+1, endBlockNum))
		var allFetchBlockNums = srv.GetDao().GetBlockNumArr(c, latestBlockNum+1, endBlockNum)

		if uint(len(allFetchBlockNums)) < endBlockNum-latestBlockNum {
			for i := latestBlockNum + 1; i <= endBlockNum; i++ {
				if !util.IntInSlice(int(i), allFetchBlockNums) {
					util.Logger().Info(fmt.Sprintf("Missing block %d", i))
					fillBlock(i)
				}
			}
		}
		latestBlockNum = endBlockNum
	}
}
