package service

import (
	"context"
	"github.com/itering/subscan/internal/model"
)

func (s *Service) GetActiveAccountCount(c context.Context) int {
	count := s.dao.GetActiveAccountCount(c)
	return count
}

func (s *Service) GetAccountList(c context.Context, queryWhere ...string) []*model.ChainAccount {
	list, _ := s.dao.GetAccountList(c, 0, 1000000, "desc", "id", queryWhere...)
	return list
}

func (s *Service) RefreshAccount(account *model.ChainAccount, u map[string]interface{}) error {
	return s.dao.RefreshAccount(account, u)
}

func (s *Service) UpdateAccountAllBalance(address string) {
	c := context.TODO()
	if account, err := s.dao.TouchAccount(c, address); err == nil {
		_, _, _ = s.dao.UpdateAccountBalance(c, account, "balances")
		_ = s.dao.UpdateAccountLock(c, account.Address, "ring")
	}
}
