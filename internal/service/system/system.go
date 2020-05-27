package system

import (
	"context"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/util"
)

type System interface {
	ExtrinsicFailed()
	NewAccount()
	KilledAccount()
	SetCode()
}

type system struct {
	dao            *dao.Dao
	e              *model.ChainExtrinsic
	event          *model.ChainEvent
	eventParams    []model.EventParam
	blockHash      string
	blockTimestamp int
	spec           int
}

func New(d *dao.Dao, e *model.ChainExtrinsic) System {
	s := system{
		dao: d,
		e:   e,
	}
	return &s
}

func NewEvent(d *dao.Dao, e *model.ChainEvent, params []model.EventParam, hash string, blockTimestamp, spec int) System {
	s := system{
		dao:            d,
		event:          e,
		eventParams:    params,
		blockHash:      hash,
		blockTimestamp: blockTimestamp,
		spec:           spec,
	}
	return &s
}

func EmitExtrinsic(d System, method string) {
	switch method {
	case "set_code":
		d.SetCode()
	}
}

func EmitEvent(s System, eventId string) {
	switch eventId {
	case "ExtrinsicFailed":
		s.ExtrinsicFailed()
	case "NewAccount":
		s.NewAccount()
	case "KilledAccount":
		s.KilledAccount()
	}
}

func (s *system) ExtrinsicFailed() {
	for _, param := range s.eventParams {
		if param.Type == "DispatchError" {
			dr := model.ParsingExtrinsicErrorParam(param.Value)
			if _, ok := dr["Error"]; ok {
				_ = s.dao.CreateExtrinsicError(s.event.ExtrinsicHash, s.dao.CheckExtrinsicError(s.spec, util.IntFromInterface(dr["Module"]), util.IntFromInterface(dr["Error"])))
			} else if _, ok := dr["Module"]; ok {
				var module model.DispatchErrorModule
				util.UnmarshalToAnything(&module, dr["Module"])
				_ = s.dao.CreateExtrinsicError(s.event.ExtrinsicHash, s.dao.CheckExtrinsicError(s.spec, module.Index, module.Error))
			} else if _, ok := dr["BadOrigin"]; ok {
				_ = s.dao.CreateExtrinsicError(s.event.ExtrinsicHash, &substrate.MetadataModuleError{Name: "BadOrigin"})
			} else if _, ok := dr["CannotLookup"]; ok {
				_ = s.dao.CreateExtrinsicError(s.event.ExtrinsicHash, &substrate.MetadataModuleError{Name: "CannotLookup"})
			} else if _, ok := dr["Other"]; ok {
				_ = s.dao.CreateExtrinsicError(s.event.ExtrinsicHash, &substrate.MetadataModuleError{Name: "Other"})
			}
			break
		}
	}

}
func (s *system) NewAccount() {
	c := context.TODO()
	if account, err := s.dao.TouchAccount(c, util.TrimHex(util.InterfaceToString(s.eventParams[0].Value))); err == nil {
		_, _, _ = s.dao.UpdateAccountBalance(c, account, "balances")
	}
}

func (s *system) KilledAccount() {
	s.dao.ResetAccountNonce(context.TODO(), util.TrimHex(util.InterfaceToString(s.eventParams[0].Value)))
}

func (s *system) SetCode() {}
