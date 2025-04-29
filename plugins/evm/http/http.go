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
	}
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
