package service

import (
	"github.com/itering/substrate-api-rpc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_EmitLog(t *testing.T) {
	txn := testSrv.dao.DbBegin()
	defer testSrv.dao.DbRollback(txn)

	blockHash := "0x5f3d69b39b3634913965bdedea920f9156eaed992851f30fd6fc19d798ac764b"
	logs, err := substrate.DecodeLogDigest([]string{
		"0x0642414245b5010102000000efa6cd0f000000004618a29aeb02e8ae7bb2360d8f5f13828c3c2f9fd15bc674be6e2c64be17a00ebb8fa2449c7b19b5988d6110e0f03a44693f246597e7bdf1a4b48aa4c50b600e6252c08951731c00e11a7f5a6b26d7c6bdf421145c575a03c23420bd76decd06",
		"0x00904d4d5252aec4a1a273aca92e65330af40d9b06447427454910e0e1b9fc9e2157b670a30f",
		"0x054241424501019e89556620e6f4ed93cf9a939349d6928b38e5688ad0abb7cd3b6f8d9c3016021ac1b30fbf4aec0de00d9a288b261da9e4ed4921f64ed6393309ddc230c9cf8d",
	})
	assert.NoError(t, err)
	validatorList := []string{"v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9", "v10"}
	validator, err := testSrv.EmitLog(txn, blockHash, 300000, logs, true, validatorList)
	assert.NoError(t, err)
	assert.Equal(t, "v3", validator)

}
