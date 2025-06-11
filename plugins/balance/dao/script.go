package dao

import (
	"context"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/model"
	bModel "github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/subscan/share/substrate"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
	"github.com/panjf2000/ants/v2"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"log"
	"sync"
)

func InitAccount(sg *Storage) {
	ctx := context.Background()
	wg := new(sync.WaitGroup)
	bp, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		wg.Add(1)
		defer wg.Done()
		params := i.([]interface{})
		addr := params[0].(string)
		info := params[1].(*bModel.AccountData)
		sg.AddOrUpdateItem(ctx, &bModel.Account{
			Address:  addr,
			Nonce:    info.Nonce,
			Balance:  info.Data.Free.Add(info.Data.Reserved),
			Locked:   decimal.Max(info.Data.MiscFrozen, info.Data.FeeFrozen),
			Reserved: info.Data.Reserved,
		}, []string{"address"}, "nonce", "balance", "locked", "reserved")
	})
	defer bp.Release()

	// refresh account balance
	if err := substrate.BatchReadKeysPaged(ctx, "System", "Account", "", func(keys []string, scaleType string) error {
		r, _ := substrate.BatchStorageByKey(ctx, keys, scaleType, "")
		for key, v := range r {
			val, _ := substrate.ParseStorageKey(key)
			addr := address.Format(val[0].ToString())
			accountData := new(bModel.AccountData)
			v.ToAny(accountData)
			util.Logger().Error(bp.Invoke([]interface{}{addr, accountData}))
		}
		return nil
	}); err != nil {
		log.Panic(err)
	}
	wg.Wait()
}

func RefreshAllAccount(_ *Storage) {

}

func InitTransfer(sg *Storage) {
	c := context.TODO()
	db := sg.Dao.GetDbInstance().(*gorm.DB)

	blockNum, _ := sg.Dao.GetCurrentBlockNum(c)
	for i := int(blockNum); i >= 0; i -= int(model.SplitTableBlockNum) {

		tableName := model.TableNameFromInterface(&model.ChainEvent{BlockNum: uint(i)}, db)
		var events []*model.ChainEvent

		query := db.Table(tableName).
			Where("module_id = ?", "balances").
			Where("event_id = ?", "Transfer")
		query.FindInBatches(&events, 50000, func(tx *gorm.DB, batch int) error {
			var blocks = make(map[int]*storage.Block)
			var blockNums []uint

			for _, e := range events {
				blockNums = append(blockNums, e.BlockNum)
			}

			for _, b := range sg.Dao.GetBlocksByNums(c, blockNums, "id,block_num,block_timestamp") {
				blocks[b.BlockNum] = b
			}

			var extrinsicIds []string
			for _, e := range events {
				extrinsicIds = append(extrinsicIds, e.ExtrinsicIndex)
			}
			for index := range events {
				event := events[index]
				_ = EmitEvent(c, sg, event.AsPlugin(), blocks[int(event.BlockNum)])
			}
			return nil
		})
	}

}
