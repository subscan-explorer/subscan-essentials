package service

import (
	"github.com/itering/subscan/lib/substrate"
	internalModel "github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins/storage"
	"github.com/itering/subscan/plugins/system/dao"
	"github.com/itering/subscan/plugins/system/model"
	"github.com/itering/subscan/util"
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
	return dao.ExtrinsicError(s.d.DB(), hash)
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
				_ = dao.CreateExtrinsicError(s.d.DB(), event.ExtrinsicHash, dao.CheckExtrinsicError(s.d.RuntimeVersionRaw(spec), util.IntFromInterface(dr["Module"]), util.IntFromInterface(dr["Error"])))

			} else if _, ok := dr["Module"]; ok {
				var module DispatchErrorModule
				util.UnmarshalToAnything(&module, dr["Module"])
				_ = dao.CreateExtrinsicError(s.d.DB(), event.ExtrinsicHash, dao.CheckExtrinsicError(s.d.RuntimeVersionRaw(spec), module.Index, module.Error))

			} else if _, ok := dr["BadOrigin"]; ok {
				_ = dao.CreateExtrinsicError(s.d.DB(), event.ExtrinsicHash, &substrate.MetadataModuleError{Name: "BadOrigin"})

			} else if _, ok := dr["CannotLookup"]; ok {
				_ = dao.CreateExtrinsicError(s.d.DB(), event.ExtrinsicHash, &substrate.MetadataModuleError{Name: "CannotLookup"})

			} else if _, ok := dr["Other"]; ok {
				_ = dao.CreateExtrinsicError(s.d.DB(), event.ExtrinsicHash, &substrate.MetadataModuleError{Name: "Other"})

			}
			break
		}
	}
}
