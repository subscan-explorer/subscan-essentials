package delegateProxy

import (
	"context"
	"github.com/itering/subscan/plugins/evm/abi"
)

// Events
// event Upgraded(address indexed implementation);

var (
	EventUpgraded = abi.EncodingMethod("Upgraded(address)")
)

type IDelegateProxy interface {
	Implementation(context.Context) (string, error)
	Standard() string
}
