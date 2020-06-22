package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/substrate/rpc"
	"github.com/itering/subscan/internal/util"
	"github.com/itering/subscan/internal/util/ss58"
	"sort"
)

func (d *Dao) CreateBlock(c context.Context, txn *GormDB, hash string, block *rpc.Block, event, log, validator string, eventCount, extrinsicsCount, blockTimestamp int, codecError bool, version int, finalized bool) (err error) {
	extrinsicB, _ := json.Marshal(block.Extrinsics)
	cb := model.ChainBlock{
		Hash:            hash,
		BlockNum:        util.StringToInt(util.HexToNumStr(block.Header.Number)),
		BlockTimestamp:  blockTimestamp,
		ParentHash:      block.Header.ParentHash,
		StateRoot:       block.Header.StateRoot,
		ExtrinsicsRoot:  block.Header.ExtrinsicsRoot,
		Logs:            log,
		Extrinsics:      string(extrinsicB),
		Event:           event,
		SpecVersion:     version,
		ExtrinsicsCount: extrinsicsCount,
		EventCount:      eventCount,
		Validator:       validator,
		CodecError:      codecError,
		Finalized:       finalized,
	}
	query := txn.Create(&cb)
	if !d.Db.HasTable(model.ChainBlock{BlockNum: cb.BlockNum + model.SplitTableBlockNum}) {
		go d.blockMigrate(cb.BlockNum + model.SplitTableBlockNum)
	}
	return d.checkDBError(query.Error)
}

func (d *Dao) SaveFillAlreadyBlockNum(c context.Context, blockNum int) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if num, _ := redis.Int(conn.Do("GET", RedisFillAlreadyBlockNum)); blockNum > num {
		_, err = conn.Do("SET", RedisFillAlreadyBlockNum, blockNum)
	}
	return
}

func (d *Dao) SaveFillAlreadyFinalizedBlockNum(c context.Context, blockNum int) (err error) {
	conn := d.redis.Get(c)
	defer func() {
		conn.Close()
	}()

	if num, _ := redis.Int(conn.Do("GET", RedisFillFinalizedBlockNum)); blockNum > num {
		_, err = conn.Do("SET", RedisFillFinalizedBlockNum, blockNum)
	}
	return
}

func (d *Dao) GetFillAlreadyBlockNum(c context.Context) (num int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	num, err = redis.Int(conn.Do("GET", RedisFillAlreadyBlockNum))
	return
}

func (d *Dao) GetFillFinalizedBlockNum(c context.Context) (num int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	num, err = redis.Int(conn.Do("GET", RedisFillFinalizedBlockNum))
	return
}

func (d *Dao) SaveRepairBlockBlockNum(c context.Context, blockNum int) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("SET", RedisRepairBlockKey, blockNum)
	return
}

func (d *Dao) GetRepairBlockBlockNum(c context.Context) (num int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	num, err = redis.Int(conn.Do("GET", RedisRepairBlockKey))
	return
}

func (d *Dao) GetBlockNumArr(c context.Context, start, end int) []int {
	var blockNums []int
	d.Db.Model(&model.ChainBlock{BlockNum: end}).Where("block_num BETWEEN ? AND ?", start, end).Order("block_num asc").Pluck("block_num", &blockNums)
	return blockNums
}

func (d *Dao) GetBlockList(page, row int) []model.ChainBlock {
	var blocks []model.ChainBlock
	blockNum, _ := d.GetFillAlreadyBlockNum(context.TODO())
	head := blockNum - page*row
	if head < 0 {
		return nil
	}
	end := head - row
	if end < 0 {
		end = 0
	}

	d.Db.Model(model.ChainBlock{BlockNum: head}).
		Joins(fmt.Sprintf("JOIN (SELECT id,block_num from %s where block_num BETWEEN %d and %d order by block_num desc ) as t on %s.id=t.id",
			model.ChainBlock{BlockNum: head}.TableName(),
			end, head,
			model.ChainBlock{BlockNum: head}.TableName(),
		)).
		Order("block_num desc").Scan(&blocks)

	if head/model.SplitTableBlockNum != end/model.SplitTableBlockNum {
		var endBlocks []model.ChainBlock
		d.Db.Model(model.ChainBlock{BlockNum: blockNum - model.SplitTableBlockNum}).
			Joins(fmt.Sprintf("JOIN (SELECT id,block_num from %s order by block_num desc limit %d) as t on %s.id=t.id",
				model.ChainBlock{BlockNum: blockNum - model.SplitTableBlockNum}.TableName(),
				row-(head%model.SplitTableBlockNum+1),
				model.ChainBlock{BlockNum: blockNum - model.SplitTableBlockNum}.TableName(),
			)).
			Order("block_num desc").Scan(&endBlocks)
		blocks = append(blocks, endBlocks...)
	}

	return blocks
}

func (d *Dao) GetBlockByHash(c context.Context, hash string) *model.ChainBlock {
	var block model.ChainBlock
	blockNum, _ := d.GetCurrentBlockNum(context.TODO())
	for index := int(blockNum / uint64(model.SplitTableBlockNum)); index >= 0; index-- {
		query := d.Db.Model(&model.ChainBlock{BlockNum: index * model.SplitTableBlockNum}).Where("hash = ?", hash).Scan(&block)
		if query != nil && !query.RecordNotFound() {
			return &block
		}
	}
	return nil
}

func (d *Dao) GetBlockByNum(c context.Context, blockNum int) *model.ChainBlock {
	var block model.ChainBlock
	query := d.Db.Model(&model.ChainBlock{BlockNum: blockNum}).Where("block_num = ?", blockNum).Scan(&block)
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return nil
	}
	return &block
}

func (d *Dao) BlockAsJson(c context.Context, block *model.ChainBlock) *model.ChainBlockJson {
	bj := model.ChainBlockJson{
		BlockNum:        block.BlockNum,
		BlockTimestamp:  block.BlockTimestamp,
		Hash:            block.Hash,
		ParentHash:      block.ParentHash,
		StateRoot:       block.StateRoot,
		EventCount:      block.EventCount,
		ExtrinsicsCount: block.ExtrinsicsCount,
		ExtrinsicsRoot:  block.ExtrinsicsRoot,
		Extrinsics:      d.GetExtrinsicsByBlockNum(c, block.BlockNum),
		Events:          d.GetEventByBlockNum(c, block.BlockNum),
		Logs:            d.GetLogByBlockNum(c, block.BlockNum),
		Validator:       substrate.SS58Address(block.Validator),
		Finalized:       block.Finalized,
	}
	return &bj
}

func (d *Dao) UpdateEventAndExtrinsic(c context.Context, txn *GormDB, block *model.ChainBlock, eventCount, extrinsicsCount, blockTimestamp int, validator string, codecError bool, finalized bool) error {
	query := txn.Model(block).UpdateColumn(map[string]interface{}{
		"event_count":      eventCount,
		"extrinsics_count": extrinsicsCount,
		"block_timestamp":  blockTimestamp,
		"validator":        validator,
		"codec_error":      codecError,
		"hash":             block.Hash,
		"parent_hash":      block.ParentHash,
		"state_root":       block.StateRoot,
		"extrinsics_root":  block.ExtrinsicsRoot,
		"extrinsics":       block.Extrinsics,
		"event":            block.Event,
		"logs":             block.Logs,
		"finalized":        finalized,
	})
	d.delCacheBlock(context.TODO(), block)
	return query.Error
}

func (d *Dao) BlockAsSampleJson(c context.Context, block *model.ChainBlock) *model.SampleBlockJson {
	b := model.SampleBlockJson{
		BlockNum:        block.BlockNum,
		BlockTimestamp:  block.BlockTimestamp,
		Hash:            block.Hash,
		EventCount:      block.EventCount,
		ExtrinsicsCount: block.ExtrinsicsCount,
		Validator:       ss58.Encode(block.Validator, substrate.AddressType),
		Finalized:       block.Finalized,
	}
	return &b
}

func (d *Dao) GetNearBlock(c context.Context, blockNum int) *model.ChainBlock {
	var block model.ChainBlock
	query := d.Db.Model(&model.ChainBlock{BlockNum: blockNum}).Where("block_num > ?", blockNum).Order("block_num desc").First(&block)
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return nil
	}
	return &block
}

func (d *Dao) SetBlockFinalized(block *model.ChainBlock) {
	d.delCacheBlock(context.TODO(), block)
	d.Db.Model(block).UpdateColumn(model.ChainBlock{Finalized: true})
}

func (d *Dao) BlocksReverseByNum(c context.Context, blockNums []int) map[int]model.ChainBlock {
	var blocks []model.ChainBlock
	if len(blockNums) == 0 {
		return nil
	}
	sort.Ints(blockNums)
	lastNum := blockNums[len(blockNums)-1]
	for index := lastNum / model.SplitTableBlockNum; index >= 0; index-- {
		var tableData []model.ChainBlock
		query := d.Db.Model(model.ChainBlock{BlockNum: index * model.SplitTableBlockNum}).Where("block_num in (?)", blockNums).Scan(&tableData)
		if query == nil || query.Error != nil || query.RecordNotFound() {
			continue
		}
		blocks = append(blocks, tableData...)
	}

	toMap := make(map[int]model.ChainBlock)
	for _, block := range blocks {
		toMap[block.BlockNum] = block
	}

	return toMap
}
