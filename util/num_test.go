package util

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IntToString(t *testing.T) {
	testCase := []struct {
		i int
		r string
	}{
		{1, "1"},
		{-1, "-1"},
		{0, "0"},
		{2 << 32, "8589934592"},
	}

	for _, test := range testCase {
		assert.Equal(t, test.r, IntToString(test.i))
	}
}

func Test_StringToInt(t *testing.T) {
	testCase := []struct {
		s string
		r int
	}{
		{"1", 1},
		{"-1", -1},
		{"abc", 0},
	}

	for _, test := range testCase {
		assert.Equal(t, test.r, StringToInt(test.s))
	}
}

func Test_InsertInts(t *testing.T) {
	ori := []int{1, 2, 3, 4, 6}
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, InsertInts(ori, 4, 5))
	assert.Equal(t, []int{1, 2, 3, 4, 6, 7}, InsertInts(ori, 7, 7))
}

func Test_IntInSlice(t *testing.T) {
	ori := []int{1, 2, 3, 4, 6}
	assert.Equal(t, true, IntInSlice(1, ori))
	assert.Equal(t, false, IntInSlice(5, ori))
}

func Test_IntFromInterface(t *testing.T) {
	testCase := []struct {
		i64 interface{}
		r   int
	}{
		{"1", 1},
		{1, 1},
		{-1, -1},
		{"abc", 0},
		{float64(1.1), 1},
		{int64(1000), 1000},
		{uint64(2 >> 32), 2 >> 32},
		{[]byte{}, 0},
	}
	for _, test := range testCase {
		assert.Equal(t, test.r, IntFromInterface(test.i64))
	}
}

func Test_Int64FromInterface(t *testing.T) {
	testCase := []struct {
		i64 interface{}
		r   int64
	}{
		{"1", 1},
		{1, 1},
		{-1, -1},
		{"abc", 0},
		{int64(1000), 1000},
		{float64(1.1), 1},
		{uint64(2 >> 32), 2 >> 32},
		{[]byte{}, 0},
	}
	for _, test := range testCase {
		assert.Equal(t, test.r, Int64FromInterface(test.i64))
	}
}

func Test_DecimalFromInterface(t *testing.T) {
	testCase := []struct {
		i interface{}
		r decimal.Decimal
	}{
		{"1", decimal.New(1, 0)},
		{1, decimal.New(1, 0)},
		{-1, decimal.New(-1, 0)},
		{"abc", decimal.Decimal{}},
		{float64(1.1), decimal.NewFromFloat(1.1)},
		{uint64(2 >> 32), decimal.New(int64(2>>32), 0)},
		{int64(2 >> 32), decimal.New(int64(2>>32), 0)},
		{[]byte{}, decimal.Zero},
	}
	for _, test := range testCase {
		assert.Equal(t, test.r, DecimalFromInterface(test.i))
	}
}
