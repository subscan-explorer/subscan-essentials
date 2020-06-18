package substrate

import (
	"fmt"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/subscan/internal/substrate/metadata"
	"github.com/itering/subscan/internal/util"
	"github.com/shopspring/decimal"
	"math"
)

// Extrinsic decode
func DecodeExtrinsic(rawList []string, metadata *metadata.MetadataType, spec int) (r []map[string]interface{}, err error) {
	defer func() {
		if fatal := recover(); fatal != nil {
			err = fmt.Errorf("Recovering from panic in DecodeExtrinsic error is: %v \n", fatal)
		}
	}()
	m := types.MetadataStruct(*metadata)
	for _, extrinsicRaw := range rawList {
		e := scalecodec.ExtrinsicDecoder{}
		option := types.ScaleDecoderOption{Metadata: &m, Spec: spec}
		e.Init(types.ScaleBytes{Data: util.HexToBytes(extrinsicRaw)}, &option)
		e.Process()
		r = append(r, e.Value.(map[string]interface{}))
	}
	return
}

type Mortal struct {
	Period uint64
	Phase  uint64
}

func DecodeMortal(era string) *Mortal {
	if era == "" || era == "00" {
		return nil
	}
	eraU8a := util.HexToBytes(era)
	first := uint64(eraU8a[0])
	second := uint64(eraU8a[1])
	encoded := first + (second << 8)
	var period uint64 = 2 << (encoded % (1 << 4))
	quantizeFactor := math.Max(float64(period>>12), 1)
	phase := (encoded >> 4) * uint64(quantizeFactor)
	fmt.Println(period, phase)
	return &Mortal{
		Period: period,
		Phase:  phase,
	}
}

func (m *Mortal) Birth(current uint64) uint64 {
	s := (decimal.Max(decimal.New(int64(current), 0), decimal.New(int64(m.Phase), 0)).
		Sub(decimal.New(int64(m.Phase), 0))).Div(decimal.New(int64(m.Period), 0)).IntPart()
	return uint64(s)*m.Period + m.Phase
}

func (m *Mortal) Death(current uint64) uint64 {
	return m.Birth(current) + m.Period
}
