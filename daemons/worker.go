package daemons

import (
	"encoding/json"
	"fmt"
	"github.com/freehere107/go-workers"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan/util"
	"reflect"
)

func RunWorker() {
	regWorker()
	go workers.StatsServer(8080)
	workers.Run()
}

func regWorker() {
	workers.Process("balance", balanceUpdate, 10)
	workers.Process("block", blockWorker, 3)
}

func balanceUpdate(message *workers.Msg) {
	accountList, _ := message.Get("args").StringArray()
	for _, address := range accountList {
		srv.UpdateAccountAllBalance(address)
	}
}

func blockWorker(message *workers.Msg) {
	args, err := message.Get("args").Map()
	if err != nil {
		log.Error("worker blockWorker get args error %v", err)
	}
	blockNum := args["block_num"]
	finalized := args["finalized"]
	if reflect.TypeOf(finalized).Kind().String() != "bool" || reflect.TypeOf(blockNum).Kind().String() != "string" {
		log.Error("worker blockWorker args %v FillBlockData error %v", args, err)
		return
	}
	if err = srv.FillBlockData(util.StringToInt(blockNum.(json.Number).String()), finalized.(bool)); err == nil {
		srv.SetHeartBeat(fmt.Sprintf("%s:heartBeat:%s", util.NetworkNode, "substrate"))
	} else {
		panic(fmt.Sprintf("blockWorker FillBlockData get err %v", err))
	}
}
