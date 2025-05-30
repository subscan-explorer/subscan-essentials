package model

import (
	"github.com/itering/subscan/util"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Option = func(tx *gorm.DB) *gorm.DB

func Select(query interface{}, args ...interface{}) Option {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Select(query, args...)
	}
}

func Omit(args ...string) Option {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Omit(args...)
	}
}

func Where(query interface{}, args ...interface{}) Option {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(query, args...)
	}
}

func GroupBy(field string) Option {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Group(field)
	}
}

func WhereOr(query interface{}, args ...interface{}) Option {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Or(query, args...)
	}
}

func Order(value interface{}) Option {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Order(value)
	}
}

func Conditions(conds []string, params []interface{}) Option {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(strings.Join(conds, " AND "), params...)
	}
}

func ForUpdate() Option {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Clauses(clause.Locking{Strength: "UPDATE"})
	}
}

func Offset(page, rows int) Option {
	return WithLimit(page*rows, rows)
}

func Nothing() Option {
	return func(tx *gorm.DB) *gorm.DB {
		return tx
	}
}

func WithLimit(offset, rows int) Option {
	return func(tx *gorm.DB) *gorm.DB {
		if offset > -1 {
			tx = tx.Offset(offset)
		}
		if rows > 0 {
			tx = tx.Limit(rows)
		}
		return tx
	}
}

func IgnoreDuplicate(tx *gorm.DB) *gorm.DB {
	return tx.Clauses(clause.OnConflict{DoNothing: true})
}

func TableNameFunc(c interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(TableNameFromInterface(c, db))
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

func RedisKeyPrefix() string {
	return util.NetworkNode + ":"
}

func MetadataCacheKey() string {
	return RedisKeyPrefix() + "metadata"
}
