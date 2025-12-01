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

	CreateBlock(context.Context, *GormDB, *model.ChainBlock) (err error)
	UpdateEventAndExtrinsic(*GormDB, *model.ChainBlock, int, int, int, string, bool, bool) error
	GetNearBlock(uint) *model.ChainBlock
	SplitBlockTable(blockNum uint)
	BlocksReverseByNum([]uint) map[uint]model.ChainBlock
	GetBlockByHash(context.Context, string) *model.ChainBlock
	GetBlockByNum(context.Context, uint) *model.ChainBlock
	SaveFillAlreadyBlockNum(context.Context, int) error
	SaveFillAlreadyFinalizedBlockNum(c context.Context, blockNum int) (err error)
	GetFillBestBlockNum(c context.Context) (num int, err error)
	GetBlockNumArr(ctx context.Context, start, end uint) []int
	GetFillFinalizedBlockNum(c context.Context) (num int, err error)

	GetBlockListCursor(ctx context.Context, limit int, before, after uint) (list []model.ChainBlock, hasPrev, hasNext bool)
	BlockAsJson(c context.Context, block *model.ChainBlock) *model.ChainBlockJson

	CreateEvent(txn *GormDB, event []model.ChainEvent) error
	GetEventListCursor(ctx context.Context, limit int, order string, fixedTableIndex int, beforeId uint, afterId uint, where ...model.Option) (list []model.ChainEvent, hasPrev, hasNext bool)
	GetEventsByIndex(extrinsicIndex string) []model.ChainEvent
	GetEventByIdx(ctx context.Context, index string) *model.ChainEvent

	CreateExtrinsic(c context.Context, txn *GormDB, extrinsic []model.ChainExtrinsic, u int) error
	GetExtrinsicListCursor(c context.Context, limit int, fixedTableIndex int, beforeId, afterId uint, accountId string, queryWhere ...model.Option) (list []model.ChainExtrinsic, hasPrev, hasNext bool)
	GetExtrinsicsByHash(c context.Context, hash string) *model.ChainExtrinsic
	GetExtrinsicsByIndex(c context.Context, index string) *model.ChainExtrinsic
	GetExtrinsicsDetailByHash(c context.Context, hash string) *model.ExtrinsicDetail
	GetExtrinsicsDetailByIndex(c context.Context, index string) *model.ExtrinsicDetail
	ExtrinsicsAsJson(e *model.ChainExtrinsic) *model.ChainExtrinsicJson
	GetExtrinsicCount(ctx context.Context, queryWhere ...model.Option) int64

	CreateLog(txn *GormDB, ce []model.ChainLog) error
	GetLogByBlockNum(ctx context.Context, blockNum uint) []model.ChainLogJson

	SetMetadata(c context.Context, metadata map[string]interface{}) (err error)
	IncrMetadata(c context.Context, filed string, incrNum int) (err error)
	GetMetadata(c context.Context) (ms map[string]string, err error)
	GetBestBlockNum(c context.Context) (uint64, error)
	GetFinalizedBlockNum(c context.Context) (uint64, error)

	CreateRuntimeVersion(c context.Context, name string, specVersion int, blockNum uint) bool
	SetRuntimeData(specVersion int, modules string, rawData string) int64
	RuntimeVersionList() []model.RuntimeVersion
	RuntimeVersionRaw(spec int) *metadata.RuntimeRaw
	RuntimeVersionRecent() *model.RuntimeVersion

	GetSessionValidatorsById(ctx context.Context, sessionId uint) []string
	CreateNewSession(ctx context.Context, sessionId uint, validators []string) error
}
