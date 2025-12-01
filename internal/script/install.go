package script

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins"
	"github.com/itering/subscan/plugins/balance"
	"github.com/itering/subscan/plugins/evm"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/mq"
	"gorm.io/gorm"
	"io"
	"os"
	"sort"
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

func RefreshMetadata() {
	ctx := context.TODO()
	srv := service.New()
	defer srv.Close()
	d := srv.GetDao()
	u := make(map[string]interface{})
	// extrinsic
	{
		u["count_extrinsic"] = d.GetExtrinsicCount(ctx)
		u["count_signed_extrinsic"] = d.GetExtrinsicCount(ctx, model.Where("is_signed = ?", true))
	}
	util.Logger().Error(srv.GetDao().SetMetadata(ctx, u))
	// balance plugin
	b := plugins.RegisteredPlugins["balance"].(*balance.Balance)
	b.RefreshMetadata()
	// evm plugin
	e := plugins.RegisteredPlugins["evm"].(*evm.EVM)
	e.RefreshMetadata()
}

func MigrateAccountExtrinsicMapping() error {
	srv := service.New()
	defer srv.Close()
	db := srv.GetDbStorage().GetDbInstance().(*gorm.DB)
	// collect existing chain_extrinsics* tables

	mapping := map[string]map[int]struct{}{}

	tableName := fmt.Sprintf("chain_extrinsics")
	idx := 0
	for {
		if idx > 0 {
			tableName = fmt.Sprintf("chain_extrinsics_%d", idx)
		}
		if !db.Migrator().HasTable(tableName) {
			break
		}
		var accounts []string
		if err := db.Table(tableName).Distinct("account_id").Where("account_id <> ''").Pluck("account_id", &accounts).Error; err != nil {
			return err
		}
		for _, a := range accounts {
			if a == "" {
				continue
			}
			if _, ok := mapping[a]; !ok {
				mapping[a] = map[int]struct{}{}
			}
			mapping[a][idx] = struct{}{}
		}
		idx++
	}
	// upsert into account_extrinsic_mapping
	for acc, idxMap := range mapping {
		var m model.AccountExtrinsicMapping
		if err := db.Where("account_id = ?", acc).First(&m).Error; err != nil {
			idxs := model.IntSlice{}
			for k := range idxMap {
				idxs = append(idxs, k)
			}
			sort.Ints(idxs)
			m = model.AccountExtrinsicMapping{AccountId: acc, ExtrinsicTable: idxs}
			if err = db.Scopes(model.IgnoreDuplicate).Create(&m).Error; err != nil {
				return err
			}
			continue
		}

		exist := map[int]struct{}{}
		for _, v := range m.ExtrinsicTable {
			exist[v] = struct{}{}
		}

		changed := false
		for k := range idxMap {
			if _, ok := exist[k]; !ok {
				m.ExtrinsicTable = append(m.ExtrinsicTable, k)
				changed = true
			}
		}
		if changed {
			sort.Ints(m.ExtrinsicTable)
			if err := db.Model(&model.AccountExtrinsicMapping{}).Where("id = ?", m.Id).Update("extrinsic_table", m.ExtrinsicTable).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
