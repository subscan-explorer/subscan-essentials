package dao

import (
	"context"
	"sort"

	"github.com/gomodule/redigo/redis"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util/address"
)

// CreateBlock, mysql db transaction
func (d *Dao) CreateBlock(txn *GormDB, cb *model.ChainBlock) (err error) {
	query := txn.Create(cb)
	return query.Error
}

func (d *Dao) SaveFillAlreadyBlockNum(c context.Context, blockNum int) (err error) {
	conn, _ := d.redis.GetContext(c)
	defer conn.Close()
	if num, _ := redis.Int(conn.Do("GET", RedisFillAlreadyBlockNum)); blockNum > num {
		_, err = conn.Do("SET", RedisFillAlreadyBlockNum, blockNum)
	}
	return
}

func (d *Dao) SaveFillAlreadyFinalizedBlockNum(c context.Context, blockNum int) (err error) {
	conn, _ := d.redis.GetContext(c)
	defer func() {
		conn.Close()
	}()

	if num, _ := redis.Int(conn.Do("GET", RedisFillFinalizedBlockNum)); blockNum > num {
		_, err = conn.Do("SET", RedisFillFinalizedBlockNum, blockNum)
	}
	return
}

func (d *ReadOnlyDao) GetFillBestBlockNum(c context.Context) (num int, err error) {
	conn, _ := d.redis.GetContext(c)
	defer conn.Close()
	num, err = redis.Int(conn.Do("GET", RedisFillAlreadyBlockNum))
	return
}

func (d *ReadOnlyDao) GetFillFinalizedBlockNum(c context.Context) (num int, err error) {
	conn, _ := d.redis.GetContext(c)
	defer conn.Close()
	num, err = redis.Int(conn.Do("GET", RedisFillFinalizedBlockNum))
	return
}

func (d *ReadOnlyDao) GetBlockList(page, row int) []model.ChainBlock {
	var blocks []model.ChainBlock
	blockNum, _ := d.GetFillBestBlockNum(context.TODO())
	head := blockNum - page*row
	if head < 0 {
		return nil
	}
	end := head - row
	if end < 0 {
		end = 0
	}

	d.db.Model(model.ChainBlock{BlockNum: head}).
		Select("id", "block_num").
		Where("block_num BETWEEN ? AND ?", end, head).
		Order("block_num desc").Scan(&blocks)

	return blocks
}

func (d *ReadOnlyDao) GetBlockByHash(c context.Context, hash string) *model.ChainBlock {
	var block model.ChainBlock
	query := d.db.Model(&block).Where("hash = ?", hash).Scan(&block)
	if query != nil && !RecordNotFound(query) {
		return &block
	}
	return nil
}

func (d *ReadOnlyDao) GetBlockByNum(blockNum int) *model.ChainBlock {
	res, _ := findOne[model.ChainBlock](d, "*", where("block_num = ?", blockNum), nil)
	return res
}

func (d *ReadOnlyDao) BlockAsJson(c context.Context, block *model.ChainBlock) *model.ChainBlockJson {
	bj := model.ChainBlockJson{
		BlockNum:        block.BlockNum,
		BlockTimestamp:  block.BlockTimestamp,
		Hash:            block.Hash,
		ParentHash:      block.ParentHash,
		StateRoot:       block.StateRoot,
		EventCount:      block.EventCount,
		ExtrinsicsCount: block.ExtrinsicsCount,
		ExtrinsicsRoot:  block.ExtrinsicsRoot,
		Extrinsics:      d.GetExtrinsicsByBlockNum(block.BlockNum),
		Events:          d.GetEventByBlockNum(block.BlockNum),
		Logs:            d.GetLogByBlockNum(block.BlockNum),
		Validator:       address.SS58AddressFromHex(block.Validator),
		Finalized:       block.Finalized,
	}
	return &bj
}

func (d *Dao) UpdateEventAndExtrinsic(txn *GormDB, block *model.ChainBlock, eventCount, extrinsicsCount, blockTimestamp int, validator string, codecError bool, finalized bool) error {
	query := txn.Where("block_num = ?", block.BlockNum).Model(block).UpdateColumns(map[string]interface{}{
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
	return query.Error
}

func (d *ReadOnlyDao) GetNearBlock(blockNum int) *model.ChainBlock {
	var block model.ChainBlock
	query := d.db.Model(&model.ChainBlock{BlockNum: blockNum}).Where("block_num > ?", blockNum).Order("block_num desc").Scan(&block)
	if query == nil || query.Error != nil || RecordNotFound(query) {
		return nil
	}
	return &block
}

func (d *Dao) SetBlockFinalized(block *model.ChainBlock) {
	d.db.Model(block).Updates(model.ChainBlock{Finalized: true})
}

func (d *ReadOnlyDao) BlocksReverseByNum(blockNums []int) map[int]model.ChainBlock {
	if len(blockNums) == 0 {
		return nil
	}
	sort.Ints(blockNums)
	var blocks []model.ChainBlock
	query := d.db.Model(&model.ChainBlock{}).Where("block_num in (?)", blockNums).Scan(&blocks)

	if query == nil || query.Error != nil || RecordNotFound(query) {
		return nil
	}

	toMap := make(map[int]model.ChainBlock)
	for _, block := range blocks {
		toMap[block.BlockNum] = block
	}

	return toMap
}

func (d *ReadOnlyDao) GetBlockNumArr(start, end int) []int {
	var blockNums []int
	d.db.Model(model.ChainBlock{BlockNum: end}).Where("block_num BETWEEN ? AND ?", start, end).Order("block_num asc").Pluck("block_num", &blockNums)
	return blockNums
}
