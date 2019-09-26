package service

import (
	"context"
	"fmt"
	"subscan-end/internal/model"
	"subscan-end/utiles/ss58"
)

func (s *Service) GetTransfersByAccount(account string, page, row int) (*[]model.TransferJson, int) {
	c := context.TODO()
	accountHex := ss58.Decode(account)
	if accountHex == "" {
		return nil, 0
	}
	fromQuery := fmt.Sprintf("from_hex = '%s'", accountHex)
	return s.dao.GetTransfers(c, page, row, fromQuery), s.dao.GetTransferCount(c, fromQuery)
}
