package model

import (
	"context"
	"fmt"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/metadata"
	"gorm.io/gorm"
	"sync"
	"time"
)

type RuntimeType struct {
	BlockNum uint
	Spec     int
}

type runtimeVersionCacheWithExpiry struct {
	runtimeVersions []RuntimeType
	sync.Mutex
	expiry int64
}

var (
	runtimeVersionCacheValue runtimeVersionCacheWithExpiry
	// runtime version cache, block number -> spec version
)

func reloadBlockRuntimeVersion(_ context.Context, tx *gorm.DB) {
	runtimeVersionCacheValue.Lock()
	defer runtimeVersionCacheValue.Unlock()
	var list []RuntimeVersion
	tx.Model(RuntimeVersion{}).Select("block_num,spec_version").Order("spec_version asc").Find(&list)
	var r []RuntimeType
	for _, item := range list {
		r = append(r, RuntimeType{BlockNum: item.BlockNum, Spec: item.SpecVersion})
	}
	runtimeVersionCacheValue.runtimeVersions = r
	runtimeVersionCacheValue.expiry = time.Now().Unix()
}

func GetBlockRuntimeVersion(_ context.Context, tx *gorm.DB, blockNum uint) *int {
	if len(runtimeVersionCacheValue.runtimeVersions) == 0 || time.Now().Unix()-runtimeVersionCacheValue.expiry > 60*10 {
		reloadBlockRuntimeVersion(context.Background(), tx)
	}
	var spec = -1
	for _, value := range runtimeVersionCacheValue.runtimeVersions {
		if value.BlockNum <= blockNum {
			spec = value.Spec
		} else {
			break
		}
	}
	if spec == -1 {
		util.Logger().Error(fmt.Errorf("runtime version not found for block number %d", blockNum))
		var block ChainBlock
		tx.Model(ChainBlock{}).Select("spec_version").Where("block_num = ?", blockNum).Find(&block)
		// expect panic if block not found
		return &block.SpecVersion
	}
	return &spec
}

func GetMetadataInstant(_ context.Context, tx *gorm.DB, spec *int) *metadata.Instant {
	if spec == nil || *spec <= -1 {
		return metadata.Latest(nil)
	}

	metadataInstant, ok := metadata.RuntimeMetadata[*spec]
	if !ok {
		var raw metadata.RuntimeRaw
		err := tx.Model(RuntimeVersion{}).
			Select("spec_version as spec ,raw_data as raw").
			Where("spec_version = ?", spec).
			Scan(&raw).Error
		if err != nil {
			util.Logger().Error(fmt.Errorf("runtime version raw data missing for spec %d", *spec))
			return nil
		}
		metadataInstant = metadata.Process(&raw)
	}
	return metadataInstant
}
