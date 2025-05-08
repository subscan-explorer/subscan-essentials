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

// @Summary Current network metadata
// @Description get metadata info, include chain customer info, runtime info, etc.
// @Tags metadata
// @Produce json
// @Success 200 {object} http.J{data=map[string]string}
// @Router /api/scan/metadata [post]
func metadataHandle(c *gin.Context) {
	m, err := svc.Metadata(c.Request.Context())
	toJson(c, m, err)
}

// @Summary Token list
// @Tags tokens
// @Accept json
// @Produce json
// @Success 200 {object} http.J{data=object{token=[]string,detail=map[string]token.Token}}
// @Router /api/scan/token [post]
func tokenHandle(c *gin.Context) {
	toJson(c, token.GetDefaultToken(), nil)
}

type BlocksParams struct {
	Pagination
}

// @Summary Blocks list
// @Tags block
// @Accept json
// @Produce json
// @Param params body BlocksParams true "params"
// @Success 200 {object} http.J{data=object{blocks=[]model.SampleBlockJson,count=int}}
// @Router /api/scan/blocks [post]
func blocksHandle(c *gin.Context) {
	p := new(BlocksParams)
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

type BlockParams struct {
	BlockNum  uint   `json:"block_num" binding:"omitempty,min=0"`
	BlockHash string `json:"block_hash" binding:"omitempty,len=66"`
}

// @Summary Get block details
// @Tags block
// @Accept json
// @Produce json
// @Param params body BlockParams true "params"
// @Success 200 {object} http.J{data=model.ChainBlockJson}
// @Router /api/scan/block [post]
func blockHandle(c *gin.Context) {
	p := new(BlockParams)
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

type extrinsicsParams struct {
	Pagination
	Signed   string `json:"signed" binding:"omitempty"`
	Address  string `json:"address" binding:"omitempty"`
	Module   string `json:"module" binding:"omitempty"`
	Call     string `json:"call" binding:"omitempty"`
	BlockNum uint   `json:"block_num" binding:"omitempty"`
}

// extrinsicsHandle handler get extrinsics list
// @Summary Get extrinsics list
// @Tags extrinsics
// @Accept json
// @Produce json
// @Param params body extrinsicsParams true "params"
// @Success 200 {object} http.J{data=object{extrinsics=[]model.ChainExtrinsicJson,count=int}}
// @Router /api/scan/extrinsics [post]
func extrinsicsHandle(c *gin.Context) {
	p := new(extrinsicsParams)
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
	if p.BlockNum > 0 {
		query = append(query, model.Where("block_num = ?", p.BlockNum))
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

type extrinsicParams struct {
	ExtrinsicIndex string `json:"extrinsic_index" binding:"omitempty"`
	Hash           string `json:"hash" binding:"omitempty,len=66"`
}

// extrinsicHandle handler get extrinsic info by extrinsic index or extrinsic hash
// @Summary Get extrinsic details
// @Tags extrinsics
// @Accept json
// @Produce json
// @Param params body extrinsicParams true "params"
// @Success 200 {object} http.J{data=model.ExtrinsicDetail}
// @Router /api/scan/extrinsic [post]
func extrinsicHandle(c *gin.Context) {
	p := new(extrinsicParams)
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

type eventsParams struct {
	Row      int    `json:"row" binding:"min=1,max=100"`
	Page     int    `json:"page" binding:"min=0"`
	Module   string `json:"module" binding:"omitempty"`
	Event    string `json:"event" binding:"omitempty"`
	BlockNum uint   `json:"block_num" binding:"omitempty"`
}

// eventsHandle handler get events list
// @Summary Get events list
// @Tags events
// @Accept json
// @Produce json
// @Param params body eventsParams true "params"
// @Success 200 {object} http.J{data=object{events=[]model.ChainEventJson,count=int}}
// @Router /api/scan/events [post]
func eventsHandle(c *gin.Context) {
	p := new(eventsParams)
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
	if p.BlockNum > 0 {
		query = append(query, model.Where("block_num = ?", p.BlockNum))
	}

	events, count := svc.EventsList(ctx, p.Page, p.Row, query...)
	toJson(c, map[string]interface{}{"events": events, "count": count}, nil)
}

// logsHandle handler get logs list
// @Summary Get logs list
// @Tags logs
// @Accept json
// @Produce json
// @Param block_num body uint true "Block number"
// @Success 200 {object} http.J{data=[]model.ChainLogJson}
// @Router /api/scan/logs [post]
func logsHandle(c *gin.Context) {
	p := new(struct {
		BlockNum uint `json:"block_num" binding:"required"`
	})
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}
	ctx := c.Request.Context()
	toJson(c, svc.LogsList(ctx, p.BlockNum), nil)
}

type checkSearchParams struct {
	Hash string `json:"hash" binding:"len=66"`
}

// checkSearchHashHandle handler check hash type, block or extrinsic or evm tx hash
// @Summary Check hash type
// @Tags hash
// @Accept json
// @Produce json
// @Param params body checkSearchParams true "params"
// @Success 200 {object} http.J{data=map[string]string}
// @Router/api/scan/check_hash[post]
func checkSearchHashHandle(c *gin.Context) {
	p := new(checkSearchParams)
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

// @Summary Get runtime list
// @Description runtimeListHandler  get runtime list
// @Tags runtime
// @Produce json
// @Success 200 {object} http.J{data=object{list=[]model.RuntimeVersion}}
// @Router /api/scan/runtime/list [post]
func runtimeListHandler(c *gin.Context) {
	toJson(c, map[string]interface{}{
		"list": svc.SubstrateRuntimeList(),
	}, nil)
}

type runtimeMetadataParams struct {
	Spec int `json:"spec"`
}

// runtimeMetadataHandle get runtime metadata info by spec version
// @Summary Get runtime metadata
// @Tags runtime
// @Accept json
// @Produce json
// @Param params body runtimeMetadataParams true "params"
// @Success 200 {object} http.J{data=object{info=metadata.Instant}}
// @Router /api/scan/runtime/metadata [post]
func runtimeMetadataHandle(c *gin.Context) {
	p := new(runtimeMetadataParams)
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}

	if info := svc.SubstrateRuntimeInfo(p.Spec); info != nil {
		toJson(c, map[string]interface{}{"info": info.Metadata.Modules}, nil)
		return
	}

	toJson(c, map[string]interface{}{"info": nil}, nil)
}
