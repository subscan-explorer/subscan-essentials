package staking

import (
	"os"
	"testing"

	scanModel "github.com/itering/subscan/model"
	"github.com/itering/subscan/util/address"
	"github.com/lmittmann/tint"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func Test_castArg(t *testing.T) {
	slog.New(tint.NewHandler(os.Stderr, nil))
	arg := scanModel.CallArg{
		Name:  "validator_stash",
		Value: "0xbe5ddb1579b72e84524fc29e78609e3caf42e85aa118ebfe0b0ad404b5bdd25f",
	}
	acct, err := CastArg[address.SS58Address](arg, "validator_stash")
	assert.NoError(t, err)
	assert.Equal(t, address.SS58Address("5GNJqTPyNqANBkUVMN1LPPrxXnFouWXoe2wNSmmEoLctxiZY"), acct)

	acctStr, err := CastArg[string](arg, "validator_stash")
	assert.NoError(t, err)
	assert.Equal(t, "0xbe5ddb1579b72e84524fc29e78609e3caf42e85aa118ebfe0b0ad404b5bdd25f", acctStr)

	arg = scanModel.CallArg{
		Name:  "era",
		Value: float64(1),
	}
	era, err := CastArg[uint32](arg, "era")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), era)
}
