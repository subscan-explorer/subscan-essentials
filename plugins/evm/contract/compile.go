package contract

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/itering/subscan/share/web3"
	"github.com/itering/subscan/util"
	"os"
)

var (
	verifyServer = os.Getenv("VERIFY_SERVER")
)

type SmartContractCompile struct {
	ContractName      string                 `json:"contract_name"`
	SourceCode        string                 `json:"source_code"`
	CompilerVersion   string                 `json:"compiler_version"`
	EvmVersion        string                 `json:"evm_version"`
	ExternalLibraries map[string]interface{} `json:"external_libraries"`
	Optimize          bool                   `json:"optimize"`
	OptimizationRuns  uint                   `json:"optimization_runs"`
}

type SolcMetadata struct {
	Language string              `json:"language"`
	Sources  SourcesCode         `json:"sources"`
	Compiler map[string]string   `json:"compiler"`
	Output   SolcMetadataOutput  `json:"output"`
	Settings SolcMetadataSetting `json:"settings"`
	Version  float64             `json:"version"`
}

type SolcMetadataOutput struct {
	Abi                    []interface{} `json:"abi"`
	Devdoc                 interface{}   `json:"devdoc"`
	Userdoc                interface{}   `json:"userdoc"`
	CreationBytecodeLength int           `json:"creationBytecodeLength"` // custom, sourcify return
}

type SolcSources struct {
	Keccak256 string   `json:"keccak256"`
	License   string   `json:"license"`
	Content   string   `json:"content"`
	Urls      []string `json:"urls"`
}

type SolcMetadataSetting struct {
	Remappings []string `json:"remappings,omitempty"`
	Optimizer  struct {
		Enabled bool `json:"enabled"`
		Runs    int  `json:"runs"`
	} `json:"optimizer"`
	EvmVersion        string                 `json:"evmVersion"`
	Libraries         map[string]interface{} `json:"libraries,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	CompilationTarget map[string]string      `json:"compilationTarget"`
	OutputSelection   interface{}            `json:"outputSelection,omitempty"`
	ReviveVersion     string                 `json:"revive_version,omitempty"`
}

func NewSmartContractCompile(ContractName, SourceCode, CompilerVersion, EvmVersion string, ExternalLibraries map[string]interface{}, Optimize bool, OptimizationRuns uint) *SmartContractCompile {
	return &SmartContractCompile{
		ContractName:      ContractName,
		SourceCode:        SourceCode,
		CompilerVersion:   CompilerVersion,
		EvmVersion:        EvmVersion,
		ExternalLibraries: ExternalLibraries,
		Optimize:          Optimize,
		OptimizationRuns:  OptimizationRuns,
	}
}

func SolcVersions(_ context.Context, release bool) (list []string, err error) {
	if release {
		return soljsonReleases, nil
	}
	return soljsonSources, nil
}

type SourcifyRes struct {
	Result []SourcifyResResult `json:"result"`
}

type SourcifyResResult struct {
	Status string `json:"status"`
}

func (s *SmartContractCompile) AsInput(_ context.Context, compilationTarget string) *CompilerJSONInput {
	var metadataValue CompilerJSONInput

	metadataValue.Language = "Solidity"
	target := s.ContractName
	metadataValue.Settings = SolcMetadataSetting{
		Optimizer: struct {
			Enabled bool `json:"enabled"`
			Runs    int  `json:"runs"`
		}{
			Enabled: s.Optimize,
			Runs:    int(s.OptimizationRuns),
		},
		EvmVersion: s.EvmVersion,
		Libraries:  s.ExternalLibraries,
		// Metadata:          map[string]interface{}{"bytecodeHash": "ipfs"},
		CompilationTarget: map[string]string{s.ContractName: s.ContractName},
		OutputSelection: map[string]interface{}{
			"*": map[string][]string{"*": {"evm.bytecode", "metadata", "abi"}},
		},
	}
	if compilationTarget != "" {
		target = compilationTarget
		delete(metadataValue.Settings.CompilationTarget, s.ContractName)
		metadataValue.Settings.CompilationTarget[compilationTarget] = s.ContractName

	}
	metadataValue.Compiler.Version = s.CompilerVersion
	metadataValue.Sources = map[string]SolcSources{target: {Content: s.SourceCode, Keccak256: common.BytesToHash(crypto.Keccak256([]byte(s.SourceCode))).Hex()}}
	return &metadataValue
}

type VerificationRes struct {
	VerifiedStatus         string        `json:"verified_status"`
	Message                string        `json:"message"`
	Abi                    []interface{} `json:"abi"`
	CreationBytecodeLength int           `json:"creation_bytecode_length"`
	ReviveVersion          string        `json:"revive_version,omitempty"`
	ContractName           string        `json:"contract_name,omitempty"`
}

func (metadataValue *CompilerJSONInput) VerifyFromJsonInput(_ context.Context, address string) (*VerificationRes, error) {
	if verifyServer != "" {
		type Input struct {
			Address         string `json:"address"`
			CompilerVersion string `json:"compilerVersion"`
			Metadata        string `json:"metadata"`
			Chain           int    `json:"chain"`
		}
		data, err := util.PostWithJson(util.ToBytes(Input{
			Address:         address,
			CompilerVersion: metadataValue.Compiler.Version,
			Metadata:        util.ToString(metadataValue),
			Chain:           int(web3.CHAIN_ID),
		}), fmt.Sprintf("%s/verify", verifyServer))
		if err != nil {
			return nil, err
		}
		var vr VerificationRes
		if err = json.Unmarshal(data, &vr); err != nil {
			return nil, err
		}
		const mismatch = "mismatch"
		if vr.VerifiedStatus == mismatch {
			return nil, errors.New(vr.Message)
		}
		return &vr, nil
	}
	return nil, errors.New("verify server disabled")
}
