package service

import (
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/system/dao"
	"github.com/itering/subscan/plugins/system/model"
	"github.com/itering/subscan/util"
)

type Service struct {
	dao storage.Dao
}

func New(d storage.Dao) *Service {
	return &Service{
		dao: d,
	}
}

func (s *Service) GetExtrinsicError(hash string) *model.ExtrinsicError {
	return dao.ExtrinsicError(s.dao, hash)
}

func (s *Service) ExtrinsicFailed(spec int, event *storage.Event, paramEvent []storage.EventParam) {

	type DispatchErrorModule struct {
		Index int `json:"index"`
		Error int `json:"error"`
	}

	for _, param := range paramEvent {

		if param.Type == "DispatchError" {

			var dr map[string]interface{}
			util.UnmarshalAny(&dr, param.Value)

			if _, ok := dr["Error"]; ok {
				_ = dao.CreateExtrinsicError(s.dao,
					event.ExtrinsicHash,
					dao.CheckExtrinsicError(spec, s.dao.SpecialMetadata(spec), util.IntFromInterface(dr["Module"]), util.IntFromInterface(dr["Error"])))

			} else if _, ok := dr["Module"]; ok {
				var module DispatchErrorModule
				util.UnmarshalAny(&module, dr["Module"])

				_ = dao.CreateExtrinsicError(s.dao,
					event.ExtrinsicHash,
					dao.CheckExtrinsicError(spec, s.dao.SpecialMetadata(spec), module.Index, module.Error))

			} else if _, ok := dr["BadOrigin"]; ok {
				_ = dao.CreateExtrinsicError(s.dao, event.ExtrinsicHash,
					&model.MetadataModuleError{Name: "BadOrigin"})

			} else if _, ok := dr["CannotLookup"]; ok {
				_ = dao.CreateExtrinsicError(s.dao, event.ExtrinsicHash,
					&model.MetadataModuleError{Name: "CannotLookup"})

			} else if _, ok := dr["Other"]; ok {
				_ = dao.CreateExtrinsicError(s.dao, event.ExtrinsicHash,
					&model.MetadataModuleError{Name: "Other"})

			}
			break
		}
	}
}
