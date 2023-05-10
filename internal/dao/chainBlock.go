package dao

import (
	"context"
	"fmt"
	"sort"

	"github.com/gomodule/redigo/redis"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util/address"
)

// CreateBlock, mysql db transaction
// Check if you need to create a new table(block, extrinsic, event, log ) after created
func (d *Dao) CreateBlock(txn *GormDB, cb *model.ChainBlock) (err error) {
	query := txn.Save(cb)
	if !d.db.Migrator().HasTable(model.ChainBlock{BlockNum: cb.BlockNum + model.SplitTableBlockNum}) {
		go func() {
			_ = d.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
				d.InternalTables(cb.BlockNum + model.SplitTableBlockNum)...)
			d.AddIndex(cb.BlockNum + model.SplitTableBlockNum)
		}()
	}
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

func (d *Dao) GetFillBestBlockNum(c context.Context) (num int, err error) {
	conn, _ := d.redis.GetContext(c)
	defer conn.Close()
	num, err = redis.Int(conn.Do("GET", RedisFillAlreadyBlockNum))
	return
}

func (d *Dao) GetFillFinalizedBlockNum(c context.Context) (num int, err error) {
	conn, _ := d.redis.GetContext(c)
	defer conn.Close()
	num, err = redis.Int(conn.Do("GET", RedisFillFinalizedBlockNum))
	return
}

func (d *Dao) GetBlockList(page, row int) []model.ChainBlock {
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
		Joins(fmt.Sprintf("JOIN (SELECT id,block_num from %s where block_num BETWEEN %d and %d order by block_num desc ) as t on %s.id=t.id",
			model.ChainBlock{BlockNum: head}.TableName(),
			end, head,
			model.ChainBlock{BlockNum: head}.TableName(),
		)).
		Order("block_num desc").Scan(&blocks)

	if head/model.SplitTableBlockNum != end/model.SplitTableBlockNum {
		var endBlocks []model.ChainBlock
		d.db.Model(model.ChainBlock{BlockNum: blockNum - model.SplitTableBlockNum}).
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
	blockNum, _ := d.GetBestBlockNum(context.TODO())
	for index := int(blockNum / uint64(model.SplitTableBlockNum)); index >= 0; index-- {
		query := d.db.Model(&model.ChainBlock{BlockNum: index * model.SplitTableBlockNum}).Where("hash = ?", hash).Scan(&block)
		if query != nil && !RecordNotFound(query) {
			return &block
		}
	}
	return nil
}

func (d *Dao) GetBlockByNum(blockNum int) *model.ChainBlock {
	var block model.ChainBlock
	query := d.db.Model(&model.ChainBlock{BlockNum: blockNum}).Where("block_num = ?", blockNum).Scan(&block)
	if query == nil || query.Error != nil || RecordNotFound(query) {
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

func (d *Dao) GetNearBlock(blockNum int) *model.ChainBlock {
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

func (d *Dao) BlocksReverseByNum(blockNums []int) map[int]model.ChainBlock {
	var blocks []model.ChainBlock
	if len(blockNums) == 0 {
		return nil
	}
	sort.Ints(blockNums)
	lastNum := blockNums[len(blockNums)-1]
	for index := lastNum / model.SplitTableBlockNum; index >= 0; index-- {
		var tableData []model.ChainBlock
		query := d.db.Model(model.ChainBlock{BlockNum: index * model.SplitTableBlockNum}).Where("block_num in (?)", blockNums).Scan(&tableData)
		if query == nil || query.Error != nil || RecordNotFound(query) {
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

func (d *Dao) GetBlockNumArr(start, end int) []int {
	var blockNums []int
	d.db.Model(model.ChainBlock{BlockNum: end}).Where("block_num BETWEEN ? AND ?", start, end).Order("block_num asc").Pluck("block_num", &blockNums)
	return blockNums
}
