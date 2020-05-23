package metadata

import (
	"encoding/json"
	"fmt"
	"github.com/freehere107/go-scale-codec"
	"github.com/freehere107/go-scale-codec/types"
	"io/ioutil"
	"strings"
	"github.com/itering/subscan/util"
)

const (
	CrabMetadata   = "crab"
	KusamaMetadata = "kusama"
)

type MetadataType types.MetadataStruct

type RuntimeRaw struct {
	Spec int
	Raw  string
}

var (
	latestSpec      = -1
	isInit          bool
	RuntimeMetadata = make(map[int]*MetadataType)
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
	latestSpec = runtime.Spec
	instant := MetadataType(m.Metadata)
	RuntimeMetadata[latestSpec] = &instant
	d := RuntimeMetadata[latestSpec]
	return d
}

func Init(runtime *[]RuntimeRaw, spec ...int) *MetadataType {
	if isInit && len(spec) == 0 {
		d := RuntimeMetadata[latestSpec]
		return d
	}

	var processMetadata = func(specVersionInt int, value string) {
		if _, ok := RuntimeMetadata[specVersionInt]; ok {
			return
		}
		m := scalecodec.MetadataDecoder{}
		m.Init(util.HexToBytes(value))
		_ = m.Process()
		instant := MetadataType(m.Metadata)
		RuntimeMetadata[specVersionInt] = &instant
	}

	for _, value := range *runtime {
		specVersionInt := value.Spec

		if len(spec) > 0 {
			if spec[0] == specVersionInt {
				processMetadata(specVersionInt, value.Raw)
				d := RuntimeMetadata[specVersionInt]
				return d
			}
			continue
		}

		processMetadata(specVersionInt, value.Raw)

		if specVersionInt > latestSpec {
			latestSpec = specVersionInt
		}
	}
	isInit = true
	d := RuntimeMetadata[latestSpec]
	return d
}

func RegNewMetadataType(spec int, coded string) *MetadataType {
	m := scalecodec.MetadataDecoder{}
	m.Init(util.HexToBytes(coded))
	_ = m.Process()
	instant := MetadataType(m.Metadata)
	RuntimeMetadata[spec] = &instant
	if latestSpec == -1 {
		latestSpec = spec
	}
	return RuntimeMetadata[spec]
}

func (m *MetadataType) GetModuleStorageMapType(section, method string) (string, *types.StorageType) {
	modules := m.Metadata.Modules
	for _, value := range modules {
		if strings.ToLower(value.Name) == strings.ToLower(section) {
			for _, storage := range value.Storage {
				if strings.ToLower(storage.Name) == strings.ToLower(method) {
					return value.Prefix, &storage.Type
				}
			}
		}
	}
	return "", nil
}

func LocalMetadata() map[string]string {
	b, err := ioutil.ReadFile(fmt.Sprintf("../configs/metadata/%s.json", NetworkRelateName()))
	if err != nil {
		return nil
	}
	var metadataFile map[string]string
	_ = json.Unmarshal(b, &metadataFile)

	return metadataFile
}

func NetworkRelateName() string {
	var fileName string
	switch util.NetworkNode {
	case util.KusamaNetwork:
		fileName = KusamaMetadata
	case util.CrabNetwork:
		fileName = CrabMetadata
	default:
		fileName = util.NetworkNode
	}
	return fileName
}

func LatestLocal() (int, string) {
	list := LocalMetadata()
	var latest string
	for spec := range list {
		if specInt := util.StringToInt(spec); specInt > util.StringToInt(latest) {
			latest = spec
		}
	}
	return util.StringToInt(latest), list[latest]
}
