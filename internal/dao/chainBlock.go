package dao

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
)

// CreateBlock mysql db transaction
func (d *Dao) CreateBlock(txn *GormDB, cb *model.ChainBlock) (err error) {
	query := txn.Scopes().Scopes(d.TableNameFunc(cb), model.IgnoreDuplicate).Create(cb)

	// Check if you need to create a new table(block, extrinsic, event, log) after created block
	if maxTableBlockNum < cb.BlockNum+model.SplitTableBlockNum {
		if !d.db.Migrator().HasTable(model.ChainBlock{BlockNum: cb.BlockNum + model.SplitTableBlockNum}) {
			go func() {
				_ = d.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
					d.InternalTables(cb.BlockNum + model.SplitTableBlockNum)...)
				d.AddIndex(cb.BlockNum + model.SplitTableBlockNum)
			}()
		}
		maxTableBlockNum = cb.BlockNum + model.SplitTableBlockNum
	}

	return query.Error
}

func (d *Dao) SaveFillAlreadyBlockNum(c context.Context, blockNum int) (err error) {
	conn, _ := d.redis.Redis().GetContext(c)
	defer conn.Close()
	if num, _ := redis.Int(conn.Do("GET", RedisFillAlreadyBlockNum)); blockNum > num {
		_, err = conn.Do("SET", RedisFillAlreadyBlockNum, blockNum)
	}
	return
}

func (d *Dao) SaveFillAlreadyFinalizedBlockNum(c context.Context, blockNum int) (err error) {
	conn, _ := d.redis.Redis().GetContext(c)
	defer func() {
		conn.Close()
	}()

	if num, _ := redis.Int(conn.Do("GET", RedisFillFinalizedBlockNum)); blockNum > num {
		_, err = conn.Do("SET", RedisFillFinalizedBlockNum, blockNum)
	}
	return
}

func (d *Dao) GetFillBestBlockNum(c context.Context) (num int, err error) {
	conn, _ := d.redis.Redis().GetContext(c)
	defer conn.Close()
	num, err = redis.Int(conn.Do("GET", RedisFillAlreadyBlockNum))
	return
}

func (d *Dao) GetFillFinalizedBlockNum(c context.Context) (num int, err error) {
	conn, _ := d.redis.Redis().GetContext(c)
	defer conn.Close()
	num, err = redis.Int(conn.Do("GET", RedisFillFinalizedBlockNum))
	return
}

func (d *Dao) GetBlockList(ctx context.Context, page, row int) []model.ChainBlock {
	var blocks []model.ChainBlock
	blockNum, _ := d.GetFillBestBlockNum(ctx)
	head := blockNum - page*row
	if head < 0 {
		return nil
	}
	end := head - row + 1
	if end < 0 {
		end = 0
	}
	bestNum := uint(blockNum)
	headBlock := uint(head)
	endBlock := uint(end)
	d.db.Scopes(model.TableNameFunc(model.ChainBlock{BlockNum: headBlock})).
		Joins(fmt.Sprintf("JOIN (SELECT id from %s where block_num BETWEEN %d and %d order by block_num desc ) as t on %s.id=t.id",
			model.ChainBlock{BlockNum: headBlock}.TableName(),
			end, head,
			model.ChainBlock{BlockNum: headBlock}.TableName(),
		)).
		Order("block_num desc").Scan(&blocks)

	if headBlock/model.SplitTableBlockNum != endBlock/model.SplitTableBlockNum {
		var endBlocks []model.ChainBlock
		d.db.Scopes(model.TableNameFunc(model.ChainBlock{BlockNum: bestNum - model.SplitTableBlockNum})).
			Joins(fmt.Sprintf("JOIN (SELECT id from %s order by block_num desc limit %d) as t on %s.id=t.id",
				model.ChainBlock{BlockNum: bestNum - model.SplitTableBlockNum}.TableName(),
				uint(row)-(headBlock%model.SplitTableBlockNum),
				model.ChainBlock{BlockNum: bestNum - model.SplitTableBlockNum}.TableName(),
			)).
			Order("block_num desc").Scan(&endBlocks)
		blocks = append(blocks, endBlocks...)
	}

	return blocks
}

func (d *Dao) GetBlockByHash(c context.Context, hash string) *model.ChainBlock {
	var block model.ChainBlock
	blockNum, _ := d.GetBestBlockNum(c)
	for index := int(blockNum / uint64(model.SplitTableBlockNum)); index >= 0; index-- {
		query := d.db.Model(&model.ChainBlock{BlockNum: uint(index) * (model.SplitTableBlockNum)}).Where("hash = ?", hash).Scan(&block)
		if query != nil && query.Error == nil {
			return &block
		}
	}
	return nil
}

func (d *Dao) GetBlockByNum(ctx context.Context, blockNum uint) *model.ChainBlock {
	var block model.ChainBlock
	query := d.db.WithContext(ctx).Scopes(model.TableNameFunc(&model.ChainBlock{BlockNum: blockNum})).Where("block_num = ?", blockNum).Find(&block)
	if query == nil || query.Error != nil {
		return nil
	}
	return &block
}

func (d *Dao) BlockAsJson(_ context.Context, block *model.ChainBlock) *model.ChainBlockJson {
	bj := model.ChainBlockJson{
		BlockNum:        block.BlockNum,
		BlockTimestamp:  block.BlockTimestamp,
		Hash:            block.Hash,
		ParentHash:      block.ParentHash,
		StateRoot:       block.StateRoot,
		EventCount:      block.EventCount,
		ExtrinsicsCount: block.ExtrinsicsCount,
		ExtrinsicsRoot:  block.ExtrinsicsRoot,
		Validator:       address.Encode(block.Validator),
		Finalized:       block.Finalized,
		SpecVersion:     block.SpecVersion,
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
		"finalized":        finalized,
	})
	return query.Error
}

func (d *Dao) GetNearBlock(blockNum uint) *model.ChainBlock {
	var block model.ChainBlock
	query := d.db.Model(&model.ChainBlock{BlockNum: blockNum}).Where("block_num > ?", blockNum).Order("block_num desc").Scan(&block)
	if query == nil || query.Error != nil {
		return nil
	}
	return &block
}

func (d *Dao) BlocksReverseByNum(blockNums []uint) map[uint]model.ChainBlock {
	var blocks []model.ChainBlock
	if len(blockNums) == 0 {
		return nil
	}
	util.SortUintSlice(blockNums)
	lastNum := blockNums[len(blockNums)-1]
	for index := int(lastNum / model.SplitTableBlockNum); index >= 0; index-- {
		var tableData []model.ChainBlock
		query := d.db.Model(model.ChainBlock{BlockNum: uint(index) * model.SplitTableBlockNum}).Where("block_num in (?)", blockNums).Scan(&tableData)
		if query == nil || query.Error != nil {
			continue
		}
		blocks = append(blocks, tableData...)
	}

	toMap := make(map[uint]model.ChainBlock)
	for _, block := range blocks {
		toMap[block.BlockNum] = block
	}

	return toMap
}

func (d *Dao) GetBlockNumArr(c context.Context, start, end uint) []int {
	var blockNums []int
	d.db.WithContext(c).Scopes(d.TableNameFunc(&model.ChainBlock{BlockNum: end})).Where("block_num BETWEEN ? AND ?", start, end).Order("block_num asc").Pluck("block_num", &blockNums)
	return blockNums
}
