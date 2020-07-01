package scan

import (
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/internal/substrate/metadata"
)

func (s *Service) SubstrateRuntimeList() []model.RuntimeVersion {
	return s.dao.RuntimeVersionList()
}

func (s *Service) SubstrateRuntimeInfo(spec int) *metadata.MetadataType {
	if metadataInstant, ok := metadata.RuntimeMetadata[spec]; ok {
		return metadataInstant
	}
	runtime := metadata.Process(s.dao.RuntimeVersionRaw(spec))
	if runtime == nil {
		return metadata.Latest(nil)
	}
	return runtime
}
