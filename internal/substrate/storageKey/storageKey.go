package storageKey

import (
	"fmt"
	"github.com/freehere107/go-scale-codec/types"
	"github.com/itering/subscan/internal/substrate/hasher"
	"github.com/itering/subscan/internal/substrate/metadata"
	"github.com/itering/subscan/util"
)

type StorageKey struct {
	EncodeKey string
	ScaleType string
}

func EncodeStorageKey(section, method string, args ...string) (storageKey StorageKey) {
	m := metadata.Latest(nil)
	if m == nil {
		return
	}
	method = util.UpperCamel(method)
	prefix, storageType := m.GetModuleStorageMapType(section, method)
	if storageType == nil {
		return
	}
	mapType := checkoutHasherAndType(storageType, args...)
	if mapType == nil {
		return
	}
	storageKey.ScaleType = mapType.Value
	var hash []byte
	if m.MetadataVersion >= 9 {
		method = dealLinkedMethod(method, mapType, args...)
		sectionHash := hasher.HashByCryptoName([]byte(util.UpperCamel(prefix)), "Twox128")
		methodHash := hasher.HashByCryptoName([]byte(method), "Twox128")
		hash = append(sectionHash, methodHash[:]...)
		if len(args) > 0 {
			var param []byte
			param = append(param, hasher.HashByCryptoName(util.HexToBytes(args[0]), mapType.Hasher)...)
			if len(args) >= 2 {
				param = append(param, hasher.HashByCryptoName(util.HexToBytes(args[1]), mapType.Hasher2)...)
			}
			hash = append(hash, param[:]...)
			storageKey.EncodeKey = util.BytesToHex(hash)
			return
		}
	} else {
		key := []byte(fmt.Sprintf("%s %s", []byte(util.UpperCamel(prefix)), method))
		if mapType.IsLinked && len(args) == 0 {
			key = append([]byte("head of "), key...)
		}
		if len(args) > 0 {
			param := encodeParams(args)
			key = append(key, param[:]...)
		}
		hash = hasher.HashByCryptoName(key, mapType.Hasher)
	}
	storageKey.EncodeKey = util.BytesToHex(hash)
	return
}

func encodeParams(args []string) []byte {
	if len(args) < 1 {
		return []byte{}
	}
	var data []byte
	for _, arg := range args {
		data = append(data, util.HexToBytes(arg)...)
	}
	return data
}

type storageOption struct {
	Value    string `json:"value"`
	Hasher   string `json:"hasher"`
	Hasher2  string `json:"hasher_2"`
	IsLinked bool   `json:"is_linked"`
}

func checkoutHasherAndType(t *types.StorageType, arg ...string) *storageOption {
	option := storageOption{}
	switch t.Origin {
	case "MapType":
		option.Value = t.MapType.Value
		option.Hasher = t.MapType.Hasher
		if option.IsLinked = t.MapType.IsLinked; option.IsLinked {
			if len(arg) == 0 && option.Value == "ValidatorPrefs" {
				option.Value = "AccountId" // waiting validator
			} else {
				option.Value = fmt.Sprintf("(%s, Linkage<AccountId>)", option.Value)
			}
		}
	case "DoubleMapType":
		option.Value = t.DoubleMapType.Value
		option.Hasher = t.DoubleMapType.Hasher
		option.Hasher2 = t.DoubleMapType.Key2Hasher
		option.IsLinked = t.DoubleMapType.IsLinked
	default:
		option.Value = *t.PlainType
		option.Hasher = "Twox64Concat"
	}
	return &option
}

func dealLinkedMethod(method string, t *storageOption, arg ...string) string {
	if t.IsLinked && len(arg) == 0 {
		method = fmt.Sprintf("HeadOf%s", method)
	}
	return method
}
