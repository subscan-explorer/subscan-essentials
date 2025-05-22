package http

import (
	"fmt"
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins/evm/contract"
	"github.com/itering/subscan/plugins/evm/dao"
	"github.com/itering/subscan/util/validator"
	"net/http"
)

func Router() []router.Http {
	srv = &dao.ApiSrv{}
	return []router.Http{
		{"etherscan", etherscanHandle, http.MethodGet},

		{"blocks", blocksHandle, http.MethodPost},
		{"block", blockHandle, http.MethodPost},

		{"transactions", transactionsHandle, http.MethodPost},
		{"transaction", transactionHandle, http.MethodPost},

		{"accounts", accountsHandle, http.MethodPost},

		{"contract", contractHandle, http.MethodPost},
		{"contracts", contractsHandle, http.MethodPost},
		{"contract/solcs", solcVersions, http.MethodPost},
		{"contract/resolcs", resolcVersions, http.MethodPost},

		// token holder
		{"token/holder", tokenHolderHandle, http.MethodPost},
		{"tokens", tokenListHandle, http.MethodPost},
		{"token/transfer", tokenTransferHandle, http.MethodPost},
		{"token/erc721/collectibles", collectiblesHandle, http.MethodPost},
		{"account/tokens", accountTokensHandle, http.MethodPost},
	}
}

type accountTokensParams struct {
	Address  string `json:"address" validate:"required,eth_addr"`
	Category string `json:"category" validate:"omitempty,oneof=erc20 erc721"`
}

// @Summary Get account tokens
// @Tags EVM
// @Accept json
// @Produce json
// @Param params body accountTokensParams true "params"
// @Success 200 {object} J{data=[]dao.AccountTokenJson}
// @Router /api/plugin/evm/account/tokens [post]
func accountTokensHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(accountTokensParams)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	tokens := srv.AccountTokens(r.Context(), p.Address, p.Category)
	toJson(w, 0, tokens, nil)
	return nil
}

type collectiblesParams struct {
	Address  string `json:"address" validate:"omitempty,eth_addr"`
	Contract string `json:"contract" validate:"omitempty,eth_addr"`
	Page     int    `json:"page" validate:"min=0"`
	Row      int    `json:"row" validate:"min=1,max=100"`
}

// @Summary Evm Erc721 collectibles
// @Tags EVM
// @Accept json
// @Produce json
// @Param params body collectiblesParams true "params"
// @Success 200 {object} J{data=object{list=[]dao.Erc721Holders,count=int}}
// @Router /api/plugin/evm/token/erc721/collectibles [post]
func collectiblesHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(collectiblesParams)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	if p.Address == "" && p.Contract == "" {
		toJson(w, 10001, nil, fmt.Errorf("address or contract is required"))
		return nil
	}
	collectibles, count := srv.Collectibles(r.Context(), p.Address, p.Contract, p.Page, p.Row)
	toJson(w, 0, map[string]interface{}{"list": collectibles, "count": count}, nil)
	return nil
}

type tokenListParams struct {
	Page     int    `json:"page" validate:"min=0"`
	Row      int    `json:"row" validate:"min=1,max=100"`
	Category string `json:"category" validate:"omitempty,oneof=erc20 erc721"`
	Contract string `json:"contract" validate:"omitempty,eth_addr"`
}

// @Summary Evm token list
// @Tags EVM
// @Accept json
// @Produce json
// @Param params body tokenListParams true "params"
// @Success 200 {object} J{data=object{list=[]dao.Token,count=int}}
// @Router /api/plugin/evm/tokens [post]
func tokenListHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(tokenListParams)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	list, count := srv.TokenList(r.Context(), p.Contract, p.Category, p.Page, p.Row)
	toJson(w, 0, map[string]interface{}{"list": list, "count": count}, nil)
	return nil
}

type tokenTransferParams struct {
	TokenAddress string `json:"token_address" validate:"omitempty,eth_addr"`
	Address      string `json:"address" validate:"omitempty,eth_addr"`
	Page         int    `json:"page" validate:"min=0"`
	Row          int    `json:"row" validate:"min=1,max=100"`
	Category     string `json:"category" validate:"omitempty,oneof=erc20 erc721"`
}

// @Summary Evm token transfer
// @Tags EVM
// @Accept json
// @Produce json
// @Param params body tokenTransferParams true "params"
// @Success 200 {object} J{data=object{list=[]dao.TokenTransferJson,count=int}}
// @Router /api/plugin/evm/token/transfer [post]
func tokenTransferHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(tokenTransferParams)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	if p.TokenAddress == "" && p.Address == "" {
		toJson(w, 10001, nil, fmt.Errorf("token_address or address is required"))
		return nil
	}
	transfers, count := srv.TokenTransfers(r.Context(), p.Address, p.TokenAddress, p.Category, p.Page, p.Row)
	toJson(w, 0, map[string]interface{}{"transfers": transfers, "count": count}, nil)
	return nil
}

type tokenHolderParams struct {
	TokenAddress string `json:"token_address" validate:"required,eth_addr"`
	Page         int    `json:"page" validate:"min=0"`
	Row          int    `json:"row" validate:"min=1,max=100"`
}

// @Summary Evm token holder
// @Tags EVM
// @Accept json
// @Produce json
// @Param params body tokenHolderParams true "params"
// @Success 200 {object} J{data=object{list=[]dao.TokenHolder,count=int}}
// @Router /api/plugin/evm/token/holder [post]
func tokenHolderHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(tokenHolderParams)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	holders, count := srv.TokenHolders(r.Context(), p.TokenAddress, p.Page, p.Row)
	toJson(w, 0, map[string]interface{}{"holders": holders, "count": count}, nil)
	return nil
}

type EvmBlocks struct {
	Row  int `json:"row" validate:"min=1,max=100"`
	Page int `json:"page" validate:"min=0"`
}

// @Summary Evm blocks
// @Tags EVM
// @Accept json
// @Produce json
// @Param params body EvmBlocks true "params"
// @Success 200 {object} J{data=object{list=[]dao.EvmBlockJson,count=int}}
// @Router /api/plugin/evm/blocks [post]
func blocksHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(EvmBlocks)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}

	list, count := srv.Blocks(r.Context(), p.Page, p.Row)
	toJson(w, 0, map[string]interface{}{
		"list": list, "count": count,
	}, nil)
	return nil
}

type EvmBlockParams struct {
	Hash     string `json:"hash" validate:"omitempty,len=66"`
	BlockNum uint   `json:"block_num" validate:"omitempty,min=0"`
}

// @Summary Evm block info
// @Tags EVM
// @Accept json
// @Produce json
// @Param params body EvmBlockParams true "params"
// @Success 200 {object} J{data=dao.EvmBlock}
// @Router /api/plugin/evm/block [post]
func blockHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(EvmBlockParams)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	if p.Hash == "" && p.BlockNum == 0 {
		toJson(w, 10001, nil, fmt.Errorf("hash or block_num is required"))
		return nil
	}
	if p.Hash != "" {
		block := srv.BlockByHash(r.Context(), p.Hash)
		toJson(w, 0, block, nil)
		return nil
	}
	block := srv.BlockByNum(r.Context(), p.BlockNum)
	toJson(w, 0, block, nil)
	return nil
}

type transactionParam struct {
	Hash string `json:"hash" validate:"required,len=66"`
}

// @Summary Evm transaction info
// @Tags EVM
// @Accept json
// @Produce json
// @Param params body transactionParam true "params"
// @Success 200 {object} J{data=dao.Transaction}
// @Router /api/plugin/evm/transaction [post]
func transactionHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(transactionParam)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	transaction := srv.GetTransactionByHash(r.Context(), p.Hash)
	if transaction == nil {
		toJson(w, 10002, nil, fmt.Errorf("transaction not found"))
		return nil
	}
	toJson(w, 0, transaction, nil)
	return nil
}

type transactionsParams struct {
	Page     int    `json:"page" validate:"min=0"`
	Row      int    `json:"row" validate:"min=1,max=100"`
	BlockNum uint   `json:"block_num" validate:"omitempty,min=0"`
	Address  string `json:"address" validate:"omitempty,eth_addr"`
}

// @Summary Evm transactions
// @Tags EVM
// @Accept json
// @Produce json
// @Param params body transactionsParams true "params"
// @Success 200 {object} J{data=object{list=[]dao.TransactionSampleJson,count=int}}
// @Router /api/plugin/evm/transactions [post]
func transactionsHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(transactionsParams)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	var opts []model.Option
	if p.Address != "" {
		opts = append(opts, model.Where("from_address = ? or to_address = ?", p.Address, p.Address))
	}
	if p.BlockNum > 0 {
		opts = append(opts, model.Where("block_num = ?", p.BlockNum))
	}
	list, count := srv.TransactionsJson(r.Context(), model.WithLimit(p.Row*p.Page, p.Row), opts...)
	toJson(w, 0, map[string]interface{}{"list": list, "count": count}, nil)
	return nil
}

type EvmAccountParams struct {
	Page    int    `json:"page" validate:"min=0"`
	Row     int    `json:"row" validate:"min=1,max=100"`
	Address string `json:"address" validate:"omitempty,eth_addr"`
}

// @Summary Evm accounts list
// @Tags EVM
// @Accept json
// @Produce json
// @Param params body EvmAccountParams true "params"
// @Success 200 {object} J{data=object{list=[]dao.AccountsJson,count=int}}
// @Router /api/plugin/evm/accounts [post]
func accountsHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(EvmAccountParams)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	list, count := srv.Accounts(r.Context(), p.Address, p.Page, p.Row)
	toJson(w, 0, map[string]interface{}{"list": list, "count": count}, nil)
	return nil
}

type contractParams struct {
	Address string `json:"address" validate:"required,eth_addr"`
}

// @Summary Evm contract info
// @Tags EVM
// @Accept json
// @Produce json
// @Param params body contractParams true "params"
// @Success 200 {object} J{data=dao.Contract}
// @Router /api/plugin/evm/contract [post]
func contractHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(contractParams)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	contract := srv.ContractsByAddr(r.Context(), p.Address)
	if contract == nil {
		toJson(w, 10002, nil, fmt.Errorf("contract not found"))
		return nil
	}
	toJson(w, 0, contract, nil)
	return nil
}

type contractsParams struct {
	Page int `json:"page" validate:"min=0"`
	Row  int `json:"row" validate:"min=1,max=100"`
}

// @Summary Evm contract list
// @Tags EVM
// @Accept json
// @Produce json
// @Param params body contractsParams true "params"
// @Success 200 {object} J{data=object{list=[]dao.ContractsJson,count=int}}
// @Router /api/plugin/evm/contracts [post]
func contractsHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(contractsParams)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	list, count := srv.Contracts(r.Context(), p.Page, p.Row)
	toJson(w, 0, map[string]interface{}{"list": list, "count": count}, nil)
	return nil
}

// @Summary Polkadot pvm resolc versions
// @Tags EVM
// @Accept json
// @Produce json
// @Success 200 {object} J{data=[]string}
// @Router /api/plugin/evm/contract/solcs [post]
func resolcVersions(w http.ResponseWriter, r *http.Request) error {
	toJson(w, 0, contract.ReviveVersion, nil)
	return nil
}

type EvmContractSolcVersionsParam struct {
	Releases bool `json:"releases" binding:"omitempty"`
}

// @Summary EVM contract solc versions
// @Tags EVM
// @Accept json
// @Produce json
// @Param param body EvmContractSolcVersionsParam true "param"
// @Success 200 {object} J{data=[]string}
// @Router /api/scan/evm/contract/solcs [post]
func solcVersions(w http.ResponseWriter, r *http.Request) error {
	p := new(EvmContractSolcVersionsParam)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}

	list, err := contract.SolcVersions(r.Context(), p.Releases)
	toJson(w, 0, list, err)
	return nil
}
