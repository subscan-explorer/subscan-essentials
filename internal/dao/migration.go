package dao

import (
	"context"

	"github.com/itering/subscan/model"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
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
	models = append(models, model.RuntimeConstant{})
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
		tablesName = append(tablesName, d.GetModelTableName(m))
	}
	protectedTables = tablesName
	return models
}

func addIndex(db *gorm.DB, model interface{}, indexName string) {
	if !db.Migrator().HasIndex(model, indexName) {
		if err := db.Migrator().CreateIndex(model, indexName); err != nil {
			slog.Error("failed to add index", "indexName", indexName, "err", err, "model", model)
		}
	}
}

func (d *Dao) AddIndex(blockNum int) {
	db := d.db

	if blockNum == 0 {
		addIndex(db, model.RuntimeVersion{}, "spec_version")
	}

	blockModel := model.ChainBlock{BlockNum: blockNum}
	eventModel := model.ChainEvent{BlockNum: blockNum}
	extrinsicModel := model.ChainExtrinsic{BlockNum: blockNum}
	logModel := model.ChainLog{BlockNum: blockNum}

	addIndex(db, blockModel, "hash")

	addIndex(db, blockModel, "block_num")
	addIndex(db, blockModel, "codec_error")

	addIndex(db, extrinsicModel, "extrinsic_hash")
	addIndex(db, extrinsicModel, "extrinsic_index")
	addIndex(db, extrinsicModel, "block_num")
	addIndex(db, extrinsicModel, "is_signed")
	addIndex(db, extrinsicModel, "account_id")
	addIndex(db, extrinsicModel, "call_module")
	addIndex(db, extrinsicModel, "call_module_function")

	addIndex(db, eventModel, "block_num")
	addIndex(db, eventModel, "type")
	addIndex(db, eventModel, "event_index")
	addIndex(db, eventModel, "event_id")
	addIndex(db, eventModel, "module_id")
	addIndex(db, eventModel, "event_idx")

	addIndex(db, logModel, "block_num")
	addIndex(db, logModel, "log_index")
}
