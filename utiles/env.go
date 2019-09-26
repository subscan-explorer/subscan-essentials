package utiles

import (
	"os"
)

const (
	DarwiniaNetwork = "darwinia"
	KusamaNetwork   = "kusama"
)

var (
	Environment      string
	Dev              = "dev"
	ProviderEndPoint = GetEnv("CHAIN_WS_ENDPOINT", "wss://crayfish.darwinia.network/")
	NetworkNode      = GetEnv("NETWORK_NODE", DarwiniaNetwork)
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}

	return value
}
