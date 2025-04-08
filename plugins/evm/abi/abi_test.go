package abi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodingMethod(t *testing.T) {
	assert.Equal(t, "ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		EncodingMethod("Transfer(address,address,uint256)"))
	assert.Equal(t, "e1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c",
		EncodingMethod("Deposit(address,uint256)"))
	assert.Equal(t, "7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65",
		EncodingMethod("Withdrawal(address,uint256)"))
}

func TestDecodeAddress(t *testing.T) {
	assert.Equal(t, "2fdeec247037a343b26f5daa0d3af15fe23a28c5",
		DecodeAddress("0x0000000000000000000000002fdeec247037a343b26f5daa0d3af15fe23a28c5"))
	assert.Equal(t, "f3c1444cd449bd66ef6da7ca6c3e7884840a3995",
		DecodeAddress("0x000000000000000000000000f3c1444cd449bd66ef6da7ca6c3e7884840a3995"))
}

func TestDecodeStaticType(t *testing.T) {
	assert.Equal(t, "2fdeec247037a343b26f5daa0d3af15fe23a28c5", DecodeStaticType("0x0000000000000000000000002fdeec247037a343b26f5daa0d3af15fe23a28c5"))
	assert.Equal(t, "f3c1444cd449bd66ef6da7ca6c3e7884840a3995", DecodeStaticType("0x000000000000000000000000f3c1444cd449bd66ef6da7ca6c3e7884840a3995"))
	assert.Equal(t, "", DecodeStaticType("0x"))
	assert.Equal(t, "123", DecodeStaticType("0x0000000000000000000000000000000000000123"))
}
