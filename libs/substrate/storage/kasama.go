package storage

import "github.com/shopspring/decimal"

type RawBabePreDigest struct {
	Primary   *RawBabePreDigestPrimary   `json:"primary,omitempty"`
	Secondary *RawBabePreDigestSecondary `json:"secondary,omitempty"`
}

type RawBabePreDigestPrimary struct {
	AuthorityIndex uint   `json:"authorityIndex"`
	SlotNumber     uint64 `json:"slotNumber"`
	Weight         uint   `json:"weight"`
	VrfOutput      string `json:"vrfOutput"`
	VrfProof       string `json:"vrfProof"`
}

type RawBabePreDigestSecondary struct {
	AuthorityIndex uint   `json:"authorityIndex"`
	SlotNumber     uint64 `json:"slotNumber"`
	Weight         uint   `json:"weight"`
}

type LinkageAccountId struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type ValidatorPrefsLinkage struct {
	ValidatorPrefs *LowerVersionValidatorPrefsLegacy `json:"col1,omitempty"`
	Linkage        *LinkageAccountId                 `json:"col2,omitempty"`
	Commission     decimal.Decimal                   `json:"commission,omitempty"`
}

type LowerVersionValidatorPrefsLegacy struct {
	Commission            decimal.Decimal `json:"commission"`
	ValidatorPaymentRatio decimal.Decimal `json:"validator_payment_ratio,omitempty"`
	NodeName              string          `json:"node_name,omitempty"`
}

type EraPoints struct {
	Total      decimal.Decimal   `json:"total"`
	Individual []decimal.Decimal `json:"individual"`
}
