package http

import (
	"fmt"
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan/model"
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

		// token holder
		{"token/holder", tokenHolderHandle, http.MethodPost},
		{"tokens", tokenListHandle, http.MethodPost},
		{"token/transfer", tokenTransferHandle, http.MethodPost},
		{Router: "token/erc721/collectibles", Handle: collectiblesHandle, Method: http.MethodPost},
		{Router: "account/tokens", Handle: accountTokensHandle, Method: http.MethodPost},
	}
}

func accountTokensHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(struct {
		Address string `json:"address" validate:"required,eth_addr"`
	})
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	tokens := srv.AccountTokens(r.Context(), p.Address)
	toJson(w, 0, tokens, nil)
	return nil
}

func collectiblesHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(struct {
		Address  string `json:"address" validate:"omitempty,eth_addr"`
		Contract string `json:"contract" validate:"omitempty,eth_addr"`
		Page     int    `json:"page" validate:"min=0"`
		Row      int    `json:"row" validate:"min=1,max=100"`
	})
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

func tokenListHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(struct {
		Page     int    `json:"page" validate:"min=0"`
		Row      int    `json:"row" validate:"min=1,max=100"`
		Category string `json:"category" validate:"omitempty,oneof=erc20 erc721"`
	})
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	list, count := srv.TokenList(r.Context(), p.Category, p.Page, p.Row)
	toJson(w, 0, map[string]interface{}{"list": list, "count": count}, nil)
	return nil
}

func tokenTransferHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(struct {
		TokenAddress string `json:"token_address" validate:"omitempty,eth_addr"`
		Address      string `json:"address" validate:"omitempty,eth_addr"`
		Page         int    `json:"page" validate:"min=0"`
		Row          int    `json:"row" validate:"min=1,max=100"`
	})
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	if p.TokenAddress == "" && p.Address == "" {
		toJson(w, 10001, nil, fmt.Errorf("token_address or address is required"))
		return nil
	}
	transfers, count := srv.TokenTransfers(r.Context(), p.Address, p.TokenAddress, p.Page, p.Row)
	toJson(w, 0, map[string]interface{}{"transfers": transfers, "count": count}, nil)
	return nil
}

func tokenHolderHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(struct {
		TokenAddress string `json:"token_address" validate:"required,eth_addr"`
		Page         int    `json:"page" validate:"min=0"`
		Row          int    `json:"row" validate:"min=1,max=100"`
	})
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	holders, count := srv.TokenHolders(r.Context(), p.TokenAddress, p.Page, p.Row)
	toJson(w, 0, map[string]interface{}{"holders": holders, "count": count}, nil)
	return nil
}

func blocksHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(struct {
		Row  int `json:"row" validate:"min=1,max=100"`
		Page int `json:"page" validate:"min=0"`
	})
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

func blockHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(struct {
		Hash     string `json:"hash" validate:"omitempty,len=66"`
		BlockNum uint   `json:"block_num" validate:"omitempty,min=0"`
	})
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

func transactionHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(struct {
		Hash string `json:"hash" validate:"required,len=66"`
	})
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

func transactionsHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(struct {
		Page     int    `json:"page" validate:"min=0"`
		Row      int    `json:"row" validate:"min=1,max=100"`
		BlockNum uint   `json:"block_num" validate:"omitempty,min=0"`
		Address  string `json:"address" validate:"omitempty,eth_addr"`
	})
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
	opts = append(opts, model.WithLimit(p.Row*p.Page, p.Row))
	transactions := srv.TransactionsJson(r.Context(), opts...)
	toJson(w, 0, transactions, nil)
	return nil
}

func accountsHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(struct {
		Page int `json:"page" validate:"min=0"`
		Row  int `json:"row" validate:"min=1,max=100"`
	})
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	list, count := srv.Accounts(r.Context(), p.Page, p.Row)
	toJson(w, 0, map[string]interface{}{"list": list, "count": count}, nil)
	return nil
}

func contractHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(struct {
		Address string `json:"address" validate:"required,eth_addr"`
	})
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

func contractsHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(struct {
		Page int `json:"page" validate:"min=0"`
		Row  int `json:"row" validate:"min=1,max=100"`
	})
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	list, count := srv.Contracts(r.Context(), p.Page, p.Row)
	toJson(w, 0, map[string]interface{}{"list": list, "count": count}, nil)
	return nil
}
