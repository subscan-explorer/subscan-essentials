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

type DarwiniaStakingLedgers struct {
	Stash             string          `json:"stash"`
	ActiveRing        decimal.Decimal `json:"active_ring"`
	ActiveKton        decimal.Decimal `json:"active_kton"`
	ActiveDepositRing decimal.Decimal `json:"active_deposit_ring"`
}

type IcefrogPrefsLegacy struct {
	ValidatorPaymentRatio decimal.Decimal `json:"validator_payment_ratio"`
	NodeName              string          `json:"node_name"`
}

type IceValidatorPrefs struct {
	ValidatorPrefs *IcefrogPrefsLegacy `json:"col1"`
	Linkage        *LinkageAccountId   `json:"col2"`
}
