package model

import (
	"fmt"
	"time"

	"github.com/itering/subscan/libs/substrate"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/es"
	"github.com/shopspring/decimal"
)

// SplitTableBlockNum
var SplitTableBlockNum = 1000000

type ChainMetadata struct {
	Version int    `json:"version"`
	Type    string `json:"type"` // storage/call/event/constant/error
	Name    string `json:"name"`
	Prefix  string `json:"prefix"`
}

type ChainBlock struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	BlockNum        int       `json:"block_num" es:"type:keyword"`
	BlockTimestamp  int       `json:"block_timestamp"`
	CreatedAt       time.Time `json:"created_at" es:"type:date"`
	Hash            string    `sql:"default: null;size:100" json:"hash" es:"type:keyword"`
	ParentHash      string    `sql:"default: null;size:100" json:"parent_hash" es:"type:keyword"`
	StateRoot       string    `sql:"default: null;size:100" json:"state_root" es:"type:keyword"`
	ExtrinsicsRoot  string    `sql:"default: null;size:100" json:"extrinsics_root" es:"type:keyword"`
	Logs            string    `json:"logs" sql:"type:text;"`
	Extrinsics      string    `json:"extrinsics" sql:"type:MEDIUMTEXT;"`
	EventCount      int       `json:"event_count"`
	ExtrinsicsCount int       `json:"extrinsics_count"`
	Event           string    `json:"event" sql:"type:text;"`
	SpecVersion     int       `json:"spec_version"  es:"type:text"`
	Validator       string    `json:"validator"`
	CodecError      bool      `json:"codec_error"`
	Finalized       bool      `json:"finalized"`
}

func (c ChainBlock) TableName() string {
	if c.BlockNum/SplitTableBlockNum == 0 {
		return "chain_blocks"
	}
	return fmt.Sprintf("chain_blocks_%d", c.BlockNum/SplitTableBlockNum)
}

func (c *ChainBlock) AsOpenBlockJson() *OpenBlockJson {
	bj := OpenBlockJson{
		BlockNum:       c.BlockNum,
		BlockTimestamp: c.BlockTimestamp,
		Hash:           c.Hash,
		ParentHash:     c.ParentHash,
		StateRoot:      c.StateRoot,
		ExtrinsicsRoot: c.ExtrinsicsRoot,
		SpecVersion:    c.SpecVersion,
		Validator:      substrate.SS58Address(c.Validator),
		Finalized:      c.Finalized,
	}
	return &bj
}

func (c *ChainBlock) AfterCreate() {
	if esClient != nil {
		index := es.NewIndexTemplate()
		index.InjectIndex(*c)
		if err := esClient.CreateIndex(index, "chain_block"); err == nil {
			_ = esClient.Insert("chain_block", util.IntToString(int(c.ID)), c)
		}
	}
}

type ChainEvent struct {
	ID            uint        `gorm:"primary_key" json:"-"`
	CreatedAt     time.Time   `json:"-" es:"type:date"`
	EventIndex    string      `sql:"default: null;size:100;" json:"event_index" es:"type:keyword"`
	BlockNum      int         `json:"block_num" es:"type:keyword"`
	Phase         int         `json:"-"`
	ExtrinsicIdx  int         `json:"extrinsic_idx"`
	Type          string      `json:"-"`
	ModuleId      string      `json:"module_id" es:"type:keyword"`
	EventId       string      `json:"event_id" es:"type:keyword"`
	Params        interface{} `json:"params" sql:"type:text;" es:"type:text"`
	ExtrinsicHash string      `json:"extrinsic_hash" sql:"default: null" es:"type:keyword"`
	EventIdx      int         `json:"event_idx"`
	Finalized     bool        `json:"finalized"`
}

func (c ChainEvent) TableName() string {
	if c.BlockNum/SplitTableBlockNum == 0 {
		return "chain_events"
	}
	return fmt.Sprintf("chain_events_%d", c.BlockNum/SplitTableBlockNum)
}

func (c *ChainEvent) AfterCreate() {
	if esClient != nil {
		index := es.NewIndexTemplate()
		index.InjectIndex(*c)
		if err := esClient.CreateIndex(index, "chain_event"); err == nil {
			_ = esClient.Insert("chain_event", util.IntToString(int(c.ID)), c)
		}
	}
}

type ChainExtrinsic struct {
	ID                 uint            `gorm:"primary_key"`
	CreatedAt          time.Time       `json:"created_at"`
	ExtrinsicIndex     string          `json:"extrinsic_index" sql:"default: null;size:100" es:"type:keyword"`
	BlockNum           int             `json:"block_num" es:"type:keyword"`
	BlockTimestamp     int             `json:"block_timestamp"`
	ValueRaw           string          `json:"value_raw"`
	ExtrinsicLength    string          `json:"extrinsic_length"`
	VersionInfo        string          `json:"version_info"`
	CallCode           string          `json:"call_code"`
	CallModuleFunction string          `json:"call_module_function" es:"type:keyword" sql:"size:100"`
	CallModule         string          `json:"call_module" es:"type:keyword" sql:"size:100"`
	Params             interface{}     `json:"params" sql:"type:MEDIUMTEXT;" es:"type:text"`
	AccountLength      string          `json:"account_length"`
	AccountId          string          `json:"account_id"`
	AccountIndex       string          `json:"account_index"`
	Signature          string          `json:"signature"`
	Nonce              int             `json:"nonce"`
	Era                string          `json:"era"`
	ExtrinsicHash      string          `json:"extrinsic_hash" sql:"default: null" es:"type:keyword"`
	IsSigned           bool            `json:"is_signed"`
	Success            bool            `json:"success"`
	Fee                decimal.Decimal `json:"fee" sql:"type:decimal(30,0);"`
	Finalized          bool            `json:"finalized"`
	BatchIndex         int             `json:"-" gorm:"-"`
}

func (c ChainExtrinsic) TableName() string {
	if c.BlockNum/SplitTableBlockNum == 0 {
		return "chain_extrinsics"
	}
	return fmt.Sprintf("chain_extrinsics_%d", c.BlockNum/SplitTableBlockNum)
}

func (c *ChainExtrinsic) AfterCreate() {
	if esClient != nil {
		index := es.NewIndexTemplate()
		index.InjectIndex(*c)
		if err := esClient.CreateIndex(index, "chain_extrinsic"); err == nil {
			_ = esClient.Insert("chain_extrinsic", util.IntToString(int(c.ID)), c)
		}
	}
}

type ChainTransaction struct {
	ID                 uint            `gorm:"primary_key"`
	CreatedAt          time.Time       `json:"created_at"`
	FromHex            string          `json:"from_hex" es:"type:keyword"`
	Destination        string          `json:"destination" es:"type:keyword"`
	ExtrinsicIndex     string          `sql:"default: null;size:100" json:"extrinsic_index" es:"type:keyword"`
	Signature          string          `json:"signature"`
	Success            bool            `json:"success"`
	Hash               string          `sql:"default: null;size:100" json:"hash" es:"type:keyword"`
	BlockNum           int             `json:"block_num" es:"type:keyword"`
	BlockTimestamp     int             `json:"block_timestamp"`
	CallCode           string          `json:"call_code"`
	CallModuleFunction string          `json:"call_module_function" es:"type:keyword" sql:"size:100"`
	CallModule         string          `json:"call_module" es:"type:keyword" sql:"size:100"`
	Params             interface{}     `json:"params" sql:"type:MEDIUMTEXT;" es:"type:text"`
	Amount             decimal.Decimal `json:"amount" sql:"type:decimal(30,15);"`
	Fee                decimal.Decimal `json:"fee" sql:"type:decimal(30,0);"`
	Finalized          bool            `json:"finalized"`
}

func (c ChainTransaction) TableName() string {
	if c.BlockNum/SplitTableBlockNum == 0 {
		return "chain_transactions"
	}
	return fmt.Sprintf("chain_transactions_%d", c.BlockNum/SplitTableBlockNum)
}

func (c *ChainTransaction) AfterCreate() {
	if esClient != nil {
		index := es.NewIndexTemplate()
		index.InjectIndex(*c)
		if err := esClient.CreateIndex(index, "chain_transaction"); err == nil {
			_ = esClient.Insert("chain_transaction", util.IntToString(int(c.ID)), c)
		}
	}
}

type ChainAccount struct {
	ID               uint            `gorm:"primary_key"`
	CreatedAt        time.Time       `json:"created_at"`
	Address          string          `sql:"default: null;size:100" json:"address"`
	Nickname         string          `sql:"default: null;size:100" json:"nickname"`
	AccountIndex     int             `json:"account_index"`
	Nonce            int             `json:"nonce"`
	Balance          decimal.Decimal `json:"balance" sql:"type:decimal(30,15);"`
	KtonBalance      decimal.Decimal `json:"kton_balance" sql:"type:decimal(30,15);"`
	CountExtrinsic   int             `json:"count_extrinsic"`
	RingLock         decimal.Decimal `json:"ring_lock" sql:"type:decimal(30,0);"`
	KtonLock         decimal.Decimal `json:"kton_lock" sql:"type:decimal(30,0);"`
	OutputBlockCount int             `json:"output_block_count" sql:"default: 0"`
}

type ExtrinsicParam struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
	ValueRaw string      `json:"valueRaw"`
}

type EventParam struct {
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
	ValueRaw string      `json:"valueRaw"`
}

type RuntimeVersion struct {
	Id                    int    `json:"-"`
	Name                  string `json:"-"`
	SpecVersion           int    `json:"spec_version"`
	CountCallFunctions    int    `json:"-"`
	CountEvents           int    `json:"-"`
	CountStorageFunctions int    `json:"-"`
	CountConstants        int    `json:"-"`
	CountErrorType        int    `json:"-"`
	Modules               string `json:"modules"`
	RawData               string `json:"-" sql:"type:MEDIUMTEXT;"`
}

type ChainLog struct {
	ID         uint      `gorm:"primary_key"`
	CreatedAt  time.Time `json:"created_at"`
	BlockNum   int       `json:"block_num" es:"type:keyword"`
	LogIndex   string    `json:"log_index" sql:"default: null;size:100" es:"type:keyword"`
	LogType    string    `json:"log_type" es:"type:keyword"`
	OriginType string    `json:"origin_type"`
	Data       string    `json:"data" sql:"type:text;"`
	Finalized  bool      `json:"finalized"`
}

func (c ChainLog) TableName() string {
	if c.BlockNum/SplitTableBlockNum == 0 {
		return "chain_logs"
	}
	return fmt.Sprintf("chain_logs_%d", c.BlockNum/SplitTableBlockNum)
}

type ChainSession struct {
	SessionId       uint            `json:"session_id" gorm:"primary_key;auto_increment:false"`
	StartBlock      int             `json:"start_block"`
	EndBlock        int             `json:"end_block"`
	Era             int             `json:"era"`
	CountValidators int             `json:"count_validators"`
	CountNominators int             `json:"count_nominators"`
	ValidatorsBond  decimal.Decimal `json:"validators_bond" sql:"type:decimal(65,0);"`
	TotalBond       decimal.Decimal `json:"validators_total_bond"  sql:"type:decimal(65,0);"`
}

type SessionValidator struct {
	ID                  uint            `gorm:"primary_key"`
	CreatedAt           time.Time       `json:"created_at"`
	SessionId           uint            `json:"session_id"`
	RankValidator       int             `json:"rank_validator"`
	ValidatorStash      string          `json:"validator_stash" sql:"default: null;size:100"`
	ValidatorController string          `json:"validator_controller"`
	BondedTotal         decimal.Decimal `json:"bonded_total" sql:"type:decimal(65,0);"`
	BondedActive        decimal.Decimal `json:"bonded_active" sql:"type:decimal(65,0);"`
	BondedNominators    decimal.Decimal `json:"bonded_nominators" sql:"type:decimal(65,0);"`
	BondedOwner         decimal.Decimal `json:"bonded_owner" sql:"type:decimal(65,0);"`
	Unlocking           string          `json:"unlocking" sql:"type:text;"`
	CountNominators     int             `json:"count_nominators"`
	ValidatorPrefsValue int             `json:"validator_prefs_value"`
}

type SessionNominator struct {
	ID             uint            `gorm:"primary_key"`
	CreatedAt      time.Time       `json:"created_at"`
	SessionID      uint            `json:"session_id"`
	ValidatorStash string          `json:"validator_stash"`
	RankValidator  int             `json:"rank_validator"`
	RankNominator  int             `json:"rank_nominator"`
	NominatorStash string          `json:"nominator_stash"`
	Bonded         decimal.Decimal `json:"bonded" sql:"type:decimal(65,0);"`
}

type ValidatorInfo struct {
	Id                  uint            `gorm:"primary_key"`
	ValidatorStash      string          `sql:"default: null;size:100" json:"validator_stash"`
	ValidatorController string          `sql:"default: null;size:100" json:"validator_controller"`
	ValidatorReward     string          `sql:"default: null;size:100" json:"validator_reward"`
	NodeName            string          `json:"node_name"`
	BondedTotal         decimal.Decimal `json:"bonded_total" sql:"type:decimal(65,0);"`
	BondedNominators    decimal.Decimal `json:"bonded_nominators" sql:"type:decimal(65,0);"`
	BondedOwner         decimal.Decimal `json:"bonded_owner" sql:"type:decimal(65,0);"`
	CountNominators     int             `json:"count_nominators"`
	ValidatorPrefsValue int             `json:"validator_prefs_value"`
	Selected            bool            `json:"selected"`
}

type StakingReward struct {
	Id             uint            `gorm:"primary_key"`
	Account        string          `sql:"default: null;size:100" json:"account"`
	Era            int             `json:"era"`
	SessionId      uint            `json:"session_id"`
	BlockTimestamp int             `json:"block_timestamp"`
	RewardType     string          `sql:"default: null;size:50" json:"reward_type"`
	Reward         decimal.Decimal `json:"reward" sql:"type:decimal(30,15);"`
	Validator      string          `sql:"default: null;size:100" json:"validator"`
}

type BondRecord struct {
	Id                      uint            `gorm:"primary_key"`
	Account                 string          `sql:"default: null;size:100" json:"account"`
	ExtrinsicIndex          string          `json:"extrinsic_index" sql:"default: null;size:100"`
	StartAt                 int64           `json:"start_at"`
	Month                   int             `json:"month"`
	Amount                  decimal.Decimal `json:"amount" sql:"type:decimal(60,0);"`
	Status                  string          `json:"status"`
	ExpiredAt               int64           `json:"expired_at"`
	UnbondingExtrinsicIndex string          `json:"unbonding_extrinsic_index" sql:"default: null;size:100"`
	UnbondingAt             int64           `json:"unbonding_at"`
	UnbondingEnd            int64           `json:"unbonding_end"`
	Currency                string          `json:"currency"`
	PunishUnlock            bool            `json:"unlock"`
}

type Nominator struct {
	ID             uint            `gorm:"primary_key"`
	CreatedAt      time.Time       `json:"created_at"`
	ValidatorStash string          `sql:"default: null;size:100" json:"validator_stash"`
	RankNominator  int             `json:"rank_nominator"`
	NominatorStash string          `sql:"default: null;size:100" json:"nominator_stash"`
	Bonded         decimal.Decimal `json:"bonded" sql:"type:decimal(65,0);"`
}

type SlashRecord struct {
	Id             uint            `gorm:"primary_key"`
	CreatedAt      time.Time       `json:"created_at"`
	Account        string          `json:"account" sql:"default: null;size:100"`
	Amount         decimal.Decimal `json:"amount" sql:"type:decimal(30,0);"`
	ExtrinsicIndex string          `json:"extrinsic_index" sql:"default: null;size:100"`
	BlockNum       int             `json:"block_num" es:"type:keyword"`
}

type OutputBlockRecord struct {
	Id        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	Account   string    `json:"account" sql:"default: null;size:100"`
	BlockNum  int       `json:"block_num"`
}

type ChainMappingRecord struct {
	Id             uint            `gorm:"primary_key"`
	Account        string          `json:"account" sql:"default: null;size:100"`
	BlockNum       int             `json:"block_num"`
	ExtrinsicIndex string          `json:"extrinsic_index" sql:"default: null;size:100"`
	Amount         decimal.Decimal `json:"amount" sql:"type:decimal(30,0);"`
	MappingType    string          `json:"mapping_type"`
	FromTx         string          `json:"from_tx"`
	MappingAt      int             `json:"mapping_at"`
}

type ValidatorStat struct {
	Id            uint            `gorm:"primary_key"`
	Account       string          `json:"account" sql:"default: null;size:100"`
	Era           uint            `json:"era"`
	StartBlockNum uint            `json:"start_block_num"`
	EndBlockNum   uint            `json:"end_block_num"`
	Reward        decimal.Decimal `json:"reward" sql:"type:decimal(65,0);" `
	Slash         decimal.Decimal `json:"slash" sql:"type:decimal(65,0);"`
	BlockProduced string          `json:"block_produced" sql:"type:text;"`
}

type AccountIdentityInfo struct {
	AccountId      uint   `json:"account_id" gorm:"primary_key;auto_increment:false"`
	Account        string `json:"account" sql:"default: null;size:100"`
	Display        string `json:"display"`
	Legal          string `json:"legal"`
	Web            string `json:"web"`
	Riot           string `json:"riot"`
	Email          string `json:"email"`
	Image          string `json:"image"`
	PgpFingerprint string `json:"pgpFingerprint"`
	Twitter        string `json:"twitter"`
}

type DemocracyProposal struct {
	ProposalID   uint            `json:"proposal_id"  gorm:"primary_key;auto_increment:false" sql:"default: 0"`
	CreatedBlock uint            `json:"created_block"`
	UpdatedBlock uint            `json:"updated_block"`
	Status       string          `json:"status" sql:"default: null;size:100"`
	ProposalHash string          `json:"proposal_hash" sql:"default: null;size:100"`
	Value        decimal.Decimal `json:"reward" sql:"type:decimal(65,0);"`
	CallModule   string          `json:"call_module"`
	CallName     string          `json:"call_name"`
	Params       interface{}     `json:"params" sql:"type:text;" es:"type:text"`
}

type DemocracyReferendum struct {
	ReferendumIndex uint            `json:"referendum_index" gorm:"primary_key;auto_increment:false" sql:"default: 0"`
	CreatedBlock    uint            `json:"created_block"`
	UpdatedBlock    uint            `json:"updated_block"`
	Delay           uint            `json:"delay"`
	End             uint            `json:"end"`
	VoteThreshold   string          `json:"vote_threshold" sql:"default: null;size:100"`
	Status          string          `json:"status" sql:"default: null;size:100"`
	Value           decimal.Decimal `json:"value" sql:"type:decimal(65,0);"`
	PreImage        string          `json:"pre_image"`
}

type DemocracyPreImage struct {
	ID           uint            `gorm:"primary_key"`
	Hash         string          `json:"hash" sql:"default: null;size:100"`
	CreatedBlock uint            `json:"created_block"`
	UpdatedBlock uint            `json:"updated_block"`
	Status       string          `json:"status" sql:"default: null;size:100"`
	Deposit      decimal.Decimal `json:"amount" sql:"type:decimal(65,0);"`
	CallModule   string          `json:"call_module"`
	CallName     string          `json:"call_name"`
	Params       interface{}     `json:"params" sql:"type:MEDIUMTEXT;" es:"type:text"`
	Author       string          `json:"author"`
}

type DemocracyVote struct {
	ID              uint            `gorm:"primary_key"`
	BlockNum        uint            `json:"block_num"`
	ReferendumIndex uint            `json:"referendum_index"`
	Account         string          `json:"account" sql:"default: null;size:100"`
	Amount          decimal.Decimal `json:"amount" sql:"type:decimal(65,0);"`
	Passed          bool            `json:"passed"`
	ExtrinsicHash   string          `json:"extrinsic_hash" sql:"default: null;size:100"`
}

type TechcommProposal struct {
	ProposalID   uint            `json:"proposal_id"  gorm:"primary_key;auto_increment:false" sql:"default: 0"`
	CreatedBlock uint            `json:"created_block"`
	UpdatedBlock uint            `json:"updated_block"`
	Status       string          `json:"status" sql:"default: null;size:100"`
	ProposalHash string          `json:"proposal_hash" sql:"default: null;size:100"`
	Value        decimal.Decimal `json:"reward" sql:"type:decimal(65,0);"`
	CallModule   string          `json:"call_module"`
	CallName     string          `json:"call_name"`
	Proposer     string          `json:"proposer"`
	Params       interface{}     `json:"params" sql:"type:text;" es:"type:text"`
	MemberCount  uint            `json:"member_count"`
	AyeVotes     uint            `json:"aye_votes"`
	NayVotes     uint            `json:"nay_votes"`
	PreImage     string          `json:"pre_image"`
}

// 技术委员会投票
type TechcommVote struct {
	ID            uint   `gorm:"primary_key"`
	BlockNum      uint   `json:"block_num"`
	ProposalID    uint   `json:"proposal_id"`
	ProposalHash  string `json:"proposal_hash" sql:"default: null;size:100"`
	Account       string `json:"account" sql:"default: null;size:100"`
	Passed        bool   `json:"passed"`
	ExtrinsicHash string `json:"extrinsic_hash" sql:"default: null;size:100"`
}

// 议会
type CouncilProposal struct {
	ProposalID   uint            `json:"proposal_id"  gorm:"primary_key;auto_increment:false" sql:"default: 0"`
	CreatedBlock uint            `json:"created_block"`
	UpdatedBlock uint            `json:"updated_block"`
	Status       string          `json:"status" sql:"default: null;size:100"`
	ProposalHash string          `json:"proposal_hash" sql:"default: null;size:100"`
	Value        decimal.Decimal `json:"reward" sql:"type:decimal(65,0);"`
	CallModule   string          `json:"call_module"`
	CallName     string          `json:"call_name"`
	Proposer     string          `json:"proposer"`
	Params       interface{}     `json:"params" sql:"type:text;" es:"type:text"`
	MemberCount  uint            `json:"member_count"`
	Executed     bool            `json:"executed"`
	AyeVotes     uint            `json:"aye_votes"`
	NayVotes     uint            `json:"nay_votes"`
	PreImage     string          `json:"pre_image"`
}

// 财政提案
type TreasuryProposal struct {
	ProposalID   uint            `json:"proposal_id"  gorm:"primary_key;auto_increment:false" sql:"default: 0"`
	CreatedBlock uint            `json:"created_block"`
	UpdatedBlock uint            `json:"updated_block"`
	Status       string          `json:"status" sql:"default: null;size:100"`
	Value        decimal.Decimal `json:"reward" sql:"type:decimal(65,0);"`
	Beneficiary  string          `json:"beneficiary" sql:"default: null;size:100"`
	Proposer     string          `json:"proposer" sql:"default: null;size:100"`
}

// 议会投票
type CouncilVote struct {
	ID            uint   `gorm:"primary_key"`
	BlockNum      uint   `json:"block_num"`
	ProposalID    uint   `json:"proposal_id"`
	ProposalHash  string `json:"proposal_hash" sql:"default: null;size:100"`
	Account       string `json:"account" sql:"default: null;size:100"`
	Passed        bool   `json:"passed"`
	ExtrinsicHash string `json:"extrinsic_hash" sql:"default: null;size:100"`
}

type TokenClaim struct {
	ID      uint            `gorm:"primary_key"`
	Account string          `json:"account" sql:"size:100"`
	Target  string          `json:"target" sql:"size:100"`
	Amount  decimal.Decimal `json:"amount" sql:"type:decimal(65,0);"`
	ClaimAt int             `json:"claim_at"`
}

type AccountEvent struct {
	ID            uint            `gorm:"primary_key" json:"-"`
	Account       string          `json:"-" sql:"size:100;"`
	EventIndex    string          `sql:"default: null;size:100;" json:"event_index"`
	BlockNum      int             `json:"block_num" `
	ExtrinsicIdx  int             `json:"extrinsic_idx"`
	ModuleId      string          `json:"module_id" sql:"size:100;"`
	EventId       string          `json:"event_id" sql:"size:100;"`
	Params        interface{}     `json:"params" sql:"type:text;"`
	ExtrinsicHash string          `json:"extrinsic_hash" sql:"size:100;"`
	EventIdx      int             `json:"event_idx"`
	Amount        decimal.Decimal `json:"amount" sql:"type:decimal(65,0);"`
	SlashKton     decimal.Decimal `json:"slash_kton"  sql:"type:decimal(65,0);"`
}

type ExtrinsicError struct {
	ID            uint   `gorm:"primary_key" json:"-"`
	ExtrinsicHash string `json:"-" sql:"size:100;"`
	Module        string `json:"module"`
	Name          string `json:"name"`
	Doc           string `json:"doc"`
}

type DispatchError struct {
	Other        *string     `json:"Other,omitempty"`
	CannotLookup *string     `json:"CannotLookup,omitempty"`
	BadOrigin    *string     `json:"BadOrigin,omitempty"`
	Module       interface{} `json:"Module,omitempty"`
	Error        *int        `json:"error,omitempty"`
}

type DispatchErrorModule struct {
	Index int `json:"index"`
	Error int `json:"error"`
}

// from transfer event
type TransferHistory struct {
	ID             uint            `gorm:"primary_key" json:"-"`
	From           string          `sql:"default: null;size:100" json:"from"`
	To             string          `sql:"default: null;size:100" json:"to"`
	ExtrinsicIndex string          `sql:"default: null;size:100" json:"extrinsic_index"`
	EventIndex     string          `sql:"default: null;size:100" json:"-"`
	Success        bool            `json:"success"`
	Hash           string          `sql:"default: null;size:100" json:"hash"`
	BlockNum       int             `json:"block_num"`
	BlockTimestamp int             `json:"block_timestamp"`
	Module         string          `sql:"default: null;size:100" json:"module"`
	Params         interface{}     `json:"-" sql:"type:MEDIUMTEXT;"`
	Amount         decimal.Decimal `json:"amount" sql:"type:decimal(30,15);"`
	Fee            decimal.Decimal `json:"fee" sql:"type:decimal(30,0);"`
	Finalized      bool            `json:"-"`
}
