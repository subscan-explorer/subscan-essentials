package model

import "github.com/shopspring/decimal"

type ChainBlockJson struct {
	BlockNum          int                   `json:"block_num"`
	BlockTimestamp    int                   `json:"block_timestamp"`
	Hash              string                `json:"hash"`
	ParentHash        string                `json:"parent_hash"`
	StateRoot         string                `json:"state_root"`
	ExtrinsicsRoot    string                `json:"extrinsics_root"`
	Extrinsics        []*ChainExtrinsicJson `json:"extrinsics"`
	Events            []ChainEventJson      `json:"events"`
	Logs              *[]ChainLogJson       `json:"logs"`
	EventCount        int                   `json:"event_count"`
	ExtrinsicsCount   int                   `json:"extrinsics_count"`
	SpecVersion       int                   `json:"spec_version"`
	Validator         string                `json:"validator"`
	ValidatorName     string                `json:"validator_name"`
	ValidatorIndexIds string                `json:"validator_index_ids"`
	Finalized         bool                  `json:"finalized"`
}

type SampleBlockJson struct {
	BlockNum          int    `json:"block_num"`
	BlockTimestamp    int    `json:"block_timestamp"`
	Hash              string `json:"hash"`
	EventCount        int    `json:"event_count"`
	ExtrinsicsCount   int    `json:"extrinsics_count"`
	Validator         string `json:"validator"`
	ValidatorName     string `json:"validator_name"`
	ValidatorIndexIds string `json:"validator_index_ids"`
	Finalized         bool   `json:"finalized"`
}

type ChainExtrinsicJson struct {
	BlockTimestamp     int             `json:"block_timestamp"`
	BlockNum           int             `json:"block_num"`
	ExtrinsicIndex     string          `json:"extrinsic_index"`
	CallModuleFunction string          `json:"call_module_function"`
	CallModule         string          `json:"call_module"`
	Params             string          `json:"params"`
	AccountId          string          `json:"account_id"`
	AccountIndex       string          `json:"account_index"`
	Signature          string          `json:"signature"`
	Nonce              int             `json:"nonce"`
	ExtrinsicHash      string          `json:"extrinsic_hash"`
	Success            bool            `json:"success"`
	Fee                decimal.Decimal `json:"fee"`
}

type ExtrinsicDetail struct {
	BlockTimestamp     int               `json:"block_timestamp"`
	BlockNum           int               `json:"block_num"`
	ExtrinsicIndex     string            `json:"extrinsic_index"`
	CallModuleFunction string            `json:"call_module_function"`
	CallModule         string            `json:"call_module"`
	AccountId          string            `json:"account_id"`
	Signature          string            `json:"signature"`
	Nonce              int               `json:"nonce"`
	ExtrinsicHash      string            `json:"extrinsic_hash"`
	Success            bool              `json:"success"`
	Params             *[]ExtrinsicParam `json:"params"`
	Transfer           *TransferJson     `json:"transfer"`
	Event              *[]ChainEvent     `json:"event"`
	Fee                decimal.Decimal   `json:"fee"`
	Error              *ExtrinsicError   `json:"error"`
	Finalized          bool              `json:"finalized"`
}

type ChainEventJson struct {
	EventIndex     string `json:"event_index"`
	BlockNum       int    `json:"block_num"`
	ExtrinsicIdx   int    `json:"extrinsic_idx"`
	ModuleId       string `json:"module_id"`
	EventId        string `json:"event_id"`
	Params         string `json:"params"`
	EventIdx       int    `json:"event_idx"`
	ExtrinsicHash  string `json:"extrinsic_hash"`
	BlockTimestamp int    `json:"block_timestamp"`
}

type ChainLogJson struct {
	BlockNum   int    `json:"block_num" es:"type:keyword"`
	LogIndex   string `json:"log_index" sql:"default: null;size:100" es:"type:keyword"`
	LogType    string `json:"log_type" es:"type:keyword"`
	OriginType string `json:"origin_type"`
	Data       string `json:"data"`
}

type ChainTransactionJson struct {
	BlockTimestamp     int    `json:"block_timestamp"`
	ExtrinsicIndex     string `json:"extrinsic_index"`
	Signature          string `json:"signature"`
	FromHex            string `json:"from_hex"`
	Destination        string `json:"destination"`
	Hash               string `json:"hash"`
	BlockNum           int    `json:"block_num"`
	CallCode           string `json:"call_code"`
	CallModuleFunction string `json:"call_module_function"`
	CallModule         string `json:"call_module"`
	Params             string `json:"params"`
	Success            bool   `json:"success"`
}

type TransferJson struct {
	From           string          `json:"from"`
	To             string          `json:"to"`
	Module         string          `json:"module"`
	Amount         decimal.Decimal `json:"amount"`
	Hash           string          `json:"hash"`
	BlockTimestamp int             `json:"block_timestamp"`
	BlockNum       int             `json:"block_num"`
	ExtrinsicIndex string          `json:"extrinsic_index"`
	Success        bool            `json:"success"`
	Fee            decimal.Decimal `json:"fee"`
}

type AccountJson struct {
	*AccountSampleJson
	*AccountIdentityJson
	Nonce            int             `json:"nonce"`
	Power            decimal.Decimal `json:"power"`
	Role             string          `json:"role"`
	AccountIndex     string          `json:"account_index"`
	Stash            string          `json:"stash"`
	OutputBlockCount int             `json:"output_block_count"`
}

type AccountSampleJson struct {
	Address     string          `json:"address"`
	Balance     decimal.Decimal `json:"balance"`
	KtonBalance decimal.Decimal `json:"kton_balance"`
	RingLock    decimal.Decimal `json:"ring_lock"`
	KtonLock    decimal.Decimal `json:"kton_lock"`
	Nickname    string          `json:"nickname"`
}

type ExtrinsicsJson struct {
	FromHex            string          `json:"from_hex"`
	BlockNum           int             `json:"block_num"`
	BlockTimestamp     int             `json:"block_timestamp"`
	ExtrinsicIndex     string          `json:"extrinsic_index"`
	Hash               string          `json:"extrinsic_hash"`
	Success            bool            `json:"success"`
	CallModule         string          `json:"call_module"`
	CallModuleFunction string          `json:"call_module_function"`
	Params             string          `json:"params"`
	Fee                decimal.Decimal `json:"fee,omitempty"`
	Destination        string          `json:"destination,omitempty"`
	Amount             decimal.Decimal `json:"amount,omitempty"`
	Finalized          bool            `json:"finalized"`
}

type ValidatorJson struct {
	RankValidator       int             `json:"rank_validator"`
	Nickname            string          `json:"nickname"`
	Display             string          `json:"display"`
	ValidatorStash      string          `json:"validator_stash"`
	ValidatorController string          `json:"validator_controller"`
	BondedNominators    decimal.Decimal `json:"bonded_nominators" sql:"type:decimal(65,0);"`
	BondedOwner         decimal.Decimal `json:"bonded_owner" sql:"type:decimal(65,0);"`
	CountNominators     int             `json:"count_nominators"`
	ValidatorPrefsValue int             `json:"validator_prefs_value"`
	AccountIndex        string          `json:"account_index"`
}

type NominatorJson struct {
	RankNominator  int             `json:"rank_nominator"`
	NominatorStash string          `json:"nominator_stash"`
	Bonded         decimal.Decimal `json:"bonded" sql:"type:decimal(65,0);"`
	Hash           string          `json:"hash"`
}

type AccountNominateList struct {
	ValidatorJson
	Bonded decimal.Decimal `json:"bonded" sql:"type:decimal(65,0);"`
}

type ValidatorPrefsMap struct {
	Address             string
	ValidatorPrefsValue int
}

type StakingRewardMap struct {
	Validator string
	Reward    decimal.Decimal
	Type      string
	Account   string
}

type EventRecord struct {
	Phase        int          `json:"phase"`
	ExtrinsicIdx int          `json:"extrinsic_idx"`
	Type         string       `json:"type"`
	ModuleId     string       `json:"module_id"`
	EventId      string       `json:"event_id"`
	Params       []EventParam `json:"params"`
	Topics       []string     `json:"topics"`
	EventIdx     int          `json:"event_idx"`
}

type AccountIdentityJson struct {
	Display string `json:"display"`
	Web     string `json:"web"`
	Riot    string `json:"riot"`
	Email   string `json:"email"`
	Legal   string `json:"legal"`
	Twitter string `json:"twitter"`
}

type TreasuryProposalJson struct {
	ProposalID   uint            `json:"proposal_id"`
	CreatedBlock uint            `json:"created_block"`
	Status       string          `json:"status" `
	Value        decimal.Decimal `json:"reward"`
	Beneficiary  AccountDisplay  `json:"beneficiary"`
	Proposer     AccountDisplay  `json:"proposer"`
}

type TechcommProposalSample struct {
	ProposalID   uint   `json:"proposal_id"`
	CreatedBlock uint   `json:"created_block"`
	Status       string `json:"status"`
	MemberCount  uint   `json:"member_count"`
	AyeVotes     uint   `json:"aye_votes"`
	NayVotes     uint   `json:"nay_votes"`
}

type TechcommProposalJson struct {
	ProposalID   uint           `json:"proposal_id"`
	CreatedBlock uint           `json:"created_block"`
	UpdatedBlock uint           `json:"updated_block"`
	AyeVotes     uint           `json:"aye_votes"`
	NayVotes     uint           `json:"nay_votes"`
	Status       string         `json:"status"`
	ProposalHash string         `json:"proposal_hash" `
	Proposer     AccountDisplay `json:"proposer"`
	MemberCount  uint           `json:"member_count"`

	Value      decimal.Decimal `json:"value"`
	CallModule string          `json:"call_module"`
	CallName   string          `json:"call_name"`
	Params     string          `json:"params"`

	PreImage *PreImageJson      `json:"pre_image"`
	Votes    []TechcommVoteJson `json:"votes"`
}

type PreImageJson struct {
	Hash         string          `json:"hash"`
	CreatedBlock uint            `json:"created_block"`
	UpdatedBlock uint            `json:"updated_block"`
	Status       string          `json:"status"`
	Deposit      decimal.Decimal `json:"amount"`
	CallModule   string          `json:"call_module"`
	CallName     string          `json:"call_name"`
	Params       string          `json:"params"`
	Author       AccountDisplay  `json:"author"`
}

type AccountDisplay struct {
	Display      string `json:"display"`
	AccountIndex string `json:"account_index"`
	Address      string `json:"address"`
}

type TechcommVoteJson struct {
	Account       AccountDisplay `json:"account"`
	Passed        bool           `json:"passed"`
	ExtrinsicHash string         `json:"extrinsic_hash"`
}

type CouncilProposalSample struct {
	TechcommProposalSample
}
type CouncilProposalJson struct {
	ProposalID   uint           `json:"proposal_id"`
	CreatedBlock uint           `json:"created_block"`
	UpdatedBlock uint           `json:"updated_block"`
	AyeVotes     uint           `json:"aye_votes"`
	NayVotes     uint           `json:"nay_votes"`
	Status       string         `json:"status"`
	ProposalHash string         `json:"proposal_hash" `
	Proposer     AccountDisplay `json:"proposer"`
	MemberCount  uint           `json:"member_count"`

	Value      decimal.Decimal `json:"value"`
	CallModule string          `json:"call_module"`
	CallName   string          `json:"call_name"`
	Params     string          `json:"params"`

	PreImage *PreImageJson     `json:"pre_image"`
	Votes    []CouncilVoteJson `json:"votes"`
}

type CouncilVoteJson struct {
	Account       AccountDisplay `json:"account"`
	Passed        bool           `json:"passed"`
	ExtrinsicHash string         `json:"extrinsic_hash"`
}

type DemocracyReferendumSampleJson struct {
	ReferendumIndex uint   `json:"referendum_index"`
	CreatedBlock    uint   `json:"created_block"`
	VoteThreshold   string `json:"vote_threshold"`
	Status          string `json:"status" sql:"default: null;size:100"`
}

type ReferendumJson struct {
	ReferendumIndex uint            `json:"referendum_index"`
	CreatedBlock    uint            `json:"created_block"`
	UpdatedBlock    uint            `json:"updated_block"`
	VoteThreshold   string          `json:"vote_threshold"`
	PreImage        *PreImageJson   `json:"pre_image"`
	Value           decimal.Decimal `json:"value"`
	Status          string          `json:"status"`
	Delay           uint            `json:"delay"`
	End             uint            `json:"end"`
}

type DemocracyVoteJson struct {
	Account       AccountDisplay  `json:"account"`
	Amount        decimal.Decimal `json:"amount" sql:"type:decimal(65,0);"`
	Passed        bool            `json:"passed"`
	ExtrinsicHash string          `json:"extrinsic_hash"`
}

type AccountEventJson struct {
	EventIndex    string          ` json:"event_index"`
	BlockNum      int             `json:"block_num" `
	ExtrinsicIdx  int             `json:"extrinsic_idx"`
	ModuleId      string          `json:"module_id" `
	EventId       string          `json:"event_id"`
	Params        string          `json:"params"`
	ExtrinsicHash string          `json:"extrinsic_hash"`
	EventIdx      int             `json:"event_idx"`
	Amount        decimal.Decimal `json:"amount"`
}

type DemocracySampleJson struct {
	ProposalID   uint   `json:"proposal_id"`
	CreatedBlock uint   `json:"created_block"`
	Status       string `json:"status" sql:"default: null;size:100"`
}

type DemocracyJson struct {
	ProposalID   uint            `json:"proposal_id"`
	Status       string          `json:"status" `
	CreatedBlock uint            `json:"created_block"`
	UpdatedBlock uint            `json:"updated_block"`
	ProposalHash string          `json:"proposal_hash" `
	Value        decimal.Decimal `json:"value"`
	CallModule   string          `json:"call_module"`
	CallName     string          `json:"call_name"`
	Params       string          `json:"params"`

	PreImage *PreImageJson `json:"pre_image"`
}
