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
	BlocksCursor(ctx context.Context, limit int, before, after *uint) ([]EvmBlockJson, map[string]interface{})
	BlockByNum(ctx context.Context, blockNum uint) *EvmBlock
	BlockByHash(ctx context.Context, hash string) *EvmBlock
	TransactionsCursor(ctx context.Context, limit int, before, after *uint, opts ...model.Option) ([]TransactionSampleJson, map[string]interface{})
	AccountsCursor(ctx context.Context, address string, limit int, before, after *string) ([]AccountsJson, map[string]interface{})
	ContractsCursor(ctx context.Context, limit int, before, after *string) ([]ContractsJson, map[string]interface{})

	AccountTokens(ctx context.Context, address, category string) []AccountTokenJson
	CollectiblesCursor(ctx context.Context, address string, contract string, limit int, before, after *string) ([]Erc721Holders, map[string]interface{})
	TokenListCursor(ctx context.Context, contract, category string, limit int, before, after *string) ([]Token, map[string]interface{})
	TokenTransfersCursor(ctx context.Context, address, tokenAddress, category string, limit int, before, after *uint) ([]TokenTransferJson, map[string]interface{})
	TokenHoldersCursor(ctx context.Context, address string, limit int, before, after *string) ([]TokenHolder, map[string]interface{})
}

type IPagination interface {
	Cursor() string
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

func (a *ApiSrv) Blocks(ctx context.Context, page int, row int) ([]EvmBlockJson, int) { return nil, 0 }

func (a *ApiSrv) BlocksCursor(ctx context.Context, limit int, before, after *uint) ([]EvmBlockJson, map[string]interface{}) {
	var blocks []EvmBlock
	fetch := limit + 1
	q := sg.db.WithContext(ctx).Model(&EvmBlock{})
	if after != nil && *after > 0 {
		q = q.Where("block_num < ?", *after).Order("block_num desc")
	} else if before != nil && *before > 0 {
		q = q.Where("block_num > ?", *before).Order("block_num asc")
	} else {
		q = q.Order("block_num desc")
	}
	q = q.Limit(fetch).Find(&blocks)
	if q.Error != nil {
		return nil, nil
	}
	var hasPrev, hasNext bool
	if before != nil && *before > 0 {
		hasPrev = len(blocks) > limit
		if hasPrev {
			blocks = blocks[:limit]
		}
		for i, j := 0, len(blocks)-1; i < j; i, j = i+1, j-1 {
			blocks[i], blocks[j] = blocks[j], blocks[i]
		}
		hasNext = true
	} else {
		hasNext = len(blocks) > limit
		if hasNext {
			blocks = blocks[:limit]
		}
		hasPrev = after != nil && *after > 0
	}
	var res []EvmBlockJson
	for _, v := range blocks {
		res = append(res, EvmBlockJson{BlockNum: uint(v.BlockNum), Miner: v.Miner, Transactions: v.TransactionCount, BlockTimestamp: v.Timestamp})
	}
	var start, end *uint
	if len(blocks) > 0 {
		s := uint(blocks[0].BlockNum)
		e := uint(blocks[len(blocks)-1].BlockNum)
		start = &s
		end = &e
	}
	return res, map[string]interface{}{"start_cursor": start, "end_cursor": end, "has_previous_page": hasPrev, "has_next_page": hasNext}
}

func (a *ApiSrv) BlockByNum(ctx context.Context, blockNum uint) *EvmBlock {
	return GetBlockByNum(ctx, int(blockNum))
}

func (a *ApiSrv) BlockByHash(ctx context.Context, hash string) *EvmBlock {
	return GetBlockByHash(ctx, hash)
}

type TransactionSampleJson struct {
	TransactionId  uint64          `json:"transaction_id"`
	Hash           string          `json:"hash"`
	BlockNum       uint            `json:"block_num"`
	BlockTimestamp uint            `json:"block_timestamp"`
	FromAddress    string          `json:"from_address"`
	ToAddress      string          `json:"to_address"`
	Create         string          `json:"create"`
	Value          decimal.Decimal `json:"value"`
}

func (a *ApiSrv) TransactionsCursor(ctx context.Context, limit int, before, after *uint, opts ...model.Option) ([]TransactionSampleJson, map[string]interface{}) {
	var txs []Transaction
	fetch := limit + 1
	q := sg.db.WithContext(ctx).Model(Transaction{}).Scopes(opts...)
	if after != nil && *after > 0 {
		q = q.Where("transaction_id < ?", *after).Order("transaction_id desc")
	} else if before != nil && *before > 0 {
		q = q.Where("transaction_id > ?", *before).Order("transaction_id asc")
	} else {
		q = q.Order("transaction_id desc")
	}
	q = q.Limit(fetch).Find(&txs)
	if q.Error != nil {
		return nil, nil
	}
	var hasPrev, hasNext bool
	if before != nil && *before > 0 {
		hasPrev = len(txs) > limit
		if hasPrev {
			txs = txs[:limit]
		}
		for i, j := 0, len(txs)-1; i < j; i, j = i+1, j-1 {
			txs[i], txs[j] = txs[j], txs[i]
		}
		hasNext = true
	} else {
		hasNext = len(txs) > limit
		if hasNext {
			txs = txs[:limit]
		}
		hasPrev = after != nil && *after > 0
	}
	var res []TransactionSampleJson
	for _, v := range txs {
		res = append(res, TransactionSampleJson{Hash: v.Hash, BlockNum: v.BlockNum, BlockTimestamp: v.BlockTimestamp, FromAddress: v.FromAddress, ToAddress: v.ToAddress, Value: v.Value, Create: v.Contract, TransactionId: v.TransactionId})
	}
	var start, end *uint
	if len(txs) > 0 {
		s := uint(txs[0].TransactionId)
		e := uint(txs[len(txs)-1].TransactionId)
		start = &s
		end = &e
	}
	return res, map[string]interface{}{"start_cursor": start, "end_cursor": end, "has_previous_page": hasPrev, "has_next_page": hasNext}
}

type AccountsJson struct {
	EvmAccount string          `json:"evm_account"`
	Balance    decimal.Decimal `json:"balance"`
}

func (a AccountsJson) Cursor() string {
	return util.Base64Encode(fmt.Sprintf("%s_%s", a.Balance.String(), a.EvmAccount))
}

func (a *ApiSrv) AccountsCursor(ctx context.Context, address string, limit int, before, after *string) ([]AccountsJson, map[string]interface{}) {
	var list []AccountsJson
	fetch := limit + 1
	q := sg.db.WithContext(ctx).Select("evm_account,balance").Model(&Account{}).Joins("join balance_accounts on evm_accounts.address=balance_accounts.address")
	if address != "" {
		q.Where("evm_account = ?", address)
	}
	if cursor := cursorDecode(after); len(cursor) == 2 {
		q = q.Where("(balance,evm_account) < (?,?)", cursor[0], cursor[1]).Order("balance desc").Order("balance_accounts.address desc")
	} else if cursor = cursorDecode(after); len(cursor) == 2 {
		q = q.Where("(balance,evm_account) < (?,?)", cursor[0], cursor[1]).Order("balance asc").Order("balance_accounts.address asc")
	} else {
		q = q.Order("balance desc").Order("balance_accounts.address desc")
	}
	q.Limit(fetch).Scan(&list)
	var hasPrev, hasNext bool
	if before != nil && *before != "" {
		hasPrev = len(list) > limit
		if hasPrev {
			list = list[:limit]
		}
		for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
		hasNext = true
	} else {
		hasNext = len(list) > limit
		if hasNext {
			list = list[:limit]
		}
		hasPrev = after != nil && *after != ""
	}
	var start, end *string
	if len(list) > 0 {
		s := list[0].Cursor()
		e := list[len(list)-1].Cursor()
		start = &s
		end = &e
	}
	return list, map[string]interface{}{"start_cursor": start, "end_cursor": end, "has_previous_page": hasPrev, "has_next_page": hasNext}
}

type ContractsJson struct {
	ContractName     string `json:"contract_name"`
	Address          string `json:"address"`
	TransactionCount int    `json:"transaction_count"`
	VerifyStatus     string `json:"verify_status"`
}

func (c ContractsJson) Cursor() string {
	return util.Base64Encode(fmt.Sprintf("%d_%s", c.TransactionCount, c.Address))
}

func (a *ApiSrv) ContractsCursor(ctx context.Context, limit int, before, after *string) ([]ContractsJson, map[string]interface{}) {
	var list []ContractsJson
	fetch := limit + 1
	q := sg.db.WithContext(ctx).Model(&Contract{}).Select("contract_name,address,transaction_count,verify_status")
	if cursor := cursorDecode(after); len(cursor) == 2 {
		q = q.Where("(transaction_count,address) < (?,?)", cursor[0], cursor[1]).Order("transaction_count desc").Order("address desc")
	} else if cursor = cursorDecode(before); len(cursor) == 2 {
		q = q.Where("(transaction_count,address) > (?,?)", cursor[0], cursor[1]).Order("transaction_count asc").Order("address asc")
	} else {
		q = q.Order("transaction_count desc").Order("address desc")
	}
	q = q.Limit(fetch).Scan(&list)
	if q.Error != nil {
		return nil, nil
	}
	var hasPrev, hasNext bool
	if before != nil && *before != "" {
		hasPrev = len(list) > limit
		if hasPrev {
			list = list[:limit]
		}
		for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
		hasNext = true
	} else {
		hasNext = len(list) > limit
		if hasNext {
			list = list[:limit]
		}
		hasPrev = after != nil && *after != ""
	}
	var start, end *string
	if len(list) > 0 {
		s := list[0].Cursor()
		e := list[len(list)-1].Cursor()
		start = &s
		end = &e
	}
	return list, map[string]interface{}{"start_cursor": start, "end_cursor": end, "has_previous_page": hasPrev, "has_next_page": hasNext}
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

var cursorDecode = func(c *string) []string {
	// decode base64
	if c == nil || *c == "" {
		return nil
	}
	decoded := util.Base64Decode(*c)
	if decoded == "" {
		return nil
	}
	parts := strings.SplitN(decoded, "_", 2)
	if len(parts) != 2 {
		return nil
	}
	return []string{parts[0], parts[1]}
}

func (a *ApiSrv) CollectiblesCursor(ctx context.Context, address string, contract string, limit int, before, after *string) ([]Erc721Holders, map[string]interface{}) {
	var list []Erc721Holders
	fetch := limit + 1
	q := sg.db.WithContext(ctx).Model(&Erc721Holders{})
	if address != "" {
		q.Where("holder = ?", address)
	}
	if contract != "" {
		q.Where("contract = ?", contract)
	}
	if cursor := cursorDecode(after); len(cursor) == 2 {
		q = q.Where("(contract,token_id) < (?,?)", cursor[0], cursor[1]).Order("contract desc").Order("token_id desc")
	} else if cursor = cursorDecode(before); len(cursor) == 2 {
		q = q.Where("(contract,token_id) > (?,?)", cursor[0], cursor[1]).Order("contract asc").Order("token_id asc")
	} else {
		q = q.Order("contract desc").Order("token_id desc")
	}
	q = q.Limit(fetch).Find(&list)
	if q.Error != nil {
		return nil, nil
	}
	var hasPrev, hasNext bool
	if before != nil && *before != "" {
		hasPrev = len(list) > limit
		if hasPrev {
			list = list[:limit]
		}
		for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
		hasNext = true
	} else {
		hasNext = len(list) > limit
		if hasNext {
			list = list[:limit]
		}
		hasPrev = after != nil && *after != ""
	}
	var start, end *string
	if len(list) > 0 {
		s := list[0].Cursor()
		e := list[len(list)-1].Cursor()
		start = &s
		end = &e
	}
	return list, map[string]interface{}{
		"start_cursor":      start,
		"end_cursor":        end,
		"has_previous_page": hasPrev,
		"has_next_page":     hasNext,
	}
}

func (a *ApiSrv) TokenListCursor(ctx context.Context, contract, category string, limit int, before, after *string) ([]Token, map[string]interface{}) {
	var list []Token
	fetch := limit + 1
	q := sg.db.WithContext(ctx).Model(&Token{})
	if category != "" {
		q.Where("category = ?", category)
	}
	if contract != "" {
		q.Where("contract = ?", contract)
	}
	if cursor := cursorDecode(after); len(cursor) == 2 {
		q = q.Where("(holders,contract) < (?,?)", cursor[0], cursor[1]).Order("holders desc").Order("contract desc")
	} else if cursor = cursorDecode(before); len(cursor) == 2 {
		q = q.Where("(holders,contract) > (?,?)", cursor[0], cursor[1]).Order("holders asc").Order("contract asc")
	} else {
		q = q.Order("holders desc").Order("contract desc")
	}
	q = q.Limit(fetch).Find(&list)
	if q.Error != nil {
		return nil, nil
	}
	var hasPrev, hasNext bool
	if before != nil && *before != "" {
		hasPrev = len(list) > limit
		if hasPrev {
			list = list[:limit]
		}
		for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
		hasNext = true
	} else {
		hasNext = len(list) > limit
		if hasNext {
			list = list[:limit]
		}
		hasPrev = after != nil && *after != ""
	}
	var start, end *string
	if len(list) > 0 {
		s := list[0].Cursor()
		e := list[len(list)-1].Cursor()
		start = &s
		end = &e
	}
	return list, map[string]interface{}{
		"start_cursor":      start,
		"end_cursor":        end,
		"has_previous_page": hasPrev,
		"has_next_page":     hasNext,
	}
}

func (a *ApiSrv) TokenTransfersCursor(ctx context.Context, address, tokenAddress, category string, limit int, before, after *uint) ([]TokenTransferJson, map[string]interface{}) {
	var transfers []TokensTransfers
	fetch := limit + 1
	q := sg.db.WithContext(ctx).Model(&TokensTransfers{})
	if address != "" {
		q.Where("sender = ? or receiver = ?", address, address)
	}
	if tokenAddress != "" {
		q.Where("contract = ?", tokenAddress)
	}
	if category != "" {
		q.Where("category = ?", category)
	}
	if after != nil && *after > 0 {
		q = q.Where("transfer_id < ?", *after).Order("transfer_id desc")
	} else if before != nil && *before > 0 {
		q = q.Where("transfer_id > ?", *before).Order("transfer_id asc")
	} else {
		q = q.Order("transfer_id desc")
	}
	q = q.Limit(fetch).Find(&transfers)
	if q.Error != nil {
		return nil, nil
	}
	var hasPrev, hasNext bool
	if before != nil && *before > 0 {
		hasPrev = len(transfers) > limit
		if hasPrev {
			transfers = transfers[:limit]
		}
		for i, j := 0, len(transfers)-1; i < j; i, j = i+1, j-1 {
			transfers[i], transfers[j] = transfers[j], transfers[i]
		}
		hasNext = true
	} else {
		hasNext = len(transfers) > limit
		if hasNext {
			transfers = transfers[:limit]
		}
		hasPrev = after != nil && *after > 0
	}
	var res []TokenTransferJson
	var tokensAddress []string
	for _, v := range transfers {
		tokensAddress = append(tokensAddress, v.Contract)
	}
	addr2Token := ContractAddr2Token(ctx, tokensAddress)
	for index := range transfers {
		transfer := transfers[index]
		tj := TokenTransferJson{ID: transfer.TransferId, Contract: transfer.Contract, Hash: transfer.Hash, CreateAt: transfer.CreateAt, From: transfer.Sender, To: transfer.Receiver, Value: &transfer.Value}
		if token, ok := addr2Token[transfer.Contract]; ok {
			tj.Decimals = &token.Decimals
			tj.Symbol = token.Symbol
			tj.Name = token.Name
			tj.Category = token.Category
		}
		res = append(res, tj)
	}
	var start, end *uint64
	if len(transfers) > 0 {
		s := transfers[0].TransferId
		e := transfers[len(transfers)-1].TransferId
		start = &s
		end = &e
	}
	return res, map[string]interface{}{"start_cursor": start, "end_cursor": end, "has_previous_page": hasPrev, "has_next_page": hasNext}
}

func (a *ApiSrv) TokenHoldersCursor(ctx context.Context, address string, limit int, before, after *string) ([]TokenHolder, map[string]interface{}) {
	var list []TokenHolder
	fetch := limit + 1
	q := sg.db.WithContext(ctx).Model(&TokenHolder{}).Where("balance > 0")
	q.Where("contract = ?", address)
	if cursor := cursorDecode(after); len(cursor) == 2 {
		q = q.Where("(balance,id) < (?,?)", cursor[0], cursor[1]).Order("balance desc").Order("id desc")
	} else if cursor = cursorDecode(before); len(cursor) == 2 {
		q = q.Where("(balance,id) > (?,?)", cursor[0], cursor[1]).Order("balance asc").Order("id asc")
	} else {
		q = q.Order("balance desc").Order("id desc")
	}
	q = q.Limit(fetch).Find(&list)
	if q.Error != nil {
		return nil, nil
	}
	var hasPrev, hasNext bool
	if before != nil && *before != "" {
		hasPrev = len(list) > limit
		if hasPrev {
			list = list[:limit]
		}
		for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
		hasNext = true
	} else {
		hasNext = len(list) > limit
		if hasNext {
			list = list[:limit]
		}
		hasPrev = after != nil && *after != ""
	}
	var start, end *string
	if len(list) > 0 {
		s := list[0].Cursor()
		e := list[len(list)-1].Cursor()
		start = &s
		end = &e
	}
	return list, map[string]interface{}{"start_cursor": start, "end_cursor": end, "has_previous_page": hasPrev, "has_next_page": hasNext}
}
