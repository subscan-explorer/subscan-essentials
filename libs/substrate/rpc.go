package substrate

import (
	"encoding/json"
)

type RpcParam struct {
	Id      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	JsonRpc string      `json:"jsonrpc"`
}

func (rpc *RpcParam) structureQuery() []byte {
	rpc.JsonRpc = "2.0"
	b, _ := json.Marshal(rpc)
	return b
}

func SystemHealth(id int) []byte {
	rpc := RpcParam{Id: id, Method: "system_health", Params: []string{}}
	return rpc.structureQuery()
}

// connected peers
func SystemPeers(id int) []byte {
	rpc := RpcParam{Id: id, Method: "system_peers", Params: []string{}}
	return rpc.structureQuery()
}

func ChainGetBlock(id int, blockHash string) []byte {
	rpc := RpcParam{Id: id, Method: "chain_getBlock", Params: []string{blockHash}}
	return rpc.structureQuery()
}

func ChainGetBlockHash(id int, blockNum int) []byte {
	rpc := RpcParam{Id: id, Method: "chain_getBlockHash", Params: []int{blockNum}}
	return rpc.structureQuery()
}

func ChainGetRuntimeVersion(id int) []byte {
	rpc := RpcParam{Id: id, Method: "chain_getRuntimeVersion", Params: []string{}}
	return rpc.structureQuery()
}

func StateGetMetadata(id int, hash string) []byte {
	rpc := RpcParam{Id: id, Method: "state_getMetadata", Params: []string{hash}}
	return rpc.structureQuery()
}

func SystemProperties(id int) []byte {
	rpc := RpcParam{Id: id, Method: "system_properties", Params: []string{}}
	return rpc.structureQuery()
}

func SystemChain(id int) []byte {
	rpc := RpcParam{Id: id, Method: "system_chain", Params: []string{}}
	return rpc.structureQuery()
}

func SystemName(id int) []byte {
	rpc := RpcParam{Id: id, Method: "system_name", Params: []string{}}
	return rpc.structureQuery()
}

func SystemVersion(id int) []byte {
	rpc := RpcParam{Id: id, Method: "system_version", Params: []string{}}
	return rpc.structureQuery()
}

func ChainSubscribeNewHead(id int) []byte {
	rpc := RpcParam{Id: id, Method: "chain_subscribeNewHead", Params: []string{}}
	return rpc.structureQuery()
}

func StateSubscribeStorage(id int, storageKey []string) []byte {
	rpc := RpcParam{Id: id, Method: "state_subscribeStorage", Params: [][]string{storageKey}}
	return rpc.structureQuery()
}

func StateUnsubscribeStorage(id int, storageKey int) []byte {
	rpc := RpcParam{Id: id, Method: "state_unsubscribeStorage", Params: []int{storageKey}}
	return rpc.structureQuery()
}

func StateGetStorageAt(id int, storageKey, hash string) []byte {
	rpc := RpcParam{Id: id, Method: "state_getStorageAt", Params: []string{storageKey, hash}}
	return rpc.structureQuery()
}

func StateGetStorage(id int, storageKey string) []byte {
	rpc := RpcParam{Id: id, Method: "state_getStorage", Params: []string{storageKey}}
	return rpc.structureQuery()
}
