package plugins

import (
	"flag"
	"fmt"
	"github.com/itering/subscan-plugin"
	"io/ioutil"
	"plugin"
	"reflect"
	"strings"
)

type PluginFactory subscan_plugin.Plugin

var RegisteredPlugins = make(map[string]PluginFactory)

func Register(name string, f interface{}) {
	if f == nil {
		return
	}

	if _, ok := RegisteredPlugins[name]; ok {
		return
	}

	RegisteredPlugins[name] = f.(PluginFactory)

}

func List() []string {
	plugins := make([]string, 0, len(RegisteredPlugins))
	for name := range RegisteredPlugins {
		plugins = append(plugins, name)
	}
	return plugins
}

func init() {
	flag.Parse()
	pluginsDir := fmt.Sprintf("%s/plugins", flag.Lookup("conf").Value)
	files, err := ioutil.ReadDir(pluginsDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		p, err := plugin.Open(fmt.Sprintf("%s/%s", pluginsDir, file.Name()))
		if err != nil {
			panic(err)
		}
		if file.IsDir() {
			return
		}
		pluginName := strings.Split(file.Name(), ".")[0]
		f, err := p.Lookup("New")
		if err != nil {
			panic(err)
		}
		Register(pluginName, reflect.ValueOf(f).Call(nil)[0].Interface())
	}
}
