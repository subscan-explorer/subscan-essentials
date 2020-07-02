package storage

import "github.com/shopspring/decimal"

type StakingLedgers struct {
	Stash             string            `json:"stash"`
	TotalRing         decimal.Decimal   `json:"total_ring,omitempty"`
	TotalRegularRing  decimal.Decimal   `json:"total_deposit_ring,omitempty"`
	ActiveRing        decimal.Decimal   `json:"active_ring,omitempty"`
	ActiveRegularRing decimal.Decimal   `json:"active_regular_ring,omitempty"`
	TotalKton         decimal.Decimal   `json:"total_kton,omitempty"`
	ActiveKton        decimal.Decimal   `json:"active_kton,omitempty"`
	Total             decimal.Decimal   `json:"total,omitempty"`
	Active            decimal.Decimal   `json:"active,omitempty"`
	RegularItems      []TimeDepositItem `json:"deposit_items,omitempty"`
	Unlocking         []UnlockChunk     `json:"unlocking"`
	ActiveDepositRing decimal.Decimal   `json:"active_deposit_ring,omitempty"`
	LastReward        int               `json:"last_reward,omitempty"`
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
	Total          decimal.Decimal  `json:"total,omitempty"`
	Own            decimal.Decimal  `json:"own,omitempty"`
	OwnRingBalance decimal.Decimal  `json:"own_ring_balance,omitempty"`
	OwnKtonBalance decimal.Decimal  `json:"own_kton_balance,omitempty"`
	OwnPower       decimal.Decimal  `json:"own_power,omitempty"`
	TotalPower     decimal.Decimal  `json:"total_power,omitempty"`
	Others         []IndividualExpo `json:"others,omitempty"`
}

type IndividualExpo struct {
	Who         string          `json:"who"`
	Value       decimal.Decimal `json:"value,omitempty"`
	RingBalance decimal.Decimal `json:"ring_balance,omitempty"`
	KtonBalance decimal.Decimal `json:"kton_balance,omitempty"`
	Power       decimal.Decimal `json:"power,omitempty"`
}

type RawAuraPreDigest struct {
	SlotNumber int64 `json:"slotNumber"`
}

type ValidatorPrefsLegacy struct {
	UnstakeThreshold      int             `json:"unstake_threshold"`
	ValidatorPaymentRatio decimal.Decimal `json:"validator_payment_ratio"`
}

type BalanceLock struct {
	Id           string          `json:"id"`
	Amount       decimal.Decimal `json:"amount,omitempty"`
	Until        uint64          `json:"until,omitempty"`
	WithdrawLock *WithdrawLock   `json:"withdraw_lock,omitempty"`
	Reasons      []string        `json:"reasons,omitempty"`
	LockReasons  interface{}     `json:"lock_reasons,omitempty"`
	LockFor      *LockFor        `json:"lock_for,omitempty"`
}

type WithdrawLock struct {
	WithStaking *struct {
		StakingAmount decimal.Decimal `json:"staking_amount"`
		Unbondings    []NormalLock    `json:"unbondings"`
	} `json:"WithStaking,omitempty"`
	Normal NormalLock `json:"Normal,omitempty"`
}

type NormalLock struct {
	Amount decimal.Decimal `json:"amount"`
	Until  uint64          `json:"until"`
}

type Registration struct {
	Deposit decimal.Decimal `json:"deposit"`
	Info    IdentityInfo    `json:"info"`
}

type IdentityInfo struct {
	Additional     []IdentityInfoAdditional `json:"additional"`
	Display        Data                     `json:"display"`
	Legal          Data                     `json:"legal"`
	Web            Data                     `json:"web"`
	Riot           Data                     `json:"riot"`
	Email          Data                     `json:"email"`
	Image          Data                     `json:"image"`
	PgpFingerprint Data                     `json:"pgpFingerprint"`
	Twitter        Data                     `json:"twitter,omitempty"`
}

type Data struct {
	None        string `json:"none,omitempty"`
	Raw         string `json:"raw,omitempty"`
	BlakeTwo256 string `json:"BlakeTwo256,omitempty"`
	Sha256      string `json:"Sha256,omitempty"`
	Keccak256   string `json:"Keccak256,omitempty"`
	ShaThree256 string `json:"ShaThree256,omitempty"`
}

type IdentityInfoAdditional struct {
	Field Data `json:"field"`
	Value Data `json:"value"`
}

type AccountData struct {
	Free         decimal.Decimal `json:"free"`
	Reserved     decimal.Decimal `json:"reserved"`
	FreeKton     decimal.Decimal `json:"free_kton,omitempty"`
	ReservedKton decimal.Decimal `json:"reserved_kton,omitempty"`
	MiscFrozen   decimal.Decimal `json:"misc_frozen"`
	FeeFrozen    decimal.Decimal `json:"fee_frozen"`
}

type Proposal struct {
	CallModule string           `json:"call_module"`
	CallName   string           `json:"call_name"`
	Params     []ExtrinsicParam `json:"params"`
}
type ExtrinsicParam struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
	ValueRaw string      `json:"valueRaw"`
}

type ReferendumInfo struct {
	End          uint                    `json:"end,omitempty"`
	ProposalHash string                  `json:"proposalHash,omitempty"`
	Threshold    string                  `json:"threshold,omitempty"`
	Delay        uint                    `json:"delay,omitempty"`
	Ongoing      *ReferendumStatus       `json:"ongoing,omitempty"`
	Finished     *ReferendumInfoFinished `json:"finished,omitempty"`
}

type ReferendumStatus struct {
	End          uint   `json:"end"`
	ProposalHash string `json:"proposalHash,omitempty"`
	Threshold    string `json:"threshold,omitempty"`
	Delay        uint   `json:"delay,omitempty"`
}

type ReferendumInfoFinished struct {
	Approved bool `json:"approved"`
	End      uint `json:"end"`
}

type AccountInfo struct {
	Nonce    int             `json:"nonce"`
	Refcount int             `json:"refcount,omitempty"`
	Data     AccountData     `json:"data,omitempty"`
	Free     decimal.Decimal `json:"free,omitempty"`
}

type ActiveEraInfo struct {
	Index int   `json:"index"`
	Start int64 `json:"start"`
}

type DecoderLog struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type AccountVote struct {
	Standard *AccountVoteStandard `json:"Standard,omitempty"`
	Split    *AccountVoteSplit    `json:"Split,omitempty"`
}

type AccountVoteStandard struct {
	Vote    int             `json:"vote"`
	Balance decimal.Decimal `json:"balance"`
}

type AccountVoteSplit struct {
	Aye decimal.Decimal `json:"aye"`
	Nay decimal.Decimal `json:"nay"`
}

type Payee struct {
	Controller string       `json:"Controller,omitempty"`
	Stash      string       `json:"Stash,omitempty"`
	Staked     *PayeeStaked `json:"staked,omitempty"`
}

type PayeeStaked struct {
	PromiseMonth int `json:"promise_month"`
}

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

type RawBabePreDigest struct {
	Primary   *RawBabePreDigestPrimary      `json:"primary,omitempty"`
	Secondary *RawBabePreDigestSecondary    `json:"secondary,omitempty"`
	VRF       *RawBabePreDigestSecondaryVRF `json:"VRF,omitempty"`
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

type RawBabePreDigestSecondaryVRF struct {
	AuthorityIndex uint   `json:"authorityIndex"`
	SlotNumber     uint64 `json:"slotNumber"`
	VrfOutput      string `json:"vrfOutput"`
	VrfProof       string `json:"vrfProof"`
}

type EraPoints struct {
	Total      decimal.Decimal   `json:"total"`
	Individual []decimal.Decimal `json:"individual"`
}
