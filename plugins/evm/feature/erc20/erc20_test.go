package erc20

import (
	"context"
	"github.com/itering/subscan/pkg/go-web3"
	"github.com/itering/subscan/pkg/go-web3/providers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken_Decimals(t *testing.T) {
	ctx := context.Background()
	token := Init(web3.NewWeb3(providers.NewHTTPProvider("https://wss.api.moonbeam.network", 60, false)), "0xAcc15dC74880C9944775448304B263D191c6077F")

	decimals, err := token.Decimals(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int(decimals), 18)

	supply, err := token.TotalSupply(ctx)
	assert.NoError(t, err)
	assert.True(t, supply.IsPositive())

	name, err := token.Name(ctx)
	assert.NoError(t, err)
	assert.Equal(t, name, "Wrapped GLMR")

	symbol, err := token.Symbol(ctx)
	assert.NoError(t, err)
	assert.Equal(t, symbol, "WGLMR")

	balance, err := token.BalanceOf(ctx, "0xc295aa4287127C5776Ad7031648692659eF2ceBB")
	assert.NoError(t, err)
	assert.True(t, balance.IsPositive())
}
