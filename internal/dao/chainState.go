package dao

import (
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/substrate/metadata"
)

func (d *Dao) CheckExtrinsicError(spec, moduleIndex, errorIndex int) *substrate.MetadataModuleError {
	modules, ok := metadata.RuntimeMetadata[spec]
	if !ok {
		modules = metadata.Init(d.RuntimeVersionRaws(spec), spec)
	}

	if moduleIndex >= len(modules.Metadata.Modules) {
		return nil
	}

	module := modules.Metadata.Modules[moduleIndex]
	if errorIndex >= len(module.Errors) {
		return nil
	}

	err := module.Errors[errorIndex]
	return &substrate.MetadataModuleError{
		Module: module.Name,
		Name:   err.Name,
		Doc:    err.Doc,
	}
}
