package scalecodec

import (
	"errors"
	"subscan-end/libs/substrate/scalecodec/types"
)

type MetadataDecoder struct {
	types.ScaleDecoder
	Version  string               `json:"version"`
	Metadata types.MetadataStruct `json:"metadata"`
}

func (m *MetadataDecoder) Init(data []byte) {
	sData := types.ScaleBytes{Data: data}
	m.ScaleDecoder.Init(sData, "")
}

func (m *MetadataDecoder) Process() error {
	magicBytes := m.GetNextBytes(4)
	if string(magicBytes) == "meta" {
		m.Version = m.ProcessAndUpdateData("Enum", "MetadataV0Decoder", "MetadataV1Decoder", "MetadataV2Decoder", "MetadataV3Decoder", "MetadataV4Decoder", "MetadataV5Decoder", "MetadataV6Decoder", "MetadataV7Decoder").(string)
		m.Metadata = m.ProcessAndUpdateData(m.Version).(types.MetadataStruct)
		return nil
	} else {
		return errors.New("not metadata")
	}
}
