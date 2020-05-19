package model

import (
	"github.com/shopspring/decimal"
	"subscan-end/utiles"
	"subscan-end/utiles/es"
	"time"
)

type ChainMetadata struct {
	Version               int    `json:"version"`
	Metadata              string `json:"metadata"`
	MetadataDecoded       string `json:"metadata_decoded"`
	BlockNum              int    `json:"block_num"`
	CountCallFunctions    int    `json:"count_call_functions"`
	CountEvents           int    `json:"count_events"`
	CountModules          int    `json:"count_modules"`
	CountStorageFunctions int    `json:"count_storage_functions"`
}

type ChainBlock struct {
	ID               uint      `gorm:"primary_key" json:"id"`
	BlockNum         int       `json:"block_num" es:"type:keyword"`
	BlockTimestamp   int       `json:"block_timestamp"`
	CreatedAt        time.Time `json:"created_at" es:"type:date"`
	Hash             string    `sql:"default: null;size:100" json:"hash" es:"type:keyword"`
	ParentHash       string    `sql:"default: null;size:100" json:"parent_hash" es:"type:keyword"`
	StateRoot        string    `sql:"default: null;size:100" json:"state_root" es:"type:keyword"`
	ExtrinsicsRoot   string    `sql:"default: null;size:100" json:"extrinsics_root" es:"type:keyword"`
	Logs             string    `json:"logs" sql:"type:text;"`
	DecodeLogs       string    `json:"decode_logs" sql:"type:text;"`
	Extrinsics       string    `json:"extrinsics" sql:"type:text;"`
	EventCount       int       `json:"event_count"`
	ExtrinsicsCount  int       `json:"extrinsics_count"`
	DecodeExtrinsics string    `json:"decode_extrinsics" sql:"type:text;"`
	Event            string    `json:"event" sql:"type:text;"`
	DecodeEvent      string    `json:"decode_event" sql:"type:text;"`
	SpecVersion      int       `json:"spec_version"  es:"type:text"`
	Validator        string    `json:"validator"`
	CodecError       bool      `json:"codec_error"`
}

func (c *ChainBlock) AfterCreate() {
	if esClient != nil {
		index := es.NewIndexTemplate()
		index.InjectIndex(*c)
		if err := esClient.CreateIndex(index, "chain_block"); err == nil {
			_ = esClient.Insert("chain_block", utiles.IntToString(int(c.ID)), c)
		}
	}
}

type ChainEvent struct {
	ID            uint        `gorm:"primary_key"`
	CreatedAt     time.Time   `json:"created_at" es:"type:date"`
	EventIndex    string      `sql:"default: null;size:100;" json:"event_index" es:"type:keyword"`
	BlockNum      int         `json:"block_num" es:"type:keyword"`
	Phase         int         `json:"phase"`
	ExtrinsicIdx  int         `json:"extrinsic_idx"`
	Type          string      `json:"type"`
	ModuleId      string      `json:"module_id" es:"type:keyword"`
	EventId       string      `json:"event_id" es:"type:keyword"`
	Params        interface{} `json:"params" sql:"type:text;" es:"type:text"`
	ExtrinsicHash string      `json:"extrinsic_hash" sql:"default: null" es:"type:keyword"`
	EventIdx      int         `json:"event_idx"`
}

func (c *ChainEvent) AfterCreate() {
	if esClient != nil {
		index := es.NewIndexTemplate()
		index.InjectIndex(*c)
		if err := esClient.CreateIndex(index, "chain_event"); err == nil {
			_ = esClient.Insert("chain_event", utiles.IntToString(int(c.ID)), c)
		}
	}
}

type ChainExtrinsic struct {
	ID                 uint        `gorm:"primary_key"`
	CreatedAt          time.Time   `json:"created_at"`
	ExtrinsicIndex     string      `json:"extrinsic_index" sql:"default: null;size:100" es:"type:keyword"`
	BlockNum           int         `json:"block_num" es:"type:keyword"`
	BlockTimestamp     int         `json:"block_timestamp"`
	ValueRaw           string      `json:"value_raw"`
	ExtrinsicLength    string      `json:"extrinsic_length"`
	VersionInfo        string      `json:"version_info"`
	CallCode           string      `json:"call_code"`
	CallModuleFunction string      `json:"call_module_function" es:"type:keyword"`
	CallModule         string      `json:"call_module" es:"type:keyword"`
	Params             interface{} `json:"params" sql:"type:text;" es:"type:text"`
	AccountLength      string      `json:"account_length"`
	AccountId          string      `json:"account_id"`
	AccountIndex       string      `json:"account_index"`
	Signature          string      `json:"signature"`
	Nonce              int         `json:"nonce"`
	Era                string      `json:"era"`
	ExtrinsicHash      string      `json:"extrinsic_hash" sql:"default: null" es:"type:keyword"`
	IsSigned           bool        `json:"is_signed"`
	Success            bool        `json:"success"`
}

func (c *ChainExtrinsic) AfterCreate() {
	if esClient != nil {
		index := es.NewIndexTemplate()
		index.InjectIndex(*c)
		if err := esClient.CreateIndex(index, "chain_extrinsic"); err == nil {
			_ = esClient.Insert("chain_extrinsic", utiles.IntToString(int(c.ID)), c)
		}
	}
}

type ChainTransaction struct {
	ID                 uint            `gorm:"primary_key"`
	CreatedAt          time.Time       `json:"created_at"`
	FromHex            string          `json:"from_hex" es:"type:keyword"`
	Destination        string          `json:"destination" es:"type:keyword"`
	ExtrinsicIndex     string          `json:"extrinsic_index" es:"type:keyword"`
	Signature          string          `json:"signature"`
	Success            bool            `json:"success"`
	Hash               string          `json:"hash" es:"type:keyword"`
	BlockNum           int             `json:"block_num" es:"type:keyword"`
	BlockTimestamp     int             `json:"block_timestamp"`
	CallCode           string          `json:"call_code"`
	CallModuleFunction string          `json:"call_module_function" es:"type:keyword"`
	CallModule         string          `json:"call_module" es:"type:keyword"`
	Params             interface{}     `json:"params" sql:"type:text;" es:"type:text"`
	Amount             decimal.Decimal `json:"amount" sql:"type:decimal(30,15);"`
}

func (c *ChainTransaction) AfterCreate() {
	if esClient != nil {
		index := es.NewIndexTemplate()
		index.InjectIndex(*c)
		if err := esClient.CreateIndex(index, "chain_transaction"); err == nil {
			_ = esClient.Insert("chain_transaction", utiles.IntToString(int(c.ID)), c)
		}
	}
}

type ChainAccount struct {
	ID             uint            `gorm:"primary_key"`
	CreatedAt      time.Time       `json:"created_at"`
	Address        string          `sql:"default: null;size:100" json:"address"`
	AccountIndex   int             `json:"account_index"`
	Balance        decimal.Decimal `json:"balance" sql:"type:decimal(30,15);"`
	KtonBalance    decimal.Decimal `json:"kton_balance" sql:"type:decimal(30,15);"`
	CountExtrinsic int             `json:"count_extrinsic"`
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
	Id          int    `json:"id"`
	Name        string `json:"name"`
	SpecVersion int    `json:"spec_version"`
}

type ChainLog struct {
	ID         uint      `gorm:"primary_key"`
	CreatedAt  time.Time `json:"created_at"`
	BlockNum   int       `json:"block_num" es:"type:keyword"`
	LogIndex   string    `json:"log_index" sql:"default: null;size:100" es:"type:keyword"`
	LogType    string    `json:"log_type" es:"type:keyword"`
	OriginType string    `json:"origin_type"`
	Data       string    `json:"data" sql:"type:text;"`
}

type ChainSession struct {
	SessionId       uint `json:"session_id" gorm:"primary_key;auto_increment:false"`
	StartBlock      int  `json:"start_block"`
	EndBlock        int  `json:"end_block"`
	Era             int  `json:"era"`
	CountValidators int  `json:"count_validators"`
	CountNominators int  `json:"count_nominators"`
}

type SessionValidator struct {
	ID                  uint            `gorm:"primary_key"`
	CreatedAt           time.Time       `json:"created_at"`
	SessionId           uint            `json:"session_id"`
	RankValidator       int             `json:"rank_validator"`
	ValidatorStash      string          `json:"validator_stash"`
	ValidatorController string          `json:"validator_controller"`
	ValidatorSession    string          `json:"validator_session"`
	BondedTotal         decimal.Decimal `json:"bonded_total" sql:"type:decimal(65,0);"`
	BondedActive        decimal.Decimal `json:"bonded_active" sql:"type:decimal(65,0);"`
	BondedNominators    decimal.Decimal `json:"bonded_nominators" sql:"type:decimal(65,0);"`
	BondedOwner         decimal.Decimal `json:"bonded_own" sql:"type:decimal(65,0);"`
	Unlocking           string          `json:"unlocking" sql:"type:text;"`
	CountNominators     int             `json:"count_nominators"`
	ValidatorPrefsValue int             `json:"validator_prefs_value"`
}

type SessionNominator struct {
	ID             uint            `gorm:"primary_key"`
	CreatedAt      time.Time       `json:"created_at"`
	SessionId      uint            `json:"session_id"`
	RankValidator  int             `json:"rank_validator"`
	RankNominator  int             `json:"rank_nominator"`
	NominatorStash string          `json:"nominator_stash"`
	Bonded         decimal.Decimal `json:"bonded" sql:"type:decimal(65,0);"`
}

type ValidatorInfo struct {
	Id                  uint   `gorm:"primary_key"`
	ValidatorStash      string `sql:"default: null;size:100" json:"validator_stash"`
	ValidatorController string `sql:"default: null;size:100" json:"validator_controller"`
	ValidatorSession    string `sql:"default: null;size:100" json:"validator_session"`
	ValidatorReward     string `sql:"default: null;size:100" json:"validator_reward"`
	NodeName            string `json:"node_name"`
}
