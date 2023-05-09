package util

import (
	"encoding/json"
	"math/rand"

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
		r, err := rpc.ReadStorage(p, module, prefix, hash, arg...)
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
