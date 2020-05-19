package http

import (
	"fmt"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	"github.com/bilibili/kratos/pkg/net/http/blademaster/binding"
	"subscan-end/internal/server/websocket"
	"subscan-end/utiles"
	"subscan-end/utiles/ss58"
	"time"
)

func metadata(c *bm.Context) {
	metadata, err := svc.GetChainMetadata()
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
	metadata, err := svc.GetChainMetadata()
	blockNum := utiles.StringToInt(metadata["blockNum"])
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
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	var query []string
	if p.Signed == "signed" {
		query = append(query, "is_signed = 1")
	}
	if p.Address != "" {
		account := ss58.Decode(p.Address)
		if account == "" {
			c.JSON(nil, utiles.InvalidAccountAddress)
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
		ExtrinsicIndex string `json:"extrinsic_index"`
		Hash           string `json:"hash" validate:"omitempty,len=66"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	if p.ExtrinsicIndex == "" && p.Hash == "" {
		c.JSON(nil, utiles.ParamsError)
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
		Row  int `json:"row" validate:"min=1,max=100"`
		Page int `json:"page" validate:"min=0"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	events, count := svc.GetEventList(p.Page, p.Row)
	c.JSON(map[string]interface{}{
		"events": events, "count": count,
	}, nil)
}

func event(c *bm.Context) {
	p := new(struct {
		EventIndex string `json:"event_index" validate:"required"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	c.JSON(svc.GetEventByIndex(p.EventIndex), nil)
}

func search(c *bm.Context) {
	p := new(struct {
		Key  string `json:"key" validate:"required"`
		Row  int    `json:"row" validate:"min=1,max=100"`
		Page int    `json:"page" validate:"min=0"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	c.JSON(svc.SearchByKey(p.Key, p.Page, p.Row), nil)
}

func logs(c *bm.Context) {
	p := new(struct {
		Row  int `json:"row" validate:"min=1,max=100"`
		Page int `json:"page" validate:"min=0"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	logs, count := svc.GetLogList(p.Page, p.Row)
	c.JSON(map[string]interface{}{
		"logs": logs, "count": count,
	}, nil)
}

func logInfo(c *bm.Context) {
	p := new(struct {
		LogIndex string `json:"log_index" validate:"required"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	c.JSON(svc.GetLogByIndex(p.LogIndex), nil)
}

func wsPullHandle(c *bm.Context) {
	websocket.ServeWs(WsHub, c)
}

func transfers(c *bm.Context) {
	p := new(struct {
		Row     int    `json:"row" validate:"min=1,max=100"`
		Page    int    `json:"page" validate:"min=0"`
		Address string `json:"address" validate:"omitempty"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	if p.Address == "" {
		list := svc.GetTransferList(p.Page, p.Row)
		m, _ := svc.GetChainMetadata()
		c.JSON(map[string]interface{}{"transfers": list, "count": utiles.StringToInt(m["count_transfer"])}, nil)
	} else {
		list, count := svc.GetTransfersByAccount(p.Address, p.Page, p.Row)
		c.JSON(map[string]interface{}{"transfers": list, "count": count}, nil)
	}
}

func dailyStat(c *bm.Context) {
	p := new(struct {
		Start string `json:"start" validate:"required"`
		End   string `json:"end" validate:"required"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	if _, err := time.Parse(`2006-01-02`, p.Start); err != nil {
		c.JSON(nil, utiles.ParamsError)
		return
	}
	if _, err := time.Parse(`2006-01-02`, p.End); err != nil {
		c.JSON(nil, utiles.ParamsError)
		return
	}
	c.JSON(map[string]interface{}{
		"list": svc.DailyStat(p.Start, p.End),
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
	c.JSON(nil, utiles.RecordNotFound)
}
