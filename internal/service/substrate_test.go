package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubscribeParserMessage(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	subscribeSrv := testSrv.initSubscribeService(done)

	err := subscribeSrv.parser([]byte(``))
	assert.Error(t, err)

	testSrv.dao.(*MockDao).On("GetBestBlockNum", context.TODO()).Return(uint64(1245201), nil)
	testSrv.dao.(*MockDao).On("GetFillBestBlockNum", context.TODO()).Return(1245200, nil)
	testSrv.dao.(*MockDao).On("GetFinalizedBlockNum", context.TODO()).Return(uint64(1245201), nil)
	testSrv.dao.(*MockDao).On("GetFillFinalizedBlockNum", context.TODO()).Return(1245200, nil)

	_ = subscribeSrv.parser([]byte(`{"jsonrpc":"2.0","result":{"apis":[["0xdf6acb689907609b",3],["0x37e397fc7c91f5e4",1],["0x40fe3ad401f8959a",4],["0xd2bc9897eed08f15",2],["0xf78b278be53f454c",2],["0xed99c5acb25eedf5",2],["0xcbca25e39f142387",2],["0x687ad44ad37f03c2",1],["0xab3c0572291feb8b",1],["0xbc9d89904f5b923f",1],["0x37c8bb1350a9a2a8",1],["0x199af487d84d9847",1],["0x18ef58a3b67ba770",1]],"authoringVersion":0,"implName":"Crab","implVersion":0,"specName":"Crab","specVersion":4,"transactionVersion":1},"id":1}`))

	_ = subscribeSrv.parser([]byte(`{"jsonrpc":"2.0","method":"chain_newHead","params":{"result":{"digest":{"logs":["0x0642414245b501011d000000d360dc0f000000009a42cc30d1aa157dddb14bb360f0e18a3c5c72d6c7da9d53c753e53d8e65737e33f0953ea4c95060cb13913903341e4179a221260451daf381d7994d3fbba907f0dfac992a65f24519280e3f338ddaa7fdc026098dfa3c5a5a590b22339fe20e","0x00904d4d5252f3631c7802415955cea4b46e49fd14863a7f5577c85c1e37ccb807bf32397cd5","0x0542414245010144f0ebd1dd8acc99a57a31ef5f2c4cb9ac1705b2e1ee5be32a13e98997cc1937ee50ce79845d1a7ed11b0ad169ef822c1405043f786b405255225343e8244c8f"]},"extrinsicsRoot":"0xe60a136d4d711c2b17a2c3729a07ced8fc0b89db49a5227e5bb7748371c1f945","number":"0x130013","parentHash":"0x2f6e814b6915e7904eba800b0d765f24eff3220b02fb71e5897f38b6af74f4f1","stateRoot":"0x11116beee5d91fc1a665b2b5862ee777800ca10f2977ba4285adb2283824ea9f"},"subscription":19139}}`))

	_ = subscribeSrv.parser([]byte(`{"jsonrpc":"2.0","method":"chain_finalizedHead","params":{"result":{"digest":{"logs":["0x064241424534021c0000000a61dc0f00000000","0x00904d4d5252c32aca3b03b483429056df60abca5bb2142f34124c6d03c664c191e7178fffbf","0x054241424501013a92db271ab64cc1827ca7c76ae324c8dd978679758935b94503c0765b5d8a5fe3b2d266b78b8b8879b68200a31fb7e56b9da323069d4a76a7c52715faf54d86"]},"extrinsicsRoot":"0x10a8648c10909b1a68a80679a27f650e82b3962fbb42f021883f359bbce0ad55","number":"0x13004a","parentHash":"0x5db569d34f8018f4b49adfee7b719e19f5c192a50b4c1e680f1fd550209c988d","stateRoot":"0x798f98a70681cd1d70b39b806dcb0db6aa68ed6395d3051cd9ed1eb7c1efb6ff"},"subscription":19171}}`))

	_ = subscribeSrv.parser([]byte(`{"jsonrpc":"2.0","method":"state_storage","params":{"result":{"block":"0xcee4c91b637487d951ef4704ffe6b36de5bb2a54fe39016dafae5f118d5b8752","changes":[["0x481e203dcea218263e3a96ca9e4b193857c875e4cff74148e4628f264b974c80","0xf48667ede356681b0000000000000000"]]},"subscription":19447}}`))
}

func TestService_FillBlockData(t *testing.T) {
	testSrv.dao.(*MockDao).On("GetBlockByNum", 1245201).Return(nil, nil)
	tc := TestConn{}
	err := testSrv.FillBlockData(&tc, 1245201, true)
	assert.NoError(t, err)
}

func TestService_subscribeFetchBlock(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	testSrv.dao.(*MockDao).On("GetFinalizedBlockNum", context.TODO()).Return(uint64(1245201), nil)
	testSrv.dao.(*MockDao).On("GetFillFinalizedBlockNum", context.TODO()).Return(1245200, nil)
	sub := testSrv.initSubscribeService(done)
	go sub.subscribeFetchBlock()
	sub.newFinHead <- true
}
