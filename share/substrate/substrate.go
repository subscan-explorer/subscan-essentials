package substrate

import (
	"fmt"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/metadata"
)

func DecodeExtrinsicParams(raw string, metadata *metadata.Instant, call *types.MetadataCalls, spec int) (params []scalecodec.ExtrinsicParam, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Recovering from panic in DecodeEventParams error is: %v \n", r)
		}
	}()
	e := types.ScaleDecoder{}
	m := types.MetadataStruct(*metadata)
	option := types.ScaleDecoderOption{Metadata: &m, Spec: spec}
	e.Init(scaleBytes.ScaleBytes{Data: util.HexToBytes(raw)}, &option)
	for _, arg := range call.Args {
		value := e.ProcessAndUpdateData(arg.Type)
		param := scalecodec.ExtrinsicParam{Type: arg.Type, Value: value, Name: arg.Name, TypeName: arg.TypeName}
		params = append(params, param)
	}
	return params, err
}

// DecodeEventParams decode event params
func DecodeEventParams(raw string, argsType []string, metadata *metadata.Instant, event *types.MetadataEvents, spec int) (params []scalecodec.EventParam, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Recovering from panic in DecodeEventParams error is: %v \n", r)
		}
	}()
	e := types.ScaleDecoder{}
	m := types.MetadataStruct(*metadata)
	option := types.ScaleDecoderOption{Metadata: &m, Spec: spec}
	e.Init(scaleBytes.ScaleBytes{Data: util.HexToBytes(raw)}, &option)
	for index, argType := range argsType {
		value := e.ProcessAndUpdateData(argType)
		param := scalecodec.EventParam{Type: argType, Value: value}
		if len(event.ArgsTypeName) == len(event.Args) {
			param.TypeName = event.ArgsTypeName[index]
		}
		if len(event.ArgsName) == len(event.Args) {
			param.Name = event.ArgsName[index]
		}
		params = append(params, param)
	}
	return params, err
}
