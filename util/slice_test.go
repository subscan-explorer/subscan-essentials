package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	str = "subscan"
	sst = []string{"subscan0", "subscan", "subscan1"}
	ssf = []string{"subscan0", "subscan1", "subscan2"}
	ssm = map[string]bool{
		"subscan0": true,
		"subscan1": true,
		"subscan2": true,
	}
)

func TestLookup(t *testing.T) {
	if StringInSlice(str, sst) == false {
		t.Errorf(
			"Lookup string in string slice failed, got %v, want %v",
			false,
			true,
		)
	}

	if StringInSlice(str, ssf) == true {
		t.Errorf(
			"Lookup string in string slice failed, got %v, want %v",
			true,
			false,
		)
	}
}

func TestContinuousNums(t *testing.T) {
	assert.Nil(t, ContinuousNums(7, 0, "asc"))
	assert.Equal(t, []int{6, 5, 4, 3, 2, 1, 0}, ContinuousNums(6, 7, "desc"))
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6}, ContinuousNums(0, 7, "asc"))
	assert.Equal(t, []int{6, 5, 4, 3, 2, 1, 0}, ContinuousNums(6, 8, "desc"))

}

func TestMap(t *testing.T) {
	res := MapStringToSlice(ssm)
	resLen := len(res)
	ssfLen := len(ssf)
	if resLen != ssfLen {
		t.Errorf(
			"Map string to string slice length failed, got %v, want %v",
			resLen,
			ssfLen,
		)
	}

	for i := range res {
		if !StringInSlice(res[i], ssf) {
			t.Errorf(
				"Map string to string slice failed #%d, got %v, want %v",
				i,
				res[i],
				ssf[i],
			)
		}
	}
}
