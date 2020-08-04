package plugins

import (
	"github.com/itering/subscan-plugin"
	"github.com/itering/subscan/plugins/balance"
)

func init() {
	registerNative()
	// registerStatic()
}

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

func registerNative() {
	Register("balance", balance.New())
}

// Currently the go plugin solution is not stable yet
// func registerStatic() {
// 	flag.Parse()
//
// 	pluginsDir := "../configs/plugins"
// 	if confFlag := flag.Lookup("conf"); confFlag != nil {
// 		pluginsDir = fmt.Sprintf("%s/plugins", flag.Lookup("conf").Value)
// 	}
//
// 	files, _ := ioutil.ReadDir(pluginsDir)
// 	for _, file := range files {
// 		p, err := plugin.Open(fmt.Sprintf("%s/%s", pluginsDir, file.Name()))
// 		if err != nil {
// 			panic(err)
// 		}
// 		if file.IsDir() {
// 			return
// 		}
// 		pluginName := strings.Split(file.Name(), ".")[0]
// 		f, err := p.Lookup("New")
// 		if err != nil {
// 			panic(err)
// 		}
// 		Register(pluginName, reflect.ValueOf(f).Call(nil)[0].Interface())
// 	}
// }
