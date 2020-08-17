package service

import (
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/websocket"
	"sync"
)

var onceToken sync.Once

// Unknown token reg
func (s *Service) unknownToken() {
	websocket.SetEndpoint(util.WSEndPoint)
	onceToken.Do(func() {
		if p, _ := rpc.GetSystemProperties(nil); p != nil {
			util.AddressType = util.IntToString(p.Ss58Format)
			util.BalanceAccuracy = util.IntToString(p.TokenDecimals)
		}
	})
}
