package model_test

import (
	"github.com/itering/subscan/model"
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
	{instant: model.ChainExtrinsic{BlockNum: 10000000}, tableName: "chain_extrinsics_10"},
	{instant: model.ChainLog{BlockNum: 999999}, tableName: "chain_logs"},
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
