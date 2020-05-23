package daemons

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/freehere107/go-workers"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan/libs/substrate/rpc"
	"github.com/itering/subscan/util"
)

func RunWorker() {
	regWorker()
	go workers.StatsServer(8080)
	workers.Run()
}

func regWorker() {
	workers.Process("balance", balanceUpdate, 10)
	workers.Process("validatorStakingInfo", validatorStakingInfo, 10)
	workers.Process("block", blockWorker, 3)
	workers.Process("freshIdentityInfo", freshIdentityInfo, 10)
	workers.Process("freshWaiting", freshWaiting, 1)
}

func balanceUpdate(message *workers.Msg) {
	accountList, _ := message.Get("args").StringArray()
	for _, address := range accountList {
		srv.UpdateAccountAllBalance(address)
	}
}

func validatorStakingInfo(message *workers.Msg) {
	stash, err := message.Get("args").String()
	if err != nil {
		log.Error("worker ValidatorStakingInfo error %v", err)
		return
	}
	err = srv.UpdateValidatorStakingInfo(stash)
	if err != nil {
		log.Error("worker ValidatorStakingInfo error %v", err)
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

func freshIdentityInfo(message *workers.Msg) {
	account, err := message.Get("args").String()
	if err != nil {
		log.Error("worker freshIdentityInfo error %v", err)
		return
	}
	if err := srv.UpdateAccountIdentityInfo(account); err != nil {
		log.Error("freshIdentityInfo args %v get error %v", account, err)
	}
}

// after era, fresh all waiting validators
func freshWaiting(message *workers.Msg) {
	list, err := rpc.StakingValidators(nil)
	if err != nil {
		log.Error("worker freshWaitingValidators error %v", err)
		return
	}

	elected := srv.ElectedValidators()
	for _, validator := range list {
		if !util.StringInSlice(validator.Address, elected) {
			_ = srv.UpdateValidatorStakingInfo(validator.Address)
		}
	}
}
