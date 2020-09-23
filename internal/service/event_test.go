package service

import (
	"github.com/itering/subscan/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_AddEvent(t *testing.T) {
	txn := testSrv.dao.DbBegin()
	defer testSrv.dao.DbRollback(txn)
	affect, err := testSrv.AddEvent(txn, &testBlock, []model.ChainEvent{testEvent}, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, affect)
}

func TestService_GetEventList(t *testing.T) {
	list, count := testSrv.RenderEvents(0, 1000, "desc")
	assert.Equal(t, 1, count)
	assert.Equal(t, []model.ChainEventJson{
		{EventIndex: "947687-0",
			BlockNum:       947687,
			ModuleId:       "imonline",
			EventId:        "AllGood",
			ExtrinsicHash:  "",
			Params:         "[]",
			BlockTimestamp: 1594791900,
		}}, list)
}
