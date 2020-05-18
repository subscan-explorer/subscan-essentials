package substrate

import (
	"encoding/json"
	"subscan-end/libs/substrate/scalecodec"
	"subscan-end/libs/substrate/scalecodec/types"
	"subscan-end/utiles"
)

func (m *MetadataType) DecodeExtrinsic(extrinsic string) []map[string]interface{} {
	var extrinsics []string
	bm, _ := json.Marshal(m)
	_ = json.Unmarshal([]byte(extrinsic), &extrinsics)
	e := scalecodec.ExtrinsicsDecoder{}
	var result []map[string]interface{}
	for _, value := range extrinsics {
		e.Init(types.ScaleBytes{Data: utiles.HexToBytes(value)}, []string{"", string(bm)})
		result = append(result, e.Process())
	}
	return result
}
