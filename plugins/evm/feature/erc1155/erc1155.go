package erc1155

import (
	"context"
	eAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/itering/subscan/pkg/go-web3"
	"github.com/itering/subscan/pkg/go-web3/dto"
	"github.com/itering/subscan/plugins/evm/abi"
	"github.com/itering/subscan/plugins/evm/contract"
	"github.com/itering/subscan/util"
	"math/big"
	"strings"

	"github.com/shopspring/decimal"
)

// https://github.com/ethereum/ercs/blob/master/ERCS/erc-1155.md

// view
// function balanceOf(address _owner, uint256 _id) external view returns (uint256);
// function balanceOfBatch(address[] calldata _owners, uint256[] calldata _ids) external view returns (uint256[] memory);
// ERC165 interface
// function supportsInterface(bytes4 interfaceID) external view returns (bool);
// function uri(uint256 _id)(string)
//

// events
// event TransferSingle(address indexed _operator, address indexed _from, address indexed _to, uint256 _id, uint256 _value);
// event TransferBatch(address indexed _operator, address indexed _from, address indexed _to, uint256[] _ids, uint256[] _values);
// event URI(string _value, uint256 indexed _id);
//

var (
	EventTransferSingle = abi.EncodingMethod("TransferSingle(address,address,address,uint256,uint256)")
	EventTransferBatch  = abi.EncodingMethod("TransferBatch(address,address,address,uint256[],uint256[])")
	EventURI            = abi.EncodingMethod("URI(string,uint256)")
)

type Token struct {
	contract.Contract
}

func Init(w3 *web3.Web3, contract string) *Token {
	c, err := w3.Eth.NewContract(abi.Erc1155)
	if err != nil {
		panic(err)
	}
	t := Token{}
	t.Contract.EthContract = c
	t.Contract.TransParam = dto.TransactionParameters{To: contract, Data: ""}
	return &t
}

// Name optional
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

// Symbol optional
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

// Decimals not implemented
func (c Token) Decimals(_ context.Context) (uint, error) { return 0, nil }

// TotalSupply not implemented
func (c Token) TotalSupply(ctx context.Context) (decimal.Decimal, error) {
	return decimal.Zero, nil
}

func (c Token) BalanceOf(ctx context.Context, accountId string) (decimal.Decimal, error) {
	return decimal.Zero, nil
}

// SupportsInterface 0xd9b67a26
func (c Token) SupportsInterface(ctx context.Context) (bool, error) {
	trueValue := "0x0000000000000000000000000000000000000000000000000000000000000001"
	if value, err := c.GetStorage(ctx, "supportsInterface", "0xd9b67a26"); err == nil {
		if value.Error != nil {
			return false, nil
		}
		return value.Result.(string) == trueValue, nil
	} else {
		return false, err
	}
}

func (c Token) Uri(ctx context.Context, tokenId string) (string, error) {
	if value, err := c.GetStorage(ctx, "uri", decimal.RequireFromString(tokenId).Coefficient()); err == nil {
		if value.Error != nil {
			return "", nil
		}
		return util.AbiStringDecoder(value.Result.(string)), nil
	} else {
		return "", err
	}
}

func (c Token) BalanceOfWithTokenId(ctx context.Context, accountId string, tokenId string) (decimal.Decimal, error) {
	if value, err := c.GetStorage(ctx, "balanceOf", accountId, decimal.RequireFromString(tokenId).Coefficient()); err == nil {
		if value.Error != nil {
			return decimal.Zero, nil
		}
		return decimal.NewFromBigInt(util.U256(value.Result.(string)), 0), nil
	} else {
		return decimal.Zero, err
	}
}

func (c Token) BalanceOfBatch(ctx context.Context, accountId []string, tokenId []*big.Int) ([]decimal.Decimal, error) {

	eABI, err := eAbi.JSON(strings.NewReader(abi.Erc1155))
	if err != nil {
		panic(err)
	}

	var accounts []common.Address
	for _, v := range accountId {
		accounts = append(accounts, common.HexToAddress(v))
	}
	method := "balanceOfBatch"
	data, err := eABI.Pack(method, accounts, tokenId)
	if err != nil {
		return nil, err
	}
	if len(util.BytesToHex(data)) < 8 {
		return nil, nil
	}
	if value, err := c.GetStorage(ctx, method, util.BytesToHex(data)[8:], ""); err == nil {
		if value.Error != nil {
			return nil, nil
		}
		result, err := eABI.Unpack(method, util.HexToBytes(value.Result.(string)))
		if err != nil {
			return nil, err
		}
		if len(result) == 0 {
			return nil, nil
		}
		if _, ok := result[0].([]*big.Int); !ok {
			return nil, nil
		}
		var results []decimal.Decimal
		for _, v := range result[0].([]*big.Int) {
			results = append(results, decimal.NewFromBigInt(v, 0))
		}
		return results, nil
	}
	return nil, nil
}
