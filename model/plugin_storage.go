package model

import (
	"encoding/json"
	"github.com/itering/subscan-plugin/storage"
)

func (c *ChainBlock) AsPluginBlock() *storage.Block {
	return &storage.Block{
		BlockNum:       c.BlockNum,
		BlockTimestamp: c.BlockTimestamp,
		Hash:           c.Hash,
		SpecVersion:    c.SpecVersion,
		Validator:      c.Validator,
		Finalized:      c.Finalized,
	}
}

func (c *ChainExtrinsic) AsPluginExtrinsic() *storage.Extrinsic {
	paramBytes, _ := json.Marshal(c.Params)
	return &storage.Extrinsic{
		ExtrinsicIndex:     c.ExtrinsicIndex,
		CallModule:         c.CallModule,
		CallModuleFunction: c.CallModuleFunction,
		Params:             paramBytes,
		AccountId:          c.AccountId,
		Signature:          c.Signature,
		Nonce:              c.Nonce,
		Era:                c.Era,
		ExtrinsicHash:      c.ExtrinsicHash,
		Success:            c.Success,
		Fee:                c.Fee,
	}

}

func (c *ChainEvent) AsPluginEvent() *storage.Event {
	paramBytes, _ := json.Marshal(c.Params)
	return &storage.Event{
		BlockNum:      c.BlockNum,
		ExtrinsicIdx:  c.ExtrinsicIdx,
		ModuleId:      c.ModuleId,
		EventId:       c.EventId,
		Params:        paramBytes,
		ExtrinsicHash: c.ExtrinsicHash,
		EventIdx:      c.EventIdx,
	}
}
