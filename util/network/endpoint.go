package network

var endpoints = map[Network]string{
	Polkadot:          "wss://rpc.polkadot.io",
	KusamaNetwork:     "wss://kusama-rpc.dwellir.com",
	Westend:           "wss://westend.api.onfinality.io/public-ws",
	AssethubWestend:   "wss://westend-asset-hub-rpc.polkadot.io",
	BridgehubWestend:  "wss://westend-bridge-hub-rpc.polkadot.io",
	BridgehubPolkadot: "wss://polkadot-bridge-hub-rpc.polkadot.io",
	Moonbeam:          "wss://moonbeam-rpc.dwellir.com",
	MoonRiver:         "wss://moonriver.api.onfinality.io/public-ws",
}

func GetDefaultEndpoint(network Network) string {
	return endpoints[network]
}
