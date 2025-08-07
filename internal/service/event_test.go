package service

import (
	"context"
	"github.com/itering/subscan/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_AddEvent(t *testing.T) {
	txn := testSrv.dao.DbBegin()
	defer testSrv.dao.DbRollback(txn)
	err := testSrv.AddEvent(txn, &testBlock, []model.ChainEvent{testEvent})
	assert.NoError(t, err)
}

func TestService_GetEventList(t *testing.T) {
	list, count := testSrv.EventsList(context.TODO(), 0, 1000, -1, 0)
	assert.Equal(t, 1, count)
	assert.Equal(t, []model.ChainEventJson{
		{EventIndex: "947687-0",
			BlockNum:       947687,
			ModuleId:       "imonline",
			EventId:        "AllGood",
			Params:         model.EventParams{},
			BlockTimestamp: 1594791900,
			ExtrinsicIndex: "947687-0",
		}}, list)
}
