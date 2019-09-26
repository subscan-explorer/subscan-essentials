package types

import (
	"subscan-end/utiles"
)

type ScaleBytes struct {
	Data   []byte `json:"data"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

func (s *ScaleBytes) GetNextBytes(length int) []byte {
	data := s.Data[s.Offset : s.Offset+length]
	s.Offset = s.Offset + length
	return data
}

func (s *ScaleBytes) GetRemainingLength() int {
	return s.Length - s.Offset
}

func (s *ScaleBytes) String() string {
	return utiles.AddHex(utiles.BytesToHex(s.Data))
}

func (s *ScaleBytes) Reset() {
	s.Offset = 0
}
