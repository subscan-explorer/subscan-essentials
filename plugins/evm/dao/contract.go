package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	evmABI "github.com/itering/subscan/plugins/evm/abi"
	evmContract "github.com/itering/subscan/plugins/evm/contract"
	"github.com/itering/subscan/plugins/evm/feature/delegateProxy"
	"github.com/itering/subscan/share/web3"
	"github.com/itering/subscan/util"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Contract struct {
	Address string `json:"address" gorm:"primaryKey;autoIncrement:false;size:255"`

	Abi          datatypes.JSON `json:"abi" `
	SourceCode   string         `json:"source_code" gorm:"type:string"`
	CreationCode string         `json:"creation_code" gorm:"type:string"`

	CreationBytecode  string         `json:"creation_bytecode" gorm:"type:string"`
	MethodIdentifiers datatypes.JSON `json:"method_identifiers" `
	EventIdentifiers  datatypes.JSON `json:"event_identifiers" gorm:"-"`
	Deployer          string         `json:"deployer" gorm:"size:100"`
	BlockNum          uint           `json:"block_num" gorm:"size:32"`
	TxHash            string         `json:"tx_hash" gorm:"size:70"`
	DeployAt          uint           `json:"deploy_at" gorm:"size:32"`
	VerifyStatus      string         `json:"verify_status"  gorm:"size:32;index:verify_status"`
	VerifyType        string         `json:"verify_type"  gorm:"size:100;default:SingleFile;"`
	ContractName      string         `json:"contract_name"  gorm:"size:255"`
	CompilerVersion   string         `json:"compiler_version"  gorm:"size:100"`
	EvmVersion        string         `json:"evm_version" gorm:"size:100"`

	ExternalLibraries datatypes.JSON `json:"external_libraries"`
	Optimize          bool           `json:"optimize"`
	OptimizationRuns  uint           `json:"optimization_runs"  gorm:"size:32"`
	ExtrinsicIndex    string         `json:"extrinsic_index" gorm:"size:255"`
	VerifyTime        uint           `json:"verify_time" gorm:"size:32"`
	TransactionCount  uint           `json:"transaction_count" gorm:"size:32;index:transaction_count"`
	Precompile        uint           `json:"precompile"`
	CompileSettings   datatypes.JSON `json:"CompileSettings"`

	EipStandard          string `json:"eip_standard" gorm:"size:100"`
	ProxyImplementation  string `json:"proxy_implementation" gorm:"size:64"`
	ConstructorArguments string `json:"constructor_arguments" gorm:"type:string"`
	DeployCodeHash       string `json:"deploy_code_hash" gorm:"size:70;index:deploy_code_hash;default:'';not null"`
}

type ContractSampleJson struct {
	Address           string         `json:"address"`
	Deployer          string         `json:"deployer"`
	BlockNum          uint           `json:"block_num" `
	DeployAt          uint           `json:"deploy_at"`
	VerifyStatus      string         `json:"verify_status" `
	VerifyTime        uint           `json:"verify_time"`
	ContractName      string         `json:"contract_name"`
	MethodIdentifiers datatypes.JSON `json:"method_identifiers"`
	EventIdentifiers  datatypes.JSON `json:"event_identifiers"`
	Abi               datatypes.JSON `json:"-" `
}

type ContractListJson struct {
	Address          string          `json:"address"`
	ContractName     string          `json:"contract_name"`
	EvmVersion       string          `json:"evm_version"`
	Balances         decimal.Decimal `json:"balances"`
	VerifyStatus     string          `json:"verify_status" `
	VerifyTime       uint            `json:"verify_time"`
	TransactionCount uint            `json:"transaction_count"`
}

const (
	VerifyTypeSingleFile   = "SingleFile"
	VerifyStandardJsonFile = "StandardJson"
)

func (c *Contract) TableName() string {
	return "evm_contracts"
}

func (c *Contract) AfterCreate(*gorm.DB) (err error) {
	_ = TouchAccount(context.Background(), c.Address)
	return nil
}

func (c *Contract) VerifySuccess(ctx context.Context, verifyRes *evmContract.VerificationRes, VerifyType string, input *evmContract.CompilerJSONInput) error {
	if verifyRes == nil {
		return errors.New("verify status is nil")
	}

	abiBytes, _ := json.Marshal(verifyRes.Abi)
	c.Abi = abiBytes
	if verifyRes.CreationBytecodeLength > 0 && len(c.CreationCode) >= verifyRes.CreationBytecodeLength {
		c.CreationBytecode = c.CreationCode[:verifyRes.CreationBytecodeLength]
	}
	c.ContractName = input.FormatContractName()
	c.SourceCode = input.Sources.AsString()
	c.VerifyType = VerifyType
	c.VerifyStatus = verifyRes.VerifiedStatus
	c.CompilerVersion = input.Compiler.Version
	c.EvmVersion = input.Settings.EvmVersion
	c.Optimize = input.Settings.Optimizer.Enabled
	c.OptimizationRuns = uint(input.Settings.Optimizer.Runs)
	c.VerifyTime = uint(time.Now().Unix())
	var abiValue abi.ABI
	_ = abiValue.UnmarshalJSON(c.Abi)
	methodIdentifiers := make(map[string]string)
	for _, method := range abiValue.Methods {
		methodIdentifiers[method.Sig] = util.BytesToHex(method.ID)
	}
	MethodIdentifiers, _ := json.Marshal(methodIdentifiers)
	c.MethodIdentifiers = MethodIdentifiers
	settingsValue, _ := json.Marshal(input.Settings)
	c.CompileSettings = settingsValue
	if len(input.Settings.Libraries) > 0 {
		externalLibraries, _ := json.Marshal(input.Settings.Libraries)
		c.ExternalLibraries = externalLibraries
	}
	c.afterVerify(ctx)

	return sg.db.Model(Contract{}).Where("address = ?", c.Address).Updates(c).Error
}

func (c *Contract) afterVerify(ctx context.Context) {
	_ = c.fetchAbiMapping(context.Background())
	// check it is proxy contract
	if c.hasEvent(ctx, delegateProxy.EventUpgraded) {
		var proxy delegateProxy.IDelegateProxy
		if c.hasStorage(ctx, "implementation") {
			proxy = delegateProxy.Init897(web3.RPC, c.Address)
		} else {
			proxy = delegateProxy.Init1967(web3.RPC, c.Address)
		}
		if implementation, _ := proxy.Implementation(ctx); implementation != "" {
			c.ProxyImplementation = implementation
			c.EipStandard = proxy.Standard()
		}
	}
}

func setContractProxyImplementation(ctx context.Context, contractAddress, implementation string) {
	if contract := GetContract(ctx, contractAddress); contract != nil && contract.ProxyImplementation != "" {
		sg.db.Model(Contract{}).Debug().Where("address = ?", contractAddress).Update("proxy_implementation", implementation)
	}
}

func (c *Contract) hasEvent(_ context.Context, eventId string) bool {
	var abiValue abi.ABI
	_ = abiValue.UnmarshalJSON(c.Abi)
	eventId = util.AddHex(eventId)
	for _, event := range abiValue.Events {
		if eventId == event.ID.Hex() {
			return true
		}
	}
	return false
}

func (c *Contract) hasStorage(_ context.Context, storageName string) bool {
	var abiValue abi.ABI
	_ = abiValue.UnmarshalJSON(c.Abi)
	for name, method := range abiValue.Methods {
		if method.StateMutability == "view" && name == storageName {
			return true
		}
	}
	return false
}

func (t *Transaction) NewContract(ctx context.Context) error {
	contract := &Contract{
		Address:        t.Contract,
		CreationCode:   t.InputData,
		DeployAt:       t.BlockTimestamp,
		BlockNum:       t.BlockNum,
		Deployer:       t.FromAddress,
		ExtrinsicIndex: t.ExtrinsicIndex,
		Precompile:     t.Precompile,
	}
	return sg.AddOrUpdateItem(ctx, contract, []string{"address"}, "creation_code", "deploy_at", "block_num", "deployer", "extrinsic_index", "precompile").Error
}

func ContractAddr(ctx context.Context) (list []string) {
	sg.db.WithContext(ctx).Model(Contract{}).Pluck("address", &list)
	return
}

func IsContract(ctx context.Context, addr string) bool {
	return util.StringInSlice(addr, ContractAddr(ctx))
}

func ContractCount(ctx context.Context) int64 {
	return int64(len(ContractAddr(ctx)))
}

func VerifyContractCount(ctx context.Context) (count int64) {
	sg.db.WithContext(ctx).Model(Contract{}).Where("verify_status !=''").Count(&count)
	return
}

func incrContractTransactionCount(ctx context.Context, address string) {
	sg.db.WithContext(ctx).Model(Contract{}).
		Where("address = ?", address).
		UpdateColumns(map[string]interface{}{"transaction_count": gorm.Expr("transaction_count + 1")})
}

type ContractDisplay struct {
	Address        string `json:"address"`
	IsContract     bool   `json:"is_contract"`
	PrecompileName string `json:"precompile_name,omitempty"`
}

func ContractsList(ctx context.Context, page, row int, addresses []string, Verified bool, search, order, orderField string) (contracts []ContractListJson, count int64) {
	query := sg.db.Model(Contract{})
	if len(addresses) > 0 {
		query.Where("address in (?)", addresses)
	}
	if Verified {
		query.Where("verify_status !=''")
	}
	if search != "" {
		query.Where("contract_name like ?", "%"+search+"%")
	}
	query.Count(&count)

	if order != "" && orderField != "" {
		query.Order(fmt.Sprintf("%s %s", orderField, order))
	}
	query.Offset(page * row).Limit(row).Scan(&contracts)
	return
}

func ContractsByAddrList(ctx context.Context, addresses []string) (contracts []ContractSampleJson) {
	sg.db.Model(Contract{}).Where("address in (?)", addresses).Scan(&contracts)
	for k, v := range contracts {
		if len(v.Abi) > 0 && v.Abi.String() != "null" {
			contracts[k].EventIdentifiers = findEventIdentifiers(ctx, v.Abi)
		}
	}
	return
}

func ContractMethodList(ctx context.Context) (list []datatypes.JSON, err error) {
	err = sg.db.WithContext(ctx).Model(Contract{}).Not("verify_status = ''").Where("method_identifiers IS NOT NULL").Pluck("method_identifiers", &list).Error
	return
}

func ContractsByAddr(ctx context.Context, contracts string) (contract *Contract) {
	if q := sg.db.Model(Contract{}).Where("address = ?", contracts).First(&contract); q.Error != nil {
		return nil
	}

	if len(contract.Abi) > 0 && contract.Abi.String() != "null" {
		contract.EventIdentifiers = findEventIdentifiers(ctx, contract.Abi)
	}
	return
}

func findEventIdentifiers(_ context.Context, abiRaw []byte) []byte {
	var abiValue abi.ABI
	_ = abiValue.UnmarshalJSON(abiRaw)
	eventIdentifiers := make(map[string]string)
	for _, event := range abiValue.Events {
		eventIdentifiers[event.String()] = event.ID.Hex()
	}
	EventIdentifiers, _ := json.Marshal(eventIdentifiers)
	return EventIdentifiers
}

func FindAbiMethodIdentifiers(_ context.Context, abiRaw []byte) []byte {
	var abiValue abi.ABI
	_ = abiValue.UnmarshalJSON(abiRaw)
	methodsIdentifiers := make(map[string]string)
	for _, method := range abiValue.Methods {
		var args []string
		for _, arg := range method.Inputs {
			args = append(args, arg.Type.String())
		}
		methodName := fmt.Sprintf("%s(%s)", method.Name, strings.Join(args, ","))
		methodsIdentifiers[methodName] = evmABI.EncodingMethod(methodName)[0:8]
	}
	identifiers, _ := json.Marshal(methodsIdentifiers)
	return identifiers
}

func GetContract(ctx context.Context, address string) *Contract {
	var contract Contract
	if q := sg.db.WithContext(ctx).Model(Contract{}).Where("address = ?", address).First(&contract); q.Error != nil {
		return nil
	}
	return &contract
}

func GetContractName(ctx context.Context, addresses []string) map[string]string {
	var contractName []struct {
		Address      string `json:"address"`
		ContractName string `json:"contract_name"`
	}
	result := make(map[string]string)
	if len(addresses) == 0 {
		return result
	}
	sg.db.WithContext(ctx).Model(&Contract{}).Where("address IN ?", addresses).Find(&contractName)
	for _, c := range contractName {
		result[c.Address] = c.ContractName
	}
	return result
}

// EvmVersionSelect https://github.com/ethereum/solidity/blob/develop/Changelog.md
func EvmVersionSelect(version string) string {
	if version == "" {
		return ""
	}
	re := regexp.MustCompile(`v\d+(\.\d+){0,2}`)
	v := re.FindStringSubmatch(version)
	if len(v) == 0 {
		return ""
	}
	calcVer := func(major, minor, revision int) int {
		return major*10000 + minor*10000 + revision
	}
	verNumArr := strings.Split(strings.ReplaceAll(v[0], "v", ""), ".")
	if len(verNumArr) < 2 {
		return ""
	}
	var major, minor, revision int
	var err error
	if major, err = strconv.Atoi(verNumArr[0]); err != nil {
		return ""
	}
	if minor, err = strconv.Atoi(verNumArr[1]); err != nil {
		return ""
	}
	if len(verNumArr) > 2 {
		if revision, err = strconv.Atoi(verNumArr[2]); err != nil {
			return ""
		}
	}
	ver := calcVer(major, minor, revision)
	if ver >= calcVer(0, 8, 25) {
		return "cancun"
	}
	if ver >= calcVer(0, 8, 20) {
		return "shanghai"
	}
	if ver >= calcVer(0, 8, 18) {
		return "paris"
	}
	if ver >= calcVer(0, 8, 7) {
		return "london"
	}
	if ver >= calcVer(0, 8, 5) {
		return "berlin"
	}
	if ver >= calcVer(0, 5, 14) {
		return "istanbul"
	}
	return "default"
}
