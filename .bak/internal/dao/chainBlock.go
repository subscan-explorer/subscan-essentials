package dao

import (
	"context"
	"encoding/json"
	"github.com/bilibili/kratos/pkg/cache/redis"
	"subscan-end/internal/model"
	"subscan-end/libs/substrate"
	"subscan-end/utiles"
	"subscan-end/utiles/ss58"
)

const (
	PubTopicBlock = "block_new"
)

func (d *Dao) CreateBlock(c context.Context, txn *GormDB, hash string, block *substrate.Block, event, decodeEvent, decodeExtrinsics, log, decodeLog, validator string, eventCount, extrinsicsCount, blockTimestamp int, codecError bool) (err error) {
	extrinsicB, _ := json.Marshal(block.Extrinsics)
	cb := model.ChainBlock{
		Hash:             hash,
		BlockNum:         utiles.StringToInt(utiles.HexToNumStr(block.Header.Number)),
		BlockTimestamp:   blockTimestamp,
		ParentHash:       block.Header.ParentHash,
		StateRoot:        block.Header.StateRoot,
		ExtrinsicsRoot:   block.Header.ExtrinsicsRoot,
		Logs:             log,
		DecodeLogs:       decodeLog,
		Extrinsics:       string(extrinsicB),
		Event:            event,
		DecodeEvent:      decodeEvent,
		DecodeExtrinsics: decodeExtrinsics,
		SpecVersion:      substrate.CurrentRuntimeSpecVersion,
		ExtrinsicsCount:  extrinsicsCount,
		EventCount:       eventCount,
		Validator:        validator,
		CodecError:       codecError,
	}
	query := txn.Create(&cb)
	d.BroadCastToChanel(c, PubTopicBlock, cb)
	return query.Error
}

func (d *Dao) SaveFillAlreadyBlockNum(c context.Context, blockNum int) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if num, _ := redis.Int(conn.Do("GET", RedisFillAlreadyBlockNum)); blockNum > num {
		_, err = conn.Do("SET", RedisFillAlreadyBlockNum, blockNum)
	}
	return
}

func (d *Dao) GetFillAlreadyBlockNum(c context.Context) (num int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	num, err = redis.Int(conn.Do("GET", RedisFillAlreadyBlockNum))
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
	d.db.Model(&model.ChainBlock{}).Where("block_num BETWEEN ? AND ?", start, end).Order("block_num asc").Pluck("block_num", &blockNums)
	return blockNums
}

func (d *Dao) GetBlockList(c context.Context, page, row int) []model.ChainBlock {
	var blocks []model.ChainBlock
	query := d.db.Model(&model.ChainBlock{}).Offset(page * row).Limit(row).Order("block_num desc").Find(&blocks)
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return nil
	}
	return blocks
}

func (d *Dao) GetBlockByHash(c context.Context, hash string) *model.ChainBlock {
	var block model.ChainBlock
	query := d.db.Model(&model.ChainBlock{}).Where("hash = ?", hash).First(&block)
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return nil
	}
	return &block
}

func (d *Dao) GetBlockByNum(c context.Context, blockNum int) *model.ChainBlock {
	var block model.ChainBlock
	query := d.db.Model(&model.ChainBlock{}).Where("block_num = ?", blockNum).First(&block)
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
		Validator:       ss58.Encode(block.Validator),
	}
	return &bj
}

func (d *Dao) GetAllBlocksNeedFix(c context.Context) *[]model.ChainBlock {
	var blocks []model.ChainBlock
	query := d.db.Model(&model.ChainBlock{}).Where("codec_error=?", "1").Find(&blocks)
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return nil
	}
	return &blocks
}

func (d *Dao) UpdateEventAndExtrinsic(c context.Context, txn *GormDB, block *model.ChainBlock, decodeEvent, decodeExtrinsics, decodeLog string, eventCount, extrinsicsCount, blockTimestamp int, validator string) error {
	query := txn.Model(block).UpdateColumn(map[string]interface{}{
		"decode_event":      decodeEvent,
		"decode_extrinsics": decodeExtrinsics,
		"decode_logs":       decodeLog,
		"event_count":       eventCount,
		"extrinsics_count":  extrinsicsCount,
		"block_timestamp":   blockTimestamp,
		"validator":         validator,
		"codec_error":       false,
	})
	return query.Error
}

func (d *Dao) BlockAsSampleJson(c context.Context, block *model.ChainBlock) *model.SampleBlockJson {
	b := model.SampleBlockJson{
		BlockNum:        block.BlockNum,
		BlockTimestamp:  block.BlockTimestamp,
		Hash:            block.Hash,
		EventCount:      block.EventCount,
		ExtrinsicsCount: block.ExtrinsicsCount,
		Validator:       ss58.Encode(block.Validator),
	}
	return &b
}

func (d *Dao) GetNodeNameByControllerAddress(address string) string {
	var validator model.ValidatorInfo
	query := d.db.Model(model.ValidatorInfo{}).Where("validator_controller = ?", address).Find(&validator)
	if query.RecordNotFound() {
		return ""
	} else {
		return validator.NodeName
	}
}
