package service

import (
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/shopspring/decimal"
)

// GetExtrinsicFee
func GetExtrinsicFee(p websocket.WsConn, encodeExtrinsic string) (fee decimal.Decimal, err error) {
	paymentInfo, err := rpc.GetPaymentQueryInfo(p, encodeExtrinsic)
	if paymentInfo != nil {
		return paymentInfo.PartialFee, nil
	}
	return decimal.Zero, err
}
