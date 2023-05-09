package service

import (
	"strings"

	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"golang.org/x/exp/slog"
)

var runtimeSpecs []int

func (s *Service) SubstrateRuntimeList() []model.RuntimeVersion {
	return s.dao.RuntimeVersionList()
}

func (s *Service) SubstrateRuntimeInfo(spec int) *metadata.Instant {
	if metadataInstant, ok := metadata.RuntimeMetadata[spec]; ok {
		return metadataInstant
	}
	runtime := metadata.Process(s.dao.RuntimeVersionRaw(spec))
	if runtime == nil {
		return metadata.Latest(nil)
	}
	return runtime
}

func (s *Service) regRuntimeVersion(name string, spec int, hash ...string) error {
	if util.IntInSlice(spec, runtimeSpecs) {
		return nil
	}
	if affected := s.dao.CreateRuntimeVersion(name, spec); affected > 0 {
		if coded := s.regCodecMetadata(hash...); coded != "" {
			runtime := metadata.RegNewMetadataType(spec, coded)
			s.setRuntimeData(spec, runtime, coded)
		} else {
			panic("get runtime metadata error")
		}
	}
	runtimeSpecs = append(runtimeSpecs, spec)
	return nil
}

func (s *Service) regCodecMetadata(hash ...string) string {
	if coded, err := rpc.GetMetadataByHash(nil, hash...); err == nil {
		return coded
	}
	return ""
}

func (s *Service) createRuntimeConstants(spec int, runtime *metadata.Instant) {
	var constants []model.RuntimeConstant
	for _, module := range runtime.Metadata.Modules {
		for _, constant := range module.Constants {
			constants = append(constants, model.RuntimeConstant{
				ModuleName:   module.Name,
				ConstantName: constant.Name,
				Type:         constant.Type,
				Value:        constant.ConstantsValue,
				SpecVersion:  spec,
			})
		}
	}
	err := s.dao.CreateRuntimeConstants(spec, constants)
	if err != nil {
		slog.Error("create runtime constants failed", "error", err)
	}
}

func (s *Service) setRuntimeData(spec int, runtime *metadata.Instant, rawData string) {
	var modules []string
	for _, value := range runtime.Metadata.Modules {
		modules = append(modules, value.Name)
	}
	s.dao.SetRuntimeData(spec, strings.Join(modules, "|"), rawData)
	s.createRuntimeConstants(spec, runtime)
}

func (s *Service) getMetadataInstant(spec int, hash string) *metadata.Instant {
	metadataInstant, ok := metadata.RuntimeMetadata[spec]
	if !ok {
		raw := s.dao.RuntimeVersionRaw(spec)
		if raw.Raw == "" {
			raw.Raw = s.regCodecMetadata(hash)
		}
		metadataInstant = metadata.Process(raw)
	}
	return metadataInstant
}
