package substrate

import (
	"github.com/pkg/errors"
	"strings"
	"subscan-end/libs/substrate/scalecodec"
	"subscan-end/libs/substrate/scalecodec/types"
	"subscan-end/utiles"
)

type MetadataType types.MetadataStruct

var metadata *MetadataType

func InitMetaData() *MetadataType {
	if metadata != nil {
		return metadata
	}
	m := scalecodec.MetadataDecoder{}
	m.Init(utiles.HexToBytes(getCurrentMetadata()))
	m.Process()
	d := MetadataType(m.Metadata)
	return &d
}

func (m *MetadataType) getModuleStorageMapType(section, method string) (map[string]interface{}, error) {
	modules := m.Metadata.Modules
	metadataMap := make(map[string]types.MetadataModules)
	for _, value := range modules {
		metadataMap[strings.ToLower(value.Prefix)] = value
	}
	if _, ok := metadataMap[strings.ToLower(section)]; ok == false {
		return nil, errors.New("Get storage type error")
	}
	storageMap := make(map[string]types.MetadataStorage)
	for _, value := range metadataMap[strings.ToLower(section)].Storage {
		storageMap[strings.ToLower(value.Name)] = value
	}
	if _, ok := storageMap[strings.ToLower(method)]; ok == false {
		return nil, errors.New("Get storage type error")
	}
	return storageMap[strings.ToLower(method)].Type, nil
}
