package erc721

import (
	"context"
	"errors"
	"github.com/itering/subscan/pkg/go-web3"
	"github.com/itering/subscan/pkg/go-web3/dto"
	"github.com/itering/subscan/plugins/evm/abi"
	"github.com/itering/subscan/plugins/evm/contract"
	"github.com/itering/subscan/util"

	"github.com/shopspring/decimal"
)

// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-721.md

// view
// function balanceOf(address _owner) external view returns (uint256); MUST
// function ownerOf(uint256 _tokenId) external view returns (address); MUST
// function name() public view returns (string) OPTIONAL
// function symbol() public view returns (string) OPTIONAL
// function tokenURI(uint256 _tokenId) external view returns (string); OPTIONAL
// function totalSupply() external view returns (uint256); OPTIONAL

// methods
//  function safeTransferFrom(address _from, address _to, uint256 _tokenId, bytes data) external payable; MUST
//  function safeTransferFrom(address _from, address _to, uint256 _tokenId) external payable; MUST
//  function transferFrom(address _from, address _to, uint256 _tokenId) external payable; MUST
//  function approve(address _approved, uint256 _tokenId) external payable; MUST
//  function setApprovalForAll(address _operator, bool _approved) external; MUST
//  function getApproved(uint256 _tokenId) external view returns (address); MUST
// function isApprovedForAll(address _owner, address _operator) external view returns (bool); MUST

// events
// event Transfer(address indexed _from, address indexed _to, uint256 indexed _tokenId); MUST
// event Approval(address indexed _owner, address indexed _approved, uint256 indexed _tokenId); MUST
// event ApprovalForAll(address indexed _owner, address indexed _operator, bool _approved); MUST
//

var (
	EventTransfer = abi.EncodingMethod("Transfer(address,address,uint256)")

	SetBaseTokenURIMethodId = abi.EncodingMethod("setBaseTokenURI(string)")[0:8]
	SetURIMethodId          = abi.EncodingMethod("setURI(string)")[0:8]
	SetTokenURIMethodId     = abi.EncodingMethod("setTokenURI(uint256,string)")[0:8]
	// EventApproval = abi.EncodingMethod("Approval(address,address,uint256)")
)

type Token struct {
	contract.Contract
}

func Init(w3 *web3.Web3, contract string) *Token {
	c, err := w3.Eth.NewContract(abi.Erc721)
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

func (c Token) Decimals(_ context.Context) (uint, error) { return 0, nil }

func (c Token) TokenURI(ctx context.Context, tokenId string) (string, error) {
	if value, err := c.GetStorage(ctx, "tokenURI", decimal.RequireFromString(tokenId).Coefficient()); err == nil {
		if value.Error != nil {
			return "", nil
		}
		return util.AbiStringDecoder(value.Result.(string)), nil
	} else {
		return "", err
	}
}

func (c Token) BaseURI(ctx context.Context) (string, error) {
	if value, err := c.GetStorage(ctx, "baseURI"); err == nil {
		if value.Error != nil {
			return c.BaseTokenURI(ctx)
		}
		return util.AbiStringDecoder(value.Result.(string)), nil
	} else {
		return "", err
	}
}

func (c Token) BaseTokenURI(ctx context.Context) (string, error) {
	if value, err := c.GetStorage(ctx, "baseTokenURI"); err == nil {
		if value.Error != nil {
			return "", nil
		}
		return util.AbiStringDecoder(value.Result.(string)), nil
	} else {
		return "", err
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

func (c Token) OwnerOf(ctx context.Context, tokenId string) (string, error) {
	if value, err := c.GetStorage(ctx, "ownerOf", decimal.RequireFromString(tokenId).Coefficient()); err == nil {
		if value.Error != nil {
			return "", nil
		}
		addr := value.Result.(string)
		if len(addr) < 66 {
			return "", errors.New("not address")
		}
		return util.AddHex(util.TrimHex(addr)[24:64]), nil
	} else {
		return "", err
	}
}
