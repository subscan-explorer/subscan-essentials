package http

import (
	"errors"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/share/token"
	"github.com/itering/subscan/util/address"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/itering/subscan/util"
)

type Pagination struct {
	Row  int `json:"row" binding:"min=1,max=100"`
	Page int `json:"page" binding:"min=0"`
}

// metadataHandle get metadata info, include chain customer info, runtime info, etc.
func metadataHandle(c *gin.Context) {
	m, err := svc.Metadata(c.Request.Context())
	toJson(c, m, err)
}

func tokenHandle(c *gin.Context) {
	toJson(c, token.GetDefaultToken(), nil)
}

// blocksHandle  get blocks list
func blocksHandle(c *gin.Context) {
	p := new(struct {
		Pagination
	})
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}

	ctx := c.Request.Context()
	blockNum, err := svc.GetFinalizedBlock(ctx)
	list := svc.GetBlocksSampleByNums(ctx, p.Page, p.Row)

	toJson(c, map[string]interface{}{
		"blocks": list, "count": blockNum,
	}, err)
}

// blockHandle get block info by block number or block hash
func blockHandle(c *gin.Context) {
	p := new(struct {
		BlockNum  uint   `json:"block_num" binding:"omitempty,min=0"`
		BlockHash string `json:"block_hash" binding:"omitempty,len=66"`
	})
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}
	ctx := c.Request.Context()

	if p.BlockHash == "" {
		toJson(c, svc.GetBlockByNum(ctx, p.BlockNum), nil)
	} else {
		toJson(c, svc.GetBlockByHashJson(ctx, p.BlockHash), nil)
	}
}

// extrinsicsHandle handler get extrinsics list
func extrinsicsHandle(c *gin.Context) {
	p := new(struct {
		Pagination
		Signed  string `json:"signed" binding:"omitempty"`
		Address string `json:"address" binding:"omitempty"`
		Module  string `json:"module" binding:"omitempty"`
		Call    string `json:"call" binding:"omitempty"`
	})
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}
	ctx := c.Request.Context()

	var query []model.Option
	if p.Module != "" {
		query = append(query, model.Where("call_module = ?", p.Module))
	}
	if p.Call != "" {
		query = append(query, model.Where("call_module_function = ?", p.Module))
	}

	if p.Signed == "signed" {
		query = append(query, model.Where("is_signed = 1"))
	}

	if p.Address != "" {
		account := address.Decode(p.Address)
		if account == "" {
			toJson(c, nil, util.InvalidAccountAddress)
			return
		}
		query = append(query, model.Where("account_id = ?", account))
	}

	list, count := svc.GetExtrinsicList(ctx, p.Page, p.Row, query...)
	toJson(c, map[string]interface{}{
		"extrinsics": list, "count": count,
	}, nil)

}

// extrinsicHandle handler get extrinsic info by extrinsic index or extrinsic hash
func extrinsicHandle(c *gin.Context) {
	p := new(struct {
		ExtrinsicIndex string `json:"extrinsic_index" binding:"omitempty"`
		Hash           string `json:"hash" binding:"omitempty,len=66"`
	})
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}
	if p.ExtrinsicIndex == "" && p.Hash == "" {
		toJson(c, nil, errors.New("extrinsic_index or hash is required"))
		return
	}

	ctx := c.Request.Context()

	if p.ExtrinsicIndex != "" {
		toJson(c, svc.GetExtrinsicByIndex(ctx, p.ExtrinsicIndex), nil)
	} else {
		toJson(c, svc.GetExtrinsicDetailByHash(ctx, p.Hash), nil)
	}
}

// eventsHandle handler get events list
func eventsHandle(c *gin.Context) {
	p := new(struct {
		Row    int    `json:"row" binding:"min=1,max=100"`
		Page   int    `json:"page" binding:"min=0"`
		Module string `json:"module" binding:"omitempty"`
		Event  string `json:"event" binding:"omitempty"`
	})
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}
	ctx := c.Request.Context()

	var query []model.Option
	if p.Module != "" {
		query = append(query, model.Where("module_id = ?", p.Module))
	}
	if p.Event != "" {
		query = append(query, model.Where("event_id = ?", p.Event))
	}

	events, count := svc.EventsList(ctx, p.Page, p.Row, query...)
	toJson(c, map[string]interface{}{"events": events, "count": count}, nil)
}

// checkSearchHashHandle handler check hash type, block or extrinsic or evm tx hash
func checkSearchHashHandle(c *gin.Context) {
	p := new(struct {
		Hash string `json:"hash" binding:"len=66"`
	})
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}

	ctx := c.Request.Context()

	if data := svc.GetBlockByHash(ctx, p.Hash); data != nil {
		toJson(c, map[string]string{"hash_type": "block"}, nil)
		return
	}
	if data := svc.GetExtrinsicByHash(ctx, p.Hash); data != nil {
		toJson(c, map[string]string{"hash_type": "extrinsic"}, nil)
		return
	}
	// todo evm tx hash
	toJson(c, nil, util.RecordNotFound)
}

// runtimeListHandler  get runtime list
func runtimeListHandler(c *gin.Context) {
	toJson(c, map[string]interface{}{
		"list": svc.SubstrateRuntimeList(),
	}, nil)
}

// runtimeMetadataHandle get runtime metadata info by spec version
func runtimeMetadataHandle(c *gin.Context) {
	p := new(struct {
		Spec int `json:"spec"`
	})
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}

	if info := svc.SubstrateRuntimeInfo(p.Spec); info != nil {
		toJson(c, map[string]interface{}{"info": info.Metadata.Modules}, nil)
		return
	}

	toJson(c, map[string]interface{}{"info": nil}, nil)
}
