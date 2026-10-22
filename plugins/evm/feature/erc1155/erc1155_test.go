package erc1155

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"

	eAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/itering/subscan/pkg/go-web3"
	"github.com/itering/subscan/pkg/go-web3/providers"
	"github.com/itering/subscan/plugins/evm/abi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_BaseValidate(t *testing.T) {
	ctx := context.Background()
	rpc := erc1155RPC(t)
	token := Init(rpc, "0x593eb66fffc499b3b871e9a4465658bd9b07d174")
	result, err := token.SupportsInterface(ctx)
	assert.NoError(t, err)
	assert.True(t, result)
	ifpsUri, err := token.Uri(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, "ipfs://QmWvKWanZcSrc8TAmWbtgvQ8E63a31JHydBye4vjHQPKSt/{id}.json", ifpsUri)
	results, err := token.BalanceOfBatch(ctx, []string{"0xfebd7eed30360729734e498a9c5fef0065b3193e",
		"0xfebd7eed30360729734e498a9c5fef0065b3193e",
		"0xfebd7eed30360729734e498a9c5fef0065b3193e",
		"0xfebd7eed30360729734e498a9c5fef0065b3193e",
		"0xfebd7eed30360729734e498a9c5fef0065b3193e",
		"0xfebd7eed30360729734e498a9c5fef0065b3193e",
		"0xfebd7eed30360729734e498a9c5fef0065b3193e",
	},
		[]*big.Int{
			big.NewInt(6),
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(3),
			big.NewInt(4),
			big.NewInt(5),
			big.NewInt(6),
		})
	assert.NoError(t, err)
	require.Len(t, results, 7)
	assert.Equal(t, "1", results[0].String())
	assert.Equal(t, "0", results[1].String())
	assert.Equal(t, "1", results[6].String())
	balance, err := token.BalanceOfWithTokenId(ctx, "0xfebd7eed30360729734e498a9c5fef0065b3193e", "6")
	assert.NoError(t, err)
	assert.Equal(t, "1", balance.String())
	// https://moonbeam.subscan.io/tx/0x3e74de05c17a391401e637b645a462cdb30bdb24ae4f5b38c93930d0ea439fb1
	assert.Equal(t, "c3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62", EventTransferSingle)
	// https://moonbeam.subscan.io/tx/0x5d2882dffcc847cd1fdbabe61e7fe38a69a3e36fc015a30e383a70a547b05cf3
	assert.Equal(t, "4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb", EventTransferBatch)
	// https://moonbeam.subscan.io/tx/0x4a64cf70010667b0b5737effbdf0b9d46f323baad6c0ef3bef67df269a32f9ff
	assert.Equal(t, "6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b", EventURI)
}

// erc1155RPC retains ABI encoding/decoding coverage while isolating mutable chain state.
func erc1155RPC(t *testing.T) *web3.Web3 {
	t.Helper()
	contractABI, err := eAbi.JSON(strings.NewReader(abi.Erc1155))
	require.NoError(t, err)
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
		if !assert.Equal(t, "eth_call", request.Method) || !assert.Len(t, request.Params, 2) {
			http.Error(w, "unexpected method or parameters", http.StatusBadRequest)
			return
		}
		var call map[string]string
		if err := json.Unmarshal(request.Params[0], &call); err != nil {
			t.Errorf("decode call parameters: %v", err)
			http.Error(w, "invalid call", http.StatusBadRequest)
			return
		}
		assert.Equal(t, "0x593eb66fffc499b3b871e9a4465658bd9b07d174", call["to"])
		data, err := hex.DecodeString(strings.TrimPrefix(call["data"], "0x"))
		if err != nil || len(data) < 4 {
			t.Errorf("invalid call data: %v", err)
			http.Error(w, "invalid call data", http.StatusBadRequest)
			return
		}
		method, err := contractABI.MethodById(data[:4])
		if err != nil {
			t.Errorf("unknown contract method: %v", err)
			http.Error(w, "unknown method", http.StatusBadRequest)
			return
		}
		args, err := method.Inputs.Unpack(data[4:])
		if err != nil {
			t.Errorf("unpack contract arguments: %v", err)
			http.Error(w, "invalid contract arguments", http.StatusBadRequest)
			return
		}
		var value any
		switch method.Name {
		case "supportsInterface":
			assert.Equal(t, [4]byte{0xd9, 0xb6, 0x7a, 0x26}, args[0])
			value = true
		case "uri":
			assert.Equal(t, big.NewInt(1), args[0])
			value = "ipfs://QmWvKWanZcSrc8TAmWbtgvQ8E63a31JHydBye4vjHQPKSt/{id}.json"
		case "balanceOfBatch":
			owners := args[0].([]common.Address)
			ids := args[1].([]*big.Int)
			assert.Len(t, owners, 7)
			assert.Equal(t, []*big.Int{big.NewInt(6), big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5), big.NewInt(6)}, ids)
			for _, owner := range owners {
				assert.Equal(t, common.HexToAddress("0xfebd7eed30360729734e498a9c5fef0065b3193e"), owner)
			}
			value = []*big.Int{big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1)}
		case "balanceOf":
			assert.Equal(t, common.HexToAddress("0xfebd7eed30360729734e498a9c5fef0065b3193e"), args[0])
			assert.Equal(t, big.NewInt(6), args[1])
			value = big.NewInt(1)
		default:
			t.Errorf("unexpected contract method: %s", method.Name)
			http.Error(w, "unexpected contract method", http.StatusBadRequest)
			return
		}
		encoded, err := method.Outputs.Pack(value)
		if err != nil {
			t.Errorf("encode contract result: %v", err)
			http.Error(w, "invalid contract result", http.StatusInternalServerError)
			return
		}
		response := map[string]any{"jsonrpc": "2.0", "id": request.ID, "result": "0x" + hex.EncodeToString(encoded)}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("encode JSON-RPC response: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	return web3.NewWeb3(providers.NewHTTPProvider(server.URL, 5, false))
}
