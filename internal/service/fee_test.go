package service

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

// GetExtrinsicFee
func Test_GetExtrinsicFee(t *testing.T) {
	ctx := context.TODO()
	raw := "0x510284002534454d30f8a028e42654d6b535e0651d1d026ddf115cef59ae1dd71bae074e003c696816e433613538cbea7ee411f2812c672d4254bf341121456b9bb4f9c13594b96c201748167ed8f810f90cc07158e7accd3b977cff98337bf9f40cec030c1501d20226000000050000d4c0c691f39a442a4022ac305d1bc270b8740602df4ef243cbb066c7e53b4523070008f6a4e8"
	fee, usedFee, err := GetExtrinsicFee(ctx, raw, "0xbbec5f0f1129efb846c46fb0d5670ff0f7cccc60760726d39306061352119587",
		-1, decimal.NewFromInt(290565000),
		decimal.NewFromInt(161305248), true)
	assert.NoError(t, err)
	assert.Equal(t, fee.IntPart(), int64(161305248))
	assert.Equal(t, usedFee.IntPart(), int64(161305248))
}
