package daemons

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bilibili/kratos/pkg/log"
	"github.com/gorilla/websocket"
	"subscan-end/internal/model"
	"subscan-end/internal/service"
	"subscan-end/utiles"
	"time"
)

func RepairBlock() {
	srv = service.New()
	defer srv.Close()
	c, _, _ := websocket.DefaultDialer.Dial(utiles.ProviderEndPoint, nil)
	defer c.Close()
	for {
		alreadyBlockNum, err := srv.GetAlreadyBlockNum()
		if err != nil {
			panic(err)
		}
		var thisRepairedBlock []int
		func() {
			log.Info("Start Repair Block: ", alreadyBlockNum)
			repairedBlockNum, _ := srv.GetRepairBlockBlockNum()
			allFetchBlockNums := srv.GetBlockNumArr(repairedBlockNum, alreadyBlockNum)
			for i := repairedBlockNum + 1; i < alreadyBlockNum; i++ {
				if len(allFetchBlockNums) <= i || allFetchBlockNums[i] != i {
					fillBlockData(c, i, srv)
					allFetchBlockNums = utiles.InsertInts(allFetchBlockNums, i, i)
					_ = srv.SetRepairBlockBlockNum(i)
					thisRepairedBlock = append(thisRepairedBlock, i)
				}
			}
			fmt.Println(allFetchBlockNums)
		}()
		log.Info("Check repair block over, repaired block ....", thisRepairedBlock)
		time.Sleep(10 * time.Second)
	}
}

func RepairBlockData() {
	srv = service.New()
	defer srv.Close()
	for {
		func() {
			blocks := srv.GetBlockFixDataList()
			for _, block := range *blocks {
				srv.RepairBlockData(&block)
			}
			log.Info("Check repair event over, repaired block data ....")
			time.Sleep(2 * time.Minute)
		}()
	}

}

func RepairValidateInfo() {
	srv = service.New()
	defer srv.Close()
	c := context.TODO()
	list, _ := srv.GetExtrinsicList(0, 10000, "asc", "is_signed = 1")
	for _, e := range *list {
		var paramsInstant []model.ExtrinsicParam
		_ = json.Unmarshal([]byte(e.Params), &paramsInstant)
		extrinsic := model.ChainExtrinsic{CallModule: e.CallModule, AccountId: e.AccountId, CallModuleFunction: e.CallModuleFunction}
		srv.AnalysisExtrinsic(c, &extrinsic, paramsInstant)
	}

}
