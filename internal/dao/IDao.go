package dao

import (
	"context"

	"github.com/itering/subscan/model"
	"github.com/itering/substrate-api-rpc/metadata"
)

type IDao interface {
	Close()
	Ping(ctx context.Context) (err error)
	DbBegin() *GormDB
	DbCommit(*GormDB)
	DbRollback(*GormDB)

	CreateBlock(*GormDB, *model.ChainBlock) (err error)
	UpdateEventAndExtrinsic(*GormDB, *model.ChainBlock, int, int, int, string, bool, bool) error
	GetNearBlock(uint) *model.ChainBlock
	BlocksReverseByNum([]uint) map[uint]model.ChainBlock
	GetBlockByHash(context.Context, string) *model.ChainBlock
	GetBlockByNum(context.Context, uint) *model.ChainBlock
	SaveFillAlreadyBlockNum(context.Context, int) error
	SaveFillAlreadyFinalizedBlockNum(c context.Context, blockNum int) (err error)
	GetFillBestBlockNum(c context.Context) (num int, err error)
	GetBlockNumArr(ctx context.Context, start, end uint) []int
	GetFillFinalizedBlockNum(c context.Context) (num int, err error)
	GetBlockList(ctx context.Context, page, row int) []model.ChainBlock
	BlockAsJson(c context.Context, block *model.ChainBlock) *model.ChainBlockJson

	CreateEvent(txn *GormDB, event []model.ChainEvent) error
	GetEventList(ctx context.Context, page, row int, order string, where ...model.Option) ([]model.ChainEvent, int)
	GetEventsByIndex(extrinsicIndex string) []model.ChainEvent
	GetEventByIdx(ctx context.Context, index string) *model.ChainEvent

	CreateExtrinsic(c context.Context, txn *GormDB, extrinsic []model.ChainExtrinsic, u int) error
	GetExtrinsicList(c context.Context, page, row int, order string, queryWhere ...model.Option) ([]model.ChainExtrinsic, int)
	GetExtrinsicsByHash(c context.Context, hash string) *model.ChainExtrinsic
	GetExtrinsicsByIndex(c context.Context, index string) *model.ChainExtrinsic
	GetExtrinsicsDetailByHash(c context.Context, hash string) *model.ExtrinsicDetail
	GetExtrinsicsDetailByIndex(c context.Context, index string) *model.ExtrinsicDetail
	ExtrinsicsAsJson(e *model.ChainExtrinsic) *model.ChainExtrinsicJson

	CreateLog(txn *GormDB, ce *model.ChainLog) error
	GetLogByBlockNum(ctx context.Context, blockNum uint) []model.ChainLogJson

	SetMetadata(c context.Context, metadata map[string]interface{}) (err error)
	IncrMetadata(c context.Context, filed string, incrNum int) (err error)
	GetMetadata(c context.Context) (ms map[string]string, err error)
	GetBestBlockNum(c context.Context) (uint64, error)
	GetFinalizedBlockNum(c context.Context) (uint64, error)

	CreateRuntimeVersion(name string, specVersion int) int64
	SetRuntimeData(specVersion int, modules string, rawData string) int64
	RuntimeVersionList() []model.RuntimeVersion
	RuntimeVersionRaw(spec int) *metadata.RuntimeRaw
	RuntimeVersionRecent() *model.RuntimeVersion
}
