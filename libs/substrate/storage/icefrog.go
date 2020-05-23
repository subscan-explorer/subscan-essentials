package storage

import "github.com/shopspring/decimal"

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
