package scan

import (
	"context"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/util"
	"github.com/itering/subscan/internal/util/ss58"
	"github.com/json-iterator/go"
	"regexp"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type JsonOption struct {
	Refresh bool
}

type Service struct {
	dao *dao.Dao
}

func New(d *dao.Dao) (s *Service) {
	s = &Service{
		dao: d,
	}
	return s
}

func (s *Service) GetExtrinsicList(page, row int, order string, query ...string) ([]*model.ChainExtrinsicJson, int) {
	c := context.TODO()
	list, count := s.dao.GetExtrinsicList(c, page, row, order, query...)
	var ejs []*model.ChainExtrinsicJson
	for _, extrinsic := range list {
		ejs = append(ejs, s.dao.ExtrinsicsAsJson(&extrinsic))
	}
	return ejs, count
}
func (s *Service) GetBlocksSampleByNums(page, row int) *[]model.SampleBlockJson {
	c := context.TODO()
	var blockJson []model.SampleBlockJson
	// var validatorList []string

	blocks := s.dao.GetBlockList(page, row)

	// for _, block := range blocks {
	// 	validatorList = append(validatorList, block.Validator)
	// }

	for _, block := range blocks {

		bj := s.dao.BlockAsSampleJson(c, &block)

		blockJson = append(blockJson, *bj)
	}
	return &blockJson
}

func (s *Service) GetExtrinsicByIndex(index string) *model.ExtrinsicDetail {
	c := context.TODO()
	return s.dao.GetExtrinsicsDetailByIndex(c, index)
}

func (s *Service) GetExtrinsicDetailByHash(hash string) *model.ExtrinsicDetail {
	c := context.TODO()
	return s.dao.GetExtrinsicsDetailByHash(c, hash)
}

func (s *Service) GetExtrinsicByHash(hash string) *model.ChainExtrinsic {
	c := context.TODO()
	return s.dao.GetExtrinsicsByHash(c, hash)
}

func (s *Service) GetBlockByHashJson(hash string) *model.ChainBlockJson {
	c := context.TODO()
	block := s.dao.BlockByHash(c, hash)
	if block == nil {
		return nil
	}
	return s.dao.BlockAsJson(c, block)
}

func (s *Service) EventByIndex(index string) *model.ChainEvent {
	return s.dao.GetEventByIdx(index)
}

func (s *Service) GetEventList(page, row int, order string, where ...string) ([]model.ChainEventJson, int) {
	c := context.TODO()
	var result []model.ChainEventJson
	var blockNums []int

	list, count := s.dao.GetEventList(c, page, row, order, where...)
	for _, event := range list {
		blockNums = append(blockNums, event.BlockNum)
	}
	blockMap := s.dao.BlocksReverseByNum(c, blockNums)

	for _, event := range list {
		ej := model.ChainEventJson{
			ExtrinsicIdx:  event.ExtrinsicIdx,
			EventIndex:    event.EventIndex,
			BlockNum:      event.BlockNum,
			ModuleId:      event.ModuleId,
			EventId:       event.EventId,
			Params:        util.InterfaceToString(event.Params),
			EventIdx:      event.EventIdx,
			ExtrinsicHash: event.ExtrinsicHash,
		}
		if block, ok := blockMap[event.BlockNum]; ok {
			ej.BlockTimestamp = block.BlockTimestamp
		}
		result = append(result, ej)
	}
	return result, count
}

func (s *Service) GetAccountListJson(page, row int, order, field string, queryWhere ...string) ([]model.AccountSampleJson, int) {
	c := context.TODO()
	list, count := s.dao.GetAccountList(c, page, row, order, field, queryWhere...)
	var j []model.AccountSampleJson
	for _, account := range list {
		j = append(j, *s.dao.AccountSampleJson(c, account))
	}
	return j, count
}

func (s *Service) SearchByKey(key string, page int, row int) interface{} {
	c := context.TODO()
	if regexp.MustCompile("^[0-9]+$").MatchString(key) { // block_num
		return s.GetBlockByNum(util.StringToInt(key))
	}
	if regexp.MustCompile(`^0x[0-9a-fA-F]{64}$`).MatchString(key) { // Extrinsic
		return s.dao.GetExtrinsicsDetailByHash(c, key)
	}

	addressHex := ss58.Decode(key, substrate.AddressType)
	if !regexp.MustCompile(`^[0-9a-fA-F]{64}$`).MatchString(addressHex) { // try account index
		accountIndex := ss58.DecodeAccountIndex(key, substrate.AddressType)
		if accountIndex >= 0 {
			if account, _ := s.dao.FindByIndex(context.TODO(), int(accountIndex)); account != nil {
				return map[string]interface{}{"account": s.dao.AccountAsJson(c, account)}
			}
		}
		return nil
	}
	account, _ := s.dao.FindByAddress(addressHex)
	if account == nil {
		// check balance, if balance positive, create account
		if amount, _, _ := s.dao.GetBalanceFromNetwork(c, addressHex, "balances"); amount.IsPositive() {
			account, _ = s.dao.TouchAccount(c, addressHex)
		} else {
			return nil
		}
	}
	return map[string]interface{}{"account": s.dao.AccountAsJson(c, account)}
}

func (s *Service) GetBlockByNum(num int) *model.ChainBlockJson {
	c := context.TODO()
	block := s.dao.Block(c, num)
	if block == nil {
		return nil
	}
	return s.dao.BlockAsJson(c, block)
}
