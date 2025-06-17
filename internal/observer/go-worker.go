package observer

import (
	"context"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/itering/go-workers"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins"
	"github.com/itering/subscan/share/metrics"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/mq"
	"github.com/shopspring/decimal"
	"time"
)

func Consumption() {
	concurrency := util.StringToInt(util.GetEnv("WORKER_GOROUTINE_COUNT", "10"))
	workers.Process("block", emitMsg, concurrency)
	workers.Process("balance", emitMsg, concurrency)
	workers.Process("plugin-block", emitMsg, concurrency)
	workers.Process("plugin-event", emitMsg, concurrency)
	workers.Process("plugin-extrinsic", emitMsg, concurrency)

	for _, plugin := range plugins.RegisteredPlugins {
		for _, queue := range plugin.ConsumptionQueue() {
			workers.Process(queue, emitMsg, concurrency)
		}
	}
	mq.Instant.Consumption()
}

func emitMsg(message *workers.Msg) {
	startTime := time.Now()
	defer func() {
		metrics.WorkerProcessCost.WithLabelValues(message.Get("queue").MustString(), message.Get("class").MustString()).Observe(time.Since(startTime).Seconds())
	}()
	var do = func(ctx context.Context, queue, class string, rawInterface interface{}) error {
		raw, ok := rawInterface.(*simplejson.Json)
		if !ok {
			err := fmt.Errorf("emitGoWorker get raw type not workers.Msg")
			return err
		}
		switch queue {
		case "block":
			return blockWorker(ctx, raw)
		case "plugin-block":
			type T struct {
				BlockNum   uint   `json:"block_num"`
				PluginName string `json:"plugin_name"`
			}
			var args T
			if err := util.UnmarshalAny(&args, raw); err != nil {
				panic(fmt.Errorf("plugin-block worker args unmarshal error: %s", err.Error()))
			}

			p, ok := plugins.RegisteredPlugins[args.PluginName]
			if !ok {
				return nil
			}
			return p.ProcessBlock(ctx, srv.GetDao().GetBlockByNum(ctx, args.BlockNum).AsPlugin())

		case "plugin-event":
			type T struct {
				EventIndex string `json:"event_index"`
				PluginName string `json:"plugin_name"`
			}
			var args T
			if err := util.UnmarshalAny(&args, raw); err != nil {
				panic(fmt.Errorf("plugin-block worker args unmarshal error: %s", err.Error()))
			}
			p, ok := plugins.RegisteredPlugins[args.PluginName]
			if !ok {
				return nil
			}

			d := srv.GetDao()
			eventIndexParse := model.ParseExtrinsicOrEventIndex(args.EventIndex)
			block := srv.GetDao().GetBlockByNum(ctx, eventIndexParse.BlockNum).AsPlugin()
			event := d.GetEventByIdx(ctx, args.EventIndex).AsPlugin()
			return p.ProcessEvent(block, event, decimal.Zero)

		case "plugin-extrinsic":
			type T struct {
				ExtrinsicIndex string `json:"event_index"`
				PluginName     string `json:"plugin_name"`
			}
			var args T
			if err := util.UnmarshalAny(&args, raw); err != nil {
				panic(fmt.Errorf("plugin-block worker args unmarshal error: %s", err.Error()))
			}
			p, ok := plugins.RegisteredPlugins[args.PluginName]
			if !ok {
				return nil
			}

			d := srv.GetDao()
			extrinsicIndexParse := model.ParseExtrinsicOrEventIndex(args.ExtrinsicIndex)
			block := srv.GetDao().GetBlockByNum(ctx, extrinsicIndexParse.BlockNum).AsPlugin()
			extrinsic := d.GetExtrinsicsByIndex(ctx, args.ExtrinsicIndex).AsPlugin()
			events := d.GetEventsByIndex(args.ExtrinsicIndex)
			var pEvents []storage.Event
			for _, event := range events {
				pEvents = append(pEvents, *event.AsPlugin())
			}
			return p.ProcessExtrinsic(block, extrinsic, pEvents)

		default:
			// Call the plugin's process function
			for _, plugin := range plugins.RegisteredPlugins {
				if err := plugin.ExecWorker(ctx, queue, class, rawInterface); err != nil {
					return err
				}
			}
		}
		return nil
	}
	err := do(context.Background(), message.Get("queue").MustString(), message.Get("class").MustString(), message.Get("args"))
	if err != nil {
		_ = mq.Instant.ForcePublish(message.Get("queue").MustString(), message.Get("class").MustString(), message.Get("args"))
	}
}

type blockArgs struct {
	BlockNum uint `json:"block_num"`
}

func blockWorker(ctx context.Context, raw interface{}) error {
	var args blockArgs
	if err := util.UnmarshalAny(&args, raw); err != nil {
		return err
	}

	if err := srv.FillBlockData(ctx, args.BlockNum, false); err != nil {
		util.Logger().Error(fmt.Errorf("fill block %d data error: %s", args.BlockNum, err.Error()))
		return err
	}
	return nil
}
