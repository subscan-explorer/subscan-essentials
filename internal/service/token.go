package service

import (
	"github.com/itering/subscan/lib/substrate"
	"github.com/itering/subscan/lib/substrate/rpc"
	"sync"
)

var onceToken sync.Once

// Unknown token reg
func (s *Service) UnknownToken() {
	onceToken.Do(func() {
		if p, _ := rpc.GetSystemProperties(); p != nil {
			substrate.AddressType = p.Ss58Format
			substrate.BalanceAccuracy = p.TokenDecimals
		}
	})
}
