package script

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/freehere107/go-workers"
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/libs/substrate"
	"github.com/itering/subscan/libs/substrate/metadata"
	"github.com/itering/subscan/libs/substrate/rpc"
	"github.com/itering/subscan/libs/substrate/storageKey"
	"github.com/itering/subscan/libs/substrate/websocket"
	"github.com/itering/subscan/pkg/recws"
	"github.com/itering/subscan/util"
	"github.com/panjf2000/ants/v2"
	"github.com/shopspring/decimal"
)

// 检测 blocks 的 完整性
func CheckCompleteness() {
	srv := service.New()
	defer srv.Close()
	for {
		alreadyBlockNum, err := srv.GetAlreadyBlockNum()
		if err != nil {
			panic(err)
		}
		var thisRepairedBlock []int
		func() {
			fmt.Println("Now: block height ", alreadyBlockNum)

			repairedBlockNum, _ := srv.GetRepairBlockBlockNum()
			endBlockNum := repairedBlockNum + 300

			if endBlockNum > alreadyBlockNum {
				endBlockNum = alreadyBlockNum
				return
			}

			if endBlockNum/model.SplitTableBlockNum != (repairedBlockNum+1)/model.SplitTableBlockNum {
				endBlockNum = (endBlockNum/model.SplitTableBlockNum)*model.SplitTableBlockNum - 1
			}

			allFetchBlockNums := srv.GetBlockNumArr(repairedBlockNum, endBlockNum)
			fmt.Printf("Start checkout %d ,end %d \r\n", repairedBlockNum, endBlockNum)

			for i := repairedBlockNum + 1; i < endBlockNum; i++ {
				if util.IntInSlice(i, allFetchBlockNums) == false {
					fmt.Printf("Find block %d \r\n", i)
					workers.EnqueueWithOptions("block", "block", map[string]interface{}{"block_num": i, "finalized": true}, workers.EnqueueOptions{RetryCount: 2})
					_ = srv.SetRepairBlockBlockNum(i)
					thisRepairedBlock = append(thisRepairedBlock, i)
				}
			}
			_ = srv.SetRepairBlockBlockNum(endBlockNum)
			fmt.Println(allFetchBlockNums)
		}()
		fmt.Println("Check repair block over, repaired block ....", thisRepairedBlock)
	}
}

// 修复 有 codec error 的 block
func RepairCodecError(option ...string) {
	srv := service.New()
	defer srv.Close()
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(3, func(i interface{}) {
		blockNum := i.(int)
		func(blockNum int) {
			block := srv.GetRawBlockByNum(blockNum)
			finalizedNum, _ := srv.GetFillFinalizedBlockNum()
			workers.EnqueueWithOptions("block", "block", map[string]interface{}{"block_num": blockNum, "finalized": block.BlockNum <= finalizedNum}, workers.EnqueueOptions{RetryCount: 2})
		}(blockNum)
		wg.Done()
	})
	defer p.Release()
	blocks := srv.GetBlockFixDataList(option...)
	for _, id := range blocks {
		wg.Add(1)
		_ = p.Invoke(id)
	}
	wg.Wait()
}

func RefreshAccountInfo(srv *service.Service, args ...string) {
	ctx := context.TODO()
	addresses := args
	var query []string

	if len(addresses) == 0 {
		// CrabNetwork not has account index
		if util.NetworkNode != util.CrabNetwork {
			query = append(query, "account_index>=0")
		}
		list := srv.GetAccountList(ctx, query...)
		for _, account := range list {
			addresses = append(addresses, account.Address)
		}
	}

	var wg sync.WaitGroup
	cp, _ := ants.NewPoolWithFunc(5, func(i interface{}) {
		address := i.(string)
		func(address string) {
			u := refreshAccount(nil, address)
			_ = srv.RefreshAccount(nil, u)
		}(address)
		wg.Done()
	})

	defer cp.Release()
	for _, account := range addresses {
		wg.Add(1)
		_ = cp.Invoke(util.TrimHex(account))
	}
	wg.Wait()
	fmt.Println("finish")
}

func refreshAccount(c *recws.RecConn, address string) map[string]interface{} {
	u := map[string]interface{}{"address": util.TrimHex(address)}

	if balance, otherBalance, err := rpc.GetFreeBalance(c, "Balances", address, ""); err != nil {
		fmt.Println("Get free balance error ", err)
	} else {
		u["balance"] = balance.Div(decimal.New(1, int32(substrate.BalanceAccuracy)))
		if util.IsDarwinia {
			u["kton_balance"] = otherBalance.Div(decimal.New(1, int32(substrate.BalanceAccuracy)))
		}
	}

	if lockBalance, err := rpc.GetAccountLock(c, address, "ring"); err != nil {
		fmt.Println("Get lock balance error ", err)
	} else {
		u["ring_lock"] = lockBalance
	}

	if nonce, err := rpc.GetAccountNonce(c, address); err != nil {
		fmt.Println("Get nonce error ", err)
	} else {
		u["nonce"] = nonce
	}

	if util.IsDarwinia {
		if lockBalance, err := rpc.GetAccountLock(c, address, "kton"); err != nil {
			fmt.Println("Get kton lock balance error ", err)
		} else {
			u["kton_lock"] = lockBalance
		}
	}

	return u
}

func InitValidators() {
	srv := service.New()
	p, err := websocket.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer p.Close()
	c := p.Conn
	list, err := rpc.StakingValidators(c)
	if err != nil {
		fmt.Println(err)
	}
	for _, stash := range list {
		controller, err := rpc.StashController(c, stash.Address)
		if err != nil {
			fmt.Println("controller error", err)
		}
		payee, err := rpc.RewardPayee(c, stash.Address)
		if err != nil {
			fmt.Println("RewardPayee error", err)
		}
		_ = srv.AddNewValidator(util.TrimHex(stash.Address), util.TrimHex(controller), payee, stash.ValidatorPrefsValue)
	}
}

func RefreshTransferStat() {
	srv := service.New()
	defer srv.Close()
	j, _ := srv.GetTransactionList(0, 1000000, "desc")
	for _, trans := range j {
		srv.SetDailyStat(trans.BlockNum, time.Unix(int64(trans.BlockTimestamp), 0).UTC(), "extrinsic")
		if trans.CallModuleFunction == "transfer" {
			srv.SetDailyStat(trans.BlockNum, time.Unix(int64(trans.BlockTimestamp), 0).UTC(), "transfer")
		}

	}
}

func RefreshAccountEvents() {
	srv := service.New()
	defer srv.Close()

	list, _ := srv.GetEventList(0, 1000000, "desc", "event_id in ('slash','reward')")
	for _, event := range list {
		srv.AnalysisEvent("", 0, event, decimal.Zero, 0, decimal.Zero)
	}

}

func SetCodecErrorForInsertError() {
	srv := service.New()
	defer srv.Close()
	alreadyBlockNum, _ := srv.GetAlreadyBlockNum()
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(100, func(i interface{}) {
		id := i.(int)
		func(id int) {
			block := srv.GetBlockByNum(id)
			if block.EventCount != len(block.Events) || block.ExtrinsicsCount != len(block.Extrinsics) {
				srv.SetBlockCodecError(id)
			}
		}(id)
		wg.Done()
	})
	defer p.Release()
	blockNum := 1
	for {
		if blockNum >= alreadyBlockNum {
			break
		}
		wg.Add(1)
		_ = p.Invoke(blockNum)
		blockNum++
	}
	wg.Wait()
}

func InitRuntimeData() {
	srv := service.New()
	defer srv.Close()
	local := metadata.LocalMetadata()
	for spec, raw := range local {
		r := metadata.RuntimeRaw{Spec: util.StringToInt(spec), Raw: raw}
		runtime := []metadata.RuntimeRaw{r}
		m := metadata.Init(&runtime, util.StringToInt(spec))
		srv.MigrateRuntimeVersion(util.StringToInt(spec), m, raw)
	}
}

func InitIdentityData() {
	srv := service.New()
	defer srv.Close()
	list, _ := srv.GetTransactionList(0, 10000, "desc", "call_module = 'identity'")
	for _, v := range list {
		srv.AnalysisExtrinsic(context.TODO(), &model.ChainExtrinsic{
			AccountId:          v.FromHex,
			CallModuleFunction: v.CallModuleFunction,
			CallModule:         v.CallModule,
			Params:             v.Params,
		}, nil)
	}
}

func InitTreasuryPropose() {
	srv := service.New()
	defer srv.Close()
	list, _ := srv.GetTransactionList(0, 10000, "asc", "call_module = 'treasury'", "success = 1")
	for _, v := range list {
		if err := srv.AnalysisExtrinsic(context.TODO(), &model.ChainExtrinsic{
			AccountId:          v.FromHex,
			BlockNum:           v.BlockNum,
			CallModuleFunction: v.CallModuleFunction,
			CallModule:         v.CallModule,
			Params:             []byte(v.Params),
		}, srv.GetEventByIndex(v.ExtrinsicIndex)); err != nil {
			fmt.Println("AnalysisExtrinsic get error", err)
		}

	}
	events, _ := srv.GetEventList(0, 1000000, "asc", "module_id='treasury'")
	for _, event := range events {
		block := srv.GetRawBlockByNum(event.BlockNum)
		srv.AnalysisEvent(block.Hash, block.BlockTimestamp, event, decimal.Zero, 0, decimal.Zero)
	}
}

func InitTechcommCouncilPropose() {
	srv := service.New()
	defer srv.Close()
	list, _ := srv.GetTransactionList(0, 10000, "asc", "call_module in ('technicalcommittee','council','sudo')", "success = 1")
	for _, v := range list {
		if err := srv.AnalysisExtrinsic(context.TODO(), &model.ChainExtrinsic{
			AccountId:          v.FromHex,
			BlockNum:           v.BlockNum,
			CallModuleFunction: v.CallModuleFunction,
			CallModule:         v.CallModule,
			Params:             []byte(v.Params),
			ExtrinsicHash:      v.Hash,
		}, srv.GetEventByIndex(v.ExtrinsicIndex)); err != nil {
			fmt.Println("AnalysisExtrinsic get error", err)
		}
	}
}

func InitDemocracyPropose() {
	srv := service.New()
	defer srv.Close()
	events, _ := srv.GetEventList(0, 1000000, "asc", "module_id='democracy'")
	for _, event := range events {
		block := srv.GetRawBlockByNum(event.BlockNum)
		srv.AnalysisEvent(block.Hash, block.BlockTimestamp, event, decimal.Zero, 0, decimal.Zero)
	}
	list, _ := srv.GetTransactionList(0, 10000, "asc", "call_module  = 'democracy'", "success = 1")
	for _, v := range list {
		if err := srv.AnalysisExtrinsic(context.TODO(), &model.ChainExtrinsic{
			AccountId:          v.FromHex,
			BlockNum:           v.BlockNum,
			CallModuleFunction: v.CallModuleFunction,
			CallModule:         v.CallModule,
			Params:             []byte(v.Params),
			ExtrinsicHash:      v.Hash,
		}, srv.GetEventByIndex(v.ExtrinsicIndex)); err != nil {
			fmt.Println("AnalysisExtrinsic get error", err)
		}
	}
}

func FixSessionsData() {
	srv := service.New()
	defer srv.Close()
	events, _ := srv.GetEventList(0, 1000000, "asc", "event_id='newSession'")
	for _, event := range events {
		if err := srv.EmitRepairAction("FixSessionsData", event); err != nil {
			fmt.Println(err)
		}
	}
}

func FixStakingRewardData() {
	srv := service.New()
	defer srv.Close()
	events, _ := srv.GetEventList(0, 1000000, "asc", "event_id='reward'")
	for _, event := range events {
		if err := srv.EmitRepairAction("FixStakingRewardData", event); err != nil {
			fmt.Println(err)
		}
	}
}

func FillTransferHistory() {
	// Transfer event
	srv := service.New()
	defer srv.Close()
	events, _ := srv.GetEventList(0, 1000000, "asc", "event_id='transfer'")
	for _, event := range events {
		if err := srv.EmitRepairAction("FillTransferHistory", event); err != nil {
			fmt.Println(err)
		}
	}
	// fail transfer
	transfers, _ := srv.GetTransactionList(0, 1000000, "asc", "call_module_function='transfer'", "success = 0")
	for _, transfer := range transfers {
		if err := srv.EmitRepairAction("FillTransferHistory", transfer); err != nil {
			fmt.Println(err)
		}
	}
}

// PrintRuntimeStorageKey
func PrintRuntimeStorageKey() {
	srv := service.New()
	defer srv.Close()
	runtime := metadata.Latest(nil)
	var keys []string
	for _, modules := range runtime.Metadata.Modules {
		for _, storage := range modules.Storage {
			name := modules.Name
			method := storage.Name
			sk := storageKey.EncodeStorageKey(name, method)
			keys = append(keys, fmt.Sprintf("%s|%s|%s", sk.EncodeKey, name, method))
		}
	}
	fmt.Println(util.InterfaceToString(keys))
}
