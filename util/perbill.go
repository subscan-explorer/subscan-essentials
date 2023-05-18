package util

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

var billion = decimal.NewFromInt(1_000_000_000)

type PerBill struct {
	value uint64
}

type Rounding = uint8

const (
	roundingDown Rounding = iota
	roundingUp
)

func divRounded(n, d decimal.Decimal, r Rounding) decimal.Decimal {
	o := n.Div(d)
	switch r {
	case roundingDown:
		return o
	case roundingUp:
		if n.Mod(d).IsZero() {
			o = o.Add(decimal.NewFromInt(1))
		}
	}
	return o
}

func (p PerBill) decimalValue() decimal.Decimal {
	v, _ := decimal.NewFromString(fmt.Sprint(p.value))
	return v
}

func (p PerBill) Mul(v decimal.Decimal) decimal.Decimal {
	return v.Mul(p.decimalValue()).Div(billion)
}

func PerBillFromRational(p, q decimal.Decimal) (PerBill, error) {
	var bad PerBill
	if q.IsZero() {
		return bad, errors.New("quotient is zero")
	}
	if p.GreaterThan(q) {
		return bad, errors.New("numerator is greater than denominator")
	}

	factor := decimal.Max(divRounded(p, billion, roundingUp), decimal.NewFromInt(1))
	qReduce := divRounded(q, factor, roundingDown)
	pReduce := divRounded(p, factor, roundingDown)
	n := pReduce.Mul(billion)
	d := qReduce
	part := divRounded(n, d, roundingDown)
	return PerBill{value: part.BigInt().Uint64()}, nil
}
