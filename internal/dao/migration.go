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

	count := int64(0)
	db.Model(&model.LastProcessedBlock{}).Count(&count)
	if count == 0 {
		db.Model(&model.LastProcessedBlock{}).Create(&model.LastProcessedBlock{Number: 0})
	}
}

func (d *Dao) InternalTables(blockNum int) (models []interface{}) {
	models = append(models, model.RuntimeVersion{})
	models = append(models, model.RuntimeConstant{})
	models = append(
		models,
		model.ChainBlock{},
		model.ChainEvent{},
		model.ChainExtrinsic{},
		model.ChainLog{},
		model.LastProcessedBlock{},
	)

	var tablesName []string
	for _, m := range models {
		tablesName = append(tablesName, d.GetModelTableName(m))
	}
	protectedTables = tablesName
	return models
}
