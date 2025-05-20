package dao

import (
	"context"
	"github.com/itering/subscan/util/address"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Contract(t *testing.T) {
	var contractAddress = address.Format("0xb2bf0bf26A4e98a6AEe1484b3bdaf50E3fb4a346")
	ctx := context.TODO()
	// contract := GetContract(ctx, contractAddress)

	t.Run("contract create will add a account", func(t *testing.T) {
		var account Account
		err := sg.db.Where("evm_account = ?", contractAddress).First(&account).Error
		assert.NoError(t, err)
		assert.Equal(t, contractAddress, account.EvmAccount)
	})

	t.Run("setContractProxyImplementation will set contract proxy implementation", func(t *testing.T) {
		setContractProxyImplementation(ctx, contractAddress, "0x1234567890abcdef1234567890abcdef12345678")
		afterSetContractProxyImplementation := GetContract(ctx, contractAddress)
		assert.Equal(t, "0x1234567890abcdef1234567890abcdef12345678", afterSetContractProxyImplementation.ProxyImplementation)
	})

}
