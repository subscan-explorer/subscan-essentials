package rpc

import (
	"encoding/json"
)

type Param struct {
	Id      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	JsonRpc string      `json:"jsonrpc"`
}

func (rpc *Param) structureQuery() []byte {
	rpc.JsonRpc = "2.0"
	b, _ := json.Marshal(rpc)
	return b
}

func SystemHealth(id int) []byte {
	rpc := Param{Id: id, Method: "system_health", Params: []string{}}
	return rpc.structureQuery()
}

func ChainGetBlock(id int, hash string) []byte {
	rpc := Param{Id: id, Method: "chain_getBlock", Params: []string{hash}}
	return rpc.structureQuery()
}

func ChainGetBlockHash(id int, blockNum int) []byte {
	rpc := Param{Id: id, Method: "chain_getBlockHash", Params: []int{blockNum}}
	return rpc.structureQuery()
}

func ChainGetRuntimeVersion(id int, hash ...string) []byte {
	rpc := Param{Id: id, Method: "chain_getRuntimeVersion", Params: hash}
	return rpc.structureQuery()
}

func StateGetMetadata(id int, hash ...string) []byte {
	rpc := Param{Id: id, Method: "state_getMetadata", Params: hash}
	return rpc.structureQuery()
}

func SystemProperties(id int) []byte {
	rpc := Param{Id: id, Method: "system_properties", Params: []string{}}
	return rpc.structureQuery()
}

func SystemChain(id int) []byte {
	rpc := Param{Id: id, Method: "system_chain", Params: []string{}}
	return rpc.structureQuery()
}

func SystemName(id int) []byte {
	rpc := Param{Id: id, Method: "system_name", Params: []string{}}
	return rpc.structureQuery()
}

func SystemVersion(id int) []byte {
	rpc := Param{Id: id, Method: "system_version", Params: []string{}}
	return rpc.structureQuery()
}

func ChainSubscribeNewHead(id int) []byte {
	rpc := Param{Id: id, Method: "chain_subscribeNewHead", Params: []string{}}
	return rpc.structureQuery()
}

func ChainSubscribeFinalizedHeads(id int) []byte {
	rpc := Param{Id: id, Method: "chain_subscribeFinalizedHeads", Params: []string{}}
	return rpc.structureQuery()
}

func StateSubscribeStorage(id int, storageKey []string) []byte {
	rpc := Param{Id: id, Method: "state_subscribeStorage", Params: [][]string{storageKey}}
	return rpc.structureQuery()
}

func AccountNonce(id int, address string) []byte {
	rpc := Param{Id: id, Method: "account_nextIndex", Params: []string{address}}
	return rpc.structureQuery()
}
func StateGetStorage(id int, storageKey string, hash string) []byte {
	rpc := Param{Id: id, Method: "state_getStorage", Params: []string{storageKey}}
	if hash != "" {
		rpc = Param{Id: id, Method: "state_getStorageAt", Params: []string{storageKey, hash}}
	}
	return rpc.structureQuery()
}

func StateGetKeysPaged(id int, storageKey string) []byte {
	rpc := Param{Id: id, Method: "state_getKeysPaged", Params: []interface{}{storageKey, 256, storageKey}}
	return rpc.structureQuery()
}

func SystemPaymentQueryInfo(id int, encodedExtrinsic string) []byte {
	rpc := Param{Id: id, Method: "payment_queryInfo", Params: []string{encodedExtrinsic}}
	return rpc.structureQuery()
}

func PowerOf(id int, address string) []byte {
	rpc := Param{Id: id, Method: "staking_powerOf", Params: []string{address}}
	return rpc.structureQuery()
}

// Query historical storage entries (by key) starting from a start block
// key Vec<StorageKey>
func StateQuerystorage(id int, key, start, end string) []byte {
	rpc := Param{Id: id, Method: "state_queryStorage", Params: []string{key, start, end}}
	return rpc.structureQuery()
}

//  Query storage entries (by key) starting at block hash given as the second parameter
// key Vec<StorageKey>
func StateQueryStorageAt(id int, key, start string) []byte {
	rpc := Param{Id: id, Method: "state_queryStorageAt", Params: []string{key, start}}
	return rpc.structureQuery()
}
