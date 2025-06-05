package service

import (
	"context"
	"fmt"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins"
	redisDao "github.com/itering/subscan/share/redis"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/mq"
	"strings"
)

var (
	subscribeExtrinsic = make(map[string][]string)
	subscribeEvent     = make(map[string][]string)
)

// registered storage
func pluginRegister(ds *dao.DbStorage, pool *redisDao.Dao) {
	for name, plugin := range plugins.RegisteredPlugins {
		db := *ds
		db.Prefix = name
		plugin.InitDao(&db)
		plugin.SetRedisPool(pool)
		for _, moduleId := range plugin.SubscribeExtrinsic() {
			subscribeExtrinsic[moduleId] = append(subscribeExtrinsic[moduleId], name)
		}
		for _, moduleId := range plugin.SubscribeEvent() {
			subscribeEvent[moduleId] = append(subscribeEvent[moduleId], name)
		}
	}
}

var ignoreEvent = []string{"system.ExtrinsicSuccess"}

// after event created, emit event data to subscribe plugins
func (s *Service) emitEvent(event *model.ChainEvent) (err error) {
	// ignore some event
	if util.StringInSliceFold(fmt.Sprintf("%s.%s", event.ModuleId, event.EventId), ignoreEvent) {
		return
	}
	for _, pluginName := range subscribeEvent[strings.ToLower(event.ModuleId)] {
		if plugins.RegisteredPlugins[pluginName].Enable() {
			if err = mq.Instant.Publish("plugin-event", "process", map[string]interface{}{"event_index": event.EventIndex(), "plugin_name": pluginName}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Service) emitBlock(_ context.Context, block *model.ChainBlock) (err error) {
	for name, plugin := range plugins.RegisteredPlugins {
		if plugin.Enable() {
			if err = mq.Instant.Publish("plugin-block", "process", map[string]interface{}{"block_num": block.BlockNum, "plugin_name": name}); err != nil {
				return err
			}
		}
	}
	return
}

// after extrinsic created, emit extrinsic data to subscribe plugins
func (s *Service) emitExtrinsic(_ context.Context, extrinsic *model.ChainExtrinsic) (err error) {
	for _, pluginName := range subscribeExtrinsic[extrinsic.CallModule] {
		if plugins.RegisteredPlugins[pluginName].Enable() {
			if err = mq.Instant.Publish("plugin-extrinsic", "process", map[string]interface{}{"extrinsic_index": extrinsic.ExtrinsicIndex, "plugin_name": pluginName}); err != nil {
				return err
			}
		}
	}
	return nil
}
