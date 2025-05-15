package http

import (
	"context"
	"github.com/itering/subscan/model"
	balanceModel "github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/subscan/plugins/evm/dao"
	"github.com/shopspring/decimal"
)

type MockServer struct {
}

func (m MockServer) AccountTokens(ctx context.Context, address, _ string) []dao.AccountTokenJson {
	return nil
}

func (m MockServer) Collectibles(ctx context.Context, address string, contract string, page, row int) ([]dao.Erc721Holders, int) {
	return nil, 0
}

func (m MockServer) TokenList(ctx context.Context, _, category string, page, row int) ([]dao.Token, int) {
	return nil, 0
}

func (m MockServer) TokenTransfers(ctx context.Context, address, tokenAddress string, page, row int) ([]dao.TokenTransferJson, int) {
	return nil, 0
}

func (m MockServer) TokenHolders(ctx context.Context, address string, page int, row int) ([]dao.TokenHolder, int) {
	return nil, 0
}

func (m MockServer) Blocks(ctx context.Context, page int, row int) ([]dao.EvmBlockJson, int) {
	return nil, 0
}

func (m MockServer) BlockByNum(ctx context.Context, blockNum uint) *dao.EvmBlock {
	return nil
}

func (m MockServer) BlockByHash(ctx context.Context, hash string) *dao.EvmBlock {
	return nil
}

func (m MockServer) TransactionsJson(ctx context.Context, page model.Option, opts ...model.Option) ([]dao.TransactionSampleJson, int) {
	return nil, 0
}

func (m MockServer) Accounts(ctx context.Context, address string, page int, row int) ([]dao.AccountsJson, int64) {
	return nil, 0
}

func (m MockServer) Contracts(ctx context.Context, page int, row int) ([]dao.ContractsJson, int64) {
	// TODO implement me
	panic("implement me")
}

func (m MockServer) GetTransactionByHash(_ context.Context, _ string) *dao.Transaction {
	return &dao.Transaction{Hash: "0xdf03f7309487778643a40a7fc4a8224f8c984f7f1821d970458cabc51c6a59b6", Success: true}
}

func (m MockServer) API_GetLogs(ctx context.Context, opts ...model.Option) (res []dao.EtherscanLogsRes) {
	return []dao.EtherscanLogsRes{
		{
			Address:          "0x1c3d21ac81860deaf7736fe87d664eeb788bacc1",
			Topics:           []string{"0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62", "0x0000000000000000000000008245637968c2e16e9c28d45067bf6dd4334e6db0", "0x0000000000000000000000008245637968c2e16e9c28d45067bf6dd4334e6db0", "0x000000000000000000000000ab528d46e35d05a50e32c86e2fe3a437b7adecea"},
			Data:             "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001",
			BlockNumber:      "0x9fa80b",
			BlockHash:        "0x146b18bb883b1667e7936347453a1fdadca1178090df9f8a2b34665bdae662a1",
			Timestamp:        "0x67fe0cba",
			GasPrice:         "0x870ab1a80",
			GasUsed:          "0xbaeefb",
			LogIndex:         "1",
			TransactionHash:  "0xdf03f7309487778643a40a7fc4a8224f8c984f7f1821d970458cabc51c6a59b6",
			TransactionIndex: "3",
		},
	}
}

func (m MockServer) API_GetAccounts(ctx context.Context, h160 []string) (map[string]balanceModel.Account, error) {
	return map[string]balanceModel.Account{
		"0x66b8c60c79dfad02fc04f1f13aab0f6feff8615b": {Balance: decimal.NewFromInt(1)},
		"0xe22d73f5dcccb31a994ad4e7ad265cf69b4e725a": {Balance: decimal.NewFromInt(2)},
	}, nil
}

func (m MockServer) API_Transactions(ctx context.Context, opts ...model.Option) (res []dao.EtherscanTxnRes) {
	return nil
}

func (m MockServer) API_TokenEventRes(ctx context.Context, opts ...model.Option) []dao.EtherscanTokenEventRes {
	return nil
}

func (m MockServer) API_ContractSourceCode(_ context.Context, c *dao.Contract) *dao.EtherscanContractSourceCodeRes {
	return nil
}

func (m MockServer) API_GetContractCreation(ctx context.Context, addresses []string) (res []dao.EtherscanContractCreationRes) {
	return []dao.EtherscanContractCreationRes{
		{
			ContractAddress: "0x1c3d21ac81860deaf7736fe87d664eeb788bacc1",
		},
	}
}

func (m MockServer) ContractsByAddr(ctx context.Context, address string) (contract *dao.Contract) {
	return &dao.Contract{Address: address, VerifyStatus: "perfect"}
}

func init() {
	srv = MockServer{}
}
