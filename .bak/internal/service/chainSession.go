package service

import (
	"context"
	"github.com/shopspring/decimal"
	"subscan-end/internal/model"
)

func (s *Service) AddSessionValidator(sessionId uint, rankValidator int, controller, session, stash string,
	BondedActive, BondedTotal, BondedNominators, BondedOwner decimal.Decimal, unlocking string, CountNominators int, validatorPrefsValue decimal.Decimal) error {
	c := context.TODO()
	cs := model.SessionValidator{
		SessionId:           sessionId,
		RankValidator:       rankValidator,
		ValidatorController: controller,
		ValidatorSession:    session,
		ValidatorStash:      stash,
		BondedActive:        BondedActive,
		BondedTotal:         BondedTotal,
		BondedNominators:    BondedNominators,
		BondedOwner:         BondedOwner,
		Unlocking:           unlocking,
		CountNominators:     CountNominators,
		ValidatorPrefsValue: int(validatorPrefsValue.IntPart()),
	}
	err := s.dao.AddSessionValidator(c, &cs)
	return err
}

func (s *Service) AddSessionNominator(sessionId uint, RankValidator, RankNominator int, NominatorStash string, bonded decimal.Decimal) error {
	c := context.TODO()
	sn := model.SessionNominator{
		SessionId:      sessionId,
		RankValidator:  RankValidator,
		RankNominator:  RankNominator,
		NominatorStash: NominatorStash,
		Bonded:         bonded,
	}
	err := s.dao.AddSessionNominator(c, &sn)
	return err
}

func (s *Service) AddSession(sessionId uint, StartBlock, EndBlock, Era, CountValidators, CountNominators int) error {
	c := context.TODO()
	cs := model.ChainSession{
		SessionId:       sessionId,
		StartBlock:      StartBlock,
		EndBlock:        EndBlock,
		Era:             Era,
		CountNominators: CountNominators,
		CountValidators: CountValidators,
	}
	err := s.dao.AddSession(c, &cs)
	return err
}
