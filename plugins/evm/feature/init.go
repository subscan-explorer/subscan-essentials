package feature

import (
	"context"
	"fmt"
	"github.com/itering/subscan/pkg/go-web3"
	"github.com/itering/subscan/plugins/evm/feature/erc1155"
	"github.com/itering/subscan/plugins/evm/feature/erc20"
	"github.com/itering/subscan/plugins/evm/feature/erc721"

	"github.com/shopspring/decimal"
)

type IToken interface {
	Name(context.Context) (string, error)
	Symbol(context.Context) (string, error)
	Decimals(context.Context) (uint, error)
	TotalSupply(context.Context) (decimal.Decimal, error)
	BalanceOf(ctx context.Context, address string) (decimal.Decimal, error)
}

func InitToken(w3 *web3.Web3, category, contract string) IToken {
	if category == "erc20" {
		return erc20.Init(w3, contract)
	}
	if category == "erc721" {
		return erc721.Init(w3, contract)
	}
	if category == "erc1155" {
		return erc1155.Init(w3, contract)
	}
	panic(fmt.Errorf("not support token %s", category))

}
