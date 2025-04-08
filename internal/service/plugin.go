package service

import (
	"context"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins"
	redisDao "github.com/itering/subscan/share/redis"
	"github.com/shopspring/decimal"
)

var (
	subscribeExtrinsic = make(map[string][]plugins.PluginFactory)
	subscribeEvent     = make(map[string][]plugins.PluginFactory)
)

// registered storage
func pluginRegister(ds *dao.DbStorage, pool *redisDao.Dao) {
	for name, plugin := range plugins.RegisteredPlugins {
		db := *ds
		db.Prefix = name
		plugin.InitDao(&db)
		plugin.SetRedisPool(pool)
		for _, moduleId := range plugin.SubscribeExtrinsic() {
			subscribeExtrinsic[moduleId] = append(subscribeExtrinsic[moduleId], plugin)
		}
		for _, moduleId := range plugin.SubscribeEvent() {
			subscribeEvent[moduleId] = append(subscribeEvent[moduleId], plugin)
		}
	}
}

// after event created, emit event data to subscribe plugins
func (s *Service) emitEvent(block *model.ChainBlock, event *model.ChainEvent) {
	pBlock := block.AsPlugin()
	pEvent := event.AsPlugin()
	for _, plugin := range subscribeEvent[event.ModuleId] {
		if plugin.Enable() {
			_ = plugin.ProcessEvent(pBlock, pEvent, decimal.Zero)
		}
	}

}

func (s *Service) emitBlock(ctx context.Context, block *model.ChainBlock) {
	pBlock := block.AsPlugin()
	for _, plugin := range plugins.RegisteredPlugins {
		if plugin.Enable() {
			_ = plugin.ProcessBlock(ctx, pBlock)
		}
	}
}

// after extrinsic created, emit extrinsic data to subscribe plugins
func (s *Service) emitExtrinsic(block *model.ChainBlock, extrinsic *model.ChainExtrinsic, events []model.ChainEvent) {
	block.BlockTimestamp = extrinsic.BlockTimestamp
	pBlock := block.AsPlugin()
	pExtrinsic := extrinsic.AsPlugin()

	var pEvents []storage.Event
	for _, event := range events {
		pEvents = append(pEvents, *event.AsPlugin())
	}

	for _, plugin := range subscribeExtrinsic[extrinsic.CallModule] {
		if plugin.Enable() {
			_ = plugin.ProcessExtrinsic(pBlock, pExtrinsic, pEvents)
		}
	}
}
