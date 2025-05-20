package erc20

import (
	"context"
	"github.com/itering/subscan/pkg/go-web3"
	"github.com/itering/subscan/pkg/go-web3/dto"
	"github.com/itering/subscan/plugins/evm/abi"
	"github.com/itering/subscan/plugins/evm/contract"
	"github.com/itering/subscan/util"

	"github.com/shopspring/decimal"
)

// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-20.md

// view
// function name() public view returns (string) OPTIONAL
// function symbol() public view returns (string) OPTIONAL
// function decimals() public view returns (uint8) OPTIONAL
// function totalSupply() public view returns (uint256)
// function balanceOf(address _owner) public view returns (uint256 balance)

// method
// function transfer(address _to, uint256 _value) public returns (bool success)
// function transferFrom(address _from, address _to, uint256 _value) public returns (bool success)
// function approve(address _spender, uint256 _value) public returns (bool success)
// function allowance(address _owner, address _spender) public view returns (uint256 remaining)

// events
// event Transfer(address indexed _from, address indexed _to, uint256 _value)
// event Approval(address indexed _owner, address indexed _spender, uint256 _value)

var (
	EventTransfer = abi.EncodingMethod("Transfer(address,address,uint256)")
	EventDeposit  = abi.EncodingMethod("Deposit(address,uint256)")
	EventWithdraw = abi.EncodingMethod("Withdrawal(address,uint256)")
	// EventApproval = abi.EncodingMethod("Approval(address,address,uint256)")
)

type Token struct {
	contract.Contract
}

func Init(w3 *web3.Web3, contract string) *Token {
	c, err := w3.Eth.NewContract(abi.Erc20)
	if err != nil {
		panic(err)
	}
	t := Token{}
	t.Contract.EthContract = c
	t.Contract.TransParam = dto.TransactionParameters{To: contract, Data: ""}
	return &t
}

func (c Token) Name(ctx context.Context) (string, error) {
	if value, err := c.GetStorage(ctx, "name"); err == nil {
		if value.Error != nil {
			return "", nil
		}
		return util.AbiStringDecoder(value.Result.(string)), nil
	} else {
		return "", err
	}
}

func (c Token) Symbol(ctx context.Context) (string, error) {
	if value, err := c.GetStorage(ctx, "symbol"); err == nil {
		if value.Error != nil {
			return "", nil
		}
		return util.AbiStringDecoder(value.Result.(string)), nil
	} else {
		return "", err
	}
}

func (c Token) Decimals(ctx context.Context) (uint, error) {
	if value, err := c.GetStorage(ctx, "decimals"); err == nil {
		if value.Error != nil {
			return 0, nil
		}
		return uint(util.U256(value.Result.(string)).Uint64()), nil
	} else {
		return 0, err
	}
}

func (c Token) TotalSupply(ctx context.Context) (decimal.Decimal, error) {
	if value, err := c.GetStorage(ctx, "totalSupply"); err == nil {
		if value.Error != nil {
			return decimal.Zero, nil
		}
		return decimal.NewFromBigInt(util.U256(value.Result.(string)), 0), nil
	} else {
		return decimal.Zero, err
	}
}

func (c Token) BalanceOf(ctx context.Context, address string) (decimal.Decimal, error) {
	if value, err := c.GetStorage(ctx, "balanceOf", address); err == nil {
		if value.Error != nil {
			return decimal.Zero, nil
		}
		return decimal.NewFromBigInt(util.U256(value.Result.(string)), 0), nil
	} else {
		return decimal.Zero, err
	}
}
