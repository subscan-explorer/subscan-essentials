package dao

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/itering/subscan/util"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"strings"
)

type TransactionReceipt struct {
	Id              uint64 `json:"id" gorm:"primaryKey;autoIncrement:false" `
	Topics          string `json:"topics" gorm:"type:TEXT" `
	Address         string `json:"address" gorm:"size:70;index:address" `
	TransactionHash string `json:"transaction_hash" gorm:"size:70" `
	Index           int    `json:"index" gorm:"size:32"  `
	Data            string `json:"data" gorm:"size:TEXT" `
	MethodHash      string `json:"method_hash" gorm:"size:70;index:method_hash"`
	BlockTimestamp  uint   `json:"block_timestamp" gorm:"size:32"  `

	Topic1 string `json:"topic1" gorm:"size:70"`
	Topic2 string `json:"topic2" gorm:"size:70"`
	Topic3 string `json:"topic3" gorm:"size:70"`

	BlockNum         uint64 `json:"block_num"   index:"block_num"`
	TransactionIndex uint64 `json:"transaction_index" gorm:"size:32"  `
}

func (t TransactionReceipt) TableName() string {
	return "evm_transaction_receipts"
}

func (t *TransactionReceipt) AfterCreate(txn *gorm.DB) (err error) {
	return t.EventProcess(txn.Statement.Context)
}

type EventLog struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNum         uint64   `json:"block_num"`
	Timestamp        uint64   `json:"timestamp"`
	TransactionHash  string   `json:"transaction_hash"`
	TransactionIndex uint     `json:"transaction_index"`
	LogIndex         uint     `json:"log_index"`
}

func BillionAddress(ctx context.Context) string {
	// d := sg.db
	minBalance := decimal.New(1, 18) // 1 ETH
	var res []AccountsJson
	sg.db.WithContext(ctx).Debug().Select("evm_account,balance").
		Model(&Account{}).Joins("left join balance_accounts on evm_accounts.address=balance_accounts.address").Where("balance>?", minBalance).Scan(&res)
	if len(res) > 0 {
		return res[0].EvmAccount
	}
	return NullAddress
}

func SplitReceiptData(abiStr string, method, data string) []interface{} {
	// split data
	eABI, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil
	}
	events, err := eABI.Unpack(method, util.HexToBytes(data))
	if err != nil {
		return nil
	}
	return events
}
