package delegateProxy

import (
	"context"
	"errors"
	"github.com/itering/subscan/pkg/go-web3"
	"github.com/itering/subscan/pkg/go-web3/dto"
	"github.com/itering/subscan/plugins/evm/contract"
	"github.com/itering/subscan/util"
)

// https://eips.ethereum.org/EIPS/eip-1967
// slot index  0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc = bytes32(uint256(keccak256('eip1967.proxy.implementation')) - 1)

const EIP1967Standard = "EIP1967"

type EIP1967 struct {
	contract.Contract
}

func Init1967(w3 *web3.Web3, contract string) *EIP1967 {
	t := EIP1967{}
	t.Eth = w3.Eth
	t.Contract.TransParam = dto.TransactionParameters{To: contract, Data: ""}
	return &t
}

func (c *EIP1967) Implementation(ctx context.Context) (string, error) {
	if addr, err := c.GetStorageByKey(ctx, c.Contract.TransParam.To, "0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"); err == nil {
		if len(addr) < 66 {
			return "", errors.New("not address")
		}
		return util.AddHex(util.TrimHex(addr)[24:64]), nil
	} else {
		return "", err
	}
}

func (c *EIP1967) Standard() string {
	return EIP1967Standard
}
