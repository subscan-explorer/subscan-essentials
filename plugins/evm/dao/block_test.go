package dao

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_EvmBlock(t *testing.T) {
	t.Run("run AddEvmBlock will be ok", func(t *testing.T) {
		ctx := context.TODO()
		var blockNum uint = 11520079
		assert.NoError(t, sg.AddEvmBlock(ctx, blockNum, false))
		block := GetBlockByNum(ctx, int(blockNum))
		assert.NotNil(t, block)
		assert.Equal(t, uint64(blockNum), GetBlockByHash(ctx, block.BlockHash).BlockNum)
	})

	t.Run("input a blockNum that is not exist", func(t *testing.T) {
		ctx := context.TODO()
		var blockNum uint = 99999999999
		assert.Error(t, sg.AddEvmBlock(ctx, blockNum, false))
		block := GetBlockByNum(ctx, int(blockNum))
		assert.Nil(t, block)
	})

	t.Run("run GetBlockNumsByRange will be ok", func(t *testing.T) {
		ctx := context.TODO()
		for i := 1; i < 10; i++ {
			blockNum := uint64(11520079 - i)
			sg.AddOrUpdateItem(ctx, EvmBlock{BlockNum: blockNum, BlockHash: fmt.Sprintf("hash%d", blockNum)}, []string{"block_num"})
		}
		blocksNum := GetBlockNumsByRange(ctx, 11520070, 11520078)
		assert.Equal(t, blocksNum, []int{11520070, 11520071, 11520072, 11520073, 11520074, 11520075, 11520076, 11520077, 11520078})
		blocks := GetBlockByNums(ctx, 11520070, 11520078)
		assert.Len(t, blocks, 9)
	})

	t.Run("GetBlocksSampleByNums will return SampleBlockJson array", func(t *testing.T) {
		ctx := context.TODO()
		blocks, count := GetBlocksSampleByNums(ctx, 0, 10)
		assert.Len(t, blocks, 10)
		assert.Greater(t, count, 0)
	})

	t.Run("Run BlockNums2Blocks will be ok", func(t *testing.T) {
		ctx := context.TODO()
		blockNums := []uint64{11520079, 11520078, 11520077}
		blocks := BlockNums2Blocks(ctx, blockNums)
		assert.Len(t, blocks, len(blockNums))
		for _, block := range blocks {
			assert.NotNil(t, block)
			assert.Contains(t, blockNums, block.BlockNum)
		}
	})
}
