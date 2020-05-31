package dao

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-sql-driver/mysql"
	"github.com/itering/subscan/internal/util"
	"github.com/jinzhu/gorm"
)

// logs
type ormLog struct{}

func (l ormLog) Print(v ...interface{}) {
	log.Info(strings.Repeat("%v ", len(v)), v...)
}

// db
type GormDB struct {
	*gorm.DB
	gdbDone bool
}

func (c *GormDB) DbCommit() {
	if c.gdbDone {
		return
	}
	tx := c.Commit()
	c.gdbDone = true
	if err := tx.Error; err != nil && err != sql.ErrTxDone {
		fmt.Println("Fatal error DbCommit", err)
	}
}

func (c *GormDB) DbRollback() {
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
func newDb(dc mysqlConf) (db *gorm.DB) {
	var err error
	if os.Getenv("TASK_MOD") == "true" {
		db, err = gorm.Open("mysql", dc.Task.DSN)
	} else if os.Getenv("TEST_MOD") == "true" {
		db, err = gorm.Open("mysql", dc.Test.DSN)
	} else {
		db, err = gorm.Open("mysql", dc.Api.DSN)
	}
	if err != nil {
		panic(err)
	}
	db.DB().SetConnMaxLifetime(5 * time.Minute)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(10)
	if util.IsProduction() {
		db.SetLogger(ormLog{})
	}
	db.LogMode(true)
	return db
}

func (d *Dao) checkDBError(err error) error {
	if err == mysql.ErrInvalidConn || err == driver.ErrBadConn {
		return err
	}
	return nil
}
