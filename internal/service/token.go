package service

import (
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/rpc"
	"sync"
)

var onceToken sync.Once

// Unknown token reg
func (s *Service) unknownToken() {
	onceToken.Do(func() {
		if p, _ := rpc.GetSystemProperties(nil); p != nil {
			util.AddressType = util.IntToString(p.Ss58Format)
			util.BalanceAccuracy = util.IntToString(p.TokenDecimals)
		}
	})
}
