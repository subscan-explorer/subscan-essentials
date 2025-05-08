package util

import (
	"os"
)

var (
	// CurrentRuntimeSpecVersion current runtime spec version
	CurrentRuntimeSpecVersion int

	// AddressType ss58 address type, default is 0(polkadot)
	AddressType = GetEnv("SUBSTRATE_ADDRESS_TYPE", "0")
	// BalanceAccuracy balance accuracy, default is 10(DOT)
	BalanceAccuracy = GetEnv("SUBSTRATE_ACCURACY", "10")
	// WSEndPoint chain rpc endpoint, default is wss://polkadot-rpc.dwellir.com
	WSEndPoint = GetEnv("CHAIN_WS_ENDPOINT", "wss://rpc.polkadot.io")
	// NetworkNode network node name, default is polkadot
	NetworkNode = GetEnv("NETWORK_NODE", "polkadot")
	// ConfDir config directory, default is ../configs
	ConfDir = GetEnv("CONF_DIR", "../configs")

	// IsEvmChain is evm chain, address type is 0x h160
	IsEvmChain = StringInSlice(NetworkNode, []string{"moonbeam", "moonriver", "moonbase"})
)

// EventStorageKey state system.events storage key
const EventStorageKey = "0x26aa394eea5630e07c48ae0c9558cef780d41e5e16056765bc8461851072c9d7"

// GetEnv get env value by key, if not exist return default value
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
