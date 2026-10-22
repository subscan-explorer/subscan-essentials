package delegateProxy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/itering/subscan/pkg/go-web3"
	"github.com/itering/subscan/pkg/go-web3/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// proxyRPC exercises the HTTP provider against stable JSON-RPC replies.
func proxyRPC(t *testing.T, result func(string, []json.RawMessage) string) *web3.Web3 {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			ID     json.RawMessage   `json:"id"`
			Method string            `json:"method"`
			Params []json.RawMessage `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode JSON-RPC request: %v", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		response := map[string]any{"jsonrpc": "2.0", "id": request.ID, "result": result(request.Method, request.Params)}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("encode JSON-RPC response: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	return web3.NewWeb3(providers.NewHTTPProvider(server.URL, 5, false))
}

func TestDelegateProxy(t *testing.T) {
	const implementation = "0x85520f613021e5db2afb40cb91ef066ccf212111"
	rpc := proxyRPC(t, func(method string, params []json.RawMessage) string {
		assert.Equal(t, "eth_call", method)
		assert.Len(t, params, 2)
		var call map[string]string
		if assert.NoError(t, json.Unmarshal(params[0], &call)) {
			assert.Equal(t, "0x14a4123da9ad21b2215dc0ab6984ec1e89842c6d", call["to"])
			assert.Equal(t, "0x5c60da1b", call["data"])
		}
		return "0x" + strings.Repeat("0", 24) + strings.TrimPrefix(implementation, "0x")
	})
	proxy := Init897(rpc, "0x14a4123da9ad21b2215dc0ab6984ec1e89842c6d")
	actual, err := proxy.Implementation(context.Background())
	require.NoError(t, err)
	assert.Equal(t, implementation, actual)
}

func TestEIP1967_Implementation(t *testing.T) {
	const implementation = "0x40f1eca9c82200428704aa555da1009ad4beb2e2"
	rpc := proxyRPC(t, func(method string, params []json.RawMessage) string {
		assert.Equal(t, "eth_getStorageAt", method)
		assert.Len(t, params, 3)
		var address, slot string
		if !assert.NoError(t, json.Unmarshal(params[0], &address)) || !assert.NoError(t, json.Unmarshal(params[1], &slot)) {
			return "0x"
		}
		assert.Equal(t, "0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc", slot)
		if address == "0x25442adf37379be90ed1f7fccd9c9417b10aa4dc" {
			return "0x" + strings.Repeat("0", 24) + strings.TrimPrefix(implementation, "0x")
		}
		assert.Equal(t, "0x22b1a40e3178fe7c7109efcc247c5bb2b34abe32", address)
		return "0x" + strings.Repeat("0", 64)
	})
	proxy := Init1967(rpc, "0x25442adf37379be90ed1f7fccd9c9417b10aa4dc")
	actual, err := proxy.Implementation(context.Background())
	require.NoError(t, err)
	assert.Equal(t, implementation, actual)

	proxy.Contract.TransParam.To = "0x22b1a40e3178fe7c7109efcc247c5bb2b34abe32"
	actual, err = proxy.Implementation(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "0x0000000000000000000000000000000000000000", actual)
}
