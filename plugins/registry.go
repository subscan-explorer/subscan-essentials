package plugins

import (
	"github.com/itering/subscan-plugin"
	"github.com/itering/subscan/plugins/balance"
	"github.com/itering/subscan/plugins/system"
	"reflect"
	"strings"
)

type PluginFactory subscan_plugin.Plugin

var RegisteredPlugins = make(map[string]PluginFactory)

// register local plugin
func init() {
	registerNative(balance.New())
	registerNative(system.New())
}

func register(name string, f interface{}) {
	name = strings.ToLower(name)
	if f == nil {
		return
	}

	if _, ok := RegisteredPlugins[name]; ok {
		return
	}

	if _, ok := f.(PluginFactory); ok {
		RegisteredPlugins[name] = f.(PluginFactory)
	}
}

func registerNative(p interface{}) {
	register(reflect.ValueOf(p).Type().Elem().Name(), p)
}

type PluginInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Ui      bool   `json:"ui"`
}

func List() []PluginInfo {
	plugins := make([]PluginInfo, 0, len(RegisteredPlugins))
	for name, plugin := range RegisteredPlugins {
		plugins = append(plugins, PluginInfo{Name: name, Version: plugin.Version(), Ui: plugin.UiConf() != nil})
	}
	return plugins
}
