package plugins

import "github.com/itering/subscan/internal/plugins/account"

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
	Register("account", func() Plugin { return account.New() })
}

func List() []string {
	plugins := make([]string, 0, len(RegisteredPlugins))
	for name := range RegisteredPlugins {
		plugins = append(plugins, name)
	}
	return plugins
}
