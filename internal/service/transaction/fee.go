package transaction

import (
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/shopspring/decimal"
)

// GetExtrinsicFee
func GetExtrinsicFee(encodeExtrinsic string) (fee decimal.Decimal) {
	paymentInfo, _ := rpc.GetPaymentQueryInfo(nil, encodeExtrinsic)
	if paymentInfo != nil {
		return paymentInfo.PartialFee
	}
	return
}
