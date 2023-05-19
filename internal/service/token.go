package service

import (
	"sync"

	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/websocket"
)

var onceToken sync.Once

// Unknown token reg
func (s *Service) unknownToken() {
	websocket.SetEndpoint(util.WSEndPoint)
	onceToken.Do(func() {
		if p, _ := util.GetSystemProperties(nil); p != nil {
			if p.Ss58Format != nil {
				util.AddressType = util.IntToString(*p.Ss58Format)
			}
			if p.TokenDecimals != nil {
				util.BalanceAccuracy = util.IntToString(*p.TokenDecimals)
			}
		}
	})
}
