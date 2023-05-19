package storage

import (
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Dao interface {
	DB
	// Find spec metadata raw
	SpecialMetadata(int) string

	// Substrate websocket rpc pool
	RPCPool() *websocket.PoolConn

	// Plugin set prefix
	SetPrefix(string)

	GetRuntimeConstant(moduleName string, constantName string) *RuntimeConstant
}

type Option struct {
	PluginPrefix string
	PageSize     int
	Page         int
	Order        string
}

// DB interface
// Every query can be found here https://gorm.io/docs/
type DB interface {
	Query(model interface{}) *gorm.DB

	// Can query database all tables data
	// Query ** no prefix ** table default, option PluginPrefix can specify other plugin model
	FindBy(record interface{}, query interface{}, option *Option) (int, bool)

	// Only can exec plugin relate tables
	// Migration
	AutoMigration(model interface{}) error
	// Add column Index
	AddIndex(model interface{}, indexName string, columns ...string) error
	// Add column unique index
	AddUniqueIndex(model interface{}, indexName string, columns ...string) error

	// Create one record
	Create(record interface{}) error
	// Update one or more column
	Update(model interface{}, query interface{}, attr map[string]interface{}) error
	// Delete one or more record
	Delete(model interface{}, query interface{}) error
}

type RuntimeConstant struct {
	SpecVersion  int    `json:"spec_version"`
	ModuleName   string `json:"module_name"`
	ConstantName string `json:"constant_name"`
	Type         string `json:"type"`
	Value        string `json:"value"`
}

type Block struct {
	BlockNum       int    `json:"block_num"`
	BlockTimestamp int    `json:"block_timestamp"`
	Hash           string `json:"hash"`
	SpecVersion    int    `json:"spec_version"`
	Validator      string `json:"validator"`
	Finalized      bool   `json:"finalized"`
}

type Extrinsic struct {
	ExtrinsicIndex     string          `json:"extrinsic_index" `
	CallCode           string          `json:"call_code"`
	CallModuleFunction string          `json:"call_module_function" `
	CallModule         string          `json:"call_module"`
	Params             []byte          `json:"params"`
	AccountId          string          `json:"account_id"`
	Signature          string          `json:"signature"`
	Nonce              int             `json:"nonce"`
	Era                string          `json:"era"`
	ExtrinsicHash      string          `json:"extrinsic_hash"`
	Success            bool            `json:"success"`
	Fee                decimal.Decimal `json:"fee"`
}

type Event struct {
	BlockNum      int    `json:"block_num"`
	ExtrinsicIdx  int    `json:"extrinsic_idx"`
	ModuleId      string `json:"module_id"`
	EventId       string `json:"event_id"`
	Params        string `json:"params"`
	ExtrinsicHash string `json:"extrinsic_hash"`
	EventIdx      int    `json:"event_idx"`
	EventIndex    string `json:"event_index"`
}

type ExtrinsicParam struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type EventParam struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type CallArg struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func (a CallArg) GetName() string {
	return a.Name
}

func (a CallArg) GetValue() interface{} {
	return a.Value
}

type Call struct {
	BlockNum       int       `json:"block_num"`
	CallIdx        int       `json:"call_idx"`
	BlockTimestamp int       `json:"block_timestamp"`
	ExtrinsicHash  string    `json:"extrinsic_hash"`
	ModuleId       string    `json:"module_id"`
	CallId         string    `json:"call_id"`
	Params         []CallArg `json:"params" sql:"type:text;"`
	Events         []Event   `json:"events"`
}
