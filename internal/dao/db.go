package dao

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/itering/subscan/configs"
	"os"
	"strings"
	"time"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-sql-driver/mysql"
	"github.com/itering/subscan/util"
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
func newDb(dc configs.MysqlConf) (db *gorm.DB) {
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
	if util.IsProduction {
		db.SetLogger(ormLog{})
	}
	if os.Getenv("TEST_MOD") != "true" {
		db.LogMode(true)
	}
	return db
}

func (d *Dao) checkDBError(err error) error {
	if err == mysql.ErrInvalidConn || err == driver.ErrBadConn {
		return err
	}
	return nil
}
