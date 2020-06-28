package ss58_test

import (
	"github.com/itering/subscan/internal/util/ss58"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecode(t *testing.T) {
	address := "5FcEGUiujfdWyf6RME1G8pCTkmkgXFDECaTSpVDWVnNiZJXR"
	assert.Equal(t, ss58.Decode(address, 42), "9cbfadc7579a27fcb3ea4bb1940aade652d1dd9a2dc69c9920f1de42d8ca0234")
}

func TestEncode(t *testing.T) {
	address := "0x88b3bfe1410ed8a12cd8a2c230e97cfd5a9fb1cc95ac859ec9c9a2ecfe7cf84f"
	assert.Equal(t, ss58.Encode(address, 2), "FfZRiEyrJwgxFZx1QsCnDjaJCHXoeUS4v4Hs1Yo8GpVveNQ")
}
