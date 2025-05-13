package dao

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDao_CreateEvent(t *testing.T) {
	txn := testDao.DbBegin()
	err := testDao.CreateEvent(txn, &testEvent)
	txn.Commit()
	assert.NoError(t, err)
}

func TestDao_GetEventList(t *testing.T) {
	events, _ := testDao.GetEventList(context.TODO(), 0, 100, "desc")
	assert.GreaterOrEqual(t, 2, len(events))
}

func TestDao_GetEventsByIndex(t *testing.T) {
	events := testDao.GetEventsByIndex("947687-0")
	assert.Equal(t, uint(947687), events[0].BlockNum)
	assert.Equal(t, uint(0), events[0].EventIdx)
}

func TestDao_GetEventByIdx(t *testing.T) {
	event := testDao.GetEventByIdx(context.TODO(), "947687-0")
	assert.Equal(t, uint(947687), event.BlockNum)
	assert.Equal(t, uint(0), event.EventIdx)
}
