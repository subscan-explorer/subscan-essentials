package substrate

import (
	"fmt"
	"subscan-end/libs/substrate/scalecodec/types"
	"subscan-end/utiles"
)

func encodeStorageKey(section, method string, arg ...string) (storageKey string, scaleType string) {
	metadata := InitMetaData()
	storageType, err := metadata.getModuleStorageMapType(section, method)
	if err != nil {
		return "", ""
	}
	mapType := chooseMapType(storageType)
	if mapType == nil {
		return "", ""
	}
	param := encodeParams(mapType["key"], arg)
	key := []byte(fmt.Sprintf("%s %s", section, method))
	if mapType["key"] != "" {
		key = append(key, param[:]...)
	}
	hash := hashBytesByHasher(key, mapType["hasher"])
	byteInstant := types.HexBytes{}
	byteInstant.Init(types.ScaleBytes{Data: compactAddLength(hash)}, "")
	byteInstant.Process()
	return byteInstant.Value.(string), convertType(mapType["value"])
}

func compactAddLength(b []byte) []byte {
	prefix := []byte{byte(len(b) << 2)}
	return append(prefix, b[:]...)
}

func encodeParams(scalecodeType string, arg []string) []byte {
	if len(arg) < 1 {
		return []byte{}
	}
	s := types.ScaleDecoder{}
	s.Init(types.ScaleBytes{Data: utiles.HexToBytes(arg[0])}, "")
	return utiles.HexToBytes(s.ProcessAndUpdateData(scalecodeType).(string))

}

func chooseMapType(t map[string]interface{}) map[string]string {
	if t["MapType"] != nil {
		mapType := t["MapType"].(map[string]interface{})
		return map[string]string{"key": mapType["key"].(string), "value": mapType["value"].(string), "hasher": mapType["hasher"].(string)}
	}

	if t["DoubleMapType"] != nil {
		doubleMapType := t["DoubleMapType"].(map[string]interface{})
		return map[string]string{"key": doubleMapType["key1"].(string), "key2": doubleMapType["key2"].(string), "value": doubleMapType["value"].(string), "hasher": doubleMapType["hasher"].(string), "key2Hasher": doubleMapType["key2Hasher"].(string)}
	}

	if t["PlainType"] != nil {
		return map[string]string{"hasher": "Twox128", "value": t["PlainType"].(string)}
	}
	return nil
}

func convertType(scaleType string) string {
	if utiles.NetworkNode == utiles.DarwiniaNetwork {
		if scaleType == "ValidatorPrefs" {
			scaleType = "ValidatorPrefsForDarwinia"
		}
	}
	return scaleType
}
