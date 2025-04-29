package dao

import (
	"context"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/pkg/go-web3/dto"
	"github.com/itering/subscan/share/web3"
	"github.com/itering/subscan/util"
	"strings"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	// TransactionIdGenerateCoefficient Transaction limit 1m 999_999
	TransactionIdGenerateCoefficient = 100_000 // max 999_999
	// TxnReceiptLimit max transaction receipt limit 100k 99_999
	TxnReceiptLimit = 10_000
)

type Transaction struct {
	Hash           string          `json:"hash"  gorm:"primaryKey;autoIncrement:false;size:255"`
	BlockNum       uint            `json:"block_num" gorm:"default: null;size:32;index:block_num" `
	BlockTimestamp uint            `json:"block_timestamp" gorm:"size:32" `
	FromAddress    string          `json:"from_address" gorm:"default: null;size:70;index:sender"`
	ToAddress      string          `json:"to_address" gorm:"default: null;size:70"`
	InputData      string          `json:"input_data" gorm:"type:string"`
	Nonce          uint            `json:"nonce" gorm:"size:32" `
	GasLimit       decimal.Decimal `json:"gas_limit" gorm:"default: 0;type:decimal(40);" `
	GasPrice       decimal.Decimal `json:"gas_price" gorm:"default: 0;type:decimal(40);" `
	GasUsed        decimal.Decimal `json:"gas_used" gorm:"default: 0;type:decimal(40);" `
	Contract       string          `json:"contract" gorm:"size:100"`
	Success        bool            `json:"success"`
	R              string          `json:"r" gorm:"size:100"`
	S              string          `json:"s" gorm:"size:100"`
	V              uint            `json:"v" gorm:"size:32"`
	Value          decimal.Decimal `json:"value" gorm:"default: 0;type:decimal(40);"`
	ExtrinsicIndex string          `json:"extrinsic_index" gorm:"default: null;size:100;"`
	// eip 1559
	EffectiveGasPrice    decimal.Decimal `json:"effective_gas_price" gorm:"default: 0;type:decimal(40);" `
	MaxPriorityFeePerGas decimal.Decimal `json:"max_priority_fee_per_gas" gorm:"default:-1;type:decimal(40)"`
	MaxFeePerGas         decimal.Decimal `json:"max_fee_per_gas" gorm:"default:0;type:decimal(40)"`
	CumulativeGasUsed    decimal.Decimal `json:"cumulative_gas_used" gorm:"default:0;type:decimal(40)"`
	Precompile           uint            `json:"precompile" gorm:"-"`
	TxnType              uint            `json:"txn_type" gorm:"size:32"`
	TransactionIndex     uint64          `json:"transaction_index" gorm:"size:32"`
	// pk
	TransactionId uint64 `json:"transaction_id" gorm:"size:64;index:transaction_id,unique" `
}

type TransactionSample struct {
	Hash              string          `json:"hash"`
	From              string          `json:"from"`
	To                string          `json:"to"`
	Value             decimal.Decimal `json:"value"`
	GasPrice          decimal.Decimal `json:"gas_price"`
	GasUsed           decimal.Decimal `json:"gas_used"`
	Success           bool            `json:"success"`
	BlockTimestamp    uint            `json:"block_timestamp"`
	Contract          string          `json:"contract"`
	ContractName      string          `json:"contract_name"`
	Method            string          `json:"method"`
	DecodeMethod      *string         `json:"decode_method,omitempty"`
	EffectiveGasPrice decimal.Decimal `json:"effective_gas_price"`
	TransactionId     uint64          `json:"transaction_id"`
}

type TransactionInputData struct {
	Hash      string `json:"hash"`
	InputData string `json:"input_data"`
}

type DetailsTokenTransfer struct {
	ID       uint            `json:"id"`
	Contract string          `json:"contract"`
	Hash     string          `json:"hash"`
	CreateAt uint            `json:"create_at"`
	From     string          `json:"from"`
	To       string          `json:"to"`
	Value    decimal.Decimal `json:"value"`
	TokenId  string          `json:"token_id"`
}

type TransactionDetail struct {
	BlockNum       uint                   `json:"block_num"`
	Hash           string                 `json:"hash"`
	BlockTimestamp uint                   `json:"block_timestamp"`
	Success        bool                   `json:"success"`
	ErrorType      string                 `json:"error_type"`
	ErrorMsg       string                 `json:"error_msg"`
	TraceErrorMsg  string                 `json:"trace_error_msg"`
	From           string                 `json:"from"`
	To             ContractDisplay        `json:"to"`
	Contract       string                 `json:"contract"`
	Value          decimal.Decimal        `json:"value" `
	GasLimit       decimal.Decimal        `json:"gas_limit"`
	GasPrice       decimal.Decimal        `json:"gas_price"`
	GasUsed        decimal.Decimal        `json:"gas_used"`
	Nonce          uint                   `json:"nonce"`
	InputData      string                 `json:"input_data"`
	DecodedData    *interface{}           `json:"decoded_data,omitempty"`
	TokenTransfers []DetailsTokenTransfer `json:"token_transfers"`
	R              string                 `json:"r"`
	S              string                 `json:"s"`
	V              uint                   `json:"v"`

	GasFee               decimal.Decimal `json:"gas_fee"`
	EffectiveGasPrice    decimal.Decimal `json:"effective_gas_price"`
	MaxPriorityFeePerGas decimal.Decimal `json:"max_priority_fee_per_gas,omitempty"`
	MaxFeePerGas         decimal.Decimal `json:"max_fee_per_gas,omitempty"`
	TreasuryFees         decimal.Decimal `json:"treasury_fees,omitempty"`
	BurntFee             decimal.Decimal `json:"burnt_fee,omitempty"`
	SavingsFee           decimal.Decimal `json:"savings_fee,omitempty"`
	BaseFee              decimal.Decimal `json:"base_fee,omitempty"`
}

func (t *Transaction) TableName() string {
	return "evm_transactions"
}

func (t *Transaction) AfterCreate(txn *gorm.DB) (err error) {
	ctx := txn.Statement.Context
	// Increase Contract transaction count
	if IsContract(ctx, t.ToAddress) {
		incrContractTransactionCount(ctx, t.ToAddress)
	}
	_, _ = sg.redis.HINCRBY(context.Background(), model.MetadataCacheKey(), "total_transaction", 1)
	return nil
}

// SetEvmAddressRelate set account evm address
// func SetEvmAddressRelate(ctx context.Context, addr ...string) {
// 	if util.IsEvmChain {
// 		return
// 	}
// 	for _, h160 := range addr {
// 		if accountId := H160ToAccountId(ctx, h160); accountId != "" {
// 			if util.ForceUseEvmAddress {
// 				_, _ = sg.db.TouchAccount(ctx, h160)
// 			}
// 			_ = sg.db.SetAccountEvmAddress(ctx, accountId, h160, false)
// 		}
// 	}
// }

func (t *Transaction) TokenTransferRecords(ctx context.Context) (data []DetailsTokenTransfer) {
	var list []TokensTransfers
	sg.db.WithContext(ctx).Model(TokensTransfers{}).Where("hash = ?", t.Hash).Find(&list)
	for _, v := range list {
		data = append(data, DetailsTokenTransfer{
			Contract: v.Contract,
			Hash:     v.Hash,
			CreateAt: v.CreateAt,
			From:     v.Sender,
			To:       v.Receiver,
			Value:    v.Value,
			TokenId:  v.TokenId,
		})
	}
	return
}

func GetTransactionByHash(c context.Context, hash string) *Transaction {
	var t Transaction
	if query := sg.db.WithContext(c).Where("hash = ?", hash).Take(&t); query.Error != nil {
		return nil
	}
	return &t
}

func (s *Storage) CreateTransactionByExecuted(ctx context.Context, blockTimestamp uint, ethTransaction *dto.BlockTransaction, extrinsicIndex string) (err error) {
	transaction := Transaction{
		BlockNum:       uint(util.U256(ethTransaction.BlockNumber).Uint64()),
		BlockTimestamp: blockTimestamp,
		ExtrinsicIndex: extrinsicIndex,
	}

	// build transaction from BlockTransaction
	transaction.Nonce = uint(util.DecimalFromU256(ethTransaction.Nonce).IntPart())
	transaction.GasLimit = util.DecimalFromU256(ethTransaction.Gas)
	transaction.GasPrice = util.DecimalFromU256(ethTransaction.GasPrice)
	transaction.InputData = ethTransaction.Input
	transaction.Value = util.DecimalFromU256(ethTransaction.Value)
	transaction.R = ethTransaction.R
	transaction.V = uint(util.U256(ethTransaction.V).Uint64())
	transaction.S = ethTransaction.S
	transaction.FromAddress = ethTransaction.From
	transaction.Hash = ethTransaction.Hash
	transaction.ToAddress = ethTransaction.To
	transaction.Contract = ethTransaction.Creates
	transaction.MaxPriorityFeePerGas = util.DecimalFromU256(ethTransaction.MaxPriorityFeePerGas)
	transaction.MaxFeePerGas = util.DecimalFromU256(ethTransaction.MaxFeePerGas)
	transaction.TxnType = uint(util.U256(ethTransaction.Type).Uint64())
	transaction.TransactionIndex = util.U256(ethTransaction.TransactionIndex).Uint64()
	transaction.TransactionId = uint64(transaction.BlockNum)*TransactionIdGenerateCoefficient + transaction.TransactionIndex

	// prevent transaction blocking
	var (
		receipts   []TransactionReceipt
		ethReceipt *dto.TransactionReceipt
	)

	//  eth_getTransactionReceipt
	if ethReceipt, err = web3.RPC.Eth.GetTransactionReceipt(ctx, transaction.Hash); err != nil {
		return err
	}

	if ethReceipt.ContractAddress != "" {
		transaction.Contract = ethReceipt.ContractAddress
	}
	// Confirm whether this transaction is to create a contract
	if transaction.Contract != "" {
		_ = transaction.NewContract(ctx)
	}
	// gas used & effective_gas_price
	transaction.GasUsed = decimal.NewFromBigInt(ethReceipt.GasUsed, 0)
	transaction.EffectiveGasPrice = decimal.NewFromBigInt(ethReceipt.EffectiveGasPrice, 0)
	transaction.Success = ethReceipt.Status

	for index, receipt := range ethReceipt.Logs {
		tr := TransactionReceipt{
			Id:              transaction.TransactionId*TxnReceiptLimit + uint64(index), // ensure unique id for each receipt
			Topics:          strings.Join(receipt.Topics, ","),
			Address:         receipt.Address,
			TransactionHash: transaction.Hash,
			Index:           index,
			Data:            util.IfEmptyElse(util.TrimHex(receipt.Data), ""),
			MethodHash:      receipt.Topics[0],
			BlockTimestamp:  transaction.BlockTimestamp,
			BlockNum:        uint64(transaction.BlockNum),
		}
		if len(receipt.Topics) > 1 {
			tr.Topic1 = receipt.Topics[1]
		}
		if len(receipt.Topics) > 2 {
			tr.Topic2 = receipt.Topics[2]
		}
		if len(receipt.Topics) > 3 {
			tr.Topic3 = receipt.Topics[3]
		}
		receipts = append(receipts, tr)
	}
	// receipt
	if err = sg.db.Scopes(model.IgnoreDuplicate).CreateInBatches(receipts, 3000).Error; err != nil {
		return
	}
	// transaction
	query := sg.AddOrUpdateItem(ctx, &transaction, []string{"hash"}, "transaction_index")
	if query.Error != nil {
		return query.Error
	}
	_ = TouchAccount(ctx, transaction.FromAddress)
	_ = TouchAccount(ctx, transaction.ToAddress)
	return nil
}

// func FillTransactionFromChain(ctx context.Context, hash string) *Transaction {
// 	var ethTransaction *dto.TransactionResponse
//
// 	_ = retry.Do(func() error {
// 		var err error
// 		ethTransaction, err = web3.RPC.Eth.GetTransactionByHash(ctx, hash)
// 		return err
// 	}, retry.Attempts(3), retry.Delay(1*time.Second))
// 	if ethTransaction == nil {
// 		return nil
// 	}
//
// 	t := Transaction{}
// 	t.Nonce = uint(ethTransaction.Nonce.Uint64())
// 	t.GasLimit = decimal.NewFromBigInt(ethTransaction.Gas, 0)
// 	t.GasPrice = decimal.NewFromBigInt(ethTransaction.GasPrice, 0)
// 	t.InputData = ethTransaction.Input
// 	t.Value = decimal.NewFromBigInt(ethTransaction.Value, 0)
// 	t.R = ethTransaction.R
// 	t.V = uint(util.U256(ethTransaction.V).Uint64())
// 	t.S = ethTransaction.S
// 	t.From = ethTransaction.From
// 	t.Hash = ethTransaction.Hash
// 	t.To = ethTransaction.To
// 	t.Contract = ethTransaction.Creates
// 	t.MaxPriorityFeePerGas = util.DecimalFromU256(ethTransaction.MaxPriorityFeePerGas)
// 	t.MaxFeePerGas = util.DecimalFromU256(ethTransaction.MaxFeePerGas)
// 	return &t
// }

// Hash2Transaction get transactions map by hash, hash=>transaction
func Hash2Transaction(ctx context.Context, hash []string) map[string]Transaction {
	var transactions []Transaction
	result := make(map[string]Transaction)

	// Query transactions by hash
	sg.db.WithContext(ctx).Where("hash IN ?", hash).Find(&transactions)

	// Map transactions by hash
	for _, t := range transactions {
		result[t.Hash] = t
	}

	return result
}
