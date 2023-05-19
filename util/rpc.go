package util

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/itering/substrate-api-rpc/rpc"
	rpcStorage "github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/storageKey"
	"github.com/itering/substrate-api-rpc/util"
	"github.com/itering/substrate-api-rpc/websocket"
	"golang.org/x/exp/slog"
)

type Result[T any] struct {
	Value T
	Error error
}

type FutureResult[T any] chan Result[T]

func (f FutureResult[T]) Wait() (T, error) {
	r := <-f
	return r.Value, r.Error
}

func StartReadStorage(p websocket.WsConn, module, prefix string, hash string, arg ...string) (ch FutureResult[rpcStorage.StateStorage]) {
	ch = make(chan Result[rpcStorage.StateStorage], 1)
	go func() {
		r, err := ReadStorage(p, module, prefix, hash, arg...)
		ch <- Result[rpcStorage.StateStorage]{r, err}
		close(ch)
	}()
	return ch
}

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
	res, err := SendWsRequest(p, StateGetKeysPagedAt(rand.Intn(10000), util.AddHex(key.EncodeKey), at))
	if err != nil {
		return
	}
	if keys, err := res.ToInterfaces(); err == nil {
		for _, k := range keys {
			r = append(r, k.(string))
		}
	}
	return r, key.ScaleType, err
}

func ReadStorage(p websocket.WsConn, module, prefix string, hash string, arg ...string) (r rpcStorage.StateStorage, err error) {
	key := storageKey.EncodeStorageKey(module, prefix, arg...)
	res, err := SendWsRequest(p, rpc.StateGetStorage(rand.Intn(10000), util.AddHex(key.EncodeKey), hash))
	if err != nil {
		return
	}
	if dataHex, err := res.ToString(); err == nil {
		if dataHex == "" {
			return "", nil
		}
		return rpcStorage.Decode(dataHex, key.ScaleType, nil)
	}
	return r, err
}

func SendWsRequest(conn websocket.WsConn, action []byte) (rpc.JsonRpcResult, error) {
	return WithRetriesAndTimeout(time.Second*5, 5, func() (rpc.JsonRpcResult, error) {
		v := &rpc.JsonRpcResult{}
		e := websocket.SendWsRequest(conn, v, action)
		return *v, e
	})
}

func ReadStorageByKey(p websocket.WsConn, key storageKey.StorageKey, hash string) (r rpcStorage.StateStorage, err error) {
	res, err := SendWsRequest(p, rpc.StateGetStorage(rand.Intn(10000), key.EncodeKey, hash))
	if err != nil {
		return
	}
	if dataHex, err := res.ToString(); err == nil {
		if dataHex == "" {
			return rpcStorage.StateStorage(""), nil
		}
		return rpcStorage.Decode(dataHex, key.ScaleType, nil)
	}
	return
}

type Properties struct {
	Ss58Format    *int    `json:"ss58Format"`
	TokenDecimals *int    `json:"tokenDecimals"`
	TokenSymbol   *string `json:"tokenSymbol"`
}

func GetSystemProperties(p websocket.WsConn) (*Properties, error) {
	var t Properties
	res, err := SendWsRequest(p, rpc.SystemProperties(rand.Intn(1000)))
	if err != nil {
		return nil, err
	}
	err = res.ToAnyThing(&t)
	return &t, err
}
