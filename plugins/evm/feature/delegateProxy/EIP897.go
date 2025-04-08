package delegateProxy

import (
	"context"
	"errors"
	"github.com/itering/subscan/pkg/go-web3"
	"github.com/itering/subscan/pkg/go-web3/dto"
	"github.com/itering/subscan/plugins/evm/abi"
	"github.com/itering/subscan/plugins/evm/contract"
	"github.com/itering/subscan/util"
)

// https://eips.ethereum.org/EIPS/eip-897

// interface ERCProxy {
//  function proxyType() public pure returns (uint256 proxyTypeId);
//  function implementation() public view returns (address codeAddr);
// }

const Eip897Standard = "EIP897"

type EIP897 struct {
	contract.Contract
}

func Init897(w3 *web3.Web3, contract string) *EIP897 {
	c, err := w3.Eth.NewContract(abi.EIP897)
	if err != nil {
		panic(err)
	}
	t := EIP897{}
	t.Eth = w3.Eth
	t.Contract.EthContract = c
	t.Contract.TransParam = dto.TransactionParameters{To: contract, Data: ""}
	return &t
}

func (c *EIP897) Implementation(ctx context.Context) (string, error) {
	if value, err := c.GetStorage(ctx, "implementation"); err == nil {
		if value.Error != nil {
			return "", errors.New(value.Error.Message)
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

func (c *EIP897) Standard() string {
	return Eip897Standard
}
