package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/hasher"
	"time"
)

func GetContractMethod(ctx context.Context, field ...string) map[string]string {
	if len(field) == 0 {
		return nil
	}
	cacheKey := fmt.Sprintf("%s:evm:contract:methodID", util.NetworkNode)
	methods := sg.redis.HMGet(ctx, cacheKey, field...)
	ttl := sg.redis.GetCacheTtl(ctx, cacheKey)
	if ttl == -2 && len(methods) == 0 {
		methods = WriteContractMethod(ctx, cacheKey)
	}
	if ttl > 0 && ttl < 3600 {
		go func() {
			// refresh cache
			c, cancel := context.WithTimeout(context.Background(), time.Second*60)
			defer cancel()
			WriteContractMethod(c, cacheKey)
		}()
	}
	return methods
}

func WriteContractMethod(ctx context.Context, key string) map[string]string {
	methodMap := make(map[string]string)
	methodID, err := ContractMethodList(ctx)
	if err != nil {
		return methodMap
	}
	for _, method := range methodID {
		var contMethod map[string]string
		_ = json.Unmarshal(method, &contMethod)
		for fn, hash := range contMethod {
			methodMap[hash] = fn
		}
	}
	_ = sg.redis.HmSetEx(ctx, key, methodMap, 3600*6)
	return methodMap
}

func DoBlake2_256(data string) string {
	return util.AddHex(util.BytesToHex(hasher.HashByCryptoName(util.HexToBytes(data), "Blake2_256")))
}
