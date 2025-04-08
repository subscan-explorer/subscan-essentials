package dao

import (
	"context"
	"crypto/sha1"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itering/subscan/plugins/evm/feature/erc1155"
	"github.com/itering/subscan/plugins/evm/feature/erc721"
	"github.com/itering/subscan/share/web3"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/ipfs"
	"github.com/shopspring/decimal"
	"strings"

	"gorm.io/gorm"
)

type Erc721Holders struct {
	Id         string   `json:"-" gorm:"primaryKey;size:100"`
	Contract   string   `json:"contract" gorm:"index:contract;index:contract_hold;index:contract_token_id;size:100"`
	Holder     string   `json:"holder" gorm:"index:hold;index:contract_hold;size:100" `
	TokenId    string   `json:"token_id" gorm:"default: 0;size:255;index:contract_token_id"`
	Metadata   Metadata `json:"metadata" gorm:"type:json"`
	StorageUrl string   `json:"storage_url" gorm:"type:text"`
}

func (t *Erc721Holders) TableName() string {
	return "evm_erc721_holders"
}

type CollectiblesJson struct {
	Contract   string           `json:"contract"`
	Holder     string           `json:"holder,omitempty"`
	TokenId    string           `json:"token_id" `
	StorageUrl string           `json:"storage_url"`
	Holders    uint             `json:"holders"`
	Balance    *decimal.Decimal `json:"balance,omitempty"`
}

func (j Metadata) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *Metadata) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		if len(v) == 0 {
			return nil
		}
	case string:
		if v == "" {
			return nil
		}
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	var result Metadata
	err := util.UnmarshalAny(&result, value)
	*j = result
	return err
}

type Metadata struct {
	Description     string                   `json:"description,omitempty"`
	ExternalUrl     string                   `json:"external_url,omitempty" validate:"omitempty,url"`
	Image           string                   `json:"image,omitempty" validate:"omitempty,url"`
	Name            string                   `json:"name,omitempty"`
	AnimationUrl    string                   `json:"animation_url,omitempty" validate:"omitempty,url"`
	YoutubeUrl      string                   `json:"youtube_url,omitempty" validate:"omitempty,url"`
	BackgroundColor string                   `json:"background_color,omitempty" validate:"omitempty,hexcolor"`
	Attributes      []map[string]interface{} `json:"attributes,omitempty"`
}

func (c *Token) RefreshErc721Holders(ctx context.Context, tokenId string) error {
	if tokenId == "" {
		return nil
	}
	to, err := erc721.Init(web3.RPC, c.Contract).OwnerOf(ctx, tokenId)
	if err != nil {
		return err
	}
	// find db
	collectible := GetCollectible(ctx, c.Contract, tokenId, "")
	db := sg.db
	// nft has burned
	if to == "" {
		if collectible != nil {
			if q := db.Where("contract = ?", c.Contract).Where("token_id =?", tokenId).Delete(Erc721Holders{}); q.RowsAffected > 0 {
				db.Model(Token{}).Where("contract =?", c.Contract).UpdateColumns(map[string]interface{}{"total_supply": gorm.Expr("total_supply - ?", 1)})
			}
		}
		return nil
	}
	id := fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%s%s%s", c.Contract, to, tokenId))))
	defer func() {
		_ = c.refreshGetMetadata(context.Background(), c.Contract, tokenId)
	}()
	if collectible == nil {
		// refresh nft metadata
		if q := db.Create(&Erc721Holders{Id: id, Contract: c.Contract, Holder: to, TokenId: tokenId}); q.RowsAffected > 0 {
			db.Model(Token{}).Where("contract =?", c.Contract).UpdateColumns(map[string]interface{}{"total_supply": gorm.Expr("total_supply + ?", 1)})
		}
	} else {
		db.Model(Erc721Holders{}).Where("contract = ? and token_id = ?", c.Contract, tokenId).UpdateColumns(Erc721Holders{
			Id:     id,
			Holder: to,
		})
	}
	return nil
}

func GetCollectible(ctx context.Context, contract, tokenId string, cols string) *Erc721Holders {
	var holder Erc721Holders
	query := sg.db
	if cols != "" {
		query = query.Select(cols)
	}
	if query = query.Where("contract = ?", contract).Where("token_id =?", tokenId).First(&holder); query.Error != nil {
		return nil
	}
	return &holder
}

func Erc721Collectibles(c context.Context, page, row int, address, contract, tokenId string) (collectibles []CollectiblesJson, count int64) {
	var holders []Erc721Holders
	query := sg.db.WithContext(c).Model(Erc721Holders{})
	if contract != "" {
		query.Where("contract = ?", contract)
	}
	if address != "" {
		query.Where("holder = ?", address)
	} else {
		query.Where("holder != ''")
	}
	if tokenId != "" {
		query.Where("token_id = ?", tokenId)
	}
	query.Count(&count)
	query.Offset(page * row).Limit(row).Find(&holders)

	for _, holder := range holders {
		collectibles = append(collectibles, CollectiblesJson{
			Contract:   holder.Contract,
			Holder:     holder.Holder,
			TokenId:    holder.TokenId,
			StorageUrl: holder.StorageUrl,
		})
	}
	return
}

// func collectiblesCount(ctx context.Context, contract string) (count int64) {
// 	sg.db.WithContext(ctx).Model(Erc721Holders{}).Where("contract = ?", contract).Where("holder != ''").Count(&count)
// 	return
// }

func ERC721TokenIdsCount(ctx context.Context, contract, accountId string) int64 {
	var count int64
	sg.db.WithContext(ctx).Model(Erc721Holders{}).Where("contract = ?", contract).Where("holder = ?", accountId).Count(&count)
	return count
}

func Collectible(ctx context.Context, contract, tokenId string) *Erc721Holders {
	var nft Erc721Holders
	if query := sg.db.WithContext(ctx).Model(Erc721Holders{}).Where("contract = ?", contract).Where("token_id = ?", tokenId).First(&nft); query.Error != nil {
		return nil
	}
	return &nft
}

func RefreshBaseTokenUri(ctx context.Context, contract string) (string, error) {
	token := GetTokenByContract(ctx, contract)
	if token == nil || token.Category != Eip721Token {
		return "", nil
	}
	ERC721 := erc721.Init(web3.RPC, contract)
	uri, err := ERC721.BaseTokenURI(ctx)
	if err != nil {
		return "", err
	}
	sg.db.Model(Token{}).Where("contract = ?", contract).Update("base_token_uri", uri)
	return uri, nil
}

func (c *Token) GetMetadata(ctx context.Context, contract, tokenId string) (*Metadata, string, error) {
	var (
		tokenUrl string
		err      error
	)
	if c.BaseTokenUri != "" {
		tokenUrl = fmt.Sprintf("%s%s", c.BaseTokenUri, tokenId)
	} else {
		// get token uri
		if c.Category == Eip721Token {
			ERC721 := erc721.Init(web3.RPC, contract)
			tokenUrl, err = ERC721.TokenURI(ctx, tokenId)
			if tokenUrl == "" {
				return nil, "", err
			}
		} else if c.Category == Eip1155Token {
			ERC1155 := erc1155.Init(web3.RPC, contract)
			tokenUrl, err = ERC1155.Uri(ctx, tokenId)
			if tokenUrl == "" {
				return nil, "", err
			}
			if strings.Contains(tokenUrl, "{id}") {
				tokenUrl = strings.ReplaceAll(tokenUrl, "{id}", tokenId)
			}
		} else {
			return nil, "", fmt.Errorf("unsupported token category %s", c.Category)
		}
	}

	var (
		metadata Metadata
		data     []byte
	)
	if strings.HasPrefix("/", tokenUrl) {
		tokenUrl = tokenUrl[1:]
	}
	// ipfs,ar,http,base64 metadata json
	switch {
	case strings.HasPrefix(tokenUrl, "ipfs://") || strings.HasPrefix(tokenUrl, "https://ipfs.io"):
		data, err = ipfs.OpenFile(ctx, ipfs.TrimMetadataUri(tokenUrl))
	case strings.HasPrefix(tokenUrl, "ar://"):
		data, err = ipfs.OpenArFile(ctx, ipfs.TrimMetadataUri(tokenUrl))
		// ignore localhost/127.0.0.1
	case strings.HasPrefix(tokenUrl, "http://localhost") || strings.HasPrefix(tokenUrl, "https://localhost") || strings.Contains(tokenUrl, "127.0.0.1"):
		return nil, "", nil
	case strings.HasPrefix(tokenUrl, "http") || strings.HasPrefix(tokenUrl, "https"):
		data, err = util.HttpGet(ctx, tokenUrl)
	case strings.HasPrefix(tokenUrl, "data:application/json;base64,"):
		raw := util.Base64Decode(strings.ReplaceAll(tokenUrl, "data:application/json;base64,", ""))
		data = []byte(raw)
	case strings.HasPrefix(tokenUrl, "data:application/json;utf8,"):
		data = []byte(strings.ReplaceAll(tokenUrl, "data:application/json;utf8,", ""))
	case strings.HasPrefix(tokenUrl, "did:dkg"):
		return nil, "", nil
	default:
		return nil, "", fmt.Errorf("unsupported token uri %s %s %s", tokenUrl, contract, tokenId)
	}

	if err != nil {
		return nil, "", err
	}
	if util.UnmarshalAny(&metadata, data) == nil {
		return &metadata, metadata.Image, nil
	}
	return nil, "", err
}

func (c *Token) refreshGetMetadata(ctx context.Context, contract string, tokenId string) error {
	metadata, storageUri, err := c.GetMetadata(ctx, contract, tokenId)
	if metadata != nil {
		sg.db.Model(Erc721Holders{}).Where("contract = ?", c.Contract).Where("token_id =?", tokenId).
			UpdateColumns(map[string]interface{}{"metadata": metadata, "storage_url": storageUri})
	}
	return err
}
