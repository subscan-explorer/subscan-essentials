package service

import (
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins"
	"github.com/itering/subscan/plugins/storage"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
)

var (
	subscribeCall      = make(map[string][]plugins.PluginFactory)
	subscribeExtrinsic = make(map[string][]plugins.PluginFactory)
	subscribeEvent     = make(map[string][]plugins.PluginFactory)
)

// registered storage
func pluginRegister(ds *dao.DbStorage, dd *dao.Dao) {
	for name, plugin := range plugins.RegisteredPlugins {
		db := *ds
		db.Prefix = name
		plugin.InitDao(&db, dd)
		for _, moduleId := range plugin.SubscribeExtrinsic() {
			subscribeExtrinsic[moduleId] = append(subscribeExtrinsic[moduleId], plugin)
		}
		for _, moduleId := range plugin.SubscribeEvent() {
			subscribeEvent[moduleId] = append(subscribeEvent[moduleId], plugin)
		}
		for _, moduleId := range plugin.SubscribeCall() {
			subscribeCall[moduleId] = append(subscribeCall[moduleId], plugin)
		}
	}
}

// after event created, emit event data to subscribe plugins
func (s *Service) emitEvent(block *model.ChainBlock, event *model.ChainEvent, fee decimal.Decimal, extrinsic *model.ChainExtrinsic) {
	for _, plugin := range subscribeEvent[event.ModuleId] {
		if err := plugin.ProcessEvent(block.AsPlugin(), event.AsPlugin(), fee, extrinsic.AsPlugin()); err != nil {
			slog.Error("plugin.ProcessEvent failed", "error", err)
		}
	}
}

// after extrinsic created, emit extrinsic data to subscribe plugins
func (s *Service) emitExtrinsic(block *model.ChainBlock, extrinsic *model.ChainExtrinsic, events []model.ChainEvent) {
	block.BlockTimestamp = extrinsic.BlockTimestamp

	for _, plugin := range subscribeExtrinsic[extrinsic.CallModule] {
		if err := plugin.ProcessExtrinsic(block.AsPlugin(), extrinsic.AsPlugin(), model.MapAsPlugin[*storage.Event](events)); err != nil {
			slog.Error("plugin.ProcessExtrinsic failed", "error", err)
		}
	}
}

func (s *Service) emitCall(block *model.ChainBlock, call *model.ChainCall, events []model.ChainEvent, extrinsic *model.ChainExtrinsic) {
	slog.Debug("emit call", "subscribeCall", subscribeCall)
	for _, plugin := range subscribeCall[call.ModuleId] {
		slog.Debug("calling plugin.ProcessCall", "plugin", plugin)
		if err := plugin.ProcessCall(block.AsPlugin(), call.AsPlugin(), model.MapAsPlugin[*storage.Event](events), extrinsic.AsPlugin()); err != nil {
			slog.Error("plugin.ProcessCall failed", "error", err)
		}
	}
}
