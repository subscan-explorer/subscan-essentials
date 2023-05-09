package util

import (
	"os"
)

var (
	CurrentRuntimeSpecVersion int
	EventStorageKey           = GetEnv("SUBSTRATE_EVENT_KEY", "0x26aa394eea5630e07c48ae0c9558cef780d41e5e16056765bc8461851072c9d7")
	AddressType               = GetEnv("SUBSTRATE_ADDRESS_TYPE", "42")
	BalanceAccuracy           = GetEnv("SUBSTRATE_ACCURACY", "10")
	CommissionAccuracy        = GetEnv("COMMISSION_ACCURACY", "9")
	WSEndPoint                = GetEnv("CHAIN_WS_ENDPOINT", "ws://127.0.0.1:9944")
	NetworkNode               = GetEnv("NETWORK_NODE", "polkadot")
	IsProduction              = os.Getenv("DEPLOY_ENV") == "prod"
	ConfDir                   = GetEnv("CONF_DIR", "../configs")
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
