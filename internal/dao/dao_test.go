package dao

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/itering/subscan/configs"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/jinzhu/gorm"
)

var (
	testDao   *Dao
	testBlock = model.ChainBlock{
		BlockNum:       947687,
		Hash:           "0xd68b38c412404a4b5d4974e6dbb4a491ed7b6200d4edc24152693804441ce99d",
		ParentHash:     "0x14b8b808939e4930703403d74e73ff7829c18680dd434e851b200982af423dea",
		StateRoot:      "0xd3adc9ed6f9e2df6a13a88a3628c01d7920fd709693120b3df75434aea3592a7",
		ExtrinsicsRoot: "0xc99ede2068646be80f2957c21667a7669539bd105bd855af37c2166a1ba43e4a",
		Logs:           `["0x0642414245b501010a000000fac3d70f000000009e335d221536deb53426c3f2529a14426a322463a844d527f8050c73f09c2d37bfe0d8f57a7b6c6e6cd6ef576d00bb97b5bcf8c87ec7a55670b03c0dfe823000d2d3bb5767274a282be5dd15f7e6ea333dc44c299f187dee4900fdf1a0b46003","0x00904d4d5252fbe5a48df0e2a689c92a630bcbb451d66e2ac0ea839096e2617c4fe1b22a635e","0x05424142450101ea06828ccb667fbaebdda98219e93700c24c6887b767680949fde8082a93673cf96bb377923751c892d37c78eaa5c8e6b453efbac656fbcac4a8b99a82287e89"]`,
		Extrinsics:     `["0x280402000b603301517301"]`,
		Event:          `0x040000000000000080e36a0900000000020000`,
		SpecVersion:    3,
		Validator:      "60e2feb892e672d5579ed10ecae0d162031fe5adc3692498ad262fb126a65732",
		Finalized:      true,
	}

	testEvent = model.ChainEvent{
		EventIdx:     0,
		BlockNum:     947687,
		ModuleId:     "imonline",
		EventId:      "AllGood",
		Params:       util.ToString([]interface{}{}),
		ExtrinsicIdx: 0,
		EventIndex:   "947687-0",
	}

	testExtrinsic = model.ChainExtrinsic{
		ID:                 1,
		ExtrinsicIndex:     "947687-0",
		BlockNum:           947687,
		BlockTimestamp:     1594791900,
		VersionInfo:        "04",
		CallModuleFunction: "set",
		CallModule:         "timestamp",
		Params: model.ExtrinsicParam{
			Name:  "now",
			Type:  "Compact<Moment>",
			Value: 1594791900,
		},
		Success: true,
	}

	testSignedExtrinsic = model.ChainExtrinsic{
		ID:                 2,
		ExtrinsicIndex:     "947689-1",
		BlockNum:           947689,
		BlockTimestamp:     1594791900,
		VersionInfo:        "04",
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
		Data:     `{"data":"0x0e4278b7e76436dc08ee4c47d83a0313ef5980dc9fc46b94ccf76318906a4c162e6d1a2b33a69184d4c662ce31176652f0fde8b87cd58e6d1347a28aa29fd58e","engine":1161969986}`,
	}
)

func init() {
	if client, err := paladin.NewFile("../../configs"); err != nil {
		panic(err)
	} else {
		paladin.DefaultClient = client
	}
	var (
		dc configs.MysqlConf
		rc configs.RedisConf
	)
	dc.MergeConf()
	rc.MergeConf()

	db, err := gorm.Open("mysql", dc.Test.DSN)
	if err != nil {
		panic(err)
	}
	db.LogMode(false)
	const testRedisDb = 1

	testDao = &Dao{
		db:    db,
		redis: redis.NewPool(rc.Config, redis.DialDatabase(testRedisDb)),
	}
	var tables []string
	err = db.Raw("show tables;").Pluck("Tables_in_subscan_test", &tables).Error
	if err != nil {
		panic(err)
	}

	testDao.Migration()
	for _, value := range tables {
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", value))
	}
	ctx := context.TODO()
	txn := testDao.DbBegin()
	_ = testDao.CreateBlock(txn, &testBlock)
	_ = testDao.CreateEvent(txn, &testEvent)
	_ = testDao.CreateExtrinsic(ctx, txn, &testExtrinsic)
	_ = testDao.CreateExtrinsic(ctx, txn, &testSignedExtrinsic)
	_ = testDao.CreateLog(txn, &testLog)
	txn.Commit()

	testDao.CreateRuntimeVersion("polkadot", 1)
	testDao.SetRuntimeData(1, "system|staking", "0x0")

	conn := testDao.redis.Get(ctx)
	_, _ = conn.Do("FLUSHALL")
	defer conn.Close()
}
