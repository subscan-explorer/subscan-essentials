package dao

import (
	"context"
	"fmt"

	"github.com/itering/subscan/configs"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
)

var (
	testDao   *Dao
	testBlock = model.ChainBlock{
		BlockNum:       947687,
		Hash:           "0xd68b38c412404a4b5d4974e6dbb4a491ed7b6200d4edc24152693804441ce99d",
		ParentHash:     "0x14b8b808939e4930703403d74e73ff7829c18680dd434e851b200982af423dea",
		StateRoot:      "0xd3adc9ed6f9e2df6a13a88a3628c01d7920fd709693120b3df75434aea3592a7",
		ExtrinsicsRoot: "0xc99ede2068646be80f2957c21667a7669539bd105bd855af37c2166a1ba43e4a",
		SpecVersion:    3,
		Validator:      "60e2feb892e672d5579ed10ecae0d162031fe5adc3692498ad262fb126a65732",
		Finalized:      true,
	}

	testEvent = model.ChainEvent{
		EventIdx:     0,
		BlockNum:     947687,
		ModuleId:     "imonline",
		EventId:      "AllGood",
		Params:       model.EventParams{},
		ExtrinsicIdx: 0,
		EventIndex:   "947687-0",
	}

	testExtrinsic = model.ChainExtrinsic{
		ID:                 1,
		ExtrinsicIndex:     "947687-0",
		BlockNum:           947687,
		BlockTimestamp:     1594791900,
		CallModuleFunction: "set",
		CallModule:         "timestamp",
		Params: model.ExtrinsicParams{model.ExtrinsicParam{
			Name:  "now",
			Type:  "Compact<Moment>",
			Value: 1594791900,
		}},
		Success: true,
	}

	testSignedExtrinsic = model.ChainExtrinsic{
		ID:                 2,
		ExtrinsicIndex:     "947689-1",
		BlockNum:           947689,
		BlockTimestamp:     1594791900,
		CallModuleFunction: "transfer",
		CallModule:         "balances",
		AccountId:          "242f0781faa44f34ddcbc9e731d0ddb51c97f5b58bb2202090a3a1c679fc4c63",
		Params: []model.ExtrinsicParam{
			{
				Name:  "dest",
				Type:  "Address",
				Value: "563d11af91b3a166d07110bb49e84094f38364ef39c43a26066ca123a8b9532b",
			},
			{
				Name:  "value",
				Type:  "Compact<Balance>",
				Value: "1000000000000000000",
			},
		},
		Success:       true,
		ExtrinsicHash: "0x368f61800f8645f67d59baf0602b236ff47952097dcaef3aa026b50ddc8dcea0",
		Signature:     "d46ec05eb03ef6904b36fd06fe7923d2a5bccf68ddb53573e821652dafd9644ae82e29c6dbe1519a5b7052c4647814f2987ad23b7c930ed7175726755e27898f",
		IsSigned:      true,
	}

	testLog = model.ChainLog{
		BlockNum: 947687,
		LogIndex: "947687-0",
		LogType:  "Seal",
		Data:     map[string]interface{}{"data": "0x0e4278b7e76436dc08ee4c47d83a0313ef5980dc9fc46b94ccf76318906a4c162e6d1a2b33a69184d4c662ce31176652f0fde8b87cd58e6d1347a28aa29fd58e", "engine": 1161969986},
	}
)

func init() {
	util.ConfDir = "../../configs"
	configs.Init()

	testDao, _, _ = New()
	var tables []string
	db := testDao.db
	err := db.Raw("show tables;").Pluck("Tables_in_subscan_test", &tables).Error
	if err != nil {
		panic(err)
	}

	testDao.Migration()
	for _, value := range tables {
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", value))
	}
	ctx := context.TODO()
	txn := testDao.DbBegin()
	err = testDao.CreateBlock(txn, &testBlock)
	if err != nil {
		panic(err)
	}
	err = testDao.CreateEvent(txn, &testEvent)
	if err != nil {
		panic(err)
	}
	err = testDao.CreateExtrinsic(ctx, txn, &testExtrinsic)
	if err != nil {
		panic(err)
	}
	err = testDao.CreateExtrinsic(ctx, txn, &testSignedExtrinsic)
	if err != nil {
		panic(err)
	}
	err = testDao.CreateLog(txn, &testLog)
	if err != nil {
		panic(err)
	}
	txn.Commit()

	testDao.CreateRuntimeVersion("polkadot", 1)
	testDao.SetRuntimeData(1, "system|staking", "0x0")

	conn, _ := testDao.redis.Redis().GetContext(ctx)
	_, _ = conn.Do("FLUSHALL")
	defer conn.Close()
}
