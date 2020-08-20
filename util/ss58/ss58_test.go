package ss58_test

import (
	"github.com/itering/subscan/util/ss58"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecode(t *testing.T) {
	assert.Equal(t, ss58.Decode("fawfafwaf", 42), "")
	assert.Equal(t, ss58.Decode(
		"5FcEGUiujfdWyf6RME1G8pCTkmkgXFDECaTSpVDWVnNiZJXR", 42),
		"9cbfadc7579a27fcb3ea4bb1940aade652d1dd9a2dc69c9920f1de42d8ca0234")
	assert.Equal(t, ss58.Decode(
		"132LXQzhfrQRAdWyfSr1hzSWrfGh7KcDrFVeyajLVcHy3voi", 0),
		"597d6fed2bdf74f6187833c49967e95f83e5c0a8e710ade206f6755f3448e769")
	assert.Equal(t, ss58.Decode(
		"FAHqreQSkzH5BsXFJN1m6touWNGHPpu11LuCDbyzVa5fnck", 2),
		"72612c619e1a5b8b2001fb484fd06882df5a41ae6e36afb38592a922429a2814")

}

func TestEncode(t *testing.T) {
	assert.Equal(t, ss58.Encode("0x1234567", 42), "")
	assert.Equal(t, ss58.Encode(
		"0x88b3bfe1410ed8a12cd8a2c230e97cfd5a9fb1cc95ac859ec9c9a2ecfe7cf84f", 2),
		"FfZRiEyrJwgxFZx1QsCnDjaJCHXoeUS4v4Hs1Yo8GpVveNQ")
	assert.Equal(t, ss58.Encode(
		"0xf2cb2711b197eef9f2803aa2f087a1cedfeae2e10f55ef9242230efe18454491", 42),
		"5HZ3o1uoA6oKYjb86YnuSU2nbz8dw1LNj6joFzguGtn2wHu2")
	assert.Equal(t, ss58.Encode(
		"0x3ccbd50810c15f4cec3462ddb73b1ba5982cfb8643b9214e715a785e1e88e500", 1),
		"8ANgaUSe4rALo2qjPYHYsDLLEGKf8ww9Y3wrpsUrSYgSE9K")
}
