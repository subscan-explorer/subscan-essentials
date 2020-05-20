package util

import (
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
	ns  = []int{0, 1, 2, 3, 4, 5, 6}
	nsd = []int{6, 5, 4, 3, 2, 1, 0}
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

func TestContinuous(t *testing.T) {
	rns := ContinuousSlice(0, 7, "desc")
	rnsd := ContinuousSlice(0, 7, "sced")

	for i := range rns {
		if rns[i] != ns[i] {
			t.Errorf(
				"Generate Continuous int arr failed #%d, got %v, want %v",
				i,
				rns[i],
				ns[i],
			)
		}
	}

	for i := range rnsd {
		if rnsd[i] != nsd[i] {
			t.Errorf(
				"Generate Continuous int arr failed #%d, got %v, want %v",
				i,
				rnsd[i],
				nsd[i],
			)
		}
	}
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
