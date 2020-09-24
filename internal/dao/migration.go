package dao

import (
	"context"
	"github.com/itering/subscan/model"
)

func (d *Dao) Migration() {
	db := d.db
	var blockNum int
	if d.redis != nil {
		blockNum, _ = d.GetFillBestBlockNum(context.TODO())
	}
	_ = db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(d.InternalTables(blockNum)...)

	for i := 0; i <= blockNum/model.SplitTableBlockNum; i++ {
		d.AddIndex(i * model.SplitTableBlockNum)
	}
}

func (d *Dao) InternalTables(blockNum int) (models []interface{}) {
	models = append(models, model.RuntimeVersion{})
	for i := 0; i <= blockNum/model.SplitTableBlockNum; i++ {
		models = append(
			models,
			model.ChainBlock{BlockNum: blockNum},
			model.ChainEvent{BlockNum: blockNum},
			model.ChainExtrinsic{BlockNum: blockNum},
			model.ChainLog{BlockNum: blockNum})
	}
	var tablesName []string
	for _, m := range models {
		tablesName = append(tablesName, d.db.Unscoped().NewScope(m).TableName())
	}
	protectedTables = tablesName
	return models
}

func (d *Dao) AddIndex(blockNum int) {
	db := d.db

	if blockNum == 0 {
		db.Model(model.RuntimeVersion{}).AddUniqueIndex("spec_version", "spec_version")
	}

	blockModel := model.ChainBlock{BlockNum: blockNum}
	eventModel := model.ChainEvent{BlockNum: blockNum}
	extrinsicModel := model.ChainExtrinsic{BlockNum: blockNum}
	logModel := model.ChainLog{BlockNum: blockNum}

	db.Model(blockModel).AddUniqueIndex("hash", "hash")
	db.Model(blockModel).AddUniqueIndex("block_num", "block_num")
	_ = db.Model(blockModel).AddIndex("codec_error", "codec_error")

	db.Model(extrinsicModel).AddIndex("extrinsic_hash", "extrinsic_hash")
	db.Model(extrinsicModel).AddUniqueIndex("extrinsic_index", "extrinsic_index")
	db.Model(extrinsicModel).AddIndex("block_num", "block_num")
	db.Model(extrinsicModel).AddIndex("is_signed", "is_signed")
	db.Model(extrinsicModel).AddIndex("account_id", "is_signed,account_id")
	db.Model(extrinsicModel).AddIndex("call_module", "call_module")
	db.Model(extrinsicModel).AddIndex("call_module_function", "call_module_function")

	db.Model(eventModel).AddIndex("block_num", "block_num")
	db.Model(eventModel).AddIndex("type", "type")
	db.Model(eventModel).AddIndex("event_index", "event_index")
	db.Model(eventModel).AddIndex("event_id", "event_id")
	db.Model(eventModel).AddIndex("module_id", "module_id")
	db.Model(eventModel).AddUniqueIndex("event_idx", "event_index", "event_idx")

	db.Model(logModel).AddUniqueIndex("log_index", "log_index")
	db.Model(logModel).AddIndex("block_num", "block_num")
}
