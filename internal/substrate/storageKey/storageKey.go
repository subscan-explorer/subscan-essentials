package storageKey

import (
	"fmt"
	"github.com/itering/scale.go/types"
	"github.com/itering/subscan/internal/substrate/hasher"
	"github.com/itering/subscan/internal/substrate/metadata"
	"github.com/itering/subscan/internal/util"
	"strings"
)

type StorageKey struct {
	EncodeKey string
	ScaleType string
}

type Storage struct {
	Prefix string
	Method string
	Type   types.StorageType
}

var (
	Sks           map[string]Storage
	TotalIssuance StorageKey
)

// PrintRuntimeStorageKey
func runtimeStorageKey() map[string]Storage {
	if Sks == nil {
		return Sks
	}
	runtime := metadata.Latest(nil)
	keys := make(map[string]Storage)
	for _, modules := range runtime.Metadata.Modules {
		for _, storage := range modules.Storage {
			prefix := modules.Prefix
			method := storage.Name
			keys[strings.ToLower(fmt.Sprintf("%s|%s", modules.Name, method))] = Storage{
				Prefix: util.UpperCamel(prefix),
				Method: util.UpperCamel(method),
				Type:   storage.Type,
			}
		}
	}
	Sks = keys
	return Sks
}

func SubscribeStorage() []string {
	TotalIssuance = EncodeStorageKey("Balances", "TotalIssuance")
	return []string{util.AddHex(TotalIssuance.EncodeKey)}

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

	method = dealLinkedMethod(method, mapType, args...)

	sectionHash := hasher.HashByCryptoName([]byte(util.UpperCamel(prefix)), "Twox128")
	methodHash := hasher.HashByCryptoName([]byte(method), "Twox128")

	hash = append(sectionHash, methodHash[:]...)

	if len(args) > 0 {
		var param []byte
		param = append(param, hasher.HashByCryptoName(util.HexToBytes(args[0]), mapType.Hasher)...)
		if len(args) == 2 {
			param = append(param, hasher.HashByCryptoName(util.HexToBytes(args[1]), mapType.Hasher2)...)
		}
		hash = append(hash, param[:]...)
	}
	storageKey.EncodeKey = util.BytesToHex(hash)
	return
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
