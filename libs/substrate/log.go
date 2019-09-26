package substrate

import (
	"encoding/json"
	"subscan-end/libs/substrate/protos/codec_protos"
	"subscan-end/libs/substrate/storage"
	"subscan-end/utiles"
)

const (
	CidAura = 0x61727561
	CidBabe = 0x45424142
)

type PreRuntime struct {
	Data   string `json:"data"`
	Engine int64  `json:"engine"`
}

func (p *PreRuntime) isAura() bool {
	return CidAura == p.Engine
}
func (p *PreRuntime) isBabe() bool {
	return CidBabe == p.Engine
}

func (p *PreRuntime) getAuraAuthor(sessionValidators []string) string {
	if decodeMsg, err := codec_protos.DecodeStorage(utiles.AddHex(p.Data), "RawAuraPreDigest"); err == nil {
		rawAuraPreDigestValue := storage.StateStorage(decodeMsg)
		modn := rawAuraPreDigestValue.ToRawAuraPreDigest().SlotNumber % int64(len(sessionValidators))
		return sessionValidators[modn]
	}
	return ""
}

func (p *PreRuntime) getBabeAuthor(sessionValidators []string) string {
	if decodeMsg, err := codec_protos.DecodeStorage(utiles.AddHex(p.Data), "RawBabePreDigest"); err == nil {
		rawBabePreDigestValue := storage.StateStorage(decodeMsg)
		digest := rawBabePreDigestValue.ToRawBabePreDigest()
		if digest != nil {
			if digest.Primary != nil {
				return sessionValidators[digest.Primary.AuthorityIndex]
			} else {
				return sessionValidators[digest.Secondary.AuthorityIndex]
			}
		}
	}
	return ""
}

func ExtractAuthor(data []byte, sessionValidators []string) string {
	var p PreRuntime
	if len(sessionValidators) == 0 {
		return ""
	}
	if err := json.Unmarshal(data, &p); err != nil {
		return ""
	}
	if p.isAura() {
		return p.getAuraAuthor(sessionValidators)
	} else if p.isBabe() {
		return p.getBabeAuthor(sessionValidators)
	} else {
		return ""
	}
}
