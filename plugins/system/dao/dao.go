package dao

import (
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/system/model"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/metadata"
	"strings"
)

func CreateExtrinsicError(db storage.DB, hash string, moduleError *model.MetadataModuleError) error {
	if moduleError == nil {
		return nil
	}
	err := db.Create(&model.ExtrinsicError{
		ExtrinsicHash: util.AddHex(hash),
		Module:        moduleError.Module,
		Name:          moduleError.Name,
		Doc:           strings.Join(moduleError.Doc, ","),
	})
	return err
}

func ExtrinsicError(db storage.DB, hash string) *model.ExtrinsicError {
	var e model.ExtrinsicError
	db.FindBy(&e, map[string]interface{}{"extrinsic_hash": hash}, nil)
	return &e
}

func CheckExtrinsicError(spec int, raw string, moduleIndex, errorIndex int) *model.MetadataModuleError {

	modules := metadata.Process(&metadata.RuntimeRaw{Raw: raw, Spec: spec})

	if moduleIndex >= len(modules.Metadata.Modules) {
		return nil
	}

	module := modules.Metadata.Modules[moduleIndex]
	if errorIndex >= len(module.Errors) {
		return nil
	}

	err := module.Errors[errorIndex]
	return &model.MetadataModuleError{
		Module: module.Name,
		Name:   err.Name,
		Doc:    err.Doc,
	}
}
