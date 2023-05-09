package staking

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strings"

	scale "github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	plugin "github.com/itering/subscan-plugin"
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan-plugin/storage"
	internalDao "github.com/itering/subscan/internal/dao"
	scanModel "github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins/staking/dao"
	"github.com/itering/subscan/plugins/staking/http"
	"github.com/itering/subscan/plugins/staking/model"
	"github.com/itering/subscan/plugins/staking/service"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
	"github.com/itering/substrate-api-rpc/rpc"
	rpcStorage "github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/storageKey"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
)

var srv *service.Service

type Staking struct {
	d  storage.Dao
	dd *internalDao.Dao
}

func New() *Staking {
	return &Staking{}
}

func (a *Staking) InitDao(d storage.Dao, dd *internalDao.Dao) {
	srv = service.New(d, dd)
	a.d = d
	a.dd = dd
	slog.Debug("staking init dao")
	a.Migrate()
}

func (a *Staking) InitHttp() []router.Http {
	return http.Router(srv)
}

func (a *Staking) ProcessExtrinsic(block *scanModel.ChainBlock, extrinsic *scanModel.ChainExtrinsic, events []scanModel.ChainEvent) error {
	slog.Debug("staking process extrinsic: %+v", extrinsic)
	return nil
}

type SS58Address string

func SS58AddressFromHex(hex string) SS58Address {
	return SS58Address(address.SS58Address(hex))
}

type GetNameValue interface {
	GetName() string
	GetValue
}

type GetValue interface {
	GetValue() interface{}
}

func CastUnnamedArg[T any](arg GetValue) (T, error) {
	var v T

	ty := reflect.TypeOf(v)
	ret := reflect.New(ty).Elem()
	switch ty.Kind() {
	case reflect.String:
		s, err := util.StringFromInterface(arg.GetValue())
		if err != nil {
			return v, err
		}

		if ty.Name() == "SS58Address" && strings.HasPrefix(s, "0x") {
			s = address.SS58Address(s)
		}
		ret.SetString(s)
		return ret.Interface().(T), nil
	case reflect.Uint32:
		switch arg.GetValue().(type) {
		case uint32:
			ret.SetUint(uint64(arg.GetValue().(uint32)))
			return ret.Interface().(T), nil
		case float64:
			ret.SetUint(uint64(arg.GetValue().(float64)))
			return ret.Interface().(T), nil
		}
		if val, ok := arg.GetValue().(uint32); ok {
			ret.SetUint(uint64(val))
			return ret.Interface().(T), nil
		}
		return v, fmt.Errorf("unexpected type. wanted uint32 but got %+v", reflect.TypeOf(arg.GetValue()))
		// u, err := util.IntFromInterface(arg.Value)
	default:
		name := ty.Name()
		switch name {
		case "Decimal":
			dec, err := util.MaybeDecimalFromInterface(arg.GetValue())
			if err != nil {
				return v, err
			}
			ret.Set(reflect.ValueOf(dec))
			return ret.Interface().(T), nil
		}
	}
	return v, fmt.Errorf("unsupported type: %s", ty)
}

func CastArg[T any](arg GetNameValue, name string) (T, error) {
	var v T
	if arg.GetName() != name {
		return v, fmt.Errorf("incorrect arg name: %s", arg.GetName())
	}
	return CastUnnamedArg[T](arg)
}

func switchName(a, b string) string {
	return fmt.Sprintf("%s.%s", strings.ToLower(a), strings.ToLower(b))

}

func (a *Staking) ProcessCall(block *scanModel.ChainBlock, call *scanModel.ChainCall, events []scanModel.ChainEvent, extrinsic *scanModel.ChainExtrinsic) error {
	slog.Info("staking process call", "call", call)

	if call == nil {
		return nil
	}

	name := switchName(call.ModuleId, call.CallId)
	slog.Info("staking process call", "name", name)
	switch name {
	case "staking.payout_stakers":
		validator, err := CastArg[SS58Address](call.Params[0], "validator_stash")
		if err != nil {
			return err
		}
		era, err := CastArg[uint32](call.Params[1], "era")
		if err != nil {
			return err
		}
		for _, e := range events {
			slog.Debug("event: %+v", e)
			var paramEvent []map[string]interface{}
			util.UnmarshalAny(&paramEvent, e.Params)
			args, err := getEventArgs(e.Params)
			if err != nil {
				return err
			}
			address, err := CastUnnamedArg[string](args[0])
			if err != nil {
				return err
			}
			amount, err := CastUnnamedArg[decimal.Decimal](args[1])
			if err != nil {
				return err
			}
			eventName := switchName(e.ModuleId, e.EventId)
			if eventName == "staking.rewarded" {
				dao.NewClaimedPayout(a.d, address, string(validator), amount, era, &e, block, extrinsic.ExtrinsicIndex)
			}
		}
		slog.Info("staking.payout_stakers", "validator", validator)
	}

	return nil
}

type EventArg struct {
	Type     string      `json:"type"`
	TypeName string      `json:"type_name"`
	Value    interface{} `json:"value"`
}

func (a EventArg) GetValue() interface{} {
	return a.Value
}

func getEventArgs(raw interface{}) ([]EventArg, error) {
	var paramEvent []map[string]interface{}
	util.UnmarshalAny(&paramEvent, raw)
	var args []EventArg
	for _, e := range paramEvent {
		arg, err := util.MapInterfaceAsStruct[EventArg](e)
		if err != nil {
			return []EventArg{}, err
		}
		args = append(args, arg)
	}
	return args, nil
}

// remove

func SendWsRequest(c websocket.WsConn, v interface{}, action []byte) (err error) {
	var p *websocket.PoolConn
	if c == nil {
		if p, err = websocket.Init(); err != nil {
			return
		}
		defer p.Close()
		c = p.Conn
	}
	if err = c.WriteMessage(1, action); err != nil {
		if p != nil {
			p.MarkUnusable()
		}
		return fmt.Errorf("websocket send error: %v", err)
	}
	if err = c.ReadJSON(v); err != nil {
		if p != nil {
			p.MarkUnusable()
		}
		slog.Error("websocket read error", "error", err)
		return
	}
	return nil
}

// 0x5f3e4907f716ac89b6347d15ececedca8bde0a0ea8864605e3b68ed9cb2da01b26ae334d66562f1490000000
// 0x5f3e4907f716ac89b6347d15ececedca8bde0a0ea8864605e3b68ed9cb2da01b26ae334d66562f1490000000
// 0x5f3e4907f716ac89b6347d15ececedca8bde0a0ea8864605e3b68ed9cb2da01b26ae334d66562f1490000000

func ReadStorage(p websocket.WsConn, module, prefix string, hash string, arg ...string) (r rpcStorage.StateStorage, err error) {
	key := storageKey.EncodeStorageKey(module, prefix, arg...)

	slog.Info("readstorage", "key", key)
	v := &rpc.JsonRpcResult{}
	if err = websocket.SendWsRequest(p, v, rpc.StateGetStorage(rand.Intn(10000), util.AddHex(key.EncodeKey), hash)); err != nil {
		return
	}
	slog.Info("readstorage", "response", fmt.Sprintf("%+v", v))
	if dataHex, err := v.ToString(); err == nil {
		if dataHex == "" {
			slog.Info("empty storage")
			return "", nil
		}

		return rpcStorage.Decode(dataHex, key.ScaleType, nil)
	}
	return r, err
}

// end
func structureQuery(param rpc.Param) []byte {
	param.JsonRpc = "2.0"
	b, _ := json.Marshal(param)
	return b
}

func StateGetKeysPagedAt(id int, storageKey string, at string) []byte {
	rpc := rpc.Param{Id: id, Method: "state_getKeysPaged", Params: []interface{}{storageKey, 256, nil, at}}
	return structureQuery(rpc)
}

func ReadKeysPaged(p websocket.WsConn, at, module, prefix string, args ...string) (r []string, scale string, err error) {
	key := storageKey.EncodeStorageKey(module, prefix, args...)
	slog.Debug("readkeys", "key", key)
	v := &rpc.JsonRpcResult{}
	if err = websocket.SendWsRequest(p, v, StateGetKeysPagedAt(rand.Intn(10000), util.AddHex(key.EncodeKey), at)); err != nil {
		return
	}
	if keys, err := v.ToInterfaces(); err == nil {
		for _, k := range keys {
			r = append(r, k.(string))
		}
	}
	return r, key.ScaleType, err
}

type EraStakerKey struct {
	StorageKey storageKey.StorageKey
	Era        uint32
	Validator  SS58Address
}

func eraStakerKey(key, scaleType string) EraStakerKey {
	keyBytes := util.HexToBytes(key)
	// first 32 bytes are the module and method hashes (both xx128 so 128 bit / 16 bytes each)
	keyBytes = keyBytes[32:]
	// era is xx64 concat, so first 64 bits / 8 bytes are the hash, then the next 4 bytes are the era
	keyBytes = keyBytes[8:]
	/*
				m := types.ScaleDecoder{}
		m.Init(scaleBytes.ScaleBytes{Data: util.HexToBytes(raw)}, option)
		return StateStorage(util.InterfaceToString(m.ProcessAndUpdateData(decodeType))), nil

	*/
	decoder := scale.ScaleDecoder{}
	eraBytes := keyBytes[:4]
	keyBytes = keyBytes[4:]
	decoder.Init(scaleBytes.ScaleBytes{Data: eraBytes}, nil)
	eraRaw := decoder.ProcessAndUpdateData("U32")
	era := eraRaw.(uint32)
	// staker is xx64 concat, so first 64 bits / 8 bytes are the hash, then the rest are the staker
	keyBytes = keyBytes[8:]
	decoder.Init(scaleBytes.ScaleBytes{Data: keyBytes}, nil)
	stakerRaw := decoder.ProcessAndUpdateData("AccountId")
	stakerString := stakerRaw.(string)
	staker := SS58AddressFromHex(stakerString)
	return EraStakerKey{StorageKey: storageKey.StorageKey{EncodeKey: key, ScaleType: scaleType}, Era: era, Validator: staker}
}

type EraInfo struct {
	Era              uint32
	TotalStake       decimal.Decimal
	Stakes           []EraStake
	TotalPoints      uint32
	TotalRewards     decimal.Decimal
	ValidatorPoints  map[SS58Address]uint32
	ValidatorRewards map[SS58Address]decimal.Decimal
	StakerRewards    map[SS58Address]decimal.Decimal
}

type EraStake struct {
	Validator      SS58Address
	Staker         SS58Address
	Amount         decimal.Decimal
	ValidatorTotal decimal.Decimal
}

/*
{
	total: 80
	individual: {
		5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY: 80
	}
}
*/

type StakeExposure struct {
	Total  string `json:"total"`
	Own    string `json:"own"`
	Others []struct {
		Who   string `json:"who"`
		Value string `json:"value"`
	} `json:"others"`
}

type EraPoints struct {
	Total      uint32 `json:"total"`
	Individual []struct {
		Col1 string `json:"col1"`
		Col2 uint32 `json:"col2"`
	} `json:"individual"`
}

func getEraInfo(era uint32, blockHash string, totalRewards decimal.Decimal) (EraInfo, error) {
	var bad EraInfo
	var stakes []EraStake
	eraEnc := scale.Encode("U32", era)
	keys, scaleType, err := ReadKeysPaged(nil, blockHash, "Staking", "ErasStakers", eraEnc)
	if err != nil {
		return bad, err
	}
	for _, key := range keys {
		k := eraStakerKey(key, scaleType)
		response, err := rpc.ReadStorageByKey(nil, k.StorageKey, blockHash)
		if err != nil {
			return bad, err
		}
		exposure := &StakeExposure{}
		err = json.Unmarshal([]byte(response.ToString()), exposure)
		if err != nil {
			return bad, err
		}
		validatorTotal, err := decimal.NewFromString(exposure.Total)
		if err != nil {
			return bad, err
		}
		ownAmount, err := decimal.NewFromString(exposure.Own)
		if err != nil {
			return bad, err
		}
		stakes = append(stakes, EraStake{Validator: k.Validator, Staker: k.Validator, Amount: ownAmount, ValidatorTotal: validatorTotal})

		for _, other := range exposure.Others {
			amount, err := decimal.NewFromString(other.Value)
			if err != nil {
				return bad, err
			}
			stakes = append(stakes, EraStake{Validator: k.Validator, Staker: SS58AddressFromHex(other.Who), Amount: amount, ValidatorTotal: validatorTotal})
		}
	}
	totalStakeRaw, err := ReadStorage(nil, "Staking", "ErasTotalStake", blockHash, eraEnc)
	if err != nil {
		return bad, err
	}
	totalStake := totalStakeRaw.ToDecimal()
	pointsRaw, err := ReadStorage(nil, "Staking", "ErasRewardPoints", blockHash, eraEnc)
	if err != nil {
		return bad, err
	}
	eraPoints := &EraPoints{}
	slog.Debug("getEraInfo", "pointsRaw", pointsRaw.ToString())
	err = json.Unmarshal([]byte(pointsRaw.ToString()), eraPoints)
	if err != nil {
		return bad, err
	}
	validatorPoints := make(map[SS58Address]uint32)
	for _, indiv := range eraPoints.Individual {
		validatorPoints[SS58AddressFromHex(indiv.Col1)] = indiv.Col2
	}

	validatorRewards := make(map[SS58Address]decimal.Decimal)
	stakerRewards := make(map[SS58Address]decimal.Decimal)
	totalPoints := decimal.NewFromInt(int64(eraPoints.Total))
	for v, points := range validatorPoints {
		share := decimal.NewFromInt(int64(points)).Div(totalPoints)
		slog.Debug("getEraInfo", "share", share, "totalRewards", totalRewards)
		validatorRewards[v] = totalRewards.Mul(share)
	}
	for _, stake := range stakes {
		share := stake.Amount.Div(stake.ValidatorTotal)
		stakerRewards[stake.Staker] = validatorRewards[stake.Validator].Mul(share)
	}
	return EraInfo{Era: era, TotalStake: totalStake, Stakes: stakes, TotalPoints: eraPoints.Total, ValidatorPoints: validatorPoints, ValidatorRewards: validatorRewards, StakerRewards: stakerRewards}, nil
}

func (a *Staking) ProcessEvent(block *scanModel.ChainBlock, event *scanModel.ChainEvent, fee decimal.Decimal, extrinsic *scanModel.ChainExtrinsic) error {
	name := switchName(event.ModuleId, event.EventId)
	args, err := getEventArgs(event.Params)
	switch name {
	case "staking.erapaid":
		if err != nil {
			return err
		}
		slog.Debug("staking.erapaid", "args", args)
		era, err := CastUnnamedArg[uint32](args[0])
		if err != nil {
			return err
		}
		reward, err := CastUnnamedArg[decimal.Decimal](args[1])
		if err != nil {
			return err
		}
		eraInfo, err := getEraInfo(era, block.Hash, reward)
		if err != nil {
			return err
		}
		slog.Info("staking.erapaid", "era", era, "reward", reward, "eraInfo", fmt.Sprintf("%+v", eraInfo))

	}
	return nil
}

func (a *Staking) SubscribeExtrinsic() []string {
	return []string{}
}

func (a *Staking) SubscribeCall() []string {
	return []string{"staking"}
}

func (a *Staking) SubscribeEvent() []string {
	return []string{"staking"}
}

func (a *Staking) Version() string {
	return "0.1"
}

func (a *Staking) UiConf() *plugin.UiConfig {
	conf := new(plugin.UiConfig)
	conf.Init()
	conf.Body.Api.Method = "post"
	conf.Body.Api.Url = "api/plugin/staking/accounts"
	conf.Body.Api.Adaptor = fmt.Sprintf(conf.Body.Api.Adaptor, "list")
	conf.Body.Columns = []plugin.UiColumns{
		{Name: "address", Label: "address"},
		{Name: "nonce", Label: "nonce"},
		{Name: "Staking", Label: "Staking"},
		{Name: "lock", Label: "lock"},
	}
	return conf
}

func (a *Staking) Migrate() {
	_ = a.d.AutoMigration(&model.Payout{})
}
