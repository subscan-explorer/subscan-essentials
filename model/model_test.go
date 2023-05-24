package model_test

import (
	"testing"

	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins/storage"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestModelPluginRender(t *testing.T) {
	block := model.ChainBlock{BlockNum: 1, BlockTimestamp: 1, Hash: "0x0", SpecVersion: 1, Validator: "0x0", Finalized: true}
	assert.Equal(t, &storage.Block{BlockNum: 1, BlockTimestamp: 1, Hash: "0x0", SpecVersion: 1, Validator: "0x0", Finalized: true}, block.AsPlugin())

	event := model.ChainEvent{BlockNum: 1, EventIdx: 1, ModuleId: "b", ExtrinsicHash: "0x0", EventId: "0", Params: `{"a":"b"}`}
	assert.Equal(t, &storage.Event{BlockNum: 1, EventIdx: 1, ModuleId: "b", ExtrinsicHash: "0x0", EventId: "0", Params: `{"a":"b"}`}, event.AsPlugin())

	extrinsic := model.ChainExtrinsic{BlockNum: 1, BlockTimestamp: 1, ExtrinsicHash: "0x0", Params: `{"a":"b"}`, Fee: decimal.New(1, 0)}
	assert.Equal(t, &storage.Extrinsic{ExtrinsicHash: "0x0", Params: []byte(`{"a":"b"}`), Fee: decimal.New(1, 0)}, extrinsic.AsPlugin())
}
