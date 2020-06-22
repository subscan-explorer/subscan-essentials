package dao

import (
	"github.com/itering/subscan/internal/plugins/system/model"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/util"
	"github.com/jinzhu/gorm"
	"strings"
)

func CreateExtrinsicError(db *gorm.DB, hash string, moduleError *substrate.MetadataModuleError) error {
	if moduleError == nil {
		return nil
	}
	query := db.Create(&model.ExtrinsicError{
		ExtrinsicHash: util.AddHex(hash),
		Module:        moduleError.Module,
		Name:          moduleError.Name,
		Doc:           strings.Join(moduleError.Doc, ","),
	})
	return query.Error
}

func ExtrinsicError(db *gorm.DB, hash string) *model.ExtrinsicError {
	var e model.ExtrinsicError
	db.Where("extrinsic_hash = ?", hash).Find(&e)
	return &e
}
