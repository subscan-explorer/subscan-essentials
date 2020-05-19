package scalecodec

import (
	"encoding/json"
	"golang.org/x/crypto/blake2b"
	. "subscan-end/libs/substrate/scalecodec/types"
	"subscan-end/utiles"
)

type ExtrinsicsDecoder struct {
	ScaleDecoder
	Metadata            MetadataCallAndEvent     `json:"metadata"`
	ExtrinsicLength     int                      `json:"extrinsic_length"`
	ExtrinsicHash       string                   `json:"extrinsic_hash"`
	VersionInfo         string                   `json:"version_info"`
	ContainsTransaction bool                     `json:"contains_transaction"`
	Address             map[string]string        `json:"address"`
	Signature           string                   `json:"signature"`
	Nonce               int                      `json:"nonce"`
	Era                 string                   `json:"era"`
	CallIndex           string                   `json:"call_index"`
	CallModule          MetadataModules          `json:"call_module"`
	Call                MetadataCalls            `json:"call"`
	Params              []map[string]interface{} `json:"params"`
}

func (e *ExtrinsicsDecoder) Init(data ScaleBytes, args []string) {
	e.TypeMapping = map[string]string{
		"extrinsic_length": "Compact<u32>",
		"version_info":     "u8",
		"address":          "Address",
		"signature":        "Signature",
		"nonce":            "Compact<u32>",
		"era":              "Era",
		"call_index":       "(u8,u8)",
	}
	var metadata MetadataCallAndEvent
	var subType string
	if len(args) > 0 {
		subType = args[0]
	}
	if len(args) > 1 {
		_ = json.Unmarshal([]byte(args[1]), &metadata)
	}
	e.Metadata = metadata
	e.ScaleDecoder.Init(data, subType)
}

func (e *ExtrinsicsDecoder) generateHash() string {
	if e.ContainsTransaction {
		var extrinsicData []byte
		if e.ExtrinsicLength > 0 {
			extrinsicData = e.Data.Data
		} else {
			extrinsicLengthType := CompactU32{}
			extrinsicLengthType.Encode(e.Data.Length)
			extrinsicData = append(extrinsicLengthType.Data.Data[:], e.Data.Data[:]...)
		}
		checksum, _ := blake2b.New(32, []byte{})
		checksum.Write(extrinsicData)
		h := checksum.Sum(nil)
		return utiles.BytesToHex(h)
	}
	return ""
}

func (e *ExtrinsicsDecoder) Process() map[string]interface{} {
	e.ExtrinsicLength = int(e.ProcessAndUpdateData("Compact<u32>").(int))
	if e.ExtrinsicLength != e.Data.GetRemainingLength() {
		e.ExtrinsicLength = 0
		e.Data.Reset()
	}
	e.VersionInfo = utiles.BytesToHex(e.GetNextBytes(1))
	e.ContainsTransaction = utiles.U256(e.VersionInfo).Int64() >= 80
	if e.ContainsTransaction {
		e.Address = e.ProcessAndUpdateData("Address").(map[string]string)
		e.Signature = e.ProcessAndUpdateData("Signature").(string)
		e.Nonce = int(e.ProcessAndUpdateData(e.TypeMapping["nonce"]).(int))
		e.Era = e.ProcessAndUpdateData("Era").(string)
		e.ExtrinsicHash = e.generateHash()
	}
	e.CallIndex = utiles.BytesToHex(e.GetNextBytes(2))
	if e.CallIndex != "" {
		if e.Metadata.CallIndex[e.CallIndex] != nil {
			callIndex := e.Metadata.CallIndex[e.CallIndex].(map[string]interface{})
			bc, _ := json.Marshal(callIndex["call"])
			var call MetadataCalls
			_ = json.Unmarshal(bc, &call)
			e.Call = call
			var CallModule MetadataModules
			bc, _ = json.Marshal(callIndex["module"])
			_ = json.Unmarshal(bc, &CallModule)
			e.CallModule = CallModule
		}
	}
	for _, arg := range e.Call.Args {
		argTypeObj := e.ProcessAndUpdateData(arg["type"].(string))
		e.Params = append(e.Params, map[string]interface{}{
			"name":     arg["name"].(string),
			"type":     arg["type"].(string),
			"value":    argTypeObj,
			"valueRaw": "",
		})
	}

	result := map[string]interface{}{
		"valueRaw":         e.RawValue,
		"extrinsic_length": e.ExtrinsicLength,
		"version_info":     e.VersionInfo,
	}

	if e.ContainsTransaction {
		result["account_length"] = e.Address["account_length"]
		result["account_id"] = e.Address["account_id"]
		result["account_index"] = e.Address["account_index"]
		result["account_idx"] = e.Address["account_idx"]
		result["signature"] = e.Signature
		result["nonce"] = e.Nonce
		result["era"] = e.Era
		result["extrinsic_hash"] = e.ExtrinsicHash
	}
	if e.CallIndex != "" {
		result["call_code"] = e.CallIndex
		result["call_module_function"] = e.Call.Name
		result["call_module"] = e.CallModule.Name
	}
	result["params"] = e.Params
	return result
}
