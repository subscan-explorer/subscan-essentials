package erc1155

import (
	"context"
	"math/big"

	"github.com/itering/subscan/pkg/go-web3"
	"github.com/itering/subscan/pkg/go-web3/providers"
	"github.com/itering/subscan/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BaseValidate(t *testing.T) {
	ctx := context.Background()
	rpc := web3.NewWeb3(providers.NewHTTPProvider(util.GetEnv("TEST_EVM_ENDPOINT", "https://wss.api.moonbeam.network"), 60, false))
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
	assert.Equal(t, 7, len(results))
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
