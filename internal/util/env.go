package util

import (
	"os"
)

const (
	CrabNetwork   = "crab"
	KusamaNetwork = "kusama"
	Edgeware      = "edgeware"
)

var (
	WSEndPoint  = GetEnv("CHAIN_WS_ENDPOINT", "wss://crab.darwinia.network")
	NetworkNode = GetEnv("NETWORK_NODE", CrabNetwork)
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func IsProduction() bool {
	return os.Getenv("DEPLOY_ENV") == "prod"
}
