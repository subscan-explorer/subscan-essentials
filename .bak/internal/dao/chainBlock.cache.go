package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bilibili/kratos/pkg/cache/redis"
	"sort"
	"subscan-end/internal/model"
	"sync"
)

func (d *Dao) Block(c context.Context, blockNum int) (b *model.ChainBlock) {
	addCache := true
	res, err := d.cacheBlock(c, blockNum)
	if err != nil {
		addCache = false
	}
	if res != nil {
		_ = json.Unmarshal(res, &b)
		return
	}
	b = d.GetBlockByNum(c, blockNum)
	if b == nil || !addCache {
		return
	}
	_ = d.cache.Do(c, func(ctx context.Context) {
		d.addCacheBlock(c, b)
	})
	return
}

func (d *Dao) cacheBlock(c context.Context, blockNum int) ([]byte, error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	cacheKey := fmt.Sprintf(blockCacheKey, blockNum)
	block, err := redis.Bytes(conn.Do("GET", cacheKey))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return nil, err
	}
	return block, nil
}

func (d *Dao) addCacheBlock(c context.Context, b *model.ChainBlock) {
	conn := d.redis.Get(c)
	defer conn.Close()
	cacheKey := fmt.Sprintf(blockCacheKey, b.BlockNum)
	bytes, _ := json.Marshal(b)
	conn.Do("SETEX", cacheKey, 86400*7, string(bytes))
}

func (d *Dao) BlockByHash(c context.Context, hash string) (b *model.ChainBlock) {
	addCache := true
	res, err := d.cacheBlockByHash(c, hash)
	if err != nil {
		addCache = false
	}
	if res != nil {
		_ = json.Unmarshal(res, &b)
		return
	}
	b = d.GetBlockByHash(c, hash)
	if b == nil || !addCache {
		return
	}
	_ = d.cache.Do(c, func(ctx context.Context) {
		d.addCacheBlockByHash(c, b)
	})
	return
}

func (d *Dao) cacheBlockByHash(c context.Context, hash string) ([]byte, error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	cacheKey := fmt.Sprintf(blockByHashCacheKey, hash)
	block, err := redis.Bytes(conn.Do("GET", cacheKey))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return nil, err
	}
	return block, nil
}

func (d *Dao) addCacheBlockByHash(c context.Context, b *model.ChainBlock) {
	conn := d.redis.Get(c)
	defer conn.Close()
	cacheKey := fmt.Sprintf(blockByHashCacheKey, b.Hash)
	bytes, _ := json.Marshal(b)
	conn.Do("SETEX", cacheKey, 86400*7, string(bytes))
}

func (d *Dao) BlocksJson(c context.Context, nums []int) (res *[]model.SampleBlockJson) {
	if len(nums) == 0 {
		return
	}
	res = d.cacheBlocks(c, nums)
	return
}

func (d *Dao) cacheBlocks(c context.Context, nums []int) *[]model.SampleBlockJson {
	wg := sync.WaitGroup{}
	var blocks []model.SampleBlockJson
	for _, blockNum := range nums {
		wg.Add(1)
		go func(blockNum int) {
			defer wg.Done()
			if b := d.Block(c, blockNum); b != nil {
				sample := d.BlockAsSampleJson(c, b)
				blocks = append(blocks, *sample)
			}
		}(blockNum)
	}
	wg.Wait()
	sort.Slice(blocks[:], func(i, j int) bool {
		return blocks[i].BlockNum > blocks[j].BlockNum
	})
	return &blocks
}
