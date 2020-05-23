package substrate

import (
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/ss58"
)

const (
	ChainNewHead       = "chain_newHead"
	ChainFinalizedHead = "chain_finalizedHead"
	StateStorage       = "state_storage"
	BlockTime          = 6
)

var (
	CurrentRuntimeSpecVersion int
	EventStorageKey           = util.GetEnv("SUBSTRATE_EVENT_KEY", "0x26aa394eea5630e07c48ae0c9558cef780d41e5e16056765bc8461851072c9d7")
	AddressType               = util.StringToInt(util.GetEnv("SUBSTRATE_ADDRESS_TYPE", "2"))
	BalanceAccuracy           = util.StringToInt(util.GetEnv("SUBSTRATE_ACCURACY", "9"))
	CommissionAccuracy        = util.GetEnv("COMMISSION_ACCURACY", "9")

	SupportToken = map[string][]string{
		util.CrabNetwork:   {"RING", "KTON", "POWER"},
		util.KusamaNetwork: {"KSM"},
		util.Plasm:         {"PLM"},
		util.Acala:         {"ACA"},
		util.Edgeware:      {"EDG"},
	}
	TokenSymbol = SupportToken[util.NetworkNode]
)

func SS58Address(address string) string {
	return ss58.Encode(address, AddressType)
}
