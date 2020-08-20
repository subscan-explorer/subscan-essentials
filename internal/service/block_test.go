package service

import (
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_CreateChainBlock(t *testing.T) {
	hash := "0x2895f79f46105d24813c97558f06edaecf993fad2358dcaf86b813859edd697d"
	event := "0x080000000000000080e36a0900000000020000000100000000000000000000000000020000"
	block := rpc.Block{
		Extrinsics: []string{"0x280402000b10449a7e7301", "0x1c0407005e8b4100"},
		Header: rpc.ChainNewHeadResult{
			ExtrinsicsRoot: "0x5a9403235d77280ad129b44eebdbdea3127ada25dcbf540e6ce38a1f770ad86f",
			Number:         "0x1062da",
			ParentHash:     "0x42838e9a502c5ba1faa5de1f19bfe5464b3374df5c6f53152d0543363a900bb7",
			StateRoot:      "0x8c9eec7854ab7cc1207a3eb9d54e272a07545b2c77dd2b99f3b85a642cd91a49",
			Digest: rpc.ChainNewHeadLog{
				Logs: []string{
					"0x064241424534021b00000007b6d90f00000000",
					"0x00904d4d5252708a1db71fe9eedba2439dab7209846247482fe486dc5e82ef76897f2e50a3a5",
					"0x05424142450101aed0a28294d357326d3b199cd06f23a43cd44412cc9450286252bb75c47fad17a4e9a7e19fde1cb90d668b2bbc769faea28f63bffc53a28c8f5a9d817bf07b83",
				},
			},
		},
	}
	err := testSrv.CreateChainBlock(nil, hash, &block, event, 4, true)
	assert.NoError(t, err)

}

func TestService_UpdateBlockData(t *testing.T) {
	block := model.ChainBlock{
		BlockNum:       1073882,
		BlockTimestamp: 1595556906,
		Hash:           "0x2895f79f46105d24813c97558f06edaecf993fad2358dcaf86b813859edd697d",
		ParentHash:     "0x42838e9a502c5ba1faa5de1f19bfe5464b3374df5c6f53152d0543363a900bb7",
		StateRoot:      "0x8c9eec7854ab7cc1207a3eb9d54e272a07545b2c77dd2b99f3b85a642cd91a49",
		ExtrinsicsRoot: "0x5a9403235d77280ad129b44eebdbdea3127ada25dcbf540e6ce38a1f770ad86f",
		Logs:           `["0x064241424534021b00000007b6d90f00000000","0x00904d4d5252708a1db71fe9eedba2439dab7209846247482fe486dc5e82ef76897f2e50a3a5","0x05424142450101aed0a28294d357326d3b199cd06f23a43cd44412cc9450286252bb75c47fad17a4e9a7e19fde1cb90d668b2bbc769faea28f63bffc53a28c8f5a9d817bf07b83"]`,
		Extrinsics:     `["0x280402000b10449a7e7301","0x1c0407005e8b4100"]`,
		Event:          "0x080000000000000080e36a0900000000020000000100000000000000000000000000020000",
	}
	err := testSrv.UpdateBlockData(nil, &block, true)
	assert.NoError(t, err)
}

func TestService_checkoutExtrinsicEvents(t *testing.T) {
	event1 := model.ChainEvent{BlockNum: 107388, ExtrinsicIdx: 1}
	event2 := model.ChainEvent{BlockNum: 107388, ExtrinsicIdx: 2}
	events := testSrv.checkoutExtrinsicEvents([]model.ChainEvent{event1, event2}, 107388)
	assert.Equal(t, map[string][]model.ChainEvent{"107388-1": {event1}, "107388-2": {event2}}, events)
}

func TestService_GetCurrentRuntimeSpecVersion(t *testing.T) {
	assert.Equal(t, 4, testSrv.GetCurrentRuntimeSpecVersion(107388))
}

func TestService_GetExtrinsicList(t *testing.T) {
	_, count := testSrv.GetExtrinsicList(0, 10, "desc")
	assert.Equal(t, 1, count)
}

func TestService_GetBlocksSampleByNums(t *testing.T) {
	util.AddressType = "42"
	blocks := testSrv.GetBlocksSampleByNums(0, 10)
	assert.Equal(t, []model.SampleBlockJson{{
		BlockNum:       947687,
		BlockTimestamp: 1594791900,
		Hash:           "0xd68b38c412404a4b5d4974e6dbb4a491ed7b6200d4edc24152693804441ce99d",
		Validator:      "5EFjtKj1r8kEvDSeDAMaAJxciDt2G33n7hPHjDDRiYPZNWCD",
		Finalized:      true,
	}}, blocks)
}
