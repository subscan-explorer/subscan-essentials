package dao

import (
	"context"
	"errors"
	"fmt"
	customerror "github.com/itering/subscan/pkg/go-web3/constants"
	"github.com/itering/subscan/pkg/go-web3/dto"
	"github.com/itering/subscan/share/web3"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/network"
	"github.com/panjf2000/ants/v2"
	"log"
	"math/big"
	"sync"
)

func Validate(startBlock int, fastMod bool, blockNum int) {
	ctx := context.TODO()
	db := sg.db
	// latest fill block num
	finalizedBlock := int(latestBlockNum(ctx))
	var err error
	const holdOnNum = 10

	util.Logger().Info(fmt.Sprintf("Now: block height %d", finalizedBlock))
	var latestUpdateBlockNum int

	var fillBlock = func(num int, force bool) *dto.Block {
		var blockRaw *dto.Block
		blockRaw, err = web3.RPC.Eth.GetBlockByNumber(ctx, big.NewInt(int64(num)), true)
		if err != nil {
			return nil
		}
		err = sg.processBlock(ctx, uint64(blockNum), blockRaw)
		if err != nil {
			if errors.Is(err, customerror.EMPTYRESPONSE) && network.CurrentIs(network.AssethubWestend) {
				// ignore empty response
				return nil
			}
			log.Println("Error processing block:", err)
			return nil
		}
		return blockRaw
	}

	if blockNum > 0 {
		util.Logger().Info(fmt.Sprintf("Start checkout block %d", blockNum))
		fillBlock(blockNum, true)
		return
	}

	wg := new(sync.WaitGroup)
	cp, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		wg.Add(1)
		defer wg.Done()
		var num = i.(int)
		fillBlock(num, false)
	})
	defer cp.Release()

	if startBlock > 0 {
		latestUpdateBlockNum = startBlock - 1
	}

	type NSum struct{ N int64 }

	for {
		if latestUpdateBlockNum >= finalizedBlock-holdOnNum {
			break
		}
		endBlockNum := latestUpdateBlockNum + 3000
		if endBlockNum > finalizedBlock-holdOnNum {
			endBlockNum = finalizedBlock - holdOnNum
		}

		util.Logger().Info(fmt.Sprintf("Start checkout block %d, end block %d", latestUpdateBlockNum+1, endBlockNum))
		var allFetchBlockNums = GetBlockNumsByRange(ctx, latestUpdateBlockNum+1, endBlockNum)

		// 检查 block 完整性
		if len(allFetchBlockNums) < endBlockNum-(latestUpdateBlockNum+1) {
			for i := latestUpdateBlockNum + 1; i <= endBlockNum; i++ {
				if !util.IntInSlice(i, allFetchBlockNums) {
					util.Logger().Warning(fmt.Sprintf("Missing block %d", i))
					if fastMod {
						_ = cp.Invoke(i)
					} else {
						fillBlock(i, false)
					}
				}
			}
		}
		if fastMod {
			latestUpdateBlockNum = endBlockNum
			continue
		}
		var n NSum
		db.WithContext(ctx).Model(EvmBlock{}).Where("block_num BETWEEN ? AND ?", latestUpdateBlockNum+1, endBlockNum).
			Select("sum(transaction_count) as n").Scan(&n)
		var nt int64
		db.WithContext(ctx).Model(Transaction{}).Where("block_num BETWEEN ? AND ?", latestUpdateBlockNum+1, endBlockNum).Count(&nt)
		if n.N == nt {
			latestUpdateBlockNum = endBlockNum
			continue
		}

		// 如果 Transactions 数目不一致，逐一检查
		for _, block := range GetBlockByNums(ctx, latestUpdateBlockNum+1, endBlockNum) {

			var transactionsCount int64
			db.Model(Transaction{}).Where("block_num = ?", block.BlockNum).Count(&transactionsCount)

			if block.TransactionCount != int(transactionsCount) {
				util.Logger().Warning(fmt.Sprintf("Invalid block %d", block.BlockNum))
				blockRaw := fillBlock(int(block.BlockNum), true)
				// reSync transaction index
				if blockRaw != nil {
					for _, ethTransaction := range blockRaw.Transactions {
						db.Model(Transaction{}).Where("hash = ?", ethTransaction.Hash).Update("transaction_index", util.U256(ethTransaction.TransactionIndex).Uint64())
					}
				}
			}

		}
		latestUpdateBlockNum = endBlockNum
	}
	wg.Wait()
}
