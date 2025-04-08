package dao

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"gorm.io/datatypes"
)

type AbiMapping struct {
	Id      string         `json:"id" gorm:"primaryKey;autoIncrement:false;size:255"`
	AbiFunc datatypes.JSON `json:"abi_func" es:"type:flattened"`
	AbiType ABIType        `json:"abi_type" gorm:"size:100"`
}

func (t *AbiMapping) TableName() string {
	return "evm_abi_mappings"
}

type ABIType string

const (
	MethodTypeEvent  ABIType = "event"
	MethodTypeMethod ABIType = "function"
)

func (c *Contract) fetchAbiMapping(ctx context.Context) error {
	var mappings []AbiMapping
	var abiValue abi.ABI
	if err := abiValue.UnmarshalJSON(c.Abi); err != nil {
		return err
	}
	type Argument struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Indexed bool   `json:"indexed"`
	}

	type Fields struct {
		Type            string     `json:"type"`
		Name            string     `json:"name"`
		Inputs          []Argument `json:"inputs"`
		Outputs         []Argument `json:"outputs"`
		StateMutability string     `json:"stateMutability"`
		Constant        bool       `json:"constant"`
		Payable         bool       `json:"payable"`
		Anonymous       *bool      `json:"anonymous,omitempty"`
	}

	checkoutType := func(funType abi.FunctionType) string {
		identity := "function"
		if funType == abi.Fallback {
			identity = "fallback"
		} else if funType == abi.Receive {
			identity = "receive"
		} else if funType == abi.Constructor {
			identity = "constructor"
		}
		return identity
	}
	convertArg := func(args []abi.Argument) []Argument {
		ret := []Argument{}
		for _, arg := range args {
			ret = append(ret, Argument{
				Name:    arg.Name,
				Type:    arg.Type.String(),
				Indexed: arg.Indexed,
			})
		}
		return ret
	}

	for _, method := range abiValue.Methods {
		f := Fields{
			Type:            checkoutType(method.Type),
			Name:            method.Name,
			Inputs:          convertArg(method.Inputs),
			Outputs:         convertArg(method.Outputs),
			StateMutability: method.StateMutability,
			Constant:        method.Constant,
			Payable:         method.Payable,
		}
		abiFunc, _ := json.Marshal(f)
		abiMapping := AbiMapping{
			Id:      util.AddHex(util.BytesToHex(method.ID)),
			AbiFunc: abiFunc,
			AbiType: MethodTypeMethod,
		}
		mappings = append(mappings, abiMapping)
	}

	for _, event := range abiValue.Events {
		e := event
		f := Fields{
			Type:      "event",
			Name:      event.Name,
			Inputs:    convertArg(event.Inputs),
			Outputs:   convertArg(nil),
			Anonymous: &e.Anonymous,
		}
		abiFunc, _ := json.Marshal(f)
		abiMapping := AbiMapping{
			Id:      event.ID.Hex(),
			AbiFunc: abiFunc,
			AbiType: MethodTypeEvent,
		}
		mappings = append(mappings, abiMapping)
	}
	q := sg.db.WithContext(ctx).Scopes(model.IgnoreDuplicate).CreateInBatches(&mappings, 200)
	return q.Error
}
