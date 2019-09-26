package substrate

import (
	"encoding/json"
	"subscan-end/libs/substrate/scalecodec"
	"subscan-end/libs/substrate/scalecodec/types"
	"subscan-end/utiles"
)

func (m *MetadataType) DecodeEvent(event string) []map[string]interface{} {
	bm, _ := json.Marshal(m)
	e := scalecodec.EventsDecoder{}
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes(event)}, []string{"", string(bm)})
	return e.Process()
}
