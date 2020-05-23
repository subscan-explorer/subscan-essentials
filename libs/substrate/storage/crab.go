package storage

import "github.com/shopspring/decimal"

type LockFor struct {
	StakingLock *StakingLock `json:"Staking"`
	Common      *Common      `json:"Common"`
}

type StakingLock struct {
	StakingAmount decimal.Decimal `json:"staking_amount"`
	Unbondings    []Unbondings    `json:"unbondings"`
}

type Common struct {
	Amount decimal.Decimal `json:"amount"`
}

type Unbondings struct {
	Amount decimal.Decimal `json:"amount"`
	Moment int             `json:"moment"`
}

type MerkleMountainRangeRootLog struct {
	MmrRoot string `json:"mmr_root"`
}
