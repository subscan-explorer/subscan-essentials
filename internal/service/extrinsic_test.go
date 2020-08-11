package service

import (
	"context"
	"github.com/itering/substrate-api-rpc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_createExtrinsic(t *testing.T) {
	txn := testSrv.dao.DbBegin()
	defer testSrv.dao.DbRollback(txn)
	ctx := context.TODO()
	encodeExtrinsic := []string{"0x280402000b10449a7e7301", "0x1c0407005e8b4100", "0xa10184623b0263bb111bbd81bd32bc893f82132d32b0a83a236da15cc40b8a893cc7160100fc86ad2eb5a10087838887fb31f714455a963ecea413cd5aa4e06150785f0d06d67644e8b95b214287f94fff4c54a3d06cf7bccfc86a1a6de919b661664f8d05031c001608"}
	metadataInstant := testSrv.getMetadataInstant(4, "")
	decodeExtrinsics, err := substrate.DecodeExtrinsic(encodeExtrinsic, metadataInstant, 4)
	assert.NoError(t, err)

	count, blockTime, txMap, feeMap, err := testSrv.createExtrinsic(ctx, txn, &testBlock, encodeExtrinsic, decodeExtrinsics, nil)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
	assert.Equal(t, 1595556906, blockTime)
	assert.Equal(t, 1, len(feeMap))
	assert.Equal(t, map[string]string{"947687-2": "5ee3f70a70d033d30755eeb1afe80520433cdc1d99348a522f2302d785ac907a"}, txMap)
}
