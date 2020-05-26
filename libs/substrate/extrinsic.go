package substrate

import (
	"fmt"
	scalecodec "github.com/freehere107/go-scale-codec"
	"github.com/freehere107/go-scale-codec/types"
	"github.com/itering/subscan/libs/substrate/metadata"
	"github.com/itering/subscan/util"
)

// Extrinsic decode
func DecodeExtrinsic(rawList []string, metadata *metadata.MetadataType, spec int) (r []interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Recovering from panic in DecodeExtrinsic error is: %v \n", r)
		}
	}()
	m := types.MetadataStruct(*metadata)
	for _, extrinsicRaw := range rawList {
		e := scalecodec.ExtrinsicDecoder{}
		option := types.ScaleDecoderOption{Metadata: &m, Spec: spec}
		e.Init(types.ScaleBytes{Data: util.HexToBytes(extrinsicRaw)}, &option)
		e.Process()
		r = append(r, e.Value)
	}
	return
}
