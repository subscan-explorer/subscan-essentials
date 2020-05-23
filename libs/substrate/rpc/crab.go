package rpc

import (
	"math/rand"

	"github.com/itering/subscan/libs/substrate/websocket"
	"github.com/shopspring/decimal"
)

type crab struct {
	query
}

func (c *crab) PowerOf(address string) (power decimal.Decimal) {
	v := &JsonRpcResult{}
	if err := websocket.SendWsRequest(c.c, v, PowerOf(rand.Intn(10000), address)); err != nil {
		return
	}
	var powerStruct struct {
		Power decimal.Decimal `json:"power"`
	}
	if err := v.ToAnyThing(&powerStruct); err == nil {
		return powerStruct.Power
	}
	return
}
