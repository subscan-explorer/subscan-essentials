package plugins

type PluginFactory Plugin

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

// func init() {
// 	pluginsDir := "../configs/plugins"
// 	files, err := ioutil.ReadDir(pluginsDir)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	for _, f := range files {
// 		p, err := plugin.Open(fmt.Sprintf("%s/%s", pluginsDir, f.Name()))
// 		if err != nil {
// 			panic(err)
// 		}
// 		f, err := p.Lookup("New")
// 		if err != nil {
// 			panic(err)
// 		}
// 		Register("account", reflect.ValueOf(f).Call(nil)[0].Interface())
// 	}
// }
