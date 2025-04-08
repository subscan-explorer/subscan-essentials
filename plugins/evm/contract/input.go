package contract

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/itering/subscan/util"
	"strings"
)

type CompilerJSONInput struct {
	Language string              `json:"language"`
	Sources  SourcesCode         `json:"sources"`
	Settings SolcMetadataSetting `json:"settings"`
	Compiler struct {
		Version string `json:"version"`
	} `json:"compiler"`
}

func (metadataValue *CompilerJSONInput) FormatContractName() string {
	path := util.EnumStringKey(metadataValue.Settings.CompilationTarget)
	if path == "" {
		return ""
	}
	// xx/xx.sol
	pathArr := strings.Split(path, "/")
	last := pathArr[len(pathArr)-1]
	return strings.ReplaceAll(last, ".sol", "")
}

func (metadataValue *CompilerJSONInput) Format() {
	// sources code
	sources := make(map[string]SolcSources)
	for k, v := range metadataValue.Sources {
		source := v
		if v.Keccak256 == "" {
			source.Keccak256 = common.BytesToHash(crypto.Keccak256([]byte(v.Content))).Hex()
		}
		sources[k] = source
	}
	metadataValue.Sources = sources
}

type SourcesCode map[string]SolcSources

func (s SourcesCode) AsString() string {
	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		for _, v := range s {
			return v.Content
		}
	}
	sources := make(map[string]string)
	for path, v := range s {
		sources[path] = v.Content
	}
	return util.ToString(sources)
}
