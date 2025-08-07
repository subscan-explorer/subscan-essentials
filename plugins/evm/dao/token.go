package dao

import (
	"context"
	"fmt"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins/evm/feature"
	"github.com/itering/subscan/plugins/evm/feature/erc721"
	"github.com/itering/subscan/share/web3"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/mq"
	"github.com/itering/subscan/util/parallel"
	"strings"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Token struct {
	Contract    string          `json:"contract" gorm:"primaryKey;auto_increment:false;size:255"`
	Name        string          `json:"name" gorm:"size:255"`
	Symbol      string          `json:"symbol" gorm:"size:255"`
	Decimals    uint            `json:"decimals" gorm:"size:32"`
	TotalSupply decimal.Decimal `json:"totalSupply" gorm:"default: 0;type:decimal(65);"`

	Holders       uint `json:"holders" gorm:"size:32"`
	TransferCount uint `json:"transfer_count" gorm:"size:32"`

	Category     string `json:"category" gorm:"size:70;index:category"`
	BaseTokenUri string `json:"base_token_uri" gorm:"size:255"`
}

func (c *Token) TableName() string {
	return "evm_tokens"
}

func (c *Token) incrTransferCount(ctx context.Context, count int) {
	sg.db.WithContext(ctx).Model(c).Updates(map[string]interface{}{"transfer_count": gorm.Expr(fmt.Sprintf("transfer_count + %d", count))})
}

func (c *Token) RefreshTokenHolder(ctx context.Context, holderCount int) {
	sg.db.Model(Token{}).WithContext(ctx).Where("contract = ?", c.Contract).Updates(map[string]interface{}{"holders": holderCount})
}

func TouchToken(ctx context.Context, address, category string) *Token {
	var token Token
	sg.db.Scopes(model.IgnoreDuplicate).WithContext(ctx).FirstOrCreate(&token, Token{Contract: address, Category: category})
	return &token
}

func GetTokenByContract(ctx context.Context, contract string) *Token {
	var token Token
	if q := sg.db.Model(Token{}).WithContext(ctx).Where("contract = ?", contract).First(&token); q.Error != nil {
		return nil
	}
	return &token
}

func TokenCount(ctx context.Context, category string) (count int64) {
	sg.db.WithContext(ctx).Model(Token{}).Where("category = ?", category).Count(&count)
	return
}

func TokensByContractAddr(ctx context.Context, page, row int, contracts []string, category, orderField, order string) (tokens []Token, count int64) {
	query := sg.db.WithContext(ctx).Model(Token{})
	if len(contracts) > 0 {
		query.Where("contract in (?)", contracts)
	}

	if category != "" {
		if category == "nft" {
			query.Where("category in ?", []string{Eip721Token, Eip1155Token})
		} else {
			query.Where("category = ?", category)
		}
	}
	query.Count(&count)

	if count == 0 {
		return
	}

	query.Offset(page * row).Limit(row)
	if orderField != "" {
		query.Order(fmt.Sprintf("%s %s", orderField, order))
	} else {
		query.Order("holders desc")
	}
	query.Find(&tokens)
	return
}

func (c *Token) AfterCreate(txn *gorm.DB) (err error) {
	// check erc20/erc721 Name,Symbol,Decimals,TotalSupply
	// check contract, account
	txn.Scopes(model.IgnoreDuplicate).Create(&Contract{Address: c.Contract})
	return c.mergeTokenInfo(txn)
}

func (c *Token) mergeTokenInfo(txn *gorm.DB) (err error) {
	token := feature.InitToken(web3.RPC, c.Category, c.Contract)
	parallel.FuncE(func(ctx context.Context) error {
		symbol, err := token.Symbol(ctx)
		if err == nil {
			c.Symbol = symbol
		}
		return err
	}).FuncE(func(ctx context.Context) error {
		name, err := token.Name(ctx)
		if err == nil {
			c.Name = name
		}
		return err
	}).FuncE(func(ctx context.Context) error {
		deciamls, err := token.Decimals(ctx)
		if err == nil {
			c.Decimals = deciamls
		}
		return err
	}).FuncE(func(ctx context.Context) error {
		totalSupply, err := token.TotalSupply(ctx)
		if err == nil {
			c.TotalSupply = totalSupply
		}
		return err
	}).FuncE(func(ctx context.Context) error {
		if c.Category == Eip721Token {
			ERC721 := token.(*erc721.Token)
			tokenUri, err := ERC721.BaseTokenURI(ctx)
			if err == nil {
				c.BaseTokenUri = tokenUri
			}
			return err
		}
		return nil
	}).ErrHandle(func(errs []error) {
		err = errs[0]
		util.Logger().Error(fmt.Errorf("merge token %s info error %s", c.Contract, err.Error()))
	}).Start(txn.Statement.Context).Wait()
	txn.Model(Token{}).Where("contract = ?", c.Contract).Updates(c)
	return
}

type TokenHolder struct {
	ID       uint            `gorm:"primaryKey;autoIncrement;size:32" json:"-"`
	Contract string          `json:"contract" gorm:"index:contract;index:contract_hold,unique;size:100"`
	Holder   string          `json:"holder" gorm:"index:hold;index:contract_hold,unique;size:100" `
	Balance  decimal.Decimal `json:"balance" gorm:"default: 0;type:decimal(65);"`
}

func (c *TokenHolder) TableName() string {
	return "evm_token_holders"
}

func TokenHolderCount(ctx context.Context, contract string) int {
	var count int64
	sg.db.Model(TokenHolder{}).WithContext(ctx).Where("contract = ?", contract).Where("balance > 0").Count(&count)
	return int(count)
}

func RefreshHolder(ctx context.Context, contract, address, category string) error {
	t := GetTokenByContract(ctx, contract)
	if t == nil {
		return nil
	}

	token := feature.InitToken(web3.RPC, category, contract)
	balance, err := token.BalanceOf(ctx, address)
	if err != nil {
		return err
	}

	if balance.IsNegative() {
		return nil
	}
	// erc1155 balance
	// if category == Eip1155Token {
	// 	balance = decimal.New(ERC1155TokenIdsCount(ctx, contract, address), 0)
	// }

	q := sg.AddOrUpdateItem(ctx, &TokenHolder{Contract: contract, Holder: address, Balance: balance}, []string{"contract", "holder"}, "balance")
	if q.RowsAffected == 1 || (q.RowsAffected == 2 && balance.IsZero()) {
		_ = TouchAccount(context.Background(), address)
		t.RefreshTokenHolder(ctx, TokenHolderCount(ctx, contract))
	}
	return q.Error
}

type TokensTransfers struct {
	Id uint `json:"id" gorm:"primaryKey;autoIncrement;size:32"`

	Contract string `gorm:"not null;default:'';size:100;index:contract" json:"contract"`
	Hash     string `gorm:"not null;default:'';index:hash;size:100" json:"hash"`
	CreateAt uint   `gorm:"not null;default:0;size:32;index:create_at"  json:"create_at"`

	Sender   string `gorm:"not null;default:0;size:70"  json:"sender"`
	Receiver string `gorm:"not null;default:0;size:70"  json:"receiver"`

	Value      decimal.Decimal `gorm:"default:0;type:decimal(65)" json:"value"`
	TokenId    string          `gorm:"default:null;size:255;index:token_id,length:50" json:"token_id"`
	TransferId uint64          `gorm:"index:transfer_id;index:batch,unique;size:64" json:"transfer_id"`
	BatchIndex uint            `gorm:"default:0;size:32;index:batch,unique" json:"batch_index"`
	Category   int             `gorm:"default:0;size:32;index:category" json:"category"`
}

func (t *TokensTransfers) BlockNum() uint64 {
	return t.TransferId / TransactionIdGenerateCoefficient / TxnReceiptLimit
}

const (
	TransferCategoryErc20 = iota
	TransferCategoryErc721
	TransferCategoryErc1155
)

func (t *TokensTransfers) TableName() string {
	return "evm_tokens_transfers"
}

func (t *TokensTransfers) AfterCreate(txn *gorm.DB) (err error) {
	ctx := txn.Statement.Context

	// burn or mint will update TotalSupply
	if t.Sender == NullAddress || t.Receiver == NullAddress {
		var supply decimal.Decimal
		if t.TokenId == "" {
			token := feature.InitToken(web3.RPC, Eip20Token, t.Contract)
			supply, err = token.TotalSupply(ctx)
			if err != nil {
				return nil
			}
			txn.Model(Token{}).Where("contract =?", t.Contract).Updates(Token{TotalSupply: supply})
		}
	}
	return nil
}

func (t *TransactionReceipt) ProcessTokenTransfer(ctx context.Context, category string) error {
	token := TouchToken(ctx, t.Address, category)
	topics := strings.Split(t.Topics, ",")
	if len(topics) < 3 {
		return nil
	}
	transfer := TokensTransfers{
		TransferId: t.Id,
		Contract:   token.Contract,
		Hash:       t.TransactionHash,
		CreateAt:   t.BlockTimestamp,
		Sender:     RemoveAddressPadded(topics[1]),
		Receiver:   RemoveAddressPadded(topics[2]),
		Value:      decimal.NewFromInt32(1),
	}

	switch category {
	case Eip20Token:
		transfer.Value = util.EvmU256Decoder(t.Data)
		transfer.Category = TransferCategoryErc20
	case Eip721Token:
		var tokenIdRaw = t.Data // above token_id not indexed
		if len(topics) == 4 {
			tokenIdRaw = topics[3]
		}
		transfer.TokenId = util.U256(tokenIdRaw).String()
		transfer.Category = TransferCategoryErc721
	default:
		// unsupported category
		return nil
	}

	query := sg.db.Scopes(model.IgnoreDuplicate).Create(&transfer)
	if query.RowsAffected > 0 {

		token.incrTransferCount(ctx, 1)
		if mq.Instant != nil {
			_ = Publish(category, "balance", []string{token.Contract, transfer.Sender})
			_ = Publish(category, "balance", []string{token.Contract, transfer.Receiver})
			if transfer.TokenId != "" {
				_ = Publish(category, "holder", []string{token.Contract, transfer.TokenId})
			}
		}
		return nil
	}
	return query.Error
}

func AccountTokens(ctx context.Context, holder string) (tokens []TokenHolder) {
	sg.db.WithContext(ctx).Model(TokenHolder{}).Where("holder = ?", holder).Where("balance > ?", 0).Find(&tokens)
	return
}

func TokenHolders(ctx context.Context, contract string, page, row int) (tokens []TokenHolder, count int64) {
	query := sg.db.WithContext(ctx).Model(TokenHolder{}).Where("contract = ?", contract).Where("balance > ?", 0)
	query.Count(&count)
	if count == 0 {
		return
	}
	query.Order("balance desc").Offset(page * row).Limit(row).Find(&tokens)
	return
}

type TokenTransferJson struct {
	ID         uint64           `json:"id"`
	Contract   string           `json:"contract"`
	Hash       string           `json:"hash"`
	CreateAt   uint             `json:"create_at"`
	From       string           `json:"from"`
	To         string           `json:"to" `
	Value      *decimal.Decimal `json:"value,omitempty"`
	TokenId    *string          `json:"token_id,omitempty"`
	StorageUrl *string          `json:"storage_url,omitempty"`
	Decimals   *uint            `json:"decimals,omitempty"`
	Symbol     string           `json:"symbol"`
	Name       string           `json:"name"`
	Category   string           `json:"category"`
}

func ContractAddr2Token(ctx context.Context, addr []string) map[string]Token {
	var tokens []Token
	sg.db.WithContext(ctx).Model(Token{}).Where("contract in ?", addr).Find(&tokens)
	var tokenMap = make(map[string]Token)
	for _, v := range tokens {
		tokenMap[v.Contract] = v
	}
	return tokenMap
}
