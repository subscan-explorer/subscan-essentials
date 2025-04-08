package dao

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/itering/subscan/util"
	"strings"

	"gorm.io/gorm"
)

type TransactionReceipt struct {
	Id              uint64 `json:"id" gorm:"primaryKey;autoIncrement:false" `
	Topics          string `json:"topics" gorm:"type:TEXT" `
	Address         string `json:"address" gorm:"size:100" `
	TransactionHash string `json:"transaction_hash" gorm:"size:100" `
	Index           int    `json:"index" gorm:"size:32"  `
	Data            string `json:"data" gorm:"size:TEXT" `
	MethodHash      string `json:"method_hash" gorm:"size:100"`
	BlockTimestamp  uint   `json:"block_timestamp" gorm:"size:32"  `

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

// BillionAddress todo
func BillionAddress(ctx context.Context) string {
	// d := sg.db
	var h160 string
	// minBalance := decimal.New(1, 18) // 1 ETH
	// d.GetDB(ctx).Model(&model.ChainAccount{}).Select("evm_address").Where("account_balance > ?", minBalance).Where("evm_address like ?", "0x%").Take(&h160)
	return h160
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
