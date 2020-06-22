package service

import (
	"github.com/itering/subscan/internal/dao"
	internalModel "github.com/itering/subscan/internal/model"
	system "github.com/itering/subscan/internal/plugins/system/dao"
	"github.com/itering/subscan/internal/plugins/system/model"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/util"
)

type Service struct {
	d *dao.Dao
}

func New(d *dao.Dao) *Service {
	return &Service{
		d: d,
	}
}

func (s *Service) GetExtrinsicError(hash string) *model.ExtrinsicError {
	return system.ExtrinsicError(s.d.Db, hash)
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
				_ = system.CreateExtrinsicError(s.d.Db, event.ExtrinsicHash, s.d.CheckExtrinsicError(spec, util.IntFromInterface(dr["Module"]), util.IntFromInterface(dr["Error"])))

			} else if _, ok := dr["Module"]; ok {
				var module DispatchErrorModule
				util.UnmarshalToAnything(&module, dr["Module"])
				_ = system.CreateExtrinsicError(s.d.Db, event.ExtrinsicHash, s.d.CheckExtrinsicError(spec, module.Index, module.Error))

			} else if _, ok := dr["BadOrigin"]; ok {
				_ = system.CreateExtrinsicError(s.d.Db, event.ExtrinsicHash, &substrate.MetadataModuleError{Name: "BadOrigin"})

			} else if _, ok := dr["CannotLookup"]; ok {
				_ = system.CreateExtrinsicError(s.d.Db, event.ExtrinsicHash, &substrate.MetadataModuleError{Name: "CannotLookup"})

			} else if _, ok := dr["Other"]; ok {
				_ = system.CreateExtrinsicError(s.d.Db, event.ExtrinsicHash, &substrate.MetadataModuleError{Name: "Other"})

			}
			break
		}
	}
}
