package service

import (
	"context"
	"fmt"
	"regexp"
	"subscan-end/internal/model"
	"subscan-end/utiles"
	"subscan-end/utiles/ss58"
)

func (s *Service) SearchByKey(key string, page int, row int) interface{} {
	c := context.TODO()
	if regexp.MustCompile("^[0-9]+$").MatchString(key) {
		return s.GetBlockByNum(utiles.StringToInt(key)) //
	} else if regexp.MustCompile(`^0x[0-9a-fA-F]{64}$`).MatchString(key) {
		return s.GetExtrinsicDetailByHash(key)
	} else {
		addressHex := ss58.Decode(key)
		if regexp.MustCompile(`^[0-9a-fA-F]{64}$`).MatchString(addressHex) == false {
			return nil
		}
		account, _ := s.dao.TouchAccount(c, addressHex)
		if account == nil {
			return nil
		}
		txs, count := s.dao.GetTransactionByAccount(c, addressHex, page, row)
		return map[string]interface{}{
			"account":    s.dao.AccountAsJson(c, account),
			"extrinsics": txs,
			"count":      count,
		}
	}
}

func (s *Service) DailyStat(start, end string) *[]model.DailyStatic {
	c := context.TODO()
	return s.dao.StatList(c, start, end)
}

func (s *Service) GetBalance() {
	c := context.TODO()
	fmt.Println(s.dao.GetBalanceFromNetwork(c, "7691596413414eea628899952adc3d407cec9794174cfc8f046ae27ce420b477", "balances"))
}
