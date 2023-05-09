package service

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
)

func (s *Service) createExtrinsic(c context.Context,
	txn *dao.GormDB,
	block *model.ChainBlock,
	encodeExtrinsics []string,
	decodeExtrinsics []map[string]interface{},
	eventMap map[string][]model.ChainEvent,
) (int, int, map[string]string, map[string]decimal.Decimal, error) {

	var (
		blockTimestamp int
		e              []model.ChainExtrinsic
		err            error
	)
	extrinsicFee := make(map[string]decimal.Decimal)

	eb, _ := json.Marshal(decodeExtrinsics)
	_ = json.Unmarshal(eb, &e)

	hash := make(map[string]string)

	for index, extrinsic := range e {
		extrinsic.CallModule = strings.ToLower(extrinsic.CallModule)
		extrinsic.BlockNum = block.BlockNum
		extrinsic.ExtrinsicIndex = fmt.Sprintf("%d-%d", extrinsic.BlockNum, index)
		extrinsic.Success = s.getExtrinsicSuccess(eventMap[extrinsic.ExtrinsicIndex])

		if tp := s.getTimestamp(&extrinsic); tp > 0 {
			blockTimestamp = tp
		}
		extrinsic.BlockTimestamp = blockTimestamp
		if extrinsic.ExtrinsicHash != "" {

			fee, _ := GetExtrinsicFee(nil, encodeExtrinsics[index])
			extrinsic.Fee = fee

			extrinsicFee[extrinsic.ExtrinsicIndex] = fee
			hash[extrinsic.ExtrinsicIndex] = extrinsic.ExtrinsicHash
		}

		if err = s.dao.CreateExtrinsic(c, txn, &extrinsic); err == nil {
			go s.emitExtrinsic(block, &extrinsic, eventMap[extrinsic.ExtrinsicIndex])
		} else {
			return 0, 0, nil, nil, err
		}

		s.handleCalls(block, &extrinsic, eventMap[extrinsic.ExtrinsicIndex])
	}
	return len(e), blockTimestamp, hash, extrinsicFee, err
}

type Call TypedCall[[]model.CallArg]

type TypedCall[T any] struct {
	CallIndex  string `json:"call_index"`
	CallModule string `json:"call_module"`
	CallName   string `json:"call_name"`
	Params     T      `json:"params"`
}

type HasNameValue interface {
	GetName() string
	GetValue() interface{}
}

func (a ExtrinsicArg) GetName() string {
	return a.Name
}
func (a ExtrinsicArg) GetValue() interface{} {
	return a.Value
}

func makeArgMap[A HasNameValue](args []A) map[string]interface{} {
	argMap := make(map[string]interface{})
	for _, arg := range args {
		argMap[arg.GetName()] = arg.GetValue()
	}
	return argMap
}

func formatModuleFunction(module, function string) string {
	return fmt.Sprintf("%s.%s", strings.ToLower(module), strings.ToLower(function))
}

func (s *Service) handleCalls(block *model.ChainBlock, extrinsic *model.ChainExtrinsic, events []model.ChainEvent) {
	switch formatModuleFunction(extrinsic.CallModule, extrinsic.CallModuleFunction) {
	case "utility.batch":
		args, err := extrinsicArgs(extrinsic)
		if err != nil {
			slog.Error("utility.batch extrinsicArgs error:", err)
			return
		}
		argMap := makeArgMap(args)
		callsList := argMap["calls"].([]interface{})
		calls := make([]Call, 0)
		for _, arg := range callsList {
			call, err := util.MapInterfaceAsStruct[Call](arg)
			if err != nil {
				slog.Error("utility.batch Call error:", err)
				return
			}
			calls = append(calls, call)
		}
		slog.Debug("utility.batch calls:", calls)

		callEvents := make([][]model.ChainEvent, len(calls))
		eventIdx := 0
		for i := 0; i < len(calls); i++ {
			for ; eventIdx < len(events); eventIdx++ {
				ev := events[eventIdx]
				if e := formatModuleFunction(ev.ModuleId, ev.EventId); e == "utility.itemcompleted" {
					eventIdx++
					break
				}
				callEvents[i] = append(callEvents[i], ev)
			}
		}

		for i := 0; i < len(calls); i++ {
			call := calls[i]
			callEvents := callEvents[i]
			slog.Debug("call:", call)
			slog.Debug("callEvents:", callEvents)
			chainCall := &model.ChainCall{
				BlockNum:       block.BlockNum,
				ExtrinsicHash:  extrinsic.ExtrinsicHash,
				CallIdx:        i,
				BlockTimestamp: block.BlockTimestamp,
				ModuleId:       strings.ToLower(call.CallModule),
				CallId:         call.CallName,
				Events:         callEvents,
				Params:         call.Params,
			}
			slog.Info("emitting call", "call", chainCall)
			s.emitCall(block, chainCall, callEvents, extrinsic)
		}
	}
}

func extrinsicArgs(extrinsic *model.ChainExtrinsic) ([]ExtrinsicArg, error) {
	var args []ExtrinsicArg
	if extrinsic.Params == nil {
		return args, nil
	}
	paramsMap, ok := extrinsic.Params.([]interface{})
	if !ok {
		return args, fmt.Errorf("invalid extrinsic params. found %+v: %+v", reflect.TypeOf(extrinsic.Params), extrinsic.Params)
	}
	for _, param := range paramsMap {
		arg, err := util.MapInterfaceAsStruct[ExtrinsicArg](param)
		if err != nil {
			return args, err
		}
		args = append(args, arg)
	}
	return args, nil
}

// [map[name:calls type:Vec<Call> type_name:Vec<<T as Config>::RuntimeCall>
// value:[map[call_index:0612 call_module:Staking call_name:payout_stakers
// params:[map[name:validator_stash type:[U8; 32] value:0xbe5ddb1579b72e84524fc29e78609e3caf42e85aa118ebfe0b0ad404b5bdd25f]
// map[name:era type:U32 value:2326]]]

type ExtrinsicArg struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	TypeName string      `json:"type_name"`
	Value    interface{} `json:"value"`
}

func (s *Service) ExtrinsicsAsJson(e *model.ChainExtrinsic) *model.ChainExtrinsicJson {
	ej := &model.ChainExtrinsicJson{
		BlockNum:           e.BlockNum,
		BlockTimestamp:     e.BlockTimestamp,
		ExtrinsicIndex:     e.ExtrinsicIndex,
		ExtrinsicHash:      e.ExtrinsicHash,
		Success:            e.Success,
		CallModule:         e.CallModule,
		CallModuleFunction: e.CallModuleFunction,
		Params:             util.ToString(e.Params),
		AccountId:          address.SS58Address(e.AccountId),
		Signature:          e.Signature,
		Nonce:              e.Nonce,
		Fee:                e.Fee,
	}
	var paramsInstant []model.ExtrinsicParam
	if err := json.Unmarshal([]byte(ej.Params), &paramsInstant); err != nil {
		for pi, param := range paramsInstant {
			if paramsInstant[pi].Type == "Address" {
				paramsInstant[pi].Value = address.SS58Address(param.Value.(string))
			}
		}
		bp, _ := json.Marshal(paramsInstant)
		ej.Params = string(bp)
	}
	return ej
}

func (s *Service) getTimestamp(extrinsic *model.ChainExtrinsic) (blockTimestamp int) {
	if extrinsic.CallModule != "timestamp" {
		return
	}

	var paramsInstant []model.ExtrinsicParam
	util.UnmarshalAny(&paramsInstant, extrinsic.Params)

	for _, p := range paramsInstant {
		if p.Name == "now" {
			if strings.EqualFold(p.Type, "compact<U64>") {
				return int(util.Int64FromInterface(p.Value) / 1000)
			}
			extrinsic.BlockTimestamp = util.IntFromInterface(p.Value)
			return extrinsic.BlockTimestamp
		}
	}
	return
}

func (s *Service) getExtrinsicSuccess(e []model.ChainEvent) bool {
	for _, event := range e {
		if strings.EqualFold(event.ModuleId, "system") {
			return strings.EqualFold(event.EventId, "ExtrinsicFailed")
		}
	}
	return true
}
