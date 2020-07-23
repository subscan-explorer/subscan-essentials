package dao

import (
	"github.com/itering/subscan/model"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDao_CreateRuntimeVersion(t *testing.T) {
	affect := testDao.CreateRuntimeVersion("polkadot", 1)
	assert.Equal(t, int64(0), affect)
}

func TestDao_SetRuntimeData(t *testing.T) {
	affect := testDao.SetRuntimeData(1, "system|staking|session", "0x0")
	assert.Equal(t, int64(1), affect)
	testDao.SetRuntimeData(1, "system|staking", "0x0")
}

func TestDao_RuntimeVersionList(t *testing.T) {
	list := testDao.RuntimeVersionList()
	assert.Equal(t, 1, len(list))
	assert.Equal(t, model.RuntimeVersion{
		SpecVersion: 1,
		Modules:     "system|staking",
	}, list[0])
}

func TestDao_RuntimeVersionRecent(t *testing.T) {
	assert.Equal(t, &model.RuntimeVersion{SpecVersion: 1, RawData: "0x0"}, testDao.RuntimeVersionRecent())
}

func TestDao_RuntimeVersionRaw(t *testing.T) {
	recent := testDao.RuntimeVersionRaw(2)
	assert.Nil(t, recent)
	assert.Equal(t, &metadata.RuntimeRaw{Spec: 1, Raw: "0x0"}, testDao.RuntimeVersionRaw(1))
}
