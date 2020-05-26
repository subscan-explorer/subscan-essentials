package dao

import (
	"context"
	"github.com/itering/subscan/internal/model"
)

func (d *Dao) Migration() {
	db := d.db
	d.splitTableMigrate()
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.ChainAccount{},
		&model.RuntimeVersion{},
		&model.ExtrinsicError{},
	)
	db.Model(model.ChainAccount{}).AddUniqueIndex("address", "address")
	db.Model(model.ChainAccount{}).AddIndex("account_index", "account_index")
	db.Model(model.RuntimeVersion{}).AddUniqueIndex("spec_version", "spec_version")
	db.Model(model.RuntimeVersion{}).ModifyColumn("modules", "text")

	db.Model(model.ExtrinsicError{}).AddUniqueIndex("extrinsic_hash", "extrinsic_hash")
}

func (d *Dao) splitTableMigrate() {
	var blockNum = 0
	if d.redis != nil {
		blockNum, _ = d.GetFillAlreadyBlockNum(context.TODO())
	}
	for i := 0; i <= blockNum/model.SplitTableBlockNum; i++ {
		d.blockMigrate(i * model.SplitTableBlockNum)
	}
}

func (d *Dao) blockMigrate(blockNum int) {
	blockModel := model.ChainBlock{BlockNum: blockNum}
	eventModel := model.ChainEvent{BlockNum: blockNum}
	extrinsicModel := model.ChainExtrinsic{BlockNum: blockNum}
	transactionModel := model.ChainTransaction{BlockNum: blockNum}
	logModel := model.ChainLog{BlockNum: blockNum}

	db := d.db
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		blockModel,
		&eventModel,
		&extrinsicModel,
		&transactionModel,
		&logModel,
	)

	db.Model(blockModel).AddUniqueIndex("hash", "hash")
	db.Model(blockModel).AddUniqueIndex("block_num", "block_num")
	db.Model(blockModel).AddIndex("codec_error", "codec_error")

	db.Model(transactionModel).AddIndex("from_hex", "from_hex")
	db.Model(transactionModel).AddIndex("destination", "destination")
	db.Model(transactionModel).AddIndex("call_module_function", "call_module_function")
	db.Model(transactionModel).AddIndex("block_num", "block_num")
	db.Model(transactionModel).AddUniqueIndex("hash", "hash")
	db.Model(transactionModel).AddUniqueIndex("extrinsic_index", "extrinsic_index")
	db.Model(extrinsicModel).AddIndex("extrinsic_hash", "extrinsic_hash")
	db.Model(extrinsicModel).AddUniqueIndex("extrinsic_index", "extrinsic_index")
	db.Model(extrinsicModel).AddIndex("block_num", "block_num")
	db.Model(extrinsicModel).AddIndex("is_signed", "is_signed")
	db.Model(extrinsicModel).AddIndex("account_id", "is_signed,account_id")
	db.Model(eventModel).AddIndex("block_num", "block_num")
	db.Model(eventModel).AddIndex("type", "type")
	db.Model(eventModel).AddIndex("event_index", "event_index")
	db.Model(eventModel).AddIndex("event_id", "event_id")
	db.Model(eventModel).AddIndex("module_id", "module_id")
	db.Model(eventModel).AddUniqueIndex("event_idx", "event_index", "event_idx")
	db.Model(logModel).AddUniqueIndex("log_index", "log_index")
	db.Model(logModel).AddIndex("block_num", "block_num")

	db.Model(transactionModel).AddIndex("call_module", "call_module")
	db.Model(extrinsicModel).AddIndex("call_module", "call_module")
	db.Model(extrinsicModel).AddIndex("call_module_function", "call_module_function")
}
