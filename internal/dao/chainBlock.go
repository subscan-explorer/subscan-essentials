package dao

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
)

var splitBlockTableCache = model.RedisKeyPrefix() + "split_block_table"

func (d *Dao) SplitBlockTable(blockNum uint) {
	ctx := context.Background()
	currentTableBlock := model.ChainBlock{BlockNum: blockNum}
	tableName := TableNameFromInterface(currentTableBlock, d.db)
	if s := d.redis.GetCacheString(ctx, splitBlockTableCache); s != tableName {
		if !d.db.Migrator().HasTable(tableName) {
			d.AddIndex(blockNum / model.SplitTableBlockNum * model.SplitTableBlockNum)
		}
		_ = d.redis.SetCache(ctx, splitBlockTableCache, tableName, 3600*24*30)
	}
}

// CreateBlock mysql db transaction
func (d *Dao) CreateBlock(ctx context.Context, txn *GormDB, cb *model.ChainBlock) (err error) {
	query := txn.WithContext(ctx).Scopes().Scopes(d.TableNameFunc(cb), model.IgnoreDuplicate).Create(cb)
	// Check if you need to create a new table(block, extrinsic, event, log) after created block
	if maxTableBlockNum < cb.BlockNum+model.SplitTableBlockNum {
		tableName := model.TableNameFromInterface(model.ChainBlock{BlockNum: cb.BlockNum + model.SplitTableBlockNum}, d.db)
		if !d.db.Migrator().HasTable(tableName) {
			go func() {
				db := d.db
				if d.DbDriver == "mysql" {
					db.Set("gorm:table_options", "ENGINE=InnoDB")
				}
				d.AddIndex(cb.BlockNum + model.SplitTableBlockNum)
			}()
		}
		maxTableBlockNum = cb.BlockNum + model.SplitTableBlockNum
	}

	return query.Error
}

func (d *Dao) SaveFillAlreadyBlockNum(c context.Context, blockNum int) (err error) {
	conn, _ := d.redis.Redis().GetContext(c)
	defer func() {
		conn.Close()
	}()
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

// GetBlockListCursor implements cursor pagination on blocks across split tables.
func (d *Dao) GetBlockListCursor(ctx context.Context, limit int, before, after uint) (list []model.ChainBlock, hasPrev, hasNext bool) {
	fetch := limit + 1
	// determine max split-table index from best block
	best, _ := d.GetFillBestBlockNum(context.TODO())
	maxIdx := uint(best) / model.SplitTableBlockNum

	// next page: block_num < after, walk tables downward
	if after > 0 {
		startIdx := int(after / model.SplitTableBlockNum)
		for idx := startIdx; idx >= 0 && len(list) < fetch; idx-- {
			var tableData []model.ChainBlock
			q := d.db.WithContext(ctx).Scopes(d.TableNameFunc(&model.ChainBlock{BlockNum: uint(idx) * model.SplitTableBlockNum}))
			if idx == startIdx {
				q = q.Where("block_num < ?", after)
			}
			q = q.Order("block_num desc").Limit(fetch - len(list))
			if err := q.Find(&tableData).Error; err != nil {
				continue
			}
			list = append(list, tableData...)
		}
		hasNext = len(list) > limit
		if hasNext {
			list = list[:limit]
		}
		hasPrev = true
		return list, hasPrev, hasNext
	}

	// previous page: block_num > before, walk tables upward (asc, then reverse)
	if before > 0 {
		startIdx := int(before / model.SplitTableBlockNum)
		for idx := startIdx; uint(idx) <= maxIdx && len(list) < fetch; idx++ {
			var tableData []model.ChainBlock
			q := d.db.WithContext(ctx).Scopes(d.TableNameFunc(&model.ChainBlock{BlockNum: uint(idx) * model.SplitTableBlockNum}))
			if idx == startIdx {
				q = q.Where("block_num > ?", before)
			}
			q = q.Order("block_num asc").Limit(fetch - len(list))
			if err := q.Find(&tableData).Error; err != nil {
				continue
			}
			list = append(list, tableData...)
		}
		hasPrev = len(list) > limit
		if hasPrev {
			list = list[:limit]
		}
		// reverse to keep response in desc order
		for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
		hasNext = true
		return list, hasPrev, hasNext
	}

	// first page: walk from newest table downward
	for idx := maxIdx; len(list) < fetch; idx-- {
		var tableData []model.ChainBlock
		q := d.db.WithContext(ctx).Scopes(d.TableNameFunc(&model.ChainBlock{BlockNum: uint(idx) * model.SplitTableBlockNum}))
		q = q.Order("block_num desc").Limit(fetch - len(list))
		if err := q.Find(&tableData).Error; err != nil {
			break
		}
		list = append(list, tableData...)
		if idx == 0 {
			break
		}
	}
	hasNext = len(list) > limit
	if hasNext {
		list = list[:limit]
	}
	return list, false, hasNext
}

func (d *Dao) GetBlockByHash(c context.Context, hash string) *model.ChainBlock {
	var block model.ChainBlock
	blockNum, _ := d.GetBestBlockNum(c)
	for index := int(blockNum / uint64(model.SplitTableBlockNum)); index >= 0; index-- {
		query := d.db.Scopes(model.TableNameFunc(model.ChainBlock{BlockNum: uint(index) * (model.SplitTableBlockNum)})).Where("hash = ?", hash).Scan(&block)
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
	query := d.db.Scopes(model.TableNameFunc(model.ChainBlock{BlockNum: blockNum})).Where("block_num > ?", blockNum).Order("block_num desc").Find(&block)
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
		query := d.db.Scopes(model.TableNameFunc(model.ChainBlock{BlockNum: uint(index) * model.SplitTableBlockNum})).Where("block_num in (?)", blockNums).Scan(&tableData)
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

func (d *Dao) GetBlocksByNums(c context.Context, blockNums []uint, columns string) (blocks []*model.ChainBlock) {
	if len(blockNums) == 0 {
		return nil
	}
	idxMap := make(map[uint][]uint)
	for _, num := range blockNums {
		idxMap[num/model.SplitTableBlockNum] = append(idxMap[num/model.SplitTableBlockNum], num)
	}
	for index, ids := range idxMap {
		var tableData []*model.ChainBlock
		if err := d.db.WithContext(c).Select(columns).Scopes(model.TableNameFunc(&model.ChainBlock{BlockNum: index * model.SplitTableBlockNum})).Where("block_num IN ?", ids).Find(&tableData).Error; err != nil {
			continue
		}
		blocks = append(blocks, tableData...)
	}
	return
}
