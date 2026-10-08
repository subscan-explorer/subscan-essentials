package dao

import (
	"context"
	"fmt"
	"github.com/huandu/xstrings"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/model"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/storageKey"
	"github.com/itering/substrate-api-rpc/websocket"
	"math/rand"
)

func RemoveAddressPadded(address string) string {
	address = util.TrimHex(address)
	if len(address) < 40 {
		return NullAddress
	}
	return util.AddHex(address[len(address)-40:])
}

func h160ToAccountIdByNetwork(ctx context.Context, h160 string, _ string) (string, error) {
	return resolveH160Account(ctx, h160, util.IsEvmChain, metadata.Latest(nil) != nil && util.StringInSlice("Revive", metadata.SupportModule()), readReviveOriginalAccount)
}

func resolveH160Account(ctx context.Context, h160 string, isEvm, isRevive bool, lookup func(context.Context, string) (string, error)) (string, error) {
	if address.VerifyEthereumAddress(h160) {
		h160 = address.Format(h160)
		switch {
		case isEvm:
			return h160, nil
		case isRevive:
			mapped, err := lookup(ctx, h160)
			if err != nil {
				return "", err
			}
			if mapped != "" {
				if !address.VerifySubstrateAddress(mapped) {
					return "", fmt.Errorf("invalid Revive OriginalAccount for %s", h160)
				}
				return address.Format(mapped), nil
			}
			return reviveAccount(h160), nil
		default:
			return address.EvmToSS58Address(h160), nil
		}
	}
	return "", nil
}

func reviveAccount(h160 string) string {
	return xstrings.LeftJustify(util.TrimHex(h160), 64, "e")
}

// Resolve registered mappings before applying the stateless Ethereum fallback.
func readReviveOriginalAccount(ctx context.Context, h160 string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	key := storageKey.EncodeStorageKey("Revive", "OriginalAccount", h160)
	if key.EncodeKey == "" {
		return "", fmt.Errorf("Revive OriginalAccount storage unavailable in metadata")
	}
	result := &model.JsonRpcResult{}
	if err := websocket.SendWsRequest(nil, result, rpc.StateGetStorage(rand.Intn(10000), util.AddHex(key.EncodeKey), "")); err != nil {
		return "", err
	}
	if err := result.CheckErr(); err != nil {
		return "", err
	}
	if result.Result == nil {
		return "", nil
	}
	raw, ok := result.Result.(string)
	if !ok || !address.VerifySubstrateAddress(raw) {
		return "", fmt.Errorf("invalid Revive OriginalAccount response for %s", h160)
	}
	return address.Format(raw), nil
}
