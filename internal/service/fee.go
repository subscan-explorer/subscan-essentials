package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/itering/scale.go/types"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/model"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/shopspring/decimal"
	"math/rand"
	"strings"
)

type PaymentQueryFeeDetails struct {
	InclusionFee *inclusionFee `json:"inclusionFee"`
}
type inclusionFee struct {
	BaseFee           decimal.Decimal `json:"baseFee"`
	LenFee            decimal.Decimal `json:"lenFee"`
	AdjustedWeightFee decimal.Decimal `json:"adjustedWeightFee"`
}

func (i *inclusionFee) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	for k, v := range m {
		if hex, ok := v.(string); ok && strings.HasPrefix(hex, "0x") {
			m[k] = util.HexToNumStr(hex)
		}
	}
	type T inclusionFee
	return util.UnmarshalAny((*T)(i), m)
}

func (p *PaymentQueryFeeDetails) basicFee() decimal.Decimal {
	return p.InclusionFee.BaseFee.Add(p.InclusionFee.LenFee)
}

func (p *PaymentQueryFeeDetails) EstimateFee() decimal.Decimal {
	return p.basicFee().Add(p.InclusionFee.AdjustedWeightFee)
}

func (p *PaymentQueryFeeDetails) ActualFee(estimateWeight, actualWeight decimal.Decimal) decimal.Decimal {
	if estimateWeight.IsZero() || actualWeight.IsZero() {
		// https://alephzero.subscan.io/extrinsic/4933476-1
		return decimal.Zero
	}
	return p.basicFee().Add(p.InclusionFee.AdjustedWeightFee.Div(estimateWeight).Mul(actualWeight)).Floor()
}

type PaymentQueryInfo struct {
	Class      string
	PartialFee decimal.Decimal
	Weight     decimal.Decimal
}

func StateCallFunction(id int, method, encodedExtrinsic, hash string) []byte {
	p := rpc.Param{Id: id, Method: method, Params: []string{encodedExtrinsic, hash}}
	p.JsonRpc = "2.0"
	b, _ := json.Marshal(p)
	return b
}

func SystemPaymentQueryInfo(id int, encodedExtrinsic, hash string) []byte {
	p := rpc.Param{Id: id, Method: "payment_queryInfo", Params: []string{encodedExtrinsic, hash}}
	p.JsonRpc = "2.0"
	b, _ := json.Marshal(p)
	return b
}

func SystemPaymentQueryFeeDetails(id int, encodedExtrinsic, hash string) []byte {
	p := rpc.Param{Id: id, Method: "payment_queryFeeDetails", Params: []string{encodedExtrinsic, hash}}
	p.JsonRpc = "2.0"
	b, _ := json.Marshal(p)
	return b
}

var InvalidValue = errors.New("invalid value")

func GetPaymentQueryInfo(_ context.Context, spec int, encodedExtrinsic, hash string, isV2Weight bool) (paymentInfo *PaymentQueryInfo, err error) {
	var result string
	v := &model.JsonRpcResult{}
	err = websocket.SendWsRequest(nil, v, StateCallFunction(rand.Intn(10000), "TransactionPaymentApi_query_info", encodedExtrinsic+types.Encode("U32", len(util.HexToBytes(encodedExtrinsic))), hash))
	if err == nil {
		result, err = v.ToString()
		if err != nil {
			return nil, err
		}
		if len(result) == 0 {
			return nil, InvalidValue
		}
		storageBytes := util.HexToBytes(result)
		var decodeMsg storage.StateStorage
		if isV2Weight {
			decodeMsg, _, err = storage.Decode(util.BytesToHex(storageBytes), "RuntimeDispatchInfo", &types.ScaleDecoderOption{Spec: spec})
		}
		if err != nil || !isV2Weight {
			if decodeMsg, _, err = storage.Decode(util.BytesToHex(storageBytes), "RuntimeDispatchInfoV1", &types.ScaleDecoderOption{Spec: spec}); err != nil {
				return nil, err
			}
		}
		decodeMsg.ToAny(&paymentInfo)
		return paymentInfo, nil
	}
	v = &model.JsonRpcResult{}
	if err = websocket.SendWsRequest(nil, v, SystemPaymentQueryInfo(rand.Intn(10000), util.AddHex(encodedExtrinsic), hash)); err != nil {
		return
	}
	if v.CheckErr() != nil {
		return nil, InvalidValue
	}
	r := &PaymentQueryInfo{}
	return r, util.UnmarshalAny(&r, v.Result)
}

func GetPaymentQueryFeeDetails(_ context.Context, encodedExtrinsic, hash string) (feeDetails *PaymentQueryFeeDetails, err error) {
	v := &model.JsonRpcResult{}
	if err = websocket.SendWsRequest(nil, v, SystemPaymentQueryFeeDetails(rand.Intn(10000), util.AddHex(encodedExtrinsic), hash)); err != nil {
		return
	}
	if v.CheckErr() != nil {
		return nil, InvalidValue
	}
	return feeDetails, util.UnmarshalAny(&feeDetails, v.Result)
}

func GetExtrinsicFee(ctx context.Context, encodeExtrinsic, hash string, spec int, actualWeight, actualFeeByEvent decimal.Decimal, isV2Weight bool) (fee, actualFee decimal.Decimal, err error) {
	var paymentInfo = new(PaymentQueryInfo)
	feeDetails, err := GetPaymentQueryFeeDetails(ctx, encodeExtrinsic, hash)
	if !actualFeeByEvent.IsPositive() {
		paymentInfo, err = GetPaymentQueryInfo(ctx, spec, encodeExtrinsic, hash, isV2Weight)
		if err != nil || paymentInfo == nil {
			return decimal.Zero, actualFeeByEvent, err
		}
		if paymentInfo.Weight.IsZero() {
			return decimal.Zero, decimal.Zero, nil
		}
		if !actualFeeByEvent.IsNegative() {
			return paymentInfo.PartialFee, actualFeeByEvent, nil
		}
		if actualWeight.Equal(paymentInfo.Weight) {
			return paymentInfo.PartialFee, paymentInfo.PartialFee, nil
		}
	}
	if feeDetails == nil || feeDetails.InclusionFee == nil {
		return decimal.Zero, decimal.Zero, err
	}
	finalFee := feeDetails.EstimateFee()
	if paymentInfo.Weight.IsPositive() {
		actualFeeByEvent = feeDetails.ActualFee(paymentInfo.Weight, actualWeight)
	}
	return finalFee, actualFeeByEvent, nil
}
