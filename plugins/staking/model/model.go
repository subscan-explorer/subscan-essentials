package model

import (
	"github.com/itering/subscan/util/address"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type Payout struct {
	ID             uint                `gorm:"primary_key" json:"-"`
	Account        address.SS58Address `gorm:"index" sql:"default: null;size:100" json:"account"`
	Amount         decimal.Decimal     `sql:"type:decimal(30,0);" json:"amount"`
	Era            uint32              `gorm:"index" json:"era"`
	Stash          address.SS58Address `gorm:"index" sql:"default: null;size:100" json:"stash"`
	ValidatorStash address.SS58Address `gorm:"index" sql:"default: null;size:100" json:"validator_stash"`
	BlockTimestamp uint64              `json:"block_timestamp"`
	ModuleId       string              `json:"module_id"`
	EventId        string              `json:"event_id"`
	SlashKton      bool                `json:"slash_kton"`
	ExtrinsicIndex string              `json:"extrinsic_index"`
	EventIndex     string              `json:"event_index"`
}

type ValidatorPrefs struct {
	ID                uint                `gorm:"primary_key" json:"-"`
	Account           address.SS58Address `gorm:"index;unique" sql:"default: null;size:100" json:"account"`
	Commission        decimal.Decimal     `sql:"type:decimal(12,11);" json:"commission"`
	BlockedNomination bool                `json:"blocked_nomination"`
	CommissionHistory string              `sql:"type:text;" json:"commission_history"`
}

type EraInfo struct {
	ID               uint            `gorm:"primary_key" json:"-"`
	Era              uint32          `gorm:"index" json:"era"`
	TotalStake       decimal.Decimal `sql:"type:decimal(30,0);" json:"total_stake"`
	Stakes           datatypes.JSONSlice[EraStake]
	TotalPoints      uint32
	TotalRewards     decimal.Decimal
	ValidatorPoints  datatypes.JSONType[map[address.SS58Address]uint32]
	ValidatorRewards datatypes.JSONType[map[address.SS58Address]decimal.Decimal]
	StakerRewards    datatypes.JSONType[map[address.SS58Address]decimal.Decimal]
}

type EraStake struct {
	Validator      address.SS58Address
	Staker         address.SS58Address
	Amount         decimal.Decimal
	ValidatorTotal decimal.Decimal
}
