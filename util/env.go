package util

import (
	"math/rand"
	"os"
	"strings"
)

const (
	CrabNetwork   = "darwinia-canary"
	KusamaNetwork = "kusama"
	Edgeware      = "edgeware"
	Acala         = "acala-test"
	Plasm         = "plasm"
)

var (
	HostName     string
	Environment  string
	WSEndPoint   string
	chainRpcUrls []string
	NetworkNode  = GetEnv("NETWORK_NODE", KusamaNetwork)
	IsDarwinia   = StringInSlice(NetworkNode, []string{CrabNetwork})
)

// TODO
//
// Config file
func init() {
	urls := GetEnv("CHAIN_WS_ENDPOINT", "")
	chainRpcUrls = strings.Split(urls, ",")
	WSEndPoint = chainRpcUrls[rand.Intn(len(chainRpcUrls))]
	HostName, _ = os.Hostname()
}

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

func IsK8s() bool {
	return os.Getenv("K8S_ENV") == "true"
}
