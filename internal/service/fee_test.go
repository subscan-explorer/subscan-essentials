package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// GetExtrinsicFee
func Test_GetExtrinsicFee(t *testing.T) {
	tc := TestConn{}
	fee, err := GetExtrinsicFee(&tc, "0x710284b62d88e3f439fe9b5ea799b27bf7c6db5e795de1784f27b1bc051553499e420f01a494b989939016bce14441b2c757cbb1c959e700220757d10811d14192b320038365b1dcb1e7bc3a5221267943bd2e54fac6fbc1805229d1f04e5611d8b2c58e750000001500b62d88e3f439fe9b5ea799b27bf7c6db5e795de1784f27b1bc051553499e420f0000a0724e180900000000000000000000000000")
	assert.NoError(t, err)
	assert.Greater(t, fee.IntPart(), int64(0))
}
