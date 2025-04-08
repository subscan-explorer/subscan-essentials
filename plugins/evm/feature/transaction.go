package feature

import (
	"context"
	"encoding/hex"
	"github.com/itering/subscan/util"
	"github.com/shopspring/decimal"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EIP1559 struct {
	ChainID              *big.Int        `json:"chain_id"`
	Nonce                decimal.Decimal `json:"nonce"`
	MaxPriorityFeePerGas decimal.Decimal `json:"max_priority_fee_per_gas"`
	MaxFeePerGas         decimal.Decimal `json:"max_fee_per_gas"`
	GasLimit             decimal.Decimal `json:"gas_limit"`
	Action               struct {
		Call *string `json:"call,omitempty"`
	} `json:"action"`
	Value      decimal.Decimal  `json:"value"`
	Input      string           `json:"input"`
	AccessList types.AccessList `json:"access_list,omitempty"`
	OddYParity bool             `json:"odd_y_parity"`
	// Signature values
	V string `json:"v"`
	R string `json:"r"`
	S string `json:"s"`
}

type Legacy struct {
	Nonce    decimal.Decimal `json:"nonce"`
	GasPrice decimal.Decimal `json:"gas_price"`
	GasLimit decimal.Decimal `json:"gas_limit"`
	Action   struct {
		Call *string `json:"call,omitempty"`
	} `json:"action"`
	Value     decimal.Decimal `json:"value"`
	Input     string          `json:"input"`
	Signature struct {
		V *big.Int `json:"v"`
		R string   `json:"r"`
		S string   `json:"s"`
	} `json:"signature"`
}

type TransactionParams struct {
	EIP1559 *EIP1559 `json:"EIP1559,omitempty"`
	Legacy  *Legacy  `json:"Legacy,omitempty"`
}

func CalTransactionHash(_ context.Context, raw []byte) string {
	var txn TransactionParams
	err := util.UnmarshalAny(&txn, raw)
	if err != nil {
		util.Logger().Error(err)
		return ""
	}

	if txn.EIP1559 != nil {
		var eip1559 = txn.EIP1559
		dy := &types.DynamicFeeTx{
			ChainID:    eip1559.ChainID,
			Nonce:      uint64(eip1559.Nonce.IntPart()),
			GasTipCap:  eip1559.MaxPriorityFeePerGas.BigInt(),
			GasFeeCap:  eip1559.MaxFeePerGas.BigInt(),
			Gas:        uint64(eip1559.GasLimit.IntPart()),
			Value:      eip1559.Value.BigInt(),
			Data:       util.HexToBytes(eip1559.Input),
			AccessList: eip1559.AccessList,
			V:          big.NewInt(0),
			R:          util.U256(eip1559.R),
			S:          util.U256(eip1559.S),
		}
		if eip1559.OddYParity {
			dy.V = big.NewInt(1)
		}
		if eip1559.Action.Call != nil {
			to := common.HexToAddress(*eip1559.Action.Call)
			dy.To = &to
		}
		eip1559Txn := types.NewTx(dy)
		return eip1559Txn.Hash().Hex()
	}

	if txn.Legacy != nil {
		var legacy = txn.Legacy
		ly := &types.LegacyTx{
			Nonce:    uint64(legacy.Nonce.IntPart()),
			GasPrice: legacy.GasPrice.BigInt(),
			Gas:      uint64(legacy.GasLimit.IntPart()),
			Value:    legacy.Value.BigInt(),
			Data:     util.HexToBytes(legacy.Input),
			V:        legacy.Signature.V,
			R:        util.U256(legacy.Signature.R),
			S:        util.U256(legacy.Signature.S),
		}
		if legacy.Action.Call != nil {
			to := common.HexToAddress(*legacy.Action.Call)
			ly.To = &to
		}
		LegacyTransaction := types.NewTx(ly)
		return LegacyTransaction.Hash().Hex()
	}
	return ""
}

func CalHashByTxRaw(_ context.Context, rawTxHex string) string {
	rawTx, err := hex.DecodeString(util.TrimHex(rawTxHex))
	if err != nil {
		panic(err)
	}

	tx := new(types.Transaction)
	err = tx.UnmarshalBinary(rawTx)
	if err != nil {
		return ""
	}
	return tx.Hash().Hex()
}
