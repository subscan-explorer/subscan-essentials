package script

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"

	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
)

func Install(conf string) {
	// create database
	// conf
	_ = fileCopy(fmt.Sprintf("%s/http.toml.example", conf), fmt.Sprintf("%s/http.toml", conf))
	_ = fileCopy(fmt.Sprintf("%s/mysql.toml.example", conf), fmt.Sprintf("%s/mysql.toml", conf))
	_ = fileCopy(fmt.Sprintf("%s/redis.toml.example", conf), fmt.Sprintf("%s/redis.toml", conf))

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
func CheckCompleteness() {
	srv := service.New()
	defer srv.Close()

	c := context.TODO()
	dao := srv.GetDao()
	// latest fill block num
	latest, err := dao.GetFillFinalizedBlockNum(c)
	if err != nil {
		panic(err)
	}

	var thisRepairedBlock []int
	repairedBlockNum := 0
	for {
		endBlockNum := repairedBlockNum + 300

		if endBlockNum > latest {
			break
		}

		if endBlockNum/model.SplitTableBlockNum != (repairedBlockNum+1)/model.SplitTableBlockNum {
			endBlockNum = (endBlockNum/model.SplitTableBlockNum)*model.SplitTableBlockNum - 1
		}

		allFetchBlockNums := dao.GetBlockNumArr(repairedBlockNum, endBlockNum)

		for i := repairedBlockNum; i < endBlockNum; i++ {
			if !util.IntInSlice(i, allFetchBlockNums) {
				// err := srv.FillBlockData(nil, i, true)
				// if err != nil {
				// 	fmt.Println("FillBlockData get error", err)
				// }
				thisRepairedBlock = append(thisRepairedBlock, i)
			}
		}
		repairedBlockNum = endBlockNum
	}

	if len(thisRepairedBlock) > 0 {
		fmt.Println("Check repair block over, repaired block ....", thisRepairedBlock, len(thisRepairedBlock))
	}
}
