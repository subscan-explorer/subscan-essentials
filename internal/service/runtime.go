package service

import (
	"context"
	"strings"

	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
)

var (
	runtimeSpecs []int
)

func (s *Service) regRuntimeVersion(ctx context.Context, name string, spec int, blockNum uint, hash ...string) error {
	if util.IntInSlice(spec, runtimeSpecs) {
		return nil
	}
	if affected := s.dao.CreateRuntimeVersion(ctx, name, spec, blockNum); affected {
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

func (s *Service) setRuntimeData(spec int, runtime *metadata.Instant, rawData string) {
	var modules []string
	for _, value := range runtime.Metadata.Modules {
		modules = append(modules, value.Name)
	}
	s.dao.SetRuntimeData(spec, strings.Join(modules, "|"), rawData)
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
