package dao

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_account(t *testing.T) {

	t.Run("evm_accounts table should be created", func(t *testing.T) {
		assert.True(t, sg.db.Migrator().HasTable(&Account{}))
		tableType, _ := sg.db.Migrator().TableType(&Account{})
		assert.Equal(t, "evm_accounts", tableType.Name())
	})

	t.Run("touch account should create a new account", func(t *testing.T) {
		ctx := context.Background()

		// Valid Ethereum address
		validH160 := "0x1234567890abcdef1234567890abcdef12345678"
		err := TouchAccount(ctx, validH160)
		assert.NoError(t, err)

		var account Account
		result := sg.db.WithContext(ctx).First(&account, "evm_account = ?", validH160)
		assert.NoError(t, result.Error)
		assert.Equal(t, validH160, account.EvmAccount)
		assert.NotEmpty(t, account.Address)

		// Invalid Ethereum address
		invalidH160 := "0xINVALIDADDRESS"
		err = TouchAccount(ctx, invalidH160)
		assert.NoError(t, err)
		var accountInvalid Account
		result = sg.db.WithContext(ctx).First(&accountInvalid, "evm_account = ?", invalidH160)
		assert.Error(t, result.Error)
		assert.Empty(t, accountInvalid.EvmAccount)
	})

}
