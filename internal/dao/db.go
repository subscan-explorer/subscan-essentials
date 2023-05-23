package dao

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/itering/subscan/configs"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins/storage"
	"github.com/itering/substrate-api-rpc/websocket"

	"github.com/itering/subscan/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbStorage struct {
	db     *gorm.DB
	Prefix string
	dao    IReadOnlyDao
}

func (d *DbStorage) SetPrefix(prefix string) {
	d.Prefix = prefix
}

func (d *DbStorage) GetPrefix() string {
	return d.Prefix
}

var protectedTables []string

func (d *DbStorage) SpecialMetadata(spec int) string {
	var raw model.RuntimeVersion
	if query := d.db.Where("spec_version = ?", spec).First(&raw); RecordNotFound(query) {
		return ""
	}
	return raw.RawData
}

func (d *DbStorage) GetModelTableName(model interface{}) string {
	stmt := &gorm.Statement{DB: d.db}
	if err := stmt.Parse(model); err != nil {
		panic(fmt.Sprintf("get model table name error: %v", err))
	}
	return stmt.Schema.Table
}

func (d *DbStorage) GetRuntimeConstant(moduleName, constantName string) *storage.RuntimeConstant {
	return d.dao.GetRuntimeConstantLatest(moduleName, constantName).AsPlugin()
}

func (d *DbStorage) checkProtected(model interface{}) error {
	if util.StringInSlice(d.GetModelTableName(model), protectedTables) {
		return errors.New("protected tables")
	}
	return nil
}

func (d *DbStorage) RPCPool() *websocket.PoolConn {
	conn, _ := websocket.Init()
	return conn
}

func (d *DbStorage) getPluginPrefixTableName(instant interface{}) string {
	tableName := d.GetModelTableName(instant)
	if util.StringInSlice(tableName, protectedTables) {
		return tableName
	}
	return fmt.Sprintf("%s_%s", d.GetPrefix(), tableName)
}

func (d *DbStorage) FindBy(record interface{}, query interface{}, option *storage.Option) (int, bool) {
	var count int64
	tx := d.db

	// where
	if reflect.ValueOf(query).IsValid() {
		tx = tx.Where(query)
	}

	// plugin prefix table
	if option != nil && option.PluginPrefix != "" {
		tx = tx.Table(fmt.Sprintf("%s_%s", option.PluginPrefix, d.GetModelTableName(record)))
		if (option.Page > 0) && (option.PageSize > 0) {
			tx = tx.Limit(option.PageSize).Offset((option.Page - 1) * option.PageSize)
		}
		if option.Order != "" {
			tx = tx.Order(option.Order)
		}
	}
	// rows count
	tx.Count(&count)

	// pagination
	if option != nil {
		// default page limit 1000
		if option.PageSize == 0 {
			option.PageSize = 1000
		}
		tx = tx.Offset(option.Page * option.PageSize).Limit(option.PageSize)
	}

	tx = tx.Find(record)
	return int(count), errors.Is(tx.Error, gorm.ErrRecordNotFound)
}

func (d *DbStorage) AutoMigration(model interface{}) error {
	if d.checkProtected(model) == nil {
		err := d.db.Table(d.getPluginPrefixTableName(model)).Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model)
		return err
	}
	return nil
}

func (d *DbStorage) AddIndex(model interface{}, indexName string, columns ...string) error {
	if d.checkProtected(model) == nil {
		err := d.db.Table(d.getPluginPrefixTableName(model)).Migrator().CreateIndex(model, indexName)
		return err
	}
	return nil
}

func (d *DbStorage) AddUniqueIndex(model interface{}, indexName string, columns ...string) error {
	return d.AddIndex(model, indexName, columns...)
}

func (d *DbStorage) Query(model interface{}) *gorm.DB {
	return d.db.Table(d.getPluginPrefixTableName(model))
}

func (d *DbStorage) Create(record interface{}) error {
	if err := d.checkProtected(record); err == nil {
		tx := d.db.Table(d.getPluginPrefixTableName(record)).Save(record)
		return tx.Error
	} else {
		return err
	}
}

func (d *DbStorage) Update(model interface{}, query interface{}, attr map[string]interface{}) error {
	if err := d.checkProtected(model); err == nil {
		tx := d.db.Table(d.getPluginPrefixTableName(model)).Where(query).Updates(attr)
		return tx.Error
	} else {
		return err
	}
}

func (d *DbStorage) Delete(model interface{}, query interface{}) error {
	if err := d.checkProtected(model); err == nil {
		tx := d.db.Table(d.getPluginPrefixTableName(model)).Where(query).Delete(model)
		return tx.Error
	} else {
		return err
	}
}

// // logs
// type ormLog struct{}

// func (l ormLog) Print(v ...interface{}) {
// 	slog.Debug(strings.Repeat("%v ", len(v)), v...)
// }

// db
type GormDB struct {
	*gorm.DB
	gdbDone bool
}

func (d *Dao) DbCommit(c *GormDB) {
	if c.gdbDone {
		return
	}
	tx := c.Commit()
	c.gdbDone = true
	if err := tx.Error; err != nil && err != sql.ErrTxDone {
		fmt.Println("Fatal error DbCommit", err)
	}
}

func (d *Dao) DbRollback(c *GormDB) {
	if c.gdbDone {
		return
	}
	tx := c.Rollback()
	c.gdbDone = true
	if err := tx.Error; err != nil && err != sql.ErrTxDone {
		fmt.Println("Fatal error DbRollback", err)
	}
}

// dao funcs
func (d *Dao) DbBegin() *GormDB {
	txn := d.db.Begin()
	if txn.Error != nil {
		panic(txn.Error)
	}
	return &GormDB{txn, false}
}

// private funcs
func newDb() (db *gorm.DB) {
	var err error

	db, err = gorm.Open(mysql.Open(configs.Boot.Database.DSN))
	if err != nil {
		panic(err)
	}
	DB, _ := db.DB()
	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(10)

	return db
}

func (d *Dao) checkDBError(err error) error {
	if errors.Is(err, driver.ErrBadConn) {
		return err
	}
	return nil
}

func RecordNotFound(result *gorm.DB) bool {
	return errors.Is(result.Error, gorm.ErrRecordNotFound)
}
