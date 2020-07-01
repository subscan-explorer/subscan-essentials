package service

import (
	internalModel "github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/internal/plugins/storage"
	system "github.com/itering/subscan/internal/plugins/system/dao"
	"github.com/itering/subscan/internal/plugins/system/model"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/util"
)

type Service struct {
	d storage.Dao
}

func New(d storage.Dao) *Service {
	return &Service{
		d: d,
	}
}

func (s *Service) GetExtrinsicError(hash string) *model.ExtrinsicError {
	return system.ExtrinsicError(s.d.DB(), hash)
}

func (s *Service) ExtrinsicFailed(spec, blockTimestamp int, blockHash string, event *internalModel.ChainEvent, paramEvent []internalModel.EventParam) {

	type DispatchErrorModule struct {
		Index int `json:"index"`
		Error int `json:"error"`
	}

	for _, param := range paramEvent {

		if param.Type == "DispatchError" {

			var dr map[string]interface{}
			util.UnmarshalToAnything(&dr, param.Value)

			if _, ok := dr["Error"]; ok {
				_ = system.CreateExtrinsicError(s.d.DB(), event.ExtrinsicHash, system.CheckExtrinsicError(s.d.RuntimeVersionRaw(spec), util.IntFromInterface(dr["Module"]), util.IntFromInterface(dr["Error"])))

			} else if _, ok := dr["Module"]; ok {
				var module DispatchErrorModule
				util.UnmarshalToAnything(&module, dr["Module"])
				_ = system.CreateExtrinsicError(s.d.DB(), event.ExtrinsicHash, system.CheckExtrinsicError(s.d.RuntimeVersionRaw(spec), module.Index, module.Error))

			} else if _, ok := dr["BadOrigin"]; ok {
				_ = system.CreateExtrinsicError(s.d.DB(), event.ExtrinsicHash, &substrate.MetadataModuleError{Name: "BadOrigin"})

			} else if _, ok := dr["CannotLookup"]; ok {
				_ = system.CreateExtrinsicError(s.d.DB(), event.ExtrinsicHash, &substrate.MetadataModuleError{Name: "CannotLookup"})

			} else if _, ok := dr["Other"]; ok {
				_ = system.CreateExtrinsicError(s.d.DB(), event.ExtrinsicHash, &substrate.MetadataModuleError{Name: "Other"})

			}
			break
		}
	}
}
