package dao

import (
	"database/sql"
	"errors"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"fmt"

	"reflect"
	"time"

	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/configs"
	"github.com/itering/subscan/model"
	"github.com/itering/substrate-api-rpc/websocket"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"

	"github.com/itering/subscan/util"
	"gorm.io/gorm"
)

type DbStorage struct {
	db       *gorm.DB
	Prefix   string
	DbDriver string
}

func (d *DbStorage) GetDbInstance() any {
	return d.db
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
	if query := d.db.Where("spec_version = ?", spec).First(&raw); query.Error != nil {
		return ""
	}
	return raw.RawData
}

func (d *DbStorage) getModelTableName(model interface{}) string {
	return TableNameFromInterface(model, d.db)
}

func (d *DbStorage) checkProtected(model interface{}) error {
	if util.StringInSlice(d.getModelTableName(model), protectedTables) {
		return errors.New("protected tables")
	}
	return nil
}

func (d *DbStorage) RPCPool() *websocket.PoolConn {
	conn, _ := websocket.Init()
	return conn
}

func (d *DbStorage) getPluginPrefixTableName(instant interface{}) string {
	tableName := d.getModelTableName(instant)
	_, implementTable := instant.(Tabler)
	if util.StringInSlice(tableName, protectedTables) || implementTable {
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
		tx = tx.Table(fmt.Sprintf("%s_%s", option.PluginPrefix, d.getModelTableName(record)))
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
	return int(count), tx.Error != nil
}

func (d *DbStorage) AutoMigration(model interface{}) error {
	if d.checkProtected(model) == nil {
		if d.DbDriver == "mysql" {
			return d.db.Table(d.getPluginPrefixTableName(model)).Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model)
		}
		return d.db.Table(d.getPluginPrefixTableName(model)).AutoMigrate(model)

	}
	return nil
}

func (d *DbStorage) AddIndex(model interface{}, indexName string, columns ...string) error {
	if d.checkProtected(model) == nil {
		return d.db.Table(d.getPluginPrefixTableName(model)).Migrator().CreateIndex(indexName, columns[0])
	}
	return nil
}

func (d *DbStorage) AddUniqueIndex(model interface{}, indexName string, columns ...string) error {
	if d.checkProtected(model) == nil {
		return d.db.Table(d.getPluginPrefixTableName(model)).Migrator().CreateIndex(indexName, columns[0])
	}
	return nil
}

func (d *DbStorage) Create(record interface{}) error {
	if err := d.checkProtected(record); err == nil {
		tx := d.db.Table(d.getPluginPrefixTableName(record)).Scopes(model.IgnoreDuplicate).Create(record)
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
	if err := tx.Error; err != nil && !errors.Is(err, sql.ErrTxDone) {
		fmt.Println("Fatal error DbCommit", err)
	}
}

func (d *Dao) DbRollback(c *GormDB) {
	if c.gdbDone {
		return
	}
	tx := c.Rollback()
	c.gdbDone = true
	if err := tx.Error; err != nil && !errors.Is(err, sql.ErrTxDone) {
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

type NamingStrategy struct {
	schema.NamingStrategy
}

func (n NamingStrategy) UniqueName(_, column string) string {
	return column
}

func (n NamingStrategy) IndexName(_, column string) string {
	return column
}

// private funcs
func newDb() (db *gorm.DB) {
	var err error
	dbDriver := util.GetEnv("DB_DRIVER", "mysql")
	util.Logger().Debug(fmt.Sprintf("Set DB_DRIVER %s", dbDriver))
	conf := &gorm.Config{
		Logger: logger.Default,
	}
	if dbDriver == "mysql" {
		conf.NamingStrategy = NamingStrategy{}
		db, err = gorm.Open(mysql.Open(configs.Boot.Database.Mysql.DSN), conf)
	} else {
		db, err = gorm.Open(postgres.Open(configs.Boot.Database.Postgres.DSN), conf)
	}

	if err != nil {
		panic(err)
	}
	sqldb, _ := db.DB()
	sqldb.SetConnMaxLifetime(5 * time.Minute)
	sqldb.SetMaxOpenConns(util.StringToInt(util.GetEnv("MAX_DB_CONN_COUNT", "200")))
	sqldb.SetMaxIdleConns(10)
	return db
}
