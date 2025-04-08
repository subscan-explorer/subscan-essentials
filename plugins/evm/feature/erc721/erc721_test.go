package erc721

import (
	"context"
	"github.com/itering/subscan/pkg/go-web3"
	"github.com/itering/subscan/pkg/go-web3/providers"
	"github.com/itering/subscan/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BaseVerify(t *testing.T) {
	ctx := context.Background()
	rpc := web3.NewWeb3(providers.NewHTTPProvider("https://moonbeam-rpc.n.dwellir.com", 60, false))
	assert.True(t, BaseVerify(ctx, rpc, "0x51737fA634e26f5687E45c6ca07604E064076350", "0xF977814e90dA44bFA03b6295A0616a897441aceC", "902d"))
	assert.False(t, BaseVerify(ctx, rpc, "0xFA36Fe1dA08C89eC72Ea1F0143a35bFd5DAea108", "0xF977814e90dA44bFA03b6295A0616a897441aceC", "01"))         // ERC20
	assert.False(t, BaseVerify(ctx, rpc, "0x51737fA634e26f5687E45c6ca07604E064076350", "0xF977814e90dA44bFA03b6295A0616a897441aceC", "1234567890")) // wrong token_id
}

func Test_MethodHash(t *testing.T) {
	assert.Equal(t, util.AddHex(SetURIMethodId), "0x02fe5305")
	assert.Equal(t, util.AddHex(SetBaseTokenURIMethodId), "0x30176e13")
}

func Test_721Token(t *testing.T) {
	ctx := context.Background()
	rpc := web3.NewWeb3(providers.NewHTTPProvider("https://moonbeam-rpc.n.dwellir.com", 60, false))
	token := Init(rpc, "0xcb13945ca8104f813992e4315f8ffefe64ac49ca")

	name, err := token.Name(ctx)
	assert.NoError(t, err)
	assert.Equal(t, name, "GLMR JUNGLE")

	symbol, err := token.Symbol(ctx)
	assert.NoError(t, err)
	assert.Equal(t, symbol, "GLMJ")

	owner, err := token.OwnerOf(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, owner, "0xb7fdd27a8df011816205a6e3caa097dc4d8c2c5d")

	balanceOf, err := token.BalanceOf(ctx, "0xb7fdd27a8df011816205a6e3caa097dc4d8c2c5d")
	assert.NoError(t, err)
	assert.Equal(t, balanceOf.String(), "296")

	uri, err := token.TokenURI(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, uri, "https://www.glmrjungle.com/nfts/1.json")
}
