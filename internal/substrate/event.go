package substrate

import (
	"fmt"
	scalecodec "github.com/freehere107/go-scale-codec"
	"github.com/freehere107/go-scale-codec/types"
	"github.com/itering/subscan/internal/substrate/metadata"
	"github.com/itering/subscan/util"
)

// Event decode
func DecodeEvent(rawList string, metadata *metadata.MetadataType, spec int) (r interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Recovering from panic in DecodeEvent error is: %v \n", r)
		}
	}()
	m := types.MetadataStruct(*metadata)
	e := scalecodec.EventsDecoder{}
	option := types.ScaleDecoderOption{Metadata: &m, Spec: spec}
	e.Init(types.ScaleBytes{Data: util.HexToBytes(rawList)}, &option)
	e.Process()
	return e.Value, nil
}
