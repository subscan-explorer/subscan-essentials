package dao

import (
	"context"
	"fmt"
	"github.com/itering/subscan/pkg/go-web3/dto"
	"github.com/itering/subscan/share/web3"
	"github.com/itering/subscan/util"
	"math/big"
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/shopspring/decimal"
)

type EvmBlock struct {
	BlockNum         uint64          `json:"block_num" gorm:"size:64;index:block_num,unique"`
	BlockHash        string          `json:"block_hash" gorm:"size:70;index:block_hash"`
	ParentHash       string          `json:"parent_hash" gorm:"size:70"`
	Sha3Uncles       string          `json:"sha3_uncles" gorm:"size:70"`
	Author           string          `json:"author" gorm:"size:70"`
	Miner            string          `json:"miner" gorm:"size:70"`
	StateRoot        string          `json:"state_root" gorm:"size:70"`
	TransactionsRoot string          `json:"transactions_root" gorm:"size:70"`
	ReceiptsRoot     string          `json:"receipts_root" gorm:"size:70"`
	GasUsed          decimal.Decimal `json:"gas_used"  gorm:"type:decimal(65);"`
	GasLimit         decimal.Decimal `json:"gas_limit" gorm:"type:decimal(65);"`
	ExtraData        string          `json:"extra_data" gorm:"size:255"`
	LogsBloom        string          `json:"logs_bloom" gorm:"type:TEXT"`
	Timestamp        uint            `json:"timestamp"  gorm:"size:64"`
	Difficulty       decimal.Decimal `json:"difficulty" gorm:"type:decimal(65);"`
	TotalDifficulty  decimal.Decimal `json:"total_difficulty" gorm:"type:decimal(65);"`
	SealFields       string          `json:"seal_fields" gorm:"size:255"`
	Uncles           string          `json:"uncles" gorm:"size:255"`
	BlockSize        decimal.Decimal `json:"block_size" gorm:"type:decimal(65);"`
	TransactionCount int             `json:"transaction_count" gorm:"size:32"`
	BaseFeePerGas    decimal.Decimal `json:"base_fee_per_gas" gorm:"type:decimal(65);"`
}

func (t *EvmBlock) TableName() string {
	return "evm_blocks"
}

func (s *Storage) AddEvmBlock(ctx context.Context, blockNum uint, force bool) error {
	if block := GetBlockByNum(ctx, int(blockNum)); block != nil && !force {
		return nil
	}
	blockRaw, err := web3.RPC.Eth.GetBlockByNumber(ctx, big.NewInt(int64(blockNum)), true)
	if err != nil {
		return err
	}
	return s.processBlock(ctx, uint64(blockNum), blockRaw)
}

func (s *Storage) processBlock(ctx context.Context, blockNum uint64, blockRaw *dto.Block) (err error) {
	block := &EvmBlock{
		BlockHash:        blockRaw.Hash,
		ParentHash:       blockRaw.ParentHash,
		Sha3Uncles:       blockRaw.Sha3Uncles,
		Author:           blockRaw.Author,
		Miner:            blockRaw.Miner,
		StateRoot:        blockRaw.StateRoot,
		TransactionsRoot: blockRaw.TransactionsRoot,
		ReceiptsRoot:     blockRaw.ReceiptsRoot,
		BlockNum:         blockNum,
		GasUsed:          util.DecimalFromU256(blockRaw.GasUsed),
		GasLimit:         util.DecimalFromU256(blockRaw.GasLimit),
		ExtraData:        blockRaw.ExtraData,
		LogsBloom:        blockRaw.LogsBloom,
		Timestamp:        uint(util.U256(blockRaw.Timestamp).Int64()),
		Difficulty:       util.DecimalFromU256(blockRaw.Difficulty),
		TotalDifficulty:  util.DecimalFromU256(blockRaw.TotalDifficulty),
		SealFields:       util.ToString(blockRaw.SealFields),
		Uncles:           util.ToString(blockRaw.Uncles),
		BlockSize:        util.DecimalFromU256(blockRaw.Size),
		TransactionCount: len(blockRaw.Transactions),
		BaseFeePerGas:    util.DecimalFromU256(blockRaw.BaseFeePerGas),
	}
	hash2ExtrinsicIndex, err := findOutSubstrateExecutedEvent(ctx, uint(blockNum), blockRaw)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup

	cp, _ := ants.NewPoolWithFunc(5, func(i interface{}) {
		transaction := i.(dto.BlockTransaction)
		// override contract address
		defer wg.Done()

		if e := s.CreateTransactionByExecuted(ctx, block.Timestamp, &transaction, hash2ExtrinsicIndex[transaction.Hash]); e != nil {
			err = e
		}
	})
	defer cp.Release()
	for _, t := range blockRaw.Transactions {
		wg.Add(1)
		_ = cp.Invoke(t)
	}
	wg.Wait()
	if err != nil {
		return err
	}
	return s.AddOrUpdateItem(ctx, block, []string{"block_num"}, "transaction_count").Error
}

func GetBlockByNum(ctx context.Context, blockNum int) *EvmBlock {
	var block EvmBlock
	txn := sg.db.WithContext(ctx)
	if query := txn.Where("block_num = ?", blockNum).First(&block); query.Error != nil {
		return nil
	}
	return &block
}

func GetBlockByHash(ctx context.Context, hash string) *EvmBlock {
	var block EvmBlock
	txn := sg.db.WithContext(ctx)
	if query := txn.Where("block_hash = ?", hash).First(&block); query.Error != nil {
		return nil
	}
	return &block
}

func GetBlockNumsByRange(ctx context.Context, start, end int) []int {
	var blockNums []int
	sg.db.WithContext(ctx).Model(EvmBlock{}).Where("block_num BETWEEN ? AND ?", start, end).Order("block_num asc").Pluck("block_num", &blockNums)
	return blockNums
}

func GetBlockByNums(ctx context.Context, start, end int) (list []EvmBlock) {
	sg.db.WithContext(ctx).Model(EvmBlock{}).Where("block_num BETWEEN ? AND ?", start, end).Order("block_num asc").Find(&list)
	return
}

type SampleBlockJson struct {
	BlockNum         uint64          `json:"block_num" gorm:"size:64;index:block_num,unique"`
	BlockHash        string          `json:"block_hash" gorm:"size:70;index:block_hash,unique"`
	Author           string          `json:"author" gorm:"size:70"`
	GasUsed          decimal.Decimal `json:"gas_used"  gorm:"type:decimal(65);"`
	GasLimit         decimal.Decimal `json:"gas_limit" gorm:"type:decimal(65);"`
	Timestamp        uint            `json:"timestamp"  gorm:"size:64"`
	TransactionCount int             `json:"transaction_count" gorm:"size:32"`
}

func GetBlocksSampleByNums(c context.Context, page, row int) ([]SampleBlockJson, int) {
	var blockJson []SampleBlockJson
	blocks, count := GetBlockList(c, page, row)
	for _, block := range blocks {
		blockJson = append(blockJson, SampleBlockJson{
			BlockNum:         block.BlockNum,
			BlockHash:        block.BlockHash,
			Author:           block.Author,
			GasUsed:          block.GasUsed,
			GasLimit:         block.GasLimit,
			Timestamp:        block.Timestamp,
			TransactionCount: block.TransactionCount,
		})
	}
	return blockJson, count
}

func GetBlockList(c context.Context, page, row int) ([]EvmBlock, int) {
	var blocks []EvmBlock
	blockNum := int(latestBlockNum(c))
	if blockNum == 0 {
		return nil, 0
	}
	head := blockNum - page*row
	if head < 0 {
		return nil, blockNum
	}
	end := head - row
	if end < 0 {
		end = 0
	}

	sg.db.WithContext(c).Model(EvmBlock{}).
		Joins(fmt.Sprintf("JOIN (SELECT block_num from %s where block_num BETWEEN %d and %d order by block_num desc ) as t on %s.block_num=t.block_num",
			"evm_blocks",
			end+1, head,
			"evm_blocks",
		)).
		Order("t.block_num desc").Scan(&blocks)
	return blocks, blockNum
}

func latestBlockNum(ctx context.Context) uint {
	var block EvmBlock
	q := sg.db.WithContext(ctx).Model(EvmBlock{}).Last(&block)
	if q.Error != nil {
		return 0
	}
	return uint(block.BlockNum)
}

func BlockNums2Blocks(ctx context.Context, blockNums []uint64) map[uint64]EvmBlock {
	var blocks []EvmBlock
	sg.db.WithContext(ctx).Model(EvmBlock{}).Where("block_num in ?", blockNums).Find(&blocks)
	blockMap := make(map[uint64]EvmBlock)
	for _, block := range blocks {
		blockMap[block.BlockNum] = block
	}
	return blockMap
}
