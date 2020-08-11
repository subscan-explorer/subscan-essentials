package model_test

import (
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCases = []struct {
	instant   interface{}
	tableName string
}{
	{instant: model.ChainBlock{}, tableName: "chain_blocks"},
	{instant: model.ChainBlock{BlockNum: 1000000}, tableName: "chain_blocks_1"},
	{instant: model.ChainEvent{BlockNum: 1000000}, tableName: "chain_events_1"},
	{instant: model.ChainEvent{BlockNum: 99}, tableName: "chain_events"},
	{instant: model.ChainExtrinsic{BlockNum: 10000000}, tableName: "chain_extrinsics_10"},
	{instant: model.ChainExtrinsic{BlockNum: 10000}, tableName: "chain_extrinsics"},
	{instant: model.ChainLog{BlockNum: 999999}, tableName: "chain_logs"},
	{instant: model.ChainLog{BlockNum: 1999999}, tableName: "chain_logs_1"},
}

func TestSplitTableName(t *testing.T) {
	for _, test := range testCases {
		switch v := test.instant.(type) {
		case model.ChainBlock:
			assert.Equal(t, v.TableName(), test.tableName)
		case model.ChainEvent:
			assert.Equal(t, v.TableName(), test.tableName)
		case model.ChainExtrinsic:
			assert.Equal(t, v.TableName(), test.tableName)
		case model.ChainLog:
			assert.Equal(t, v.TableName(), test.tableName)
		}
	}
}

func TestModelPluginRender(t *testing.T) {
	block := model.ChainBlock{BlockNum: 1, BlockTimestamp: 1, Hash: "0x0", SpecVersion: 1, Validator: "0x0", Finalized: true}
	assert.Equal(t, &storage.Block{BlockNum: 1, BlockTimestamp: 1, Hash: "0x0", SpecVersion: 1, Validator: "0x0", Finalized: true}, block.AsPlugin())

	event := model.ChainEvent{BlockNum: 1, EventIdx: 1, ModuleId: "b", ExtrinsicHash: "0x0", EventId: "0", Params: `{"a":"b"}`}
	assert.Equal(t, &storage.Event{BlockNum: 1, EventIdx: 1, ModuleId: "b", ExtrinsicHash: "0x0", EventId: "0", Params: []byte(`{"a":"b"}`)}, event.AsPlugin())

	extrinsic := model.ChainExtrinsic{BlockNum: 1, BlockTimestamp: 1, ExtrinsicHash: "0x0", Params: `{"a":"b"}`, Fee: decimal.New(1, 0)}
	assert.Equal(t, &storage.Extrinsic{ExtrinsicHash: "0x0", Params: []byte(`{"a":"b"}`), Fee: decimal.New(1, 0)}, extrinsic.AsPlugin())

}
