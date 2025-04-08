package observer

import (
	"context"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/itering/go-workers"
	"github.com/itering/subscan/plugins"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/mq"
)

func Consumption() {
	// workers.Logger = log.Logger()
	concurrency := util.StringToInt(util.GetEnv("WORKER_GOROUTINE_COUNT", "10"))
	workers.Process("block", emitMsg, concurrency)
	workers.Process("balance", emitMsg, concurrency)

	for _, plugin := range plugins.RegisteredPlugins {
		for _, queue := range plugin.ConsumptionQueue() {
			workers.Process(queue, emitMsg, concurrency)
		}
	}
	mq.Instant.Consumption()
}

func emitMsg(message *workers.Msg) {
	var do = func(ctx context.Context, queue, class string, rawInterface interface{}) error {
		raw, ok := rawInterface.(*simplejson.Json)
		if !ok {
			err := fmt.Errorf("emitGoWorker get raw type not workers.Msg")
			return err
		}
		switch queue {
		case "block":
			return blockWorker(ctx, raw)
		default:
			// Call the plugin's process function
			for _, plugin := range plugins.RegisteredPlugins {
				return plugin.ExecWorker(ctx, queue, class, rawInterface)
			}
		}
		return nil
	}
	util.Logger().Error(do(context.Background(), message.Get("queue").MustString(), message.Get("class").MustString(), message.Get("args")))
}

type blockArgs struct {
	BlockNum uint `json:"block_num"`
}

func blockWorker(ctx context.Context, raw interface{}) error {
	var args blockArgs
	if err := util.UnmarshalAny(&args, raw); err != nil {
		return err
	}

	if err := srv.FillBlockData(ctx, args.BlockNum); err != nil {
		util.Logger().Error(fmt.Errorf("fill block %d data error: %s", args.BlockNum, err.Error()))
		return err
	}
	return nil
}
