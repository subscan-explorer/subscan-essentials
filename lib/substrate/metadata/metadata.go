package metadata

import (
	"github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/subscan/util"
	"strings"
)

type MetadataType types.MetadataStruct

type RuntimeRaw struct {
	Spec int
	Raw  string
}

var (
	latestSpec      = -1
	RuntimeMetadata = make(map[int]*MetadataType)
	Decoder         *scalecodec.MetadataDecoder
)

func Latest(runtime *RuntimeRaw) *MetadataType {
	if latestSpec != -1 {
		d := RuntimeMetadata[latestSpec]
		return d
	}
	if runtime == nil {
		return nil
	}
	m := scalecodec.MetadataDecoder{}
	m.Init(util.HexToBytes(runtime.Raw))
	_ = m.Process()

	Decoder = &m
	latestSpec = runtime.Spec

	instant := MetadataType(m.Metadata)
	RuntimeMetadata[latestSpec] = &instant
	return RuntimeMetadata[latestSpec]
}

func Process(runtime *RuntimeRaw) *MetadataType {
	if runtime == nil {
		return nil
	}
	if d, ok := RuntimeMetadata[runtime.Spec]; ok {
		return d
	}

	m := scalecodec.MetadataDecoder{}
	m.Init(util.HexToBytes(runtime.Raw))
	_ = m.Process()

	instant := MetadataType(m.Metadata)
	RuntimeMetadata[runtime.Spec] = &instant

	return RuntimeMetadata[runtime.Spec]
}

func RegNewMetadataType(spec int, coded string) *MetadataType {
	m := scalecodec.MetadataDecoder{}
	m.Init(util.HexToBytes(coded))
	_ = m.Process()

	instant := MetadataType(m.Metadata)
	RuntimeMetadata[spec] = &instant

	if spec > latestSpec {
		latestSpec = spec
	}
	return RuntimeMetadata[spec]
}

func (m *MetadataType) GetModuleStorageMapType(section, method string) (string, *types.StorageType) {
	modules := m.Metadata.Modules
	for _, value := range modules {
		if strings.EqualFold(strings.ToLower(value.Name), strings.ToLower(section)) {
			for _, storage := range value.Storage {
				if strings.EqualFold(strings.ToLower(storage.Name), strings.ToLower(method)) {
					return value.Prefix, &storage.Type
				}
			}
		}
	}
	return "", nil
}
