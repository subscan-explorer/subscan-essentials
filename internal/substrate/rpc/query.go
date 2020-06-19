package rpc

import (
	"fmt"
	"github.com/itering/subscan/internal/pkg/recws"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/substrate/storage"
	"github.com/itering/subscan/internal/substrate/storageKey"
	"github.com/itering/subscan/internal/substrate/websocket"
	"github.com/shopspring/decimal"
	"math/rand"

	"github.com/itering/subscan/internal/util"
)

type Query interface {
	GetCurrentEra() (int, error)
	GetActiveEra() (int, error)
	StakingStakers(string, int) *storage.Exposures
	PowerOf(string) decimal.Decimal
	RewardPayee(string) (string, error)
}

type query struct {
	c    *recws.RecConn
	hash string
}

func New(c *recws.RecConn, hash string) Query {
	s := query{c: c, hash: hash}
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

func (r *query) PowerOf(stash string) (power decimal.Decimal) {
	return
}

func (r *query) StakingStakers(stash string, currentEra int) *storage.Exposures {
	exposure, err := ReadStorage(r.c, "Staking", "ErasStakers", "", util.U32Encode(currentEra), stash)
	if err != nil {
		return nil
	}
	return exposure.ToExposures()
}

func (r *query) RewardPayee(address string) (string, error) {
	value, err := ReadStorage(r.c, "Staking", "Payee", r.hash, address)
	if err != nil {
		return "", err
	}
	return value.ToString(), nil
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

func GetFreeBalance(c *recws.RecConn, accountId, hash string) (decimal.Decimal, decimal.Decimal, error) {
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

func GetValidatorFromSub(c *recws.RecConn, hash string) ([]string, error) {
	validators, err := ReadStorage(c, "Session", "Validators", hash)
	if err != nil {
		return []string{}, err
	}
	var r []string
	for _, address := range validators.ToStringSlice() {
		r = append(r, util.TrimHex(address))
	}
	return r, nil
}

func GetSystemProperties() (*Properties, error) {
	var t Properties
	v := &JsonRpcResult{}
	if err := websocket.SendWsRequest(nil, v, SystemProperties(rand.Intn(1000))); err != nil {
		return nil, err
	}
	err := v.ToAnyThing(&t)
	return &t, err
}
