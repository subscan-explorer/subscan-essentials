package daemons

import (
	"encoding/json"
	"fmt"
	"subscan-end/internal/service"
	"subscan-end/libs/substrate"
	"subscan-end/utiles"
)

var fetchBlockFlag = false

func parserDistribution(message []byte, srv *service.Service) {
	var j substrate.JsonRpcResult
	if err := json.Unmarshal(message, &j); err != nil {
		return
	}
	if j.Id == 1 { //runtime version
		r := j.ToRuntimeVersion()
		_ = srv.CreateRuntimeVersion(r.ImplName, r.SpecVersion)
		srv.UpdateChainMetadata(map[string]interface{}{"implName": r.ImplName, "specVersion": r.SpecVersion})
		substrate.CurrentRuntimeSpecVersion = r.SpecVersion
	}

	switch j.Method {
	case substrate.ChainNewHead:
		r := j.ToNewHead()
		metadata := map[string]interface{}{"blockNum": utiles.HexToNumStr(r.Number)}
		srv.UpdateChainMetadata(metadata)
		newHead <- true
		if fetchBlockFlag == false {
			fetchBlockFlag = true
			go fillALLBlock(srv) // TODO: Fetch block will block all this function, need fix this
		}
	case substrate.StateStorage:
		r := j.ToStorage()
		fmt.Println("StateStorage", r)
	default:
		return
	}

}
