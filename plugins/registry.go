package plugins

import (
	"github.com/itering/subscan-plugin"
	"github.com/itering/subscan/plugins/balance"
	"github.com/itering/subscan/plugins/evm"
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
	registerNative(evm.New())
}

func register(name string, f subscan_plugin.Plugin) {
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

func registerNative(p subscan_plugin.Plugin) {
	register(reflect.ValueOf(p).Type().Elem().Name(), p)
}
