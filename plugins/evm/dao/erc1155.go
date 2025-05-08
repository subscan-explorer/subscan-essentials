package dao

// type ERC1155Item struct {
// 	Id         uint     `gorm:"primaryKey;autoIncrement;size:32" json:"-"`
// 	TokenId    string   `json:"token_id" gorm:"index:contract_item,unique;default:0;size:255"`
// 	Contract   string   `json:"contract" gorm:"index:contract;index:contract_item,unique;size:100"`
// 	Holders    uint     `json:"holders" gorm:"size:32;index:holders"`
// 	Metadata   Metadata `json:"metadata" gorm:"type:json"`
// 	StorageUrl string   `json:"storage_url" gorm:"type:text"`
//
// 	Balance decimal.Decimal `gorm:"-" json:"balance,omitempty"`
// }
//
// func (e *ERC1155Item) TableName() string {
// 	return "evm_erc1155_items"
// }
//
// type ERC1155Holder struct {
// 	Id       uint            `gorm:"primaryKey;autoIncrement;size:32" json:"-"`
// 	TokenId  string          `json:"token_id" gorm:"index:contract_hold,unique;default:0;size:150"`
// 	Contract string          `json:"contract" gorm:"index:contract_hold,unique;size:100"`
// 	Holder   string          `json:"holder" gorm:"index:holder;size:100;index:contract_hold,unique"`
// 	Balance  decimal.Decimal `json:"balance" gorm:"type:decimal(65);index:balance"`
// }
//
// func (e *ERC1155Holder) TableName() string {
// 	return "evm_erc1155_holders"
// }
//
// func (e *ERC1155Item) touchHolder(ctx context.Context, accountId string, balance decimal.Decimal) *gorm.DB {
// 	holder := ERC1155Holder{
// 		TokenId:  e.TokenId,
// 		Contract: e.Contract,
// 		Holder:   accountId,
// 		Balance:  balance,
// 	}
// 	query := sg.AddOrUpdateItem(ctx, &holder, []string{"token_id", "contract", "holder"}, "balance")
//
// 	if query.RowsAffected == 1 || (query.RowsAffected == 2 && balance.IsZero()) {
// 		// refresh item holder count
// 		var incrCount int
// 		if balance.IsPositive() && query.RowsAffected == 1 {
// 			incrCount = 1
// 		}
// 		if balance.IsZero() && query.RowsAffected == 2 {
// 			incrCount = -1
// 		}
// 		// refresh token holder balance
// 		var erc1155Holder TokenHolder
// 		sg.db.Model(TokenHolder{}).Where("contract = ? and holder = ?", e.Contract, accountId).First(&erc1155Holder)
// 		if erc1155Holder.Balance.GreaterThanOrEqual(decimal.NewFromInt(10000)) {
// 			sg.db.Model(TokenHolder{}).Where("contract = ? and holder = ?", e.Contract, accountId).UpdateColumns(map[string]interface{}{"balance": gorm.Expr("balance + ?", incrCount)})
// 		}
// 	}
// 	return query
// }
//
// func (e *ERC1155Item) refreshGetMetadata(ctx context.Context, t *Token) error {
// 	metadata, storageUri, err := t.GetMetadata(ctx, e.Contract, e.TokenId)
// 	if metadata != nil {
// 		sg.db.Model(e).Where("contract = ?", e.Contract).Where("token_id =?", e.TokenId).
// 			UpdateColumns(map[string]interface{}{"metadata": metadata, "storage_url": storageUri})
// 	}
// 	return err
// }
//
// func (e *ERC1155Item) refreshHolders(ctx context.Context) error {
// 	var count int64
// 	sg.db.WithContext(ctx).Model(ERC1155Holder{}).Where("contract = ?", e.Contract).Where("token_id = ?", e.TokenId).Where("balance>0").Count(&count)
// 	if e.Holders == uint(count) {
// 		return nil
// 	}
// 	return sg.db.WithContext(ctx).Model(e).Where("contract = ?", e.Contract).Where("token_id =?", e.TokenId).
// 		Update("holders", count).Error
// }
//
// func (c *Token) RefreshErc1155Holders(ctx context.Context, accountId, tokenId string, balance decimal.Decimal) (int, error) {
// 	if tokenId == "" || accountId == NullAddress {
// 		return 0, nil
// 	}
// 	collectible, createNewItem := c.TouchErc1155Item(ctx, tokenId)
// 	if collectible == nil {
// 		return 0, nil
// 	}
//
// 	defer func() {
// 		if createNewItem {
// 			if c.TotalSupply.LessThan(decimal.NewFromInt(10000)) {
// 				sg.db.Model(c).Update("total_supply", erc1155CollectiblesCount(ctx, c.Contract))
// 			}
// 			_ = collectible.refreshGetMetadata(ctx, c)
// 		}
// 	}()
//
// 	query := collectible.touchHolder(ctx, accountId, balance)
// 	if query.Error != nil {
// 		return 0, query.Error
// 	}
//
// 	if query.RowsAffected == 1 || (query.RowsAffected == 2 && balance.IsZero()) {
// 		// refresh item holder count
// 		err := collectible.refreshHolders(ctx)
// 		if err != nil {
// 			return 0, fmt.Errorf("RefreshErc1155Holders error %w contract %s token id %s", err, c.Contract, tokenId)
// 		}
// 	}
//
// 	var incrCount int
// 	if query.RowsAffected == 1 && balance.IsPositive() {
// 		incrCount = 1
// 	}
// 	return incrCount, nil
// }
//
// func (c *Token) TouchErc1155Item(ctx context.Context, tokenId string) (*ERC1155Item, bool) {
// 	if tokenId == "" {
// 		return nil, false
// 	}
// 	var item ERC1155Item
// 	q := sg.db.WithContext(ctx).Scopes(model.IgnoreDuplicate).FirstOrCreate(&item, ERC1155Item{Contract: c.Contract, TokenId: tokenId})
// 	if q.Error != nil {
// 		return nil, false
// 	}
// 	return &item, q.RowsAffected > 0
// }
//
// // ProcessErc1155 process erc1155 event
// func (t *TransactionReceipt) ProcessErc1155(ctx context.Context) error {
// 	switch util.TrimHex(t.MethodHash) {
// 	// TransferSingle(address,address,address,uint256,uint256)
// 	case erc1155.EventTransferSingle:
// 		events := SplitReceiptData(abi.Erc1155, "TransferSingle", t.Data)
// 		if len(t.Topics) < 4 || len(events) < 2 {
// 			return nil
// 		}
// 		topics := strings.Split(t.Topics, ",")
// 		_from := RemoveAddressPadded(topics[2])
// 		_to := RemoveAddressPadded(topics[3])
// 		tokenId := events[0].(*big.Int)
// 		amount := events[1].(*big.Int)
// 		// fmt.Println("from:", _from, "to:", _to, "tokenId:", tokenId, "amount:", amount)
// 		token := TouchToken(ctx, t.Address, Eip1155Token)
// 		return token.CreateErc1155Transfer(ctx, _from, _to, []*big.Int{tokenId}, []*big.Int{amount}, t)
// 	// 	TransferBatch(address,address,address,uint256[],uint256[]
// 	case erc1155.EventTransferBatch:
// 		events := SplitReceiptData(abi.Erc1155, "TransferBatch", t.Data)
// 		topics := strings.Split(t.Topics, ",")
// 		if len(t.Topics) < 4 || len(events) < 2 {
// 			return nil
// 		}
// 		_from := RemoveAddressPadded(topics[2])
// 		_to := RemoveAddressPadded(topics[3])
// 		tokenIds := events[0].([]*big.Int)
// 		amounts := events[1].([]*big.Int)
// 		token := TouchToken(ctx, t.Address, Eip1155Token)
// 		// fmt.Println("from:", _from, "to:", _to, "tokenIds:", tokenIds, "amounts:", amounts)
// 		return token.CreateErc1155Transfer(ctx, _from, _to, tokenIds, amounts, t)
// 	// URI(string,index uint256)
// 	case erc1155.EventURI:
// 		// events := SplitReceiptData(abi.Erc1155, "URI", t.Data)
// 		topics := strings.Split(t.Topics, ",")
// 		tokenId := util.U256(topics[1]).String()
// 		token := TouchToken(ctx, t.Address, Eip1155Token)
// 		collectible, _ := token.TouchErc1155Item(ctx, tokenId)
// 		_ = collectible.refreshGetMetadata(ctx, token)
// 	}
// 	return nil
// }
//
// // CreateErc1155Transfer batch transfer
// func (c *Token) CreateErc1155Transfer(ctx context.Context, _from, _to string, tokenIds []*big.Int, amount []*big.Int, receipt *TransactionReceipt) error {
// 	if len(tokenIds) == 0 {
// 		return nil
// 	}
// 	if len(tokenIds) != len(amount) {
// 		return fmt.Errorf("token %s tokenIds and amount length not equal", c.Contract)
// 	}
//
// 	var transfers []TokensTransfers
// 	var tokenIdsStr []string
// 	for i, tokenId := range tokenIds {
// 		value := amount[i]
// 		transfer := TokensTransfers{
// 			TransferId: receipt.Id,
// 			Contract:   c.Contract,
// 			Hash:       receipt.TransactionHash,
// 			CreateAt:   receipt.BlockTimestamp,
// 			Sender:     RemoveAddressPadded(_from),
// 			Receiver:   RemoveAddressPadded(_to),
// 			Value:      decimal.NewFromBigInt(value, 0),
// 			TokenId:    tokenId.String(),
// 			BatchIndex: uint(i),
// 			Category:   TransferCategoryErc1155,
// 		}
// 		transfers = append(transfers, transfer)
// 		tokenIdsStr = append(tokenIdsStr, tokenId.String())
// 	}
// 	q := sg.db.Scopes(model.IgnoreDuplicate).CreateInBatches(transfers, 2000)
//
// 	if q.Error != nil {
// 		return q.Error
// 	}
//
// 	if q.RowsAffected > 0 {
// 		c.incrTransferCount(ctx, int(q.RowsAffected))
//
// 		if mq.Instant != nil {
// 			_ = Publish(Eip1155Token, "balance", []string{c.Contract, _from, strings.Join(tokenIdsStr, ",")})
// 			_ = Publish(Eip1155Token, "balance", []string{c.Contract, _to, strings.Join(tokenIdsStr, ",")})
// 		}
// 	}
// 	return nil
// }
//
// func Erc1155Collectibles(c context.Context, page, row int, address, contract, tokenId string) (collectibles []CollectiblesJson, count int64) {
// 	var items []CollectiblesJson
// 	query := sg.db.WithContext(c).Table("evm_erc1155_items")
//
// 	if contract != "" {
// 		query.Where("evm_erc1155_items.contract = ?", contract)
// 		if tokenId != "" {
// 			query.Where("evm_erc1155_items.token_id = ?", tokenId)
// 		}
// 	}
// 	if address != "" {
// 		query.Joins("JOIN evm_erc1155_holders ON evm_erc1155_items.contract = evm_erc1155_holders.contract AND evm_erc1155_items.token_id = evm_erc1155_holders.token_id").
// 			Where("evm_erc1155_holders.holder = ?", address).Where("evm_erc1155_holders.balance > 0")
// 	}
// 	query.Count(&count)
// 	if count == 0 {
// 		return
// 	}
// 	query.Offset(page * row).Limit(row).Scan(&items)
// 	for _, item := range items {
// 		collectibles = append(collectibles, CollectiblesJson{
// 			Contract:   item.Contract,
// 			TokenId:    item.TokenId,
// 			StorageUrl: item.StorageUrl,
// 			Holders:    item.Holders,
// 			Balance:    item.Balance,
// 		})
// 	}
// 	return
// }
//
// func Erc1155Collectible(ctx context.Context, contract, tokenId string) *ERC1155Item {
// 	var nft ERC1155Item
// 	if query := sg.db.Model(ERC1155Item{}).WithContext(ctx).Where("contract = ?", contract).Where("token_id = ?", tokenId).First(&nft); query.Error != nil {
// 		return nil
// 	}
// 	return &nft
// }
//
// func (e *ERC1155Item) HoldersList(ctx context.Context, page, row int, order, orderField string) (holders []ERC1155Holder, count int64) {
// 	q := sg.db.Model(ERC1155Holder{}).Where("contract = ?", e.Contract).Where("token_id = ?", e.TokenId).Where("balance > 0")
// 	q.Count(&count)
// 	defer func() {
// 		go func() {
// 			// update count
// 			if int64(e.Holders) != count {
// 				sg.db.WithContext(context.Background()).Model(e).Where("id = ?", e.Id).Update("holders", count)
// 			}
// 		}()
// 	}()
// 	if orderField != "" && order != "" {
// 		q.Order(fmt.Sprintf("%s %s", orderField, order))
// 	} else {
// 		q.Order("balance desc")
// 	}
// 	q.Offset(page * row).Limit(row).Find(&holders)
// 	return
// }
//
// func ERC1155TokenIdsCount(ctx context.Context, contract, accountId string) int64 {
// 	var count int64
// 	sg.db.WithContext(ctx).Model(ERC1155Holder{}).Where("contract = ?", contract).Where("holder = ?", accountId).Where("balance>0").Count(&count)
// 	return count
// }
//
// func RefreshErc1155Holder(ctx context.Context, contract, address, tokenIdStr string) error {
// 	tokenIds := strings.Split(tokenIdStr, ",")
// 	if len(tokenIds) == 0 {
// 		return nil
// 	}
// 	token := TouchToken(ctx, contract, Eip1155Token)
//
// 	var (
// 		incrSum    int
// 		accountIds []string
// 		ids        []*big.Int
// 	)
// 	for _, id := range tokenIds {
// 		accountIds = append(accountIds, address)
// 		bn := new(big.Int)
// 		n, _ := bn.SetString(id, 10)
// 		ids = append(ids, n)
// 	}
// 	balances, err := erc1155.Init(web3.RPC, contract).BalanceOfBatch(ctx, accountIds, ids)
// 	if err != nil {
// 		return err
// 	}
// 	if len(balances) != len(tokenIds) {
// 		return fmt.Errorf("RefreshErc1155Holder %s balances length not equal tokenIds length", contract)
// 	}
// 	for index, tokenId := range tokenIds {
// 		if tokenId == "" {
// 			continue
// 		}
// 		incr, err := token.RefreshErc1155Holders(ctx, address, tokenId, balances[index])
// 		if err != nil {
// 			return err
// 		}
// 		incrSum += incr
// 	}
// 	var holder TokenHolder
//
// 	if q := sg.db.Model(TokenHolder{}).Where("contract = ? and holder = ?", contract, address).First(&holder); q.Error == nil {
// 		if holder.Balance.IntPart()%100 == 0 && holder.Balance.GreaterThanOrEqual(decimal.New(1000, 0)) {
// 			sg.db.Model(TokenHolder{}).Where("contract = ? and holder = ?", contract, address).UpdateColumns(map[string]interface{}{"balance": gorm.Expr("balance + ?", incrSum)})
// 			return nil
// 		}
// 	}
//
// 	if err = RefreshHolder(ctx, contract, address, Eip1155Token); err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func erc1155CollectiblesCount(ctx context.Context, contract string) int64 {
// 	var count int64
// 	sg.db.WithContext(ctx).Model(ERC1155Item{}).Where("contract = ?", contract).Count(&count)
// 	return count
// }
