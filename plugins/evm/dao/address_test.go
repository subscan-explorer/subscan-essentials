package dao

import (
	"context"
	"github.com/itering/subscan/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveAddressPadded(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0x0000000000000000000000001234567890abcdef1234567890abcdef12345678", "0x1234567890abcdef1234567890abcdef12345678"},
		{"0x123", NullAddress},
	}

	for _, test := range tests {
		result := RemoveAddressPadded(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestH160ToAccountIdByNetwork(t *testing.T) {
	tests := []struct {
		addr       string
		expected   string
		isEvmChain bool
	}{
		{"0x1234567890abcdef1234567890abcdef12345678", "1234567890abcdef1234567890abcdef12345678eeeeeeeeeeeeeeeeeeeeeeee", false},
		{"0x1234567890abcdef1234567890abcdef12345678", "0x1234567890abcdef1234567890abcdef12345678", true},
		{"invalid_address", "", false},
	}

	for _, test := range tests {
		util.IsEvmChain = test.isEvmChain
		result := h160ToAccountIdByNetwork(context.Background(), test.addr, "")
		assert.Equal(t, test.expected, result)
	}
	assert.Equal(t, "750ea21c1e98cced0d4557196b6f4a5974ccb6f5eeeeeeeeeeeeeeeeeeeeeeee", reviveAccount("0x750ea21c1e98cced0d4557196b6f4a5974ccb6f5"))
}
