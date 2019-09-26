package storage

import (
	"encoding/binary"
	"encoding/json"
	"math/big"
	"strconv"
	"subscan-end/utiles"
)

type StateStorage string

func (s *StateStorage) bytes() []byte {
	return []byte(string(*s))
}

func (s *StateStorage) string() string {
	return string(*s)
}

func (s *StateStorage) ToStringSlice() (r []string) {
	_ = json.Unmarshal(s.bytes(), &r)
	return
}

func (s *StateStorage) ToInt() (r int) {
	r = utiles.StringToInt(s.string())
	return
}

func (s *StateStorage) ToInt64() (r int64) {
	i, _ := strconv.ParseInt(s.string(), 10, 64)
	return i
}

func (s *StateStorage) ToStakingLedgers() (r *StakingLedgers) {
	_ = json.Unmarshal(s.bytes(), &r)
	return
}

func (s *StateStorage) ToMapString() (r map[string]string) {
	_ = json.Unmarshal(s.bytes(), &r)
	return
}

func (s *StateStorage) ToExposures() (r *Exposures) {
	_ = json.Unmarshal(s.bytes(), &r)
	return
}

func (s *StateStorage) ToRawAuraPreDigest() (r *RawAuraPreDigest) {
	_ = json.Unmarshal(s.bytes(), &r)
	return
}

func (s *StateStorage) ToRawBabePreDigest() (r *RawBabePreDigest) {
	_ = json.Unmarshal(s.bytes(), &r)
	return
}

func (s *StateStorage) ToValidatorPrefsLegacy() (r *ValidatorPrefsLegacy) {
	_ = json.Unmarshal(s.bytes(), &r)
	return
}

func (s *StateStorage) ToU32FromCodec() (r uint32) {
	return uint32(binary.LittleEndian.Uint32(utiles.HexToBytes(s.string())[0:4]))
}

func (s *StateStorage) ToBigInt() (r *big.Int) {
	if b, ok := new(big.Int).SetString(s.string(), 10); ok == true {
		return b
	}
	return big.NewInt(0)
}
