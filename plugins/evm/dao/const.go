package dao

import (
	"github.com/itering/subscan/util"
)

var (
	ContractAddrKey             = util.NetworkNode + ":" + "EvmContractAddr"
	ContractCreationBytecodeKey = util.NetworkNode + ":" + "EvmContractCreationBytecode"
	TransactionCount            = util.NetworkNode + ":" + "EvmTransactionCount"
	TokenAddrKey                = util.NetworkNode + ":" + "EvmTokenAddr"
	Eip20Token                  = "erc20"
	Eip721Token                 = "erc721"
	Eip1155Token                = "erc1155"

	TransactionQueue = "Transaction"
	ConvictionClass  = "ConvictionVoting"
	TraceQueue       = "Trace"

	NullAddress = "0x0000000000000000000000000000000000000000"

	Create = "CREATE"
)

const (
	PreCompareExternal = 1
)

func HashDeployedCode(code string) string {
	return util.AddHex(DoBlake2_256(code))
}
