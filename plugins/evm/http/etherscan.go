package http

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/itering/scale.go/utiles/crypto/keccak"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins/evm/contract"
	"github.com/itering/subscan/plugins/evm/dao"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
	"net/http"
	"path"
	"strings"
)

var srv dao.ISrv

func etherscanHandle(w http.ResponseWriter, r *http.Request) error {

	ErrRecordNotFound := errors.New("No records found")
	InvalidParam := errors.New("Invalid parameter")

	// module, action params
	actionParams := new(struct {
		Module string `form:"module" binding:"required,oneof=logs transaction account contract"`
		Action string `form:"action" binding:"required,oneof=getLogs getstatus gettxreceiptstatus balance balancemulti txlist txlistinternal tokentx token1155tx tokennfttx getabi getsourcecode getcontractcreation verifysourcecode checkverifystatus"`
	})

	logsParams := new(struct {
		Offset     int    `form:"offset" binding:"omitempty,min=1,max=1000"`
		Page       int    `form:"page" binding:"omitempty,min=1"`
		ToBlock    int    `form:"toBlock" binding:"min=0"`
		FromBlock  int    `form:"fromBlock" binding:"min=0"`
		Address    string `form:"address" binding:"omitempty,eth_addr"`
		Topic0     string `form:"topic0" binding:"omitempty,len=66"`
		Topic1     string `form:"topic1" binding:"omitempty,len=66"`
		Topic2     string `form:"topic2" binding:"omitempty,len=66"`
		Topic3     string `json:"topic3" binding:"omitempty,len=66"`
		Topic01Opr string `form:"topic0_1_opr" binding:"omitempty,oneof=and or"`
		Topic020pr string `form:"topic0_2_opr" binding:"omitempty,oneof=and or"`
		Topic031pr string `form:"topic0_3_opr" binding:"omitempty,oneof=and or"`
		Topic12Opr string `form:"topic1_2_opr" binding:"omitempty,oneof=and or"`
		Topic131pr string `form:"topic1_3_opr" binding:"omitempty,oneof=and or"`
		Topic23Opr string `form:"topic2_3_opr" binding:"omitempty,oneof=and or"`
	})

	txParams := new(struct {
		TxHash string `form:"txhash" binding:"omitempty,len=66"`
	})

	tokenParams := new(struct {
		Address         string `form:"address" binding:"omitempty,eth_addr"`
		ContractAddress string `form:"contractaddress" binding:"omitempty,eth_addr"`
		StartBlock      int    `form:"startblock" binding:"min=0"`
		EndBlock        int    `form:"endblock" binding:"min=0"`
		Page            int    `form:"page" binding:"min=1"`
		ToBlock         int    `form:"toBlock" binding:"min=0"`
		Offset          int    `form:"offset" binding:"min=1,max=1000"`
		Sort            string `form:"sort" binding:"omitempty,oneof=asc desc"`
	})

	if err := binding.Query.Bind(r, actionParams); err != nil {
		toJson(w, 0, nil, err)
		return nil
	}
	switch fmt.Sprintf("%s-%s", actionParams.Module, actionParams.Action) {
	case "logs-getLogs":
		if err := binding.Query.Bind(r, logsParams); err != nil {
			toJson(w, 0, nil, err)
			return nil
		}
		var opts []model.Option
		if logsParams.Offset == 0 || logsParams.Page == 0 {
			logsParams.Offset = 1000
			logsParams.Page = 1
		}
		opts = append(opts, model.WithLimit(logsParams.Page, logsParams.Offset))
		if logsParams.ToBlock > 0 {
			opts = append(opts, model.Where("block_num <= ?", logsParams.ToBlock))
		}
		if logsParams.FromBlock > 0 {
			opts = append(opts, model.Where("block_num >= ?", logsParams.FromBlock))
		}
		if logsParams.Address != "" {
			opts = append(opts, model.Where("address = ?", logsParams.Address))
		}
		if logsParams.Topic0 != "" {
			if logsParams.Topic01Opr != "" && logsParams.Topic1 != "" {
				// topic0 and topic1
				if logsParams.Topic01Opr == "and" {
					opts = append(opts, model.Where("method_hash = ? and topic1 = ?", logsParams.Topic0, logsParams.Topic1))
				} else {
					opts = append(opts, model.Where("method_hash = ? or topic1 = ?", logsParams.Topic0, logsParams.Topic1))
				}
			}
			if logsParams.Topic020pr != "" && logsParams.Topic2 != "" {
				// topic0 and topic2
				if logsParams.Topic020pr == "and" {
					opts = append(opts, model.Where("method_hash = ? and topic2 = ?", logsParams.Topic0, logsParams.Topic2))
				} else {
					opts = append(opts, model.Where("method_hash = ? or topic2 = ?", logsParams.Topic0, logsParams.Topic2))
				}
			}
			if logsParams.Topic031pr != "" && logsParams.Topic3 != "" {
				// topic0 and topic3
				if logsParams.Topic031pr == "and" {
					opts = append(opts, model.Where("method_hash = ? and topic3 = ?", logsParams.Topic0, logsParams.Topic3))
				} else {
					opts = append(opts, model.Where("method_hash = ? or topic3 = ?", logsParams.Topic0, logsParams.Topic3))
				}
			}
			if logsParams.Topic01Opr == "" && logsParams.Topic020pr == "" && logsParams.Topic031pr == "" {
				opts = append(opts, model.Where("method_hash = ?", logsParams.Topic0))
			}
		}
		if logsParams.Topic1 != "" {
			if logsParams.Topic12Opr != "" && logsParams.Topic2 != "" {
				// topic1 and topic2
				if logsParams.Topic12Opr == "and" {
					opts = append(opts, model.Where("topic1 = ? and topic2 = ?", logsParams.Topic1, logsParams.Topic2))
				} else {
					opts = append(opts, model.Where("topic1 = ? or topic2 = ?", logsParams.Topic1, logsParams.Topic2))
				}
			}
			if logsParams.Topic131pr != "" && logsParams.Topic3 != "" {
				// topic1 and topic3
				if logsParams.Topic131pr == "and" {
					opts = append(opts, model.Where("topic1 = ? and topic3 = ?", logsParams.Topic1, logsParams.Topic3))
				} else {
					opts = append(opts, model.Where("topic1 = ? or topic3 = ?", logsParams.Topic1, logsParams.Topic3))
				}
			}
			if logsParams.Topic12Opr == "" && logsParams.Topic131pr == "" && logsParams.Topic01Opr == "" {
				opts = append(opts, model.Where("topic1 = ?", logsParams.Topic1))
			}
		}
		if logsParams.Topic2 != "" {
			if logsParams.Topic23Opr != "" && logsParams.Topic3 != "" {
				// topic2 and topic3
				if logsParams.Topic23Opr == "and" {
					opts = append(opts, model.Where("topic2 = ? and topic3 = ?", logsParams.Topic2, logsParams.Topic3))
				} else {
					opts = append(opts, model.Where("topic2 = ? or topic3 = ?", logsParams.Topic2, logsParams.Topic3))
				}
			}
			if logsParams.Topic23Opr == "" && logsParams.Topic12Opr == "" && logsParams.Topic020pr == "" {
				opts = append(opts, model.Where("topic2 = ?", logsParams.Topic2))
			}
		}
		if logsParams.Topic3 != "" && logsParams.Topic23Opr == "" && logsParams.Topic131pr == "" && logsParams.Topic031pr == "" {
			opts = append(opts, model.Where("topic3 = ?", logsParams.Topic3))
		}
		res := srv.API_GetLogs(r.Context(), opts...)
		if len(res) == 0 {
			etherscanRes(w, 0, nil, ErrRecordNotFound)
			return nil
		}
		etherscanRes(w, 1, res, nil)

		// Check Contract Execution Status

	case "transaction-getstatus":
		if err := binding.Query.Bind(r, txParams); err != nil {
			toJson(w, 0, nil, err)
			return nil
		}
		txn := srv.GetTransactionByHash(r.Context(), txParams.TxHash)
		if txn == nil {
			etherscanRes(w, 0, nil, ErrRecordNotFound)
			return nil
		}
		var isError = 0
		if !txn.Success {
			isError = 1
		}
		etherscanRes(w, 1, map[string]interface{}{"isError": isError, "errDescription": ""}, nil)
		// Check Transaction Receipt Status

	case "transaction-gettxreceiptstatus":
		if err := binding.Query.Bind(r, txParams); err != nil {
			toJson(w, 0, nil, err)
			return nil
		}
		txn := srv.GetTransactionByHash(r.Context(), txParams.TxHash)
		if txn == nil {
			etherscanRes(w, 0, nil, ErrRecordNotFound)
			return nil
		}
		var status = 1
		if !txn.Success {
			status = 0
		}
		etherscanRes(w, 1, map[string]interface{}{"status": status}, nil)

	case "account-balance":
		accountParams := new(struct {
			Address string `form:"address" binding:"required,eth_addr"`
		})
		if err := binding.Query.Bind(r, accountParams); err != nil {
			toJson(w, 0, nil, err)
			return nil
		}
		account, err := srv.API_GetAccounts(r.Context(), []string{accountParams.Address})
		if err != nil || account == nil {
			etherscanRes(w, 0, nil, ErrRecordNotFound)
			return nil
		}
		etherscanRes(w, 1, map[string]interface{}{"balance": account[accountParams.Address]}, nil)

	case "account-balancemulti":
		accountParams := new(struct {
			Address string `form:"address" binding:"required"`
		})
		if err := binding.Query.Bind(r, accountParams); err != nil {
			toJson(w, 0, nil, err)
			return nil
		}
		var addresses []string
		for _, v := range strings.Split(accountParams.Address, ",") {
			if address.VerifyEthereumAddress(v) {
				addresses = append(addresses, v)
			} else {
				etherscanRes(w, 0, nil, InvalidParam)
				return nil
			}
		}
		account, err := srv.API_GetAccounts(r.Context(), addresses)
		if err != nil || account == nil {
			etherscanRes(w, 0, nil, ErrRecordNotFound)
			return nil
		}
		var result []map[string]interface{}
		for _, v := range addresses {
			result = append(result, map[string]interface{}{
				"address": v,
				"balance": account[v],
			})
		}
		etherscanRes(w, 1, result, nil)

	case "account-txlist":
		p := new(struct {
			Address    string `form:"address" binding:"required,eth_addr"`
			StartBlock int    `form:"startblock" binding:"min=0"`
			EndBlock   int    `form:"endblock" binding:"min=0"`
			Page       int    `form:"page" binding:"omitempty,min=1"`
			ToBlock    int    `form:"toBlock" binding:"min=0"`
			Offset     int    `form:"offset" binding:"min=1,max=1000"`
			Sort       string `form:"sort" binding:"omitempty,oneof=asc desc"`
		})
		if err := binding.Query.Bind(r, p); err != nil {
			toJson(w, 0, nil, err)
			return nil
		}
		if p.Offset*p.Page > 50000 {
			toJson(w, 0, nil, errors.New("page size too large"))
			return nil
		}
		var opts []model.Option
		opts = append(opts, model.Where("from_address = ?", p.Address))
		opts = append(opts, model.WithLimit(p.Page, p.Offset))
		if p.Sort == "" {
			p.Sort = "desc"
		}
		if p.ToBlock > 0 {
			opts = append(opts, model.Where("block_num <= ?", p.ToBlock))
		}
		if p.StartBlock > 0 {
			opts = append(opts, model.Where("block_num >= ?", p.StartBlock))
		}
		opts = append(opts, model.Order(fmt.Sprintf("id %s", p.Sort)))

		results := srv.API_Transactions(r.Context(), opts...)
		etherscanRes(w, 1, results, nil)

		// Get a list of 'Internal' Transactions by Address
		// https://docs.etherscan.io/etherscan-v2/api-endpoints/accounts#get-a-list-of-internal-transactions-by-address
		// https://docs.etherscan.io/etherscan-v2/api-endpoints/accounts#get-internal-transactions-by-transaction-hash
		// https://docs.etherscan.io/etherscan-v2/api-endpoints/accounts#get-internal-transactions-by-block-range

	case "account-txlistinternal": // todo

	// https://docs.etherscan.io/etherscan-v2/api-endpoints/accounts#get-a-list-of-erc20-token-transfer-events-by-address
	// https://docs.etherscan.io/etherscan-v2/api-endpoints/accounts#get-a-list-of-erc721-token-transfer-events-by-address
	// https://docs.etherscan.io/etherscan-v2/api-endpoints/accounts#get-a-list-of-erc1155-token-transfer-events-by-address
	case "account-tokentx", "account-tokennfttx", "account-token1155tx":
		fmt.Println(actionParams.Action)
		if err := binding.Query.Bind(r, tokenParams); err != nil {
			toJson(w, 0, nil, err)
			return nil
		}
		if tokenParams.Offset*tokenParams.Page > 50000 {
			toJson(w, 0, nil, errors.New("page size too large"))
			return nil
		}
		var opts []model.Option
		opts = append(opts, model.Where("from_address = ?", tokenParams.Address))
		opts = append(opts, model.WithLimit(tokenParams.Page, tokenParams.Offset))
		if tokenParams.Sort == "" {
			tokenParams.Sort = "desc"
		}
		if tokenParams.ToBlock > 0 {
			opts = append(opts, model.Where("block_num <= ?", tokenParams.ToBlock))
		}
		if tokenParams.StartBlock > 0 {
			opts = append(opts, model.Where("block_num >= ?", tokenParams.StartBlock))
		}
		opts = append(opts, model.Order(fmt.Sprintf("id %s", tokenParams.Sort)))
		var category uint
		switch actionParams.Action {
		case "account:tokentx":
			category = dao.TransferCategoryErc20
		case "account:tokennfttx":
			category = dao.TransferCategoryErc721
		case "account:token1155tx":
			category = dao.TransferCategoryErc1155
		}
		opts = append(opts, model.Where("category = ?", category))
		results := srv.API_TokenEventRes(r.Context(), opts...)
		etherscanRes(w, 1, results, nil)

	case "contract-getabi":
		p := new(struct {
			Address string `form:"address" binding:"required,eth_addr"`
		})
		if err := binding.Query.Bind(r, p); err != nil {
			toJson(w, 0, nil, err)
			return nil
		}
		contract := srv.ContractsByAddr(r.Context(), p.Address)
		if contract == nil {
			etherscanRes(w, 0, nil, ErrRecordNotFound)
			return nil
		}
		etherscanRes(w, 1, contract.Abi.String(), nil)

	case "contract-getsourcecode":
		p := new(struct {
			Address string `form:"address" binding:"required,eth_addr"`
		})
		if err := binding.Query.Bind(r, p); err != nil {
			toJson(w, 0, nil, err)
			return nil
		}
		contract := srv.ContractsByAddr(r.Context(), p.Address)
		if contract == nil {
			etherscanRes(w, 0, nil, ErrRecordNotFound)
			return nil
		}
		etherscanRes(w, 1, srv.API_ContractSourceCode(r.Context(), contract), nil)

	case "contract-getcontractcreation":
		p := new(struct {
			ContractAddresses string `form:"contractaddresses" binding:"required"`
		})
		if err := binding.Query.Bind(r, p); err != nil {
			toJson(w, 0, nil, err)
			return nil
		}
		var addresses []string
		for _, v := range strings.Split(p.ContractAddresses, ",") {
			if address.VerifyEthereumAddress(v) {
				addresses = append(addresses, v)
			} else {
				etherscanRes(w, 0, nil, InvalidParam)
				return nil
			}
		}
		if len(addresses) > 5 {
			etherscanRes(w, 0, nil, errors.New("too many addresses"))
			return nil
		}
		etherscanRes(w, 1, srv.API_GetContractCreation(r.Context(), addresses), nil)

	case "contract-verifysourcecode":
		type SourceCode struct {
			ContractAddress  string `form:"contractaddress" binding:"required,eth_addr"`
			SourceCode       string `form:"sourceCode" binding:"required"`
			CodeFormat       string `form:"codeformat" binding:"oneof=solidity-single-file solidity-standard-json-input"`
			ContractName     string `form:"contractname" binding:""`
			CompilerVersion  string `form:"compilerversion" binding:"required"`
			OptimizationUsed int    `form:"optimizationUsed" binding:"omitempty,oneof=0 1"`
			Runs             uint   `form:"runs" binding:"omitempty,min=1"`
			EvmVersion       string `form:"evmversion" binding:"omitempty"`
			LicenseType      int    `form:"licenseType" binding:"omitempty,min=1,max=14"`
			ResolcVersion    string `form:"resolcVersion" binding:"omitempty"`
		}
		const VerifyFail = "Fail - Unable to verify"

		p := new(SourceCode)
		if err := binding.Form.Bind(r, p); err != nil {
			etherscanRes(w, 0, VerifyFail, err)
			return nil
		}
		p.ContractAddress = address.Format(p.ContractAddress)

		if p.ContractName != "" {
			p.ContractName = strings.Split(p.ContractName, ":")[0]
		}
		externalLibrary := make(map[string]interface{})
		// Optimize
		if p.OptimizationUsed == 1 && p.Runs < 1 {
			etherscanRes(w, 0, VerifyFail, errors.New("runs is invalid"))
			return nil
		}
		if p.CodeFormat == "solidity-single-file" {
			p.CodeFormat = dao.VerifyTypeSingleFile
		}
		if p.CodeFormat == "solidity-standard-json-input" {
			p.CodeFormat = dao.VerifyStandardJsonFile
		}
		p.EvmVersion = dao.EvmVersionSelect(p.CompilerVersion)
		localContract := dao.ContractsByAddr(r.Context(), p.ContractAddress)
		if localContract == nil {
			etherscanRes(w, 0, "Unable to locate ContractCode", fmt.Errorf("Unable to locate ContractCode at %s", p.ContractAddress))
			return nil
		}
		if len(localContract.VerifyStatus) != 0 {
			etherscanRes(w, 1, "Already Verified", fmt.Errorf("Contract source code already verified"))
			return nil
		}
		if p.ResolcVersion != "" && !util.StringInSlice(p.ResolcVersion, contract.ReviveVersion) {
			etherscanRes(w, 1, VerifyFail, fmt.Errorf("Fail - Invalid resolc version"))
			return nil
		}
		var input *contract.CompilerJSONInput
		if p.CodeFormat == dao.VerifyTypeSingleFile {
			compileInstance := contract.NewSmartContractCompile(p.ContractName, p.SourceCode, p.CompilerVersion, p.EvmVersion, externalLibrary, p.OptimizationUsed == 1, p.Runs)
			input = compileInstance.AsInput(r.Context(), "")
		} else {
			if err := util.UnmarshalAny(&input, p.SourceCode); err != nil {
				util.Logger().Error(err)
			} else {
				if len(input.Settings.EvmVersion) == 0 {
					if p.EvmVersion != "" {
						input.Settings.EvmVersion = p.EvmVersion
					} else {
						if v := dao.EvmVersionSelect(input.Compiler.Version); v != "" {
							input.Settings.EvmVersion = v
						}
					}
				}
				if len(input.Compiler.Version) == 0 {
					input.Compiler.Version = p.CompilerVersion
				}
				if p.OptimizationUsed == 1 {
					input.Settings.Optimizer.Enabled = true
					input.Settings.Optimizer.Runs = int(p.Runs)
				}
				if len(input.Settings.CompilationTarget) == 0 && p.ContractName != "" {
					input.Settings.CompilationTarget = make(map[string]string)
					dir, pth := path.Split(p.ContractName)
					contractName := strings.ReplaceAll(pth, ".sol", "")
					if len(dir) != 0 {
						input.Settings.CompilationTarget[p.ContractName] = contractName
					}
					for ph := range input.Sources {
						if len(input.Settings.CompilationTarget) == 0 && strings.Contains(strings.ToLower(ph), strings.ToLower(contractName)) {
							input.Settings.CompilationTarget[ph] = contractName
						}
					}
					if len(input.Settings.CompilationTarget) == 0 {
						input.Settings.CompilationTarget[fmt.Sprintf("%s.sol", contractName)] = contractName
					}
				}
				for ph, source := range input.Sources {
					if len(source.Keccak256) == 0 {
						source.Keccak256 = util.AddHex(util.BytesToHex(keccak.Keccak256([]byte(source.Content))))
						input.Sources[ph] = source
					}
				}
			}
		}
		if input == nil {
			etherscanRes(w, 0, VerifyFail, fmt.Errorf("The contract code format is wrong"))
			return nil
		}

		input.Format()
		if p.ResolcVersion != "" {
			input.ResolcVersion = p.ResolcVersion
		}
		verify, err := input.VerifyFromJsonInput(r.Context(), p.ContractAddress)
		if err != nil {
			// raise http 500 error
			etherscanRes(w, 0, VerifyFail, err)
			return nil
		}
		if err = localContract.VerifySuccess(r.Context(), verify, p.CodeFormat, input); err != nil {
			etherscanRes(w, 0, VerifyFail, err)
			return nil
		}
		etherscanRes(w, 1, util.TrimHex(localContract.Address), nil)

	case "contract-checkverifystatus":
		p := new(struct {
			Guid string `form:"guid" binding:"required,eth_addr"`
		})
		if err := binding.Query.Bind(r, p); err != nil {
			toJson(w, 0, nil, err)
			return nil
		}

		localContract := srv.ContractsByAddr(r.Context(), util.AddHex(p.Guid))

		if localContract == nil {
			etherscanRes(w, 0, nil, ErrRecordNotFound)
			return nil
		}
		if len(localContract.VerifyStatus) != 0 {
			etherscanRes(w, 1, "Pass - Verified", nil)
			return nil
		} else {
			etherscanRes(w, 0, "Fail - Unable to verify", errors.New("NOTOK"))
			return nil
		}
	}
	return nil
}
