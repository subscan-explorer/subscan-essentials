package dao

import (
	"context"
	"fmt"
	"github.com/itering/subscan/model"
	balanceModel "github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/subscan/util"
	"github.com/shopspring/decimal"
	"strings"
)

type ISrv interface {
	API_GetLogs(ctx context.Context, opts ...model.Option) (res []EtherscanLogsRes)
	API_GetAccounts(ctx context.Context, h160 []string) (map[string]balanceModel.Account, error)
	API_Transactions(ctx context.Context, opts ...model.Option) (res []EtherscanTxnRes)
	API_TokenEventRes(ctx context.Context, opts ...model.Option) []EtherscanTokenEventRes
	API_ContractSourceCode(_ context.Context, c *Contract) *EtherscanContractSourceCodeRes
	API_GetContractCreation(ctx context.Context, addresses []string) (res []EtherscanContractCreationRes)

	ContractsByAddr(ctx context.Context, address string) (contract *Contract)
	GetTransactionByHash(c context.Context, hash string) *Transaction
	Blocks(ctx context.Context, page int, row int) ([]EvmBlockJson, int)
	BlockByNum(ctx context.Context, blockNum uint) *EvmBlock
	BlockByHash(ctx context.Context, hash string) *EvmBlock
	TransactionsJson(ctx context.Context, page model.Option, opts ...model.Option) ([]TransactionSampleJson, int)
	Accounts(ctx context.Context, adress string, page int, row int) ([]AccountsJson, int64)
	Contracts(ctx context.Context, page int, row int) ([]ContractsJson, int64)

	AccountTokens(ctx context.Context, address, category string) []AccountTokenJson
	Collectibles(ctx context.Context, address string, contract string, page, row int) ([]Erc721Holders, int)
	TokenList(ctx context.Context, contract, category string, page, row int) ([]Token, int)
	TokenTransfers(ctx context.Context, address, tokenAddress string, page, row int) ([]TokenTransferJson, int)
	TokenHolders(ctx context.Context, address string, page int, row int) ([]TokenHolder, int)
}

type ApiSrv struct{}

type EtherscanLogsRes struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	BlockHash        string   `json:"blockHash"`
	Timestamp        string   `json:"timestamp"`
	GasPrice         string   `json:"gasPrice"`
	GasUsed          string   `json:"gasUsed"`
	LogIndex         string   `json:"logIndex"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}

func (a *ApiSrv) API_GetLogs(ctx context.Context, opts ...model.Option) (res []EtherscanLogsRes) {
	var list []TransactionReceipt
	sg.db.WithContext(ctx).Scopes(opts...).Order("id desc").Find(&list)

	var (
		blockNums []uint64
		hashes    []string
	)
	for _, v := range list {
		blockNums = append(blockNums, v.BlockNum)
		hashes = append(hashes, v.TransactionHash)
	}

	var blocks []EvmBlock
	sg.db.WithContext(ctx).Select("block_num,block_hash").Model(&EvmBlock{}).Where("block_num in ?", blockNums).Find(&blocks)
	var hashesMap = make(map[uint64]string)
	for _, v := range blocks {
		hashesMap[v.BlockNum] = v.BlockHash
	}

	var txns []Transaction
	sg.db.WithContext(ctx).Select("hash,gas_price,gas_used").Model(&Transaction{}).Where("hash in ?", hashes).Find(&txns)
	var txnsMap = make(map[string]Transaction)
	for _, v := range txns {
		txnsMap[v.Hash] = v
	}

	for _, v := range list {
		res = append(res, EtherscanLogsRes{
			Address:          v.Address,
			Topics:           strings.Split(v.Topics, ","),
			Data:             v.Data,
			BlockNumber:      util.IntToHexNumber(v.BlockNum),
			BlockHash:        hashesMap[v.BlockNum],
			Timestamp:        util.IntToHexNumber(uint64(v.BlockTimestamp)),
			GasPrice:         util.IntToHexNumber(uint64(txnsMap[v.TransactionHash].GasPrice.IntPart())),
			GasUsed:          util.IntToHexNumber(uint64(txnsMap[v.TransactionHash].GasUsed.IntPart())),
			LogIndex:         util.IntToHexNumber(uint64(v.Index)),
			TransactionHash:  v.TransactionHash,
			TransactionIndex: util.IntToHexNumber(v.TransactionIndex),
		})
	}
	return
}

func (a *ApiSrv) API_GetAccounts(ctx context.Context, h160 []string) (map[string]balanceModel.Account, error) {
	var addresses []string
	var addr2H160 = make(map[string]string)

	for _, v := range h160 {
		addr := h160ToAccountIdByNetwork(ctx, v, util.NetworkNode)
		if addr == "" {
			return nil, fmt.Errorf("address %s not a valid address", v)
		}
		addresses = append(addresses, addr)
		addr2H160[addr] = v
	}
	var accounts []balanceModel.Account
	if err := sg.db.WithContext(ctx).Where("address in ? ", addresses).Find(&accounts).Error; err != nil {
		return nil, err
	}
	var accountMap = make(map[string]balanceModel.Account)
	for _, v := range accounts {
		accountMap[addr2H160[v.Address]] = v
	}
	return accountMap, nil
}

type EtherscanTxnRes struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	TxreceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Confirmations     string `json:"confirmations"`
	MethodId          string `json:"methodId"`
	FunctionName      string `json:"functionName"`
}

func (a *ApiSrv) API_Transactions(ctx context.Context, opts ...model.Option) (res []EtherscanTxnRes) {
	var list []Transaction
	sg.db.WithContext(ctx).Scopes(opts...).Find(&list)
	lastBlock := latestBlockNum(ctx)
	for _, v := range list {
		var isErr = "0"
		var txreceiptStatus = "1"
		if !v.Success {
			isErr = "1"
			txreceiptStatus = "0"
		}
		var methodId string
		// 0xa9059cbb
		if len(v.InputData) > 10 {
			methodId = v.InputData[:10]
		}
		res = append(res, EtherscanTxnRes{
			BlockNumber:       fmt.Sprintf("%d", v.BlockNum),
			TimeStamp:         fmt.Sprintf("%d", v.BlockTimestamp),
			Hash:              v.Hash,
			Nonce:             fmt.Sprintf("%d", v.Nonce),
			BlockHash:         "",
			TransactionIndex:  fmt.Sprintf("%d", v.TransactionIndex),
			From:              v.FromAddress,
			To:                v.ToAddress,
			Value:             v.Value.String(),
			Gas:               v.GasLimit.String(),
			GasPrice:          v.GasPrice.String(),
			IsError:           isErr,
			TxreceiptStatus:   txreceiptStatus,
			Input:             v.InputData,
			ContractAddress:   v.Contract,
			CumulativeGasUsed: v.CumulativeGasUsed.String(),
			GasUsed:           v.GasUsed.String(),
			Confirmations:     fmt.Sprintf("%d", lastBlock-v.BlockNum),
			MethodId:          methodId,
			FunctionName:      "", // todo
		})
	}
	return
}

type EtherscanTokenEventRes struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	From              string `json:"from"`
	ContractAddress   string `json:"contractAddress"`
	To                string `json:"to"`
	TokenName         string `json:"tokenName"`
	TokenSymbol       string `json:"tokenSymbol"`
	TokenDecimal      string `json:"tokenDecimal"`
	TransactionIndex  string `json:"transactionIndex"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	GasUsed           string `json:"gasUsed"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	Confirmations     string `json:"confirmations"`

	Value      string `json:"value,omitempty"`
	TokenValue string `json:"tokenValue,omitempty"`
	TokenID    string `json:"tokenID,omitempty"`
}

func (a *ApiSrv) API_TokenEventRes(ctx context.Context, opts ...model.Option) []EtherscanTokenEventRes {
	var transfers []TokensTransfers
	sg.db.WithContext(ctx).Scopes(opts...).Find(&transfers)
	var res []EtherscanTokenEventRes
	lastBlock := latestBlockNum(ctx)

	var (
		tokens []string
		blocks []uint64
		txs    []string
	)
	for _, v := range transfers {
		tokens = append(tokens, v.Contract)
		blocks = append(blocks, v.BlockNum())
		txs = append(txs, v.Hash)
	}

	blockNums2Blocks := BlockNums2Blocks(ctx, blocks)
	address2Tokens := ContractAddr2Token(ctx, tokens)
	hash2Txn := Hash2Transaction(ctx, txs)

	for _, transfer := range transfers {
		txn := hash2Txn[transfer.Hash]
		block := blockNums2Blocks[transfer.BlockNum()]
		token := address2Tokens[transfer.Contract]
		r := EtherscanTokenEventRes{
			BlockNumber:       fmt.Sprintf("%d", transfer.BlockNum()),
			TimeStamp:         fmt.Sprintf("%d", txn.BlockTimestamp),
			Hash:              transfer.Hash,
			Nonce:             fmt.Sprintf("%d", txn.Nonce),
			BlockHash:         block.BlockHash,
			From:              transfer.Sender,
			ContractAddress:   transfer.Contract,
			To:                transfer.Receiver,
			TokenName:         token.Name,
			TokenSymbol:       token.Symbol,
			TokenDecimal:      fmt.Sprintf("%d", token.Decimals),
			TransactionIndex:  fmt.Sprintf("%d", txn.TransactionIndex),
			Gas:               txn.GasLimit.String(),
			GasPrice:          txn.GasPrice.String(),
			GasUsed:           txn.GasUsed.String(),
			CumulativeGasUsed: txn.CumulativeGasUsed.String(),
			Confirmations:     fmt.Sprintf("%d", uint64(lastBlock)-transfer.BlockNum()),
		}
		switch transfer.Category {
		case TransferCategoryErc20:
			r.Value = transfer.Value.String()
		case TransferCategoryErc721:
			r.TokenID = transfer.TokenId
		case TransferCategoryErc1155:
			r.TokenValue = transfer.Value.String()
			r.TokenID = transfer.TokenId
		}
		res = append(res, r)
	}
	return res
}

type EtherscanContractSourceCodeRes struct {
	SourceCode           string `json:"SourceCode"`
	ABI                  string `json:"ABI"`
	ContractName         string `json:"ContractName"`
	CompilerVersion      string `json:"CompilerVersion"`
	OptimizationUsed     string `json:"OptimizationUsed"`
	Runs                 string `json:"Runs"`
	ConstructorArguments string `json:"ConstructorArguments"`
	EVMVersion           string `json:"EVMVersion"`
	Library              string `json:"Library"`
	LicenseType          string `json:"LicenseType"`
	Proxy                string `json:"Proxy"`
	Implementation       string `json:"Implementation"`
	SwarmSource          string `json:"SwarmSource"`
	SimilarMatch         string `json:"SimilarMatch"`
}

func (a *ApiSrv) API_ContractSourceCode(_ context.Context, c *Contract) *EtherscanContractSourceCodeRes {
	res := &EtherscanContractSourceCodeRes{
		SourceCode:       c.SourceCode,
		ABI:              c.Abi.String(),
		ContractName:     c.ContractName,
		CompilerVersion:  c.CompilerVersion,
		OptimizationUsed: "0",
		Runs:             fmt.Sprintf("%d", c.OptimizationRuns),
		EVMVersion:       c.EvmVersion,
		Library:          c.ExternalLibraries.String(),
		Proxy:            c.VerifyType,
		// ConstructorArguments: "",
		// LicenseType:          "",
	}
	if c.Optimize {
		res.OptimizationUsed = "1"
	}
	return res
}

type EtherscanContractCreationRes struct {
	ContractAddress string `json:"contractAddress"`
	ContractCreator string `json:"contractCreator"`
	TxHash          string `json:"txHash"`
	BlockNumber     string `json:"blockNumber"`
	Timestamp       string `json:"timestamp"`
	// ContractFactory  string `json:"contractFactory"`
	CreationBytecode string `json:"creationBytecode"`
}

func (a *ApiSrv) API_GetContractCreation(ctx context.Context, addresses []string) (res []EtherscanContractCreationRes) {
	var contracts []Contract
	sg.db.WithContext(ctx).Model(&Contract{}).Where("contract_address in ?", addresses).Find(&contracts)

	for _, v := range contracts {
		res = append(res, EtherscanContractCreationRes{
			ContractAddress: v.Address,
			ContractCreator: v.Deployer,
			TxHash:          v.TxHash,
			BlockNumber:     fmt.Sprintf("%d", v.BlockNum),
			Timestamp:       fmt.Sprintf("%d", v.DeployAt),
			// ContractFactory:  "",
			CreationBytecode: v.CreationBytecode,
		})
	}
	return
}

func (a *ApiSrv) ContractsByAddr(ctx context.Context, address string) (contract *Contract) {
	return ContractsByAddr(ctx, address)
}

func (a *ApiSrv) GetTransactionByHash(c context.Context, hash string) *Transaction {
	return GetTransactionByHash(c, hash)
}

type EvmBlockJson struct {
	BlockNum       uint   `json:"block_num"`
	Miner          string `json:"miner"`
	Transactions   int    `json:"transactions"`
	BlockTimestamp uint   `json:"block_timestamp"`
}

func (a *ApiSrv) Blocks(ctx context.Context, page int, row int) ([]EvmBlockJson, int) {
	list, count := GetBlockList(ctx, page, row)
	var res []EvmBlockJson
	for _, v := range list {
		res = append(res, EvmBlockJson{
			BlockNum:       uint(v.BlockNum),
			Miner:          v.Miner,
			Transactions:   v.TransactionCount,
			BlockTimestamp: v.Timestamp,
		})
	}
	return res, count
}

func (a *ApiSrv) BlockByNum(ctx context.Context, blockNum uint) *EvmBlock {
	return GetBlockByNum(ctx, int(blockNum))
}

func (a *ApiSrv) BlockByHash(ctx context.Context, hash string) *EvmBlock {
	return GetBlockByHash(ctx, hash)
}

type TransactionSampleJson struct {
	Hash           string          `json:"hash"`
	BlockNum       uint            `json:"block_num"`
	BlockTimestamp uint            `json:"block_timestamp"`
	FromAddress    string          `json:"from_address"`
	ToAddress      string          `json:"to_address"`
	Create         string          `json:"create"`
	Value          decimal.Decimal `json:"value"`
}

func (a *ApiSrv) TransactionsJson(ctx context.Context, page model.Option, opts ...model.Option) ([]TransactionSampleJson, int) {
	var list []Transaction
	var count int64
	sg.db.WithContext(ctx).Scopes(page).Scopes(opts...).Find(&list)
	sg.db.WithContext(ctx).Model(Transaction{}).Scopes(opts...).Count(&count)
	var res []TransactionSampleJson
	for _, v := range list {
		res = append(res, TransactionSampleJson{
			Hash:           v.Hash,
			BlockNum:       v.BlockNum,
			BlockTimestamp: v.BlockTimestamp,
			FromAddress:    v.FromAddress,
			ToAddress:      v.ToAddress,
			Value:          v.Value,
			Create:         v.Contract,
		})
	}
	return res, int(count)
}

type AccountsJson struct {
	EvmAccount string          `json:"evm_account"`
	Balance    decimal.Decimal `json:"balance"`
}

func (a *ApiSrv) Accounts(ctx context.Context, address string, page int, row int) ([]AccountsJson, int64) {
	var count int64
	q := sg.db.WithContext(ctx).Model(&Account{})
	if address != "" {
		q.Where("evm_account = ?", address)
	}
	q.Count(&count)
	if count == 0 {
		return nil, 0
	}
	var res []AccountsJson
	query := sg.db.WithContext(ctx).Select("evm_account,balance").Model(&Account{}).Joins("left join balance_accounts on evm_accounts.address=balance_accounts.address")
	if address != "" {
		query.Where("evm_account = ?", address)
	}
	query.Order("balance desc").Order("evm_account desc").Limit(row).Offset((page - 1) * row).Scan(&res)
	return res, count
}

type ContractsJson struct {
	ContractName     string `json:"contract_name"`
	Address          string `json:"address"`
	TransactionCount int    `json:"transaction_count"`
	VerifyStatus     string `json:"verify_status"`
}

func (a *ApiSrv) Contracts(ctx context.Context, page int, row int) ([]ContractsJson, int64) {
	var count int64
	sg.db.WithContext(ctx).Model(&Contract{}).Count(&count)
	if count == 0 {
		return nil, 0
	}
	var res []ContractsJson
	sg.db.WithContext(ctx).Model(&Contract{}).Select("contract_name,address,transaction_count,verify_status").Limit(row).Offset((page - 1) * row).Scan(&res)
	return res, count
}

type AccountTokenJson struct {
	Name     string          `json:"name"`
	Symbol   string          `json:"symbol"`
	Balance  decimal.Decimal `json:"balance"`
	Decimals uint            `json:"decimals"`
	Category string          `json:"category"`
	Contract string          `json:"contract"`
}

func (a *ApiSrv) AccountTokens(ctx context.Context, address, category string) []AccountTokenJson {
	var tokenHolders []AccountTokenJson

	q := sg.db.WithContext(ctx).Select("evm_token_holders.contract,balance,category,decimals,symbol,name").Model(&TokenHolder{}).
		Joins("join evm_tokens on evm_token_holders.contract=evm_tokens.contract").Where("holder = ?", address)
	if category != "" {
		q.Where("category = ?", category)
	}
	q.Scan(&tokenHolders)
	return tokenHolders
}

func (a *ApiSrv) Collectibles(ctx context.Context, address string, contract string, page, row int) ([]Erc721Holders, int) {
	var holders []Erc721Holders
	q := sg.db.WithContext(ctx).Model(&Erc721Holders{})
	if address != "" {
		q.Where("holder = ?", address)
	}
	if contract != "" {
		q.Where("contract = ?", contract)
	}
	var count int64
	q.Count(&count)
	q.Offset(page * row).Limit(row).Find(&holders)
	return holders, int(count)
}

func (a *ApiSrv) TokenList(ctx context.Context, contract, category string, page, row int) ([]Token, int) {
	var tokens []Token
	var count int64
	query := sg.db.WithContext(ctx).Model(&Token{})
	if category != "" {
		query.Where("category = ?", category)
	}
	if contract != "" {
		query.Where("contract = ?", contract)
	}
	query.Count(&count)
	query.Offset(page * row).Limit(row).Find(&tokens)
	return tokens, int(count)
}

func (a *ApiSrv) TokenTransfers(ctx context.Context, address, tokenAddress string, page, row int) ([]TokenTransferJson, int) {
	var transfers []TokensTransfers
	var count int64

	query := sg.db.WithContext(ctx).Model(&TokensTransfers{})
	if address != "" {
		query.Where("sender = ? or receiver = ?", address, address)
	}
	if tokenAddress != "" {
		query.Where("contract = ?", tokenAddress)
	}
	query.Count(&count)
	query.Offset(page * row).Limit(row).Find(&transfers)

	var res []TokenTransferJson
	var tokensAddress []string
	for _, v := range transfers {
		tokensAddress = append(tokensAddress, v.Contract)
	}
	addr2Token := ContractAddr2Token(ctx, tokensAddress)
	for index := range transfers {
		transfer := transfers[index]
		tj := TokenTransferJson{
			ID:       transfer.TransferId,
			Contract: transfer.Contract,
			Hash:     transfer.Hash,
			CreateAt: transfer.CreateAt,
			From:     transfer.Sender,
			To:       transfer.Receiver,
			Value:    &transfer.Value,
		}
		if token, ok := addr2Token[transfer.Contract]; ok {
			tj.Decimals = &token.Decimals
			tj.Symbol = token.Symbol
			tj.Name = token.Name
			tj.Category = token.Category
		}
		res = append(res, tj)
	}
	return res, int(count)
}

func (a *ApiSrv) TokenHolders(ctx context.Context, address string, page int, row int) ([]TokenHolder, int) {
	tokens, count := TokenHolders(ctx, address, page, row)
	return tokens, int(count)
}
