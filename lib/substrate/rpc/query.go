package rpc

import (
	"fmt"
	"github.com/itering/subscan/lib/substrate/storage"
	"github.com/itering/subscan/lib/substrate/storageKey"
	"github.com/itering/subscan/lib/substrate/websocket"
	"github.com/itering/subscan/pkg/recws"
	"github.com/shopspring/decimal"
	"math/rand"

	"github.com/itering/subscan/util"
)

type Query interface {
	GetCurrentEra() (int, error)
	GetActiveEra() (int, error)
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

func GetAccountLock(c *recws.RecConn, address string) (balance decimal.Decimal, err error) {
	var sv storage.StateStorage
	sv, err = ReadStorage(c, "Balances", "Locks", "", util.TrimHex(address))

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
	return
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
