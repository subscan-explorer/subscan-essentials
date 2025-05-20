package dao

import (
	"context"
	"github.com/itering/subscan/model"
	"gorm.io/gorm"
)

func (d *Dao) Migration() {
	db := d.db
	var blockNum uint
	if d.redis != nil {
		num, _ := d.GetFillBestBlockNum(context.TODO())
		blockNum = uint(num)
	}
	if d.DbDriver == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}
	_ = db.AutoMigrate(d.InternalTables(blockNum)...)
	for i := 0; uint(i) <= blockNum/model.SplitTableBlockNum; i++ {
		d.AddIndex(uint(i) * model.SplitTableBlockNum)
	}
}

func (d *Dao) InternalTables(blockNum uint) (models []interface{}) {
	models = append(models, model.RuntimeVersion{})
	for i := 0; uint(i) <= blockNum/model.SplitTableBlockNum; i++ {
		models = append(
			models,
			model.ChainBlock{BlockNum: blockNum},
			model.ChainEvent{BlockNum: blockNum},
			model.ChainExtrinsic{BlockNum: blockNum},
			model.ChainLog{BlockNum: blockNum})
	}
	var tablesName []string
	for _, m := range models {
		tablesName = append(tablesName, TableNameFromInterface(m, d.db))
	}
	protectedTables = tablesName
	return models
}

func (d *Dao) AddIndex(blockNum uint) {
	db := d.db
	blockModel := model.ChainBlock{BlockNum: blockNum}
	eventModel := model.ChainEvent{BlockNum: blockNum}
	extrinsicModel := model.ChainExtrinsic{BlockNum: blockNum}
	logModel := model.ChainLog{BlockNum: blockNum}
	_ = db.Scopes(d.TableNameFunc(&blockModel)).AutoMigrate(&blockModel)
	_ = db.Scopes(d.TableNameFunc(&eventModel)).AutoMigrate(&eventModel)
	_ = db.Scopes(d.TableNameFunc(&extrinsicModel)).AutoMigrate(&extrinsicModel)

	_ = db.Scopes(d.TableNameFunc(&logModel)).AutoMigrate(&logModel)
}

func (d *Dao) TableNameFunc(c interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(TableNameFromInterface(c, d.db))
	}
}

type Tabler interface {
	TableName() string
}

func TableNameFromInterface(c interface{}, db *gorm.DB) string {
	var tableName string
	if tabler, ok := c.(Tabler); ok {
		tableName = tabler.TableName()
	} else {
		stmt := &gorm.Statement{DB: db}
		_ = stmt.Parse(c)
		tableName = stmt.Schema.Table
	}
	return tableName
}
