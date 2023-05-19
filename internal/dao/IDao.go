package dao

import (
	"context"

	"github.com/itering/subscan/model"
	"github.com/itering/substrate-api-rpc/metadata"
)

type IDao interface {
	SetHeartBeatNow(context.Context, string) error
	DbBegin() *GormDB
	DbCommit(*GormDB)
	DbRollback(*GormDB)
	CreateBlock(*GormDB, *model.ChainBlock) (err error)
	UpdateEventAndExtrinsic(*GormDB, *model.ChainBlock, int, int, int, string, bool, bool) error
	SetBlockFinalized(*model.ChainBlock)
	SaveFillAlreadyBlockNum(context.Context, int) error
	SaveFillAlreadyFinalizedBlockNum(c context.Context, blockNum int) (err error)
	CreateEvent(txn *GormDB, event *model.ChainEvent) error
	CreateExtrinsic(c context.Context, txn *GormDB, extrinsic *model.ChainExtrinsic) error
	CreateLog(txn *GormDB, ce *model.ChainLog) error
	SetMetadata(c context.Context, metadata map[string]interface{}) (err error)
	IncrMetadata(c context.Context, filed string, incrNum int) (err error)
	CreateRuntimeVersion(name string, specVersion int) int64
	SetRuntimeData(specVersion int, modules string, rawData string) int64
	CreateRuntimeConstants(specVersion int, constants []model.RuntimeConstant) error

	IReadOnlyDao
}

type IReadOnlyDao interface {
	Ping(ctx context.Context) (err error)
	DaemonHealth(context.Context) map[string]bool
	GetNearBlock(int) *model.ChainBlock
	BlocksReverseByNum([]int) map[int]model.ChainBlock
	GetBlockByHash(context.Context, string) *model.ChainBlock
	GetBlockByNum(int) *model.ChainBlock
	GetFillBestBlockNum(c context.Context) (num int, err error)
	GetBlockNumArr(start, end int) []int
	GetFillFinalizedBlockNum(c context.Context) (num int, err error)
	GetBlockList(page, row int) []model.ChainBlock
	BlockAsJson(c context.Context, block *model.ChainBlock) *model.ChainBlockJson
	GetEventByBlockNum(blockNum int, where ...string) []model.ChainEventJson
	GetEventList(page, row int, order string, where ...string) ([]model.ChainEvent, int)
	GetEventsByIndex(extrinsicIndex string) []model.ChainEvent
	GetEventByIdx(index string) *model.ChainEvent
	GetRuntimeConstantLatest(moduleName string, constantName string) *model.RuntimeConstant
	GetMetadata(c context.Context) (ms map[string]string, err error)
	GetBestBlockNum(c context.Context) (uint64, error)
	GetFinalizedBlockNum(c context.Context) (uint64, error)
	GetExtrinsicsByBlockNum(blockNum int) []model.ChainExtrinsicJson
	GetExtrinsicList(c context.Context, page, row int, order string, queryWhere ...string) ([]model.ChainExtrinsic, int)
	GetExtrinsicsByHash(c context.Context, hash string) *model.ChainExtrinsic
	GetExtrinsicsDetailByHash(c context.Context, hash string) *model.ExtrinsicDetail
	GetExtrinsicsDetailByIndex(c context.Context, index string) *model.ExtrinsicDetail
	ExtrinsicsAsJson(e *model.ChainExtrinsic) *model.ChainExtrinsicJson
	GetLogsByIndex(index string) *model.ChainLogJson
	GetLogByBlockNum(blockNum int) []model.ChainLogJson
	RuntimeVersionList() []model.RuntimeVersion
	RuntimeVersionRaw(spec int) *metadata.RuntimeRaw
	RuntimeVersionRecent() *model.RuntimeVersion
	Close()
}
