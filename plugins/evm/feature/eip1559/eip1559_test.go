package eip1559

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// func Test_TreasuryFees(t *testing.T) {
// 	assert.Equal(t, TreasuryFees(decimal.New(69618750000000000, 0)).String(), decimal.New(13923750000000000, 0).String())
// }

func TestBurntFee(t *testing.T) {
	assert.Equal(t, BurntFee(decimal.New(1, 0), decimal.New(69618750000000000, 0)).String(), decimal.New(55695000000000000, 0).String())
}

func TestSavingsFee(t *testing.T) {
	assert.Equal(t, SavingsFee(decimal.New(320723064240, 0),
		decimal.New(125000000000, 0),
		decimal.NewFromInt32(0), decimal.New(556950, 0)).String(),
		decimal.New(109007960628468000, 0).String())
}
