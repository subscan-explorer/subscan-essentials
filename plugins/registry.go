package plugins

import (
	"fmt"
	"reflect"
	"strings"

	subscanPlugin "github.com/itering/subscan-plugin"
	"github.com/itering/subscan/plugins/router"
	"github.com/itering/subscan/plugins/staking"
	"github.com/itering/subscan/plugins/storage"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
)

type Plugin interface {
	// Init storage interface
	InitDao(d storage.Dao)

	// Init http router
	InitHttp() []router.Http

	// Receive Extrinsic data when subscribe extrinsic dispatch
	ProcessExtrinsic(*storage.Block, *storage.Extrinsic, []storage.Event) error

	// Receive Extrinsic data when subscribe extrinsic dispatch
	ProcessEvent(*storage.Block, *storage.Event, decimal.Decimal, *storage.Extrinsic) error

	ProcessCall(*storage.Block, *storage.Call, []storage.Event, *storage.Extrinsic) error

	// Mysql tables schema auto migrate
	Migrate()

	// Subscribe Extrinsic with special module
	SubscribeExtrinsic() []string

	// Subscribe Events with special module
	SubscribeEvent() []string

	// Subscribe Call with special module
	SubscribeCall() []string

	// Plugins version
	Version() string

	UiConf() *subscanPlugin.UiConfig
}
type PluginFactory Plugin

var RegisteredPlugins = make(map[string]PluginFactory)

// register local plugin
func init() {
	registerNative(staking.New())
}

func register(name string, f interface{}) {
	slog.Debug("register plugin", name)
	name = strings.ToLower(name)
	if f == nil {
		return
	}

	if _, ok := RegisteredPlugins[name]; ok {
		return
	}

	if _, ok := f.(PluginFactory); ok {
		RegisteredPlugins[name] = f.(PluginFactory)
	} else {
		panic(fmt.Sprintf("plugin must implement PluginFactory interface: %s", name))
	}

	slog.Debug("Now registered plugins: %v", RegisteredPlugins)
}

func registerNative(p Plugin) {
	register(reflect.ValueOf(p).Type().Elem().Name(), p)
}

type PluginInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func List() []PluginInfo {
	plugins := make([]PluginInfo, 0, len(RegisteredPlugins))

	for name, plugin := range RegisteredPlugins {
		plugins = append(plugins, PluginInfo{Name: name, Version: plugin.Version()})
	}
	return plugins
}
