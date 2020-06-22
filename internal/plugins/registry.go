package plugins

import (
	"github.com/itering/subscan/internal/plugins/balance"
	"github.com/itering/subscan/internal/plugins/system"
)

type PluginFactory func() Plugin

var RegisteredPlugins = make(map[string]PluginFactory)

func Register(name string, f PluginFactory) {
	if f == nil {
		return
	}

	if _, ok := RegisteredPlugins[name]; ok {
		return
	}

	RegisteredPlugins[name] = f

}

func init() {
	Register("account", func() Plugin { return balance.New() })
	Register("system", func() Plugin { return system.New() })
}

func List() []string {
	plugins := make([]string, 0, len(RegisteredPlugins))
	for name := range RegisteredPlugins {
		plugins = append(plugins, name)
	}
	return plugins
}
