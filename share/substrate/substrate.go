package substrate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/hasher"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/model"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/storageKey"
	"github.com/itering/substrate-api-rpc/websocket"
	"math/rand"
)

func DecodeExtrinsicParams(raw string, metadata *metadata.Instant, call *types.MetadataCalls, spec int) (params []scalecodec.ExtrinsicParam, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Recovering from panic in DecodeEventParams error is: %v \n", r)
		}
	}()
	e := types.ScaleDecoder{}
	m := types.MetadataStruct(*metadata)
	option := types.ScaleDecoderOption{Metadata: &m, Spec: spec}
	e.Init(scaleBytes.ScaleBytes{Data: util.HexToBytes(raw)}, &option)
	for _, arg := range call.Args {
		value := e.ProcessAndUpdateData(arg.Type)
		param := scalecodec.ExtrinsicParam{Type: arg.Type, Value: value, Name: arg.Name, TypeName: arg.TypeName}
		params = append(params, param)
	}
	return params, err
}

// DecodeEventParams decode event params
func DecodeEventParams(raw string, argsType []string, metadata *metadata.Instant, event *types.MetadataEvents, spec int) (params []scalecodec.EventParam, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Recovering from panic in DecodeEventParams error is: %v \n", r)
		}
	}()
	e := types.ScaleDecoder{}
	m := types.MetadataStruct(*metadata)
	option := types.ScaleDecoderOption{Metadata: &m, Spec: spec}
	e.Init(scaleBytes.ScaleBytes{Data: util.HexToBytes(raw)}, &option)
	for index, argType := range argsType {
		value := e.ProcessAndUpdateData(argType)
		param := scalecodec.EventParam{Type: argType, Value: value}
		if len(event.ArgsTypeName) == len(event.Args) {
			param.TypeName = event.ArgsTypeName[index]
		}
		if len(event.ArgsName) == len(event.Args) {
			param.Name = event.ArgsName[index]
		}
		params = append(params, param)
	}
	return params, err
}

func StateGetKeysPaged(storageKey, start, hash string, row int) []byte {
	params := []interface{}{storageKey, row, start}
	if hash != "" {
		params = append(params, hash)
	}
	p := rpc.Param{Id: rand.Intn(10000), Method: "state_getKeysPaged", Params: params}
	p.JsonRpc = "2.0"
	b, _ := json.Marshal(p)
	return b
}

func StateQueryStorageAt(key []string, start string) []byte {
	params := []interface{}{key}
	if start != "" {
		params = append(params, start)
	}
	p := rpc.Param{Id: rand.Intn(10000), Method: "state_queryStorageAt", Params: params}
	p.JsonRpc = "2.0"
	b, _ := json.Marshal(p)
	return b
}

func BatchReadKeysPaged(_ context.Context, module, prefix string, hash string, action func(keys []string, scaleType string) error, arg ...string) (err error) {
	key := storageKey.EncodeStorageKey(module, prefix, arg...)
	if key.EncodeKey == "" {
		err = fmt.Errorf("storageKey not encode with %s %s", module, prefix)
		return
	}
	start := util.AddHex(key.EncodeKey)
	for {
		var keys []any
		v := &model.JsonRpcResult{}
		if err = websocket.SendWsRequest(nil, v, StateGetKeysPaged(util.AddHex(key.EncodeKey), start, hash, 256)); err != nil {
			break
		}
		if err = v.CheckErr(); err != nil {
			break
		}
		_ = util.UnmarshalAny(&keys, v.Result)
		result := make([]string, 0, len(keys))
		for _, k := range keys {
			result = append(result, k.(string))
			start = k.(string)
		}
		if err = action(result, key.ScaleType); err != nil {
			return err
		}
		if len(keys) < 256 {
			break
		}
	}
	return
}

type StateStorageResult struct {
	Block   string     `json:"block"`
	Changes [][]string `json:"changes"`
}

func BatchStorageByKey(_ context.Context, keys []string, scaleType string, hash string) (r map[string]storage.StateStorage, err error) {
	var data []StateStorageResult
	v := &model.JsonRpcResult{}
	if err = websocket.SendWsRequest(nil, v, StateQueryStorageAt(keys, hash)); err != nil {
		return
	}
	if err = v.CheckErr(); err != nil {
		return
	}
	_ = util.UnmarshalAny(&data, v.Result)
	r = make(map[string]storage.StateStorage)
	if len(data) > 0 {
		storageChanges := data[0].Changes
		for _, changes := range storageChanges {
			var state storage.StateStorage
			if len(changes[1]) > 0 {
				state, _, _ = storage.Decode(changes[1], scaleType, &types.ScaleDecoderOption{})
			}
			r[changes[0]] = state
		}
	}
	return
}

func ParseStorageKey(key string) (KeyStorage, error) {
	_decodeType, hashers, err := CheckoutStorageKey(key)
	if err != nil {
		return nil, err
	}
	return DecodeStorageKey(key, _decodeType, hashers, nil)
}

type KeyStorage []storage.StateStorage

type KeyValueStorage struct {
	Keys   []KeyStorage
	Values []storage.StateStorage
}

func (k *KeyValueStorage) Put(key KeyStorage, value storage.StateStorage) {
	k.Keys = append(k.Keys, key)
	k.Values = append(k.Values, value)
}

func (k *KeyValueStorage) GetValue(index int) storage.StateStorage {
	return k.Values[index]
}

func (k *KeyValueStorage) GetKey(index int) KeyStorage {
	return k.Keys[index]
}

func (k *KeyValueStorage) Foreach(fn func(key KeyStorage, value storage.StateStorage)) {
	for i := range k.Keys {
		fn(k.Keys[i], k.Values[i])
	}
}

type StorageOption struct {
	Value  string   `json:"value"`
	Keys   []string `json:"keys"`
	Hasher []string `json:"hasher"`
}

func (o *StorageOption) KeyHasher() map[string]string {
	r := make(map[string]string)
	for i, key := range o.Keys {
		r[key] = o.Hasher[i]
	}
	return r
}

func CheckoutHasherAndType(t *types.StorageType) *StorageOption {
	option := StorageOption{}
	switch t.Origin {
	case "MapType":
		option.Keys = []string{t.MapType.Key}
		option.Value = t.MapType.Value
		option.Hasher = []string{t.MapType.Hasher}
	case "DoubleMapType":
		option.Value = t.DoubleMapType.Value
		option.Keys = []string{t.DoubleMapType.Key, t.DoubleMapType.Key2}
		option.Hasher = []string{t.DoubleMapType.Hasher, t.DoubleMapType.Key2Hasher}
	case "Map":
		option.Value = t.NMapType.Value
		option.Keys = t.NMapType.KeyVec
		option.Hasher = t.NMapType.Hashers
	default:
		option.Value = *t.PlainType
		option.Hasher = []string{"Twox64Concat"}
	}
	return &option
}

func CheckoutStorageKey(raw string) (keys []string, hashers []string, err error) {
	key := util.HexToBytes(raw)
	m := metadata.Latest(nil)
	moduleKey := key[:16]
	methodKey := key[16:32]
	var module types.MetadataModules
	var method types.MetadataStorage
	for _, mm := range m.Metadata.Modules {
		if bytes.EqualFold(hasher.HashByCryptoName([]byte(mm.Prefix), "Twox128"), moduleKey) {
			module = mm
			for _, storage := range mm.Storage {
				if bytes.EqualFold(hasher.HashByCryptoName([]byte(storage.Name), "Twox128"), methodKey) {
					method = storage
					break
				}
			}
			break
		}
	}

	if module.Name == "" || method.Name == "" {
		err = fmt.Errorf("module or method not found")
		return
	}

	if mapType := CheckoutHasherAndType(&method.Type); mapType != nil {
		return mapType.Keys, mapType.Hasher, nil
	}

	err = fmt.Errorf("keys not found")
	return
}

const moduleMethodSize = 32

func DecodeStorageKey(raw string, decodeTypes []string, hashers []string, option *types.ScaleDecoderOption) (ks KeyStorage, err error) {
	if len(decodeTypes) == 0 {
		return nil, nil
	}
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(error); ok {
				err = fmt.Errorf("Recovering from panic in DecodeStorageKey key %serror is: %v \n", err, raw)
			} else {
				err = fmt.Errorf("Recovering from panic in DecodeStorageKey key %s error is: %v \n", raw, r)
			}
		}
	}()

	var preSize = moduleMethodSize
	rawHex := util.HexToBytes(raw)[preSize:]
	var offset int
	var r KeyStorage
	for i, h := range hashers {
		offset += Size(h)
		v, length, err := storage.Decode(util.BytesToHex(rawHex[offset:]), decodeTypes[i], option)
		if err != nil {
			return nil, err
		}
		r = append(r, v)
		offset += length
	}
	return r, nil
}

var hashSize = map[string]int{
	"Blake2_128":       16,
	"Blake2_256":       32,
	"Twox128":          16,
	"Twox256":          32,
	"Twox64Concat":     8,
	"Identity":         0,
	"Blake2_128Concat": 16,
}

const defaultHashSize = 8

func Size(hasher string) int {
	if hasher == "" {
		return 0
	}

	if size, ok := hashSize[hasher]; ok {
		return size
	}

	return defaultHashSize
}
