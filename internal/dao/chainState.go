package dao

import (
	"github.com/itering/subscan/libs/substrate"
	"github.com/itering/subscan/libs/substrate/metadata"
)

// func (d *Dao) getConstant(module, key string) (storage.StateStorage, error) {
// 	m := metadata.Latest(nil)
// 	modules := m.Metadata.Modules
// 	metadataMap := make(map[string]types.MetadataModules)
// 	for _, value := range modules {
// 		metadataMap[strings.ToLower(value.Prefix)] = value
// 	}
//
// 	if _, ok := metadataMap[strings.ToLower(module)]; ok == false {
// 		return "", errors.New("not found this constant")
// 	}
//
// 	constantMap := make(map[string]types.MetadataConstants)
// 	for _, value := range metadataMap[strings.ToLower(module)].Constants {
// 		constantMap[strings.ToLower(value.Name)] = value
// 	}
//
// 	if _, ok := constantMap[strings.ToLower(key)]; ok == false {
// 		return "", errors.New("get storage type error")
// 	}
// 	return storage.StateStorage(constantMap[strings.ToLower(key)].ConstantsValue), nil
// }

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
