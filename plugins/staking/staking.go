package staking

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	scale "github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	plugin "github.com/itering/subscan-plugin"
	internalDao "github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/plugins/router"
	"github.com/itering/subscan/plugins/staking/dao"
	"github.com/itering/subscan/plugins/staking/http"
	"github.com/itering/subscan/plugins/staking/model"
	"github.com/itering/subscan/plugins/staking/service"
	"github.com/itering/subscan/plugins/storage"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
	"github.com/itering/substrate-api-rpc/storageKey"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
	"gorm.io/datatypes"
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

func (a *Staking) ProcessExtrinsic(block *storage.Block, extrinsic *storage.Extrinsic, events []storage.Event) error {
	slog.Debug("staking process extrinsic: %+v", extrinsic)
	return nil
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

	value := arg.GetValue()

	// gross special casing
	var anyv any = v
	switch anyv.(type) {
	case decimal.Decimal:
		dec, err := util.MaybeDecimalFromInterface(value)
		if err != nil {
			return v, err
		}
		ret.Set(reflect.ValueOf(dec))
		return ret.Interface().(T), nil
	}
	name := ty.Name()
	switch name {
	case "Decimal":
		dec, err := util.MaybeDecimalFromInterface(value)
		if err != nil {
			return v, err
		}
		ret.Set(reflect.ValueOf(dec))
		return ret.Interface().(T), nil
	}
	// end

	switch ty.Kind() {
	case reflect.String:
		s, err := util.StringFromInterface(value)
		if err != nil {
			return v, err
		}

		if ty.Name() == "SS58Address" && strings.HasPrefix(s, "0x") {
			s = address.SS58AddressFromHex(s).String()
		}
		ret.SetString(s)
		return ret.Interface().(T), nil
	case reflect.Uint32:
		switch value.(type) {
		case uint32:
			ret.SetUint(uint64(value.(uint32)))
			return ret.Interface().(T), nil
		case float64:
			ret.SetUint(uint64(value.(float64)))
			return ret.Interface().(T), nil
		}
		if val, ok := value.(uint32); ok {
			ret.SetUint(uint64(val))
			return ret.Interface().(T), nil
		}
		return v, fmt.Errorf("unexpected type. wanted uint32 but got %+v", reflect.TypeOf(value))
	// u, err := util.IntFromInterface(arg.Value)
	case reflect.Struct:
		return util.MapInterfaceAsStruct[T](value)
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

func (a *Staking) ProcessCall(block *storage.Block, call *storage.Call, events []storage.Event, extrinsic *storage.Extrinsic) error {
	slog.Debug("staking process call", "call", call)

	if call == nil {
		return nil
	}

	name := switchName(call.ModuleId, call.CallId)
	slog.Debug("staking process call", "name", name)
	switch name {
	case "staking.payout_stakers":
		validator, err := CastArg[address.SS58Address](call.Params[0], "validator_stash")
		if err != nil {
			return err
		}
		era, err := CastArg[uint32](call.Params[1], "era")
		if err != nil {
			return err
		}
		for _, e := range events {
			eventName := switchName(e.ModuleId, e.EventId)
			if eventName == "staking.rewarded" {
				var paramEvent []map[string]interface{}
				util.UnmarshalAny(&paramEvent, e.Params)
				args, err := GetEventArgs(e.Params)
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
				slog.Debug("staking.rewarded", "address", address, "amount", amount, "era", era, "event", e)
				dao.NewClaimedPayout(a.d, address, string(validator), amount, era, &e, block, extrinsic.ExtrinsicIndex)
			}
		}
		slog.Debug("staking.payout_stakers", "validator", validator)
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

func GetEventArgs(raw interface{}) ([]EventArg, error) {
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

type EraStakerKey struct {
	StorageKey storageKey.StorageKey
	Era        uint32
	Validator  address.SS58Address
}

func eraStakerKey(key, scaleType string) EraStakerKey {
	keyBytes := util.HexToBytes(key)
	// first 32 bytes are the module and method hashes (both xx128 so 128 bit / 16 bytes each)
	keyBytes = keyBytes[32:]
	// era is xx64 concat, so first 64 bits / 8 bytes are the hash, then the next 4 bytes are the era
	keyBytes = keyBytes[8:]
	eraBytes := keyBytes[:4]
	keyBytes = keyBytes[4:]
	decoder := scale.ScaleDecoder{}
	decoder.Init(scaleBytes.ScaleBytes{Data: eraBytes}, nil)
	eraRaw := decoder.ProcessAndUpdateData("U32")
	era := eraRaw.(uint32)
	// staker is xx64 concat, so first 64 bits / 8 bytes are the hash, then the rest are the staker
	keyBytes = keyBytes[8:]
	decoder.Init(scaleBytes.ScaleBytes{Data: keyBytes}, nil)
	stakerRaw := decoder.ProcessAndUpdateData("AccountId")
	stakerString := stakerRaw.(string)
	staker := address.SS58AddressFromHex(stakerString)
	return EraStakerKey{StorageKey: storageKey.StorageKey{EncodeKey: key, ScaleType: scaleType}, Era: era, Validator: staker}
}

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

func (a *Staking) getEraInfo(era uint32, blockHash string, totalRewards decimal.Decimal) (*model.EraInfo, error) {
	var bad *model.EraInfo
	var stakes []model.EraStake
	eraEnc := scale.Encode("U32", era)
	keys, scaleType, err := util.ReadKeysPaged(nil, blockHash, "Staking", "ErasStakersClipped", eraEnc)
	if err != nil {
		return bad, err
	}
	totalStakeRes := util.StartReadStorage(nil, "Staking", "ErasTotalStake", blockHash, eraEnc)
	pointsRes := util.StartReadStorage(nil, "Staking", "ErasRewardPoints", blockHash, eraEnc)
	for _, key := range keys {
		k := eraStakerKey(key, scaleType)
		response, err := util.ReadStorageByKey(nil, k.StorageKey, blockHash)
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
		stakes = append(stakes, model.EraStake{Validator: k.Validator, Staker: k.Validator, Amount: ownAmount, ValidatorTotal: validatorTotal})

		for _, other := range exposure.Others {
			amount, err := decimal.NewFromString(other.Value)
			if err != nil {
				return bad, err
			}
			stakes = append(stakes, model.EraStake{Validator: k.Validator, Staker: address.SS58AddressFromHex(other.Who), Amount: amount, ValidatorTotal: validatorTotal})
		}
	}
	totalStakeRaw, err := totalStakeRes.Wait()
	if err != nil {
		return bad, err
	}
	totalStake := totalStakeRaw.ToDecimal()
	pointsRaw, err := pointsRes.Wait()
	if err != nil {
		return bad, err
	}
	eraPoints := &EraPoints{}
	slog.Debug("getEraInfo", "pointsRaw", pointsRaw.ToString())
	err = json.Unmarshal([]byte(pointsRaw.ToString()), eraPoints)
	if err != nil {
		return bad, err
	}
	validatorPoints := make(map[address.SS58Address]uint32)
	for _, indiv := range eraPoints.Individual {
		validatorPoints[address.SS58AddressFromHex(indiv.Col1)] = indiv.Col2
	}

	validatorRewards := make(map[address.SS58Address]decimal.Decimal)
	validatorCuts := make(map[address.SS58Address]decimal.Decimal)
	stakerRewards := make(map[address.SS58Address]decimal.Decimal)
	totalPoints := decimal.NewFromInt(int64(eraPoints.Total))
	for v, points := range validatorPoints {
		prefs, err := dao.GetValidatorPrefs(a.d, v.String(), era)
		if err != nil {
			return bad, err
		}
		share, err := util.PerBillFromRational(decimal.NewFromInt(int64(points)), totalPoints)
		if err != nil {
			return bad, err
		}
		slog.Debug("updateEraInfo", "share", share, "totalRewards", totalRewards)

		validatorRewards[v] = share.Mul(totalRewards)
		validatorCuts[v] = validatorRewards[v].Mul(prefs.Commission.Round(9))
	}
	for _, stake := range stakes {
		totalRewardForValidator := validatorRewards[stake.Validator]
		cut := validatorCuts[stake.Validator]
		afterCut := totalRewardForValidator.Sub(cut)
		share, err := util.PerBillFromRational(stake.Amount, stake.ValidatorTotal)
		if err != nil {
			return bad, err
		}
		stakerReward := share.Mul(afterCut)
		if stake.Staker == stake.Validator {
			// validators take their cut
			stakerReward = stakerReward.Add(cut)
		}
		stakerRewards[stake.Staker] = stakerReward
	}
	return &model.EraInfo{Era: era, TotalStake: totalStake, Stakes: stakes, TotalPoints: eraPoints.Total, TotalRewards: totalRewards, ValidatorPoints: datatypes.NewJSONType(validatorPoints), ValidatorRewards: datatypes.NewJSONType(validatorRewards), StakerRewards: datatypes.NewJSONType(stakerRewards)}, nil
}

type ValidatorPrefs struct {
	Commission float64 `json:"commission"`
	Blocked    bool    `json:"blocked"`
}

func (a *Staking) ProcessEvent(block *storage.Block, event *storage.Event, fee decimal.Decimal, extrinsic *storage.Extrinsic) error {
	name := switchName(event.ModuleId, event.EventId)
	slog.Debug("process event", "module", event.ModuleId, "event", event.EventId, "block", block.Hash, "name", name, "event", fmt.Sprintf("%+v", event))
	switch name {
	case "staking.erapaid":
		slog.Debug("process event before get args", "event", fmt.Sprintf("%+v", event))
		args, err := GetEventArgs(event.Params)
		slog.Debug("process event after get args", "event", fmt.Sprintf("%+v", event))
		if len(args) < 2 {
			slog.Error("staking.erapaid: not enough arguments", "len", len(args), "args", args, "name", name, "event", fmt.Sprintf("%+v", event), "block", block.Hash)
			return fmt.Errorf("staking.erapaid: not enough arguments. got %d, expected at least 2", len(args))
		}
		if err != nil {
			return err
		}
		era, err := CastUnnamedArg[uint32](args[0])
		if err != nil {
			return err
		}
		reward, err := CastUnnamedArg[decimal.Decimal](args[1])
		if err != nil {
			return err
		}
		if era == 1 {
			eraInfo, err := dao.FindEraInfo(a.d, 0)
			if err != nil {
				slog.Error("error while completing era 0", "error", err)
				return err
			}
			if err = dao.CompleteEraInfo(a.d, eraInfo); err != nil {
				slog.Error("error while completing era 0", "error", err)
				return err
			}
		}
		eraInfo, err := a.getEraInfo(era, block.Hash, reward)
		if err != nil {
			return err
		}
		eraInfo.EndBlock = uint(block.BlockNum)
		dao.CompleteEraInfo(a.d, eraInfo)
		for _, stake := range eraInfo.Stakes {
			rewardAmount := eraInfo.StakerRewards.Data()[stake.Staker]
			if err := dao.NewUnclaimedPayout(a.d, stake.Staker, stake.Validator, rewardAmount, era); err != nil {
				slog.Error("failed to create new unclaimed payout", "error", err, "staker", stake.Staker.String(), "validator", stake.Validator.String(), "amount", rewardAmount.String(), "era", era)
			}
		}
		slog.Debug("staking.erapaid", "era", era, "reward", reward, "eraInfo", fmt.Sprintf("%+v", eraInfo))
		if err = dao.StartEraInfo(a.d, era+1, uint(block.BlockNum)); err != nil {
			return err
		}
		slog.Debug("new era", "era", era, "block", block.BlockNum)
	case "staking.validatorprefsset":
		args, err := GetEventArgs(event.Params)
		if err != nil {
			return err
		}
		if len(args) < 2 {
			return fmt.Errorf("staking.validatorprefsset: not enough arguments. got %d, expected at least 2", len(args))
		}
		slog.Debug("staking.validatorprefsset", "args", args)
		account, err := CastUnnamedArg[address.SS58Address](args[0])
		if err != nil {
			return err
		}
		prefs, err := CastUnnamedArg[ValidatorPrefs](args[1])
		if err != nil {
			return err
		}
		slog.Info("staking.validatorprefsset", "account", account, "prefs", prefs, "blockNum", block.BlockNum)
		// the commission is in parts per billion
		commission := decimal.NewFromFloat(prefs.Commission).Div(decimal.NewFromInt(1_000_000_000))
		if err := dao.NewValidatorPrefs(a.d, account, commission, prefs.Blocked, uint32(block.BlockNum)); err != nil {
			return err
		}
	case "nominationpools.paidout":
		args, err := GetEventArgs(event.Params)
		if err != nil {
			return err
		}
		if len(args) < 3 {
			return fmt.Errorf("nominationpools.paidout: not enough arguments. got %d, expected at least 3", len(args))
		}
		member, err := CastUnnamedArg[address.SS58Address](args[0])
		if err != nil {
			return err
		}
		poolId, err := CastUnnamedArg[uint32](args[1])
		if err != nil {
			return err
		}
		amount, err := CastUnnamedArg[decimal.Decimal](args[2])
		if err != nil {
			return err
		}
		if err := dao.NewPoolPayout(a.d, member, amount, poolId, event, block, extrinsic.ExtrinsicIndex); err != nil {
			return err
		}
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
	return []string{"staking", "nominationpools"}
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
	_ = a.d.AutoMigration(&model.PoolPayout{})
	_ = a.d.AutoMigration(&model.ValidatorPrefs{})
	_ = a.d.AutoMigration(&model.EraInfo{})

	a.d.Create(&model.EraInfo{Era: 0, StartBlock: 0})
}
