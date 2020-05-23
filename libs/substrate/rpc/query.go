package rpc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/libs/substrate"
	"github.com/itering/subscan/libs/substrate/storage"
	"github.com/itering/subscan/libs/substrate/storageKey"
	"github.com/itering/subscan/libs/substrate/websocket"
	"github.com/itering/subscan/pkg/recws"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/ss58"
	"github.com/shopspring/decimal"
)

type Query interface {
	GetCurrentEra() (int, error)
	GetActiveEra() (int, error)
	ValidatorPrefsValue(string) (int64, error)
	StakingStakers(string, int) *storage.Exposures
	PowerOf(string) decimal.Decimal
}

type query struct {
	c    *recws.RecConn
	hash string
}

func New(c *recws.RecConn, hash string) Query {
	s := query{c: c, hash: hash}
	switch util.NetworkNode {
	case util.Plasm:
		return &plasm{query{c: c, hash: hash}}
	case util.Edgeware:
		return &edg{query{c: c, hash: hash}}
	case util.CrabNetwork:
		return &crab{query{c: c, hash: hash}}
	}
	return &s
}

func (r *query) GetCurrentEra() (int, error) {
	eraIndex, err := ReadStorage(r.c, "Staking", "CurrentEra", r.hash)
	if err != nil {
		return 0, err
	}
	return eraIndex.ToInt(), nil
}

func (r *query) GetActiveEra() (int, error) {
	eraIndex, err := ReadStorage(r.c, "Staking", "ActiveEra", r.hash)
	if err != nil {
		return 0, err
	}
	if era := eraIndex.ToActiveEraInfo(); era != nil {
		return era.Index, nil
	}
	return 0, fmt.Errorf("decode ActiveEra error")
}

func (r *query) ValidatorPrefsValue(stash string) (int64, error) {
	validatorPrefs, err := ReadStorage(nil, "Staking", "Validators", "", stash)
	if err != nil {
		return 0, nil
	}
	validatorPrefsValue := validatorPrefs.ToValidatorPrefsLinkage()
	if validatorPrefsValue != nil {
		if validatorPrefsValue.ValidatorPrefs != nil {
			return validatorPrefsValue.ValidatorPrefs.Commission.IntPart(), nil
		}
		return validatorPrefsValue.Commission.IntPart(), nil
	}
	return 0, nil
}

func (r *query) PowerOf(stash string) (power decimal.Decimal) {
	return
}

func (r *query) StakingStakers(stash string, currentEra int) *storage.Exposures {
	exposure, err := ReadStorage(nil, "Staking", "ErasStakers", "", util.U32Encode(currentEra), stash)
	if err != nil {
		return nil
	}
	return exposure.ToExposures()
}

// Read substrate storage
func ReadStorage(c *recws.RecConn, module, prefix string, hash string, arg ...string) (r storage.StateStorage, err error) {
	key := storageKey.EncodeStorageKey(module, prefix, arg...)
	v := &JsonRpcResult{}
	if err = websocket.SendWsRequest(c, v, StateGetStorage(rand.Intn(10000), util.AddHex(key.EncodeKey), hash)); err != nil {
		return
	}
	if dataHex, err := v.ToString(); err == nil {
		if dataHex == "" {
			return storage.StateStorage(""), nil
		}
		return storage.Decode(dataHex, key.ScaleType, nil)
	}
	return r, err

}

func ReadKeysPaged(c *recws.RecConn, module, prefix string) (r []string, scale string, err error) {
	key := storageKey.EncodeStorageKey(module, prefix)
	v := &JsonRpcResult{}
	if err = websocket.SendWsRequest(c, v, StateGetKeysPaged(rand.Intn(10000), util.AddHex(key.EncodeKey))); err != nil {
		return
	}
	if keys, err := v.ToInterfaces(); err == nil {
		for _, k := range keys {
			r = append(r, k.(string))
		}
	}
	return r, key.ScaleType, err
}

func GetPaymentQueryInfo(c *recws.RecConn, encodedExtrinsic string) (paymentInfo *PaymentQueryInfo, err error) {
	v := &JsonRpcResult{}
	if err = websocket.SendWsRequest(c, v, SystemPaymentQueryInfo(rand.Intn(10000), util.AddHex(encodedExtrinsic))); err != nil {
		return
	}
	paymentInfo = v.ToPaymentQueryInfo()
	if paymentInfo == nil {
		return nil, fmt.Errorf("get PaymentQueryInfo error")
	}
	return
}

func ReadStorageByKey(c *recws.RecConn, key storageKey.StorageKey, hash string) (r storage.StateStorage, err error) {
	v := &JsonRpcResult{}
	if err = websocket.SendWsRequest(c, v, StateGetStorage(rand.Intn(10000), key.EncodeKey, hash)); err != nil {
		return
	}
	if dataHex, err := v.ToString(); err == nil {
		if dataHex == "" {
			return storage.StateStorage(""), nil
		}
		return storage.Decode(dataHex, key.ScaleType, nil)
	}
	return
}

func GetMetadataByHash(hash ...string) (string, error) {
	v := &JsonRpcResult{}
	if err := websocket.SendWsRequest(nil, v, StateGetMetadata(rand.Intn(10), hash...)); err != nil {
		return "", err
	}
	return v.ToString()
}

func GetFreeBalance(c *recws.RecConn, module, accountId, hash string) (decimal.Decimal, decimal.Decimal, error) {
	var accountValue storage.StateStorage
	var err error
	switch util.NetworkNode {
	case util.Edgeware:
		accountValue, err = ReadStorage(c, "Balances", "Account", hash, util.TrimHex(accountId))
		if err == nil {
			if account := accountValue.ToAccountData(); account != nil {
				return account.Free.Add(account.Reserved), decimal.Zero, nil
			}
		}
	case util.CrabNetwork:
		accountValue, err = ReadStorage(c, "System", "Account", hash, util.TrimHex(accountId))
		if err == nil {
			if account := accountValue.ToAccountInfo(); account != nil {
				return account.Data.Free.Add(account.Data.Reserved),
					account.Data.FreeKton.Add(account.Data.ReservedKton), nil
			}
		}
	default:
		accountValue, err = ReadStorage(c, "System", "Account", hash, util.TrimHex(accountId))
		if err == nil {
			if account := accountValue.ToAccountInfo(); account != nil {
				return account.Data.Free.Add(account.Data.Reserved), decimal.Zero, nil
			}
		}
	}
	return decimal.Zero, decimal.Zero, err
}

func GetAccountLock(c *recws.RecConn, address, currency string) (balance decimal.Decimal, err error) {
	var sv storage.StateStorage
	m := map[string]string{"ring": "Balances", "kton": "Kton", "": "Balances"}
	if module, ok := m[currency]; ok {
		sv, err = ReadStorage(c, module, "Locks", "", util.TrimHex(address))
		if err == nil {
			if locks := sv.ToBalanceLock(); len(locks) > 0 {
				for _, lock := range locks {
					switch util.NetworkNode {

					case util.Edgeware:
						return lock.Amount, nil

					case util.CrabNetwork:
						if lock.LockFor != nil {
							if lock.LockFor.StakingLock != nil {
								for _, unbonding := range lock.LockFor.StakingLock.Unbondings {
									balance = balance.Add(unbonding.Amount)
								}
								balance = balance.Add(lock.LockFor.StakingLock.StakingAmount).Add(balance)
							}
							if lock.LockFor.Common != nil {
								balance = balance.Add(lock.LockFor.Common.Amount)
							}
						}

					default: // kusama or other
						if lock.Amount.GreaterThanOrEqual(balance) {
							balance = lock.Amount
						}
					}
				}
				return balance, nil
			}
		}
	}
	return
}

func GetAccountNonce(c *recws.RecConn, address string) (int, error) {
	if util.StringInSlice(util.NetworkNode, []string{util.KusamaNetwork, util.CrabNetwork}) {
		v := &JsonRpcResult{}
		if err := websocket.SendWsRequest(c, v, AccountNonce(rand.Intn(10000), substrate.SS58Address(address))); err != nil {
			return 0, err
		}
		return util.IntFromInterface(v.Result), nil
	}
	nonce, err := ReadStorage(c, "System", "AccountNonce", "", util.TrimHex(address))
	if err != nil {
		return 0, err
	}
	return nonce.ToInt(), nil
}

func TotalIssuance(c *recws.RecConn, module string) (decimal.Decimal, error) {
	balanceValue, err := ReadStorage(c, module, "TotalIssuance", "")
	if err != nil {
		return decimal.Zero, err
	}
	return balanceValue.ToDecimal(), nil
}

func Nickname(c *recws.RecConn, address string) (string, error) {
	value, err := ReadStorage(c, "Nicks", "NameOf", "", util.TrimHex(address))
	if err != nil {
		return "", err
	}
	m := value.ToMapInterface()
	if m != nil {
		return m["col1"].(string), nil
	}
	return "", nil
}

func StakingValidatorCount(c *recws.RecConn) (int, error) {
	value, err := ReadStorage(c, "Staking", "ValidatorCount", "")
	if err != nil {
		return 0, err
	}
	return value.ToInt(), nil
}

func EpochIndex(c *recws.RecConn) (int, error) {
	value, err := ReadStorage(c, "Babe", "EpochIndex", "")
	if err != nil {
		return 0, err
	}
	return value.ToInt(), nil
}

func CurrentSlot(c *recws.RecConn) (int64, error) {
	value, err := ReadStorage(c, "Babe", "CurrentSlot", "")
	if err != nil {
		return 0, err
	}
	return value.ToInt64(), nil
}

func GenesisStartSlot(c *recws.RecConn) (int64, error) {
	value, err := ReadStorage(c, "Babe", "GenesisSlot", "")
	if err != nil {
		return 0, err
	}
	return value.ToInt64(), nil
}

func SessionCurrentIndex(c *recws.RecConn) (int, error) {
	value, err := ReadStorage(c, "Session", "CurrentIndex", "")
	if err != nil {
		return 0, err
	}
	return value.ToInt(), nil
}

func CurrentEraStartSessionIndex(c *recws.RecConn) (int, error) {
	switch util.NetworkNode {
	case util.KusamaNetwork, util.CrabNetwork, util.Acala:
		era, err := GetActiveEraEra(c, "")
		if err != nil {
			return 0, nil
		}
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(era.Index))
		value, err := ReadStorage(c, "Staking", "ErasStartSessionIndex", "", util.BytesToHex(bs))
		if err != nil {
			return 0, err
		}
		return value.ToInt(), nil

	default:
		value, err := ReadStorage(c, "Staking", "CurrentEraStartSessionIndex", "")
		if err != nil {
			return 0, err
		}
		return value.ToInt(), nil
	}
}

func StashController(c *recws.RecConn, address string) (string, error) {
	value, err := ReadStorage(c, "Staking", "Bonded", "", address)
	if err != nil {
		return "", err
	}
	return value.ToString(), nil
}

func RewardPayee(c *recws.RecConn, address string) (string, error) {
	value, err := ReadStorage(c, "Staking", "Payee", "", address)
	if err != nil {
		return "", err
	}
	return value.ToString(), nil
}

func IdentityOf(c *recws.RecConn, address string) (*storage.Registration, error) {
	value, err := ReadStorage(c, "Identity", "IdentityOf", "", address)
	if err != nil {
		return nil, err
	}
	return value.ToRegistration(), nil
}

func ReferendumInfoOf(c *recws.RecConn, hash string, referendumIndex uint) (*storage.ReferendumInfo, error) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, uint32(referendumIndex))
	value, err := ReadStorage(c, "Democracy", "ReferendumInfoOf", hash, util.BytesToHex(buf.Bytes()))
	if err != nil {
		return nil, err
	}
	return value.ToReferendumInfo(), nil
}

func GetActiveEraEra(c *recws.RecConn, hash string) (*storage.ActiveEraInfo, error) {
	eraIndex, err := ReadStorage(c, "Staking", "ActiveEra", hash)
	if err != nil {
		return nil, err
	}
	return eraIndex.ToActiveEraInfo(), nil
}

func StakingLedger(c *recws.RecConn, hash, controller string) *storage.StakingLedgers {
	ledger, err := ReadStorage(c, "Staking", "Ledger", hash, util.TrimHex(controller))
	if err != nil {
		return nil
	}
	return ledger.ToStakingLedgers()
}

func StakingValidators(c *recws.RecConn) ([]model.ValidatorPrefsMap, error) {
	var l []model.ValidatorPrefsMap
	if util.StringInSlice(util.NetworkNode, []string{util.KusamaNetwork, util.CrabNetwork, util.Plasm}) {
		keys, scaleType, err := ReadKeysPaged(c, "Staking", "Validators")
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			if value, err := ReadStorageByKey(c, storageKey.StorageKey{EncodeKey: key, ScaleType: scaleType}, ""); err == nil && value.ToValidatorPrefsLinkage() != nil {
				v := model.ValidatorPrefsMap{
					Address:             key[len(key)-64:],
					ValidatorPrefsValue: int(value.ToValidatorPrefsLinkage().Commission.IntPart()),
				}
				l = append(l, v)
			}
		}
	} else {
		var stakingValidator = func(c *recws.RecConn, validator string) (*storage.ValidatorPrefsLinkage, error) {
			value, err := ReadStorage(c, "Staking", "Validators", "", validator)
			if err != nil {
				return nil, err
			}
			return value.ToValidatorPrefsLinkage(), nil
		}

		value, err := ReadStorage(c, "Staking", "Validators", "")
		if err != nil {
			return nil, err
		}
		head := value.ToString()

		for {
			next, err := stakingValidator(c, head)
			if err != nil || next == nil || next.Linkage == nil || next.Linkage.Next == "0x0000000000000000000000000000000000000000000000000000000000000000" {
				break
			}
			v := model.ValidatorPrefsMap{
				Address:             head,
				ValidatorPrefsValue: int(next.ValidatorPrefs.Commission.IntPart()),
			}
			l = append(l, v)
			head = next.Linkage.Next
		}
	}
	return l, nil
}

// todo
func StakingNominators(c *recws.RecConn) error {
	keys, scaleType, err := ReadKeysPaged(c, "Staking", "Nominators")
	if err != nil {
		return err
	}
	fmt.Println(len(keys))
	for _, key := range keys {
		address := key[len(key)-64:]
		if value, err := ReadStorageByKey(c, storageKey.StorageKey{EncodeKey: key, ScaleType: scaleType}, ""); err == nil {
			fmt.Println(ss58.Encode(address, 42), value)
		}
	}
	return nil
}
