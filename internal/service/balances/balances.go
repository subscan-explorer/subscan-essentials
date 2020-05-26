package balances

import (
	"context"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/util"
	"github.com/shopspring/decimal"
)

type Balances interface {
	transfer()
	newAccount()
	reapedAccount()
}

type balances struct {
	dao            *dao.Dao
	eventParams    []model.EventParam
	e              *model.ChainEvent
	blockTimestamp int
	fee            decimal.Decimal
}

func New(d *dao.Dao, e *model.ChainEvent, params []model.EventParam, blockTimestamp int, fee decimal.Decimal) Balances {
	var service balances
	s := balances{
		dao:            d,
		eventParams:    params,
		e:              e,
		blockTimestamp: blockTimestamp,
		fee:            fee,
	}
	service = s
	return &service
}

func EmitEvent(s Balances, eventId string) {
	switch eventId {
	case "Transfer":
		s.transfer()
	case "NewAccount", "Endowed": // account created
		s.newAccount()
	case "ReapedAccount": // account reaped
		s.reapedAccount()
	}
}

func (s *balances) transfer() {
	c := context.TODO()
	if len(s.eventParams) < 3 {
		return
	}

	to := util.TrimHex(util.InterfaceToString(s.eventParams[1].Value))
	if account, err := s.dao.TouchAccount(c, to); err == nil {
		_, _, _ = s.dao.UpdateAccountBalance(c, account, s.e.ModuleId)
	}

}

func (s *balances) newAccount() {
	c := context.TODO()
	if account, err := s.dao.TouchAccount(c, util.TrimHex(util.InterfaceToString(s.eventParams[0].Value))); err == nil {
		_, _, _ = s.dao.UpdateAccountBalance(c, account, s.e.ModuleId)
	}
}

func (s *balances) reapedAccount() {
	s.dao.ResetAccountNonce(context.TODO(), util.TrimHex(util.InterfaceToString(s.eventParams[0].Value)))
}
