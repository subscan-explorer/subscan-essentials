package substrate_test

import (
	"github.com/itering/subscan/lib/substrate"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMortal(t *testing.T) {
	var current uint64 = 2497761
	era := "d501"
	mortal := substrate.DecodeMortal(era)
	assert.Equal(t, mortal.Period, uint64(64))
	assert.Equal(t, mortal.Phase, uint64(29))
	assert.Equal(t, mortal.Birth(current), uint64(2497757))
	assert.Equal(t, mortal.Death(current), uint64(2497821))
}
