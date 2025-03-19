package util

import (
	"os"
)

var (
	CurrentRuntimeSpecVersion int
	AddressType               = GetEnv("SUBSTRATE_ADDRESS_TYPE", "1")
	BalanceAccuracy           = GetEnv("SUBSTRATE_ACCURACY", "10")
	WSEndPoint                = GetEnv("CHAIN_WS_ENDPOINT", "wss://polkadot.api.onfinality.io/public-ws")
	NetworkNode               = GetEnv("NETWORK_NODE", "polkadot")
	IsProduction              = os.Getenv("DEPLOY_ENV") == "prod"
	ConfDir                   = GetEnv("CONF_DIR", "../configs")
)

const EventStorageKey = "0x26aa394eea5630e07c48ae0c9558cef780d41e5e16056765bc8461851072c9d7"

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
