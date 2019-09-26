package storage

import "github.com/shopspring/decimal"

type StakingLedgers struct {
	Stash             string            `json:"stash"`
	TotalRing         decimal.Decimal   `json:"total_ring"`
	TotalRegularRing  decimal.Decimal   `json:"total_deposit_ring"`
	ActiveRing        decimal.Decimal   `json:"active_ring"`
	ActiveRegularRing decimal.Decimal   `json:"active_deposit_ring"`
	TotalKton         decimal.Decimal   `json:"total_kton"`
	ActiveKton        decimal.Decimal   `json:"active_kton"`
	RegularItems      []TimeDepositItem `json:"deposit_items"`
	Unlocking         []UnlockChunk     `json:"unlocking"`
}

type TimeDepositItem struct {
	Value      decimal.Decimal `json:"value"`
	StartTime  int             `json:"start_time"`
	ExpireTime int             `json:"expire_time"`
}

type UnlockChunk struct {
	Value         StakingBalance `json:"value"`
	Era           int            `json:"era"`
	IsTimeDeposit bool           `json:"is_time_deposit"`
}

type StakingBalance struct {
	Ring decimal.Decimal `json:"Ring,omitempty"`
	Kton decimal.Decimal `json:"Kton,omitempty"`
}

type Exposures struct {
	Total  decimal.Decimal  `json:"total"`
	Own    decimal.Decimal  `json:"own"`
	Others []IndividualExpo `json:"others"`
}

type IndividualExpo struct {
	Who   string          `json:"who"`
	Value decimal.Decimal `json:"value"`
}

type RawAuraPreDigest struct {
	SlotNumber int64 `json:"slotNumber"`
}

type ValidatorPrefsLegacy struct {
	UnstakeThreshold      int             `json:"unstake_threshold"`
	ValidatorPaymentRatio decimal.Decimal `json:"validator_payment_ratio"`
}
