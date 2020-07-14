package service

import (
	"github.com/itering/scale.go/types"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"strings"
)

func (s *Service) regRuntimeVersion(name string, spec int) error {
	if affected := s.dao.CreateRuntimeVersion(name, spec); affected > 0 {
		if coded := s.regCodecMetadata(); coded != "" {
			runtime := metadata.RegNewMetadataType(spec, coded)
			s.setRuntimeData(spec, runtime, coded)
		}
	}
	return nil
}

func (s *Service) regCodecMetadata() string {
	if coded, err := rpc.GetMetadataByHash(nil); err == nil {
		return coded
	}
	return ""
}

func (s *Service) setRuntimeData(spec int, runtime *types.MetadataStruct, rawData string) {
	var modules []string
	for _, value := range runtime.Metadata.Modules {
		modules = append(modules, value.Name)
	}
	s.dao.SetRuntimeData(spec, strings.Join(modules, "|"), rawData)
}
