package dao

import (
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/util"
	"strings"
)

func (d *Dao) CreateExtrinsicError(hash string, moduleError *substrate.MetadataModuleError) error {
	if moduleError == nil {
		return nil
	}
	query := d.Db.Create(&model.ExtrinsicError{
		ExtrinsicHash: util.AddHex(hash),
		Module:        moduleError.Module,
		Name:          moduleError.Name,
		Doc:           strings.Join(moduleError.Doc, ","),
	})
	return query.Error
}

func (d *Dao) ExtrinsicError(hash string) *model.ExtrinsicError {
	var e model.ExtrinsicError
	d.Db.Where("extrinsic_hash = ?", hash).Find(&e)
	return &e
}
