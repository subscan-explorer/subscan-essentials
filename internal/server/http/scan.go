package http

import (
	"context"
	"fmt"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
	"github.com/itering/subscan/plugins"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/ss58"
)

func metadata(c *bm.Context) {
	metadata, err := svc.Metadata()
	c.JSON(metadata, err)
}

func blocks(c *bm.Context) {
	p := new(struct {
		Row  int `json:"row" validate:"min=1,max=100"`
		Page int `json:"page" validate:"min=0"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	blockNum, err := svc.GetCurrentBlockNum(context.TODO())
	blocks := svc.GetBlocksSampleByNums(p.Page, p.Row)
	c.JSON(map[string]interface{}{
		"blocks": blocks, "count": blockNum,
	}, err)
}

func block(c *bm.Context) {
	p := new(struct {
		BlockNum  int    `json:"block_num" validate:"omitempty,min=0"`
		BlockHash string `json:"block_hash" validate:"omitempty,len=66"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	if p.BlockHash == "" {
		c.JSON(svc.GetBlockByNum(p.BlockNum), nil)
	} else {
		c.JSON(svc.GetBlockByHashJson(p.BlockHash), nil)
	}
}

func extrinsics(c *bm.Context) {
	p := new(struct {
		Row     int    `json:"row" validate:"min=1,max=100"`
		Page    int    `json:"page" validate:"min=0"`
		Signed  string `json:"signed" validate:"omitempty"`
		Address string `json:"address" validate:"omitempty"`
		Module  string `json:"module" validate:"omitempty"`
		Call    string `json:"call" validate:"omitempty"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	var query []string
	if p.Module != "" {
		query = append(query, fmt.Sprintf("call_module = '%s'", p.Module))
	}
	if p.Call != "" {
		query = append(query, fmt.Sprintf("call_module_function = '%s'", p.Call))
	}

	if p.Signed == "signed" {
		query = append(query, "is_signed = 1")
	}
	if p.Address != "" {
		account := ss58.Decode(p.Address, util.StringToInt(util.AddressType))
		if account == "" {
			c.JSON(nil, util.InvalidAccountAddress)
			return
		}
		query = append(query, fmt.Sprintf("is_signed = 1 and account_id = '%s'", account))
	}
	extrinsics, count := svc.GetExtrinsicList(p.Page, p.Row, "desc", query...)
	c.JSON(map[string]interface{}{
		"extrinsics": extrinsics, "count": count,
	}, nil)

}

func extrinsic(c *bm.Context) {
	p := new(struct {
		ExtrinsicIndex string `json:"extrinsic_index" validate:"omitempty"`
		Hash           string `json:"hash" validate:"omitempty,len=66"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	if p.ExtrinsicIndex == "" && p.Hash == "" {
		c.JSON(nil, util.ParamsError)
		return
	}
	if p.ExtrinsicIndex != "" {
		c.JSON(svc.GetExtrinsicByIndex(p.ExtrinsicIndex), nil)
	} else {
		c.JSON(svc.GetExtrinsicDetailByHash(p.Hash), nil)
	}
}

func events(c *bm.Context) {
	p := new(struct {
		Row    int    `json:"row" validate:"min=1,max=100"`
		Page   int    `json:"page" validate:"min=0"`
		Module string `json:"module" validate:"omitempty"`
		Call   string `json:"call" validate:"omitempty"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	var query []string
	if p.Module != "" {
		query = append(query, fmt.Sprintf("module_id = '%s'", p.Module))
	}
	if p.Call != "" {
		query = append(query, fmt.Sprintf("event_id = '%s'", p.Call))
	}
	events, count := svc.RenderEvents(p.Page, p.Row, "desc", query...)
	c.JSON(map[string]interface{}{
		"events": events, "count": count,
	}, nil)
}

func checkSearchHash(c *bm.Context) {
	p := new(struct {
		Hash string `json:"hash" validate:"len=66"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	if block := svc.GetBlockByHash(p.Hash); block != nil {
		c.JSON(map[string]string{"hash_type": "block"}, nil)
		return
	}
	if extrinsic := svc.GetExtrinsicByHash(p.Hash); extrinsic != nil {
		c.JSON(map[string]string{"hash_type": "extrinsic"}, nil)
		return
	}
	c.JSON(nil, util.RecordNotFound)
}

func runtimeList(c *bm.Context) {
	c.JSON(map[string]interface{}{
		"list": svc.SubstrateRuntimeList(),
	}, nil)
}

func runtimeMetadata(c *bm.Context) {
	p := new(struct {
		Spec int `json:"spec"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	if info := svc.SubstrateRuntimeInfo(p.Spec); info == nil {
		c.JSON(map[string]interface{}{"info": nil}, nil)
	} else {
		c.JSON(map[string]interface{}{
			"info": info.Metadata.Modules,
		}, nil)
	}

}

func pluginList(c *bm.Context) {
	c.JSON(plugins.List(), nil)
}

func pluginUIConfig(c *bm.Context) {
	p := new(struct {
		Name string `json:"name" validate:"required"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	if plugin, ok := plugins.RegisteredPlugins[p.Name]; ok {
		c.JSON(plugin.UiConf(), nil)
		return
	}
	c.JSON(nil, nil)
}
