package contract

import (
	"context"
	"github.com/itering/subscan/pkg/go-web3/dto"
	"github.com/itering/subscan/pkg/go-web3/eth"
)

type Contract struct {
	Eth         *eth.Eth
	EthContract *eth.Contract
	TransParam  dto.TransactionParameters
}

func (ct *Contract) GetStorage(ctx context.Context, functionName string, arg ...interface{}) (*dto.RequestResult, error) {
	h, err := ct.EthContract.Call(ctx, &ct.TransParam, functionName, arg...)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (ct *Contract) GetStorageByKey(ctx context.Context, address string, position string) (string, error) {
	h, err := ct.Eth.GetStorageAt(ctx, address, position, "latest")
	if err != nil {
		return "", err
	}
	return h, nil
}
