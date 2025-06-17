package main

import (
	"context"
	"github.com/itering/subscan/internal/observer"
	"github.com/itering/subscan/internal/script"
	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/plugins"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/network"
	"github.com/urfave/cli"
)

func init() {
	network.SetCurrent(util.NetworkNode)
}

var commands = []cli.Command{
	{
		Name:  "start",
		Usage: "Start a daemon to subscribe to block or start a worker to process events",
		Action: func(c *cli.Context) error {
			observer.Run(c.Args().Get(0))
			return nil
		},
	},
	{
		Name:  "install",
		Usage: "Install default database and create default conf file",
		Action: func(c *cli.Context) error {
			script.Install(c.Parent().String("conf"))
			return nil
		},
	},
	{
		Name:  "CheckCompleteness",
		Usage: "Create blocks completeness",
		Flags: []cli.Flag{
			cli.UintFlag{Name: "StartBlock", Usage: "start block to check"},
			cli.BoolFlag{Name: "fast", Usage: "fast mode, will publish to worker"},
		},
		Action: func(c *cli.Context) error {
			script.CheckCompleteness(c.Uint("StartBlock"), c.Bool("fast"))
			return nil
		},
	},
	{
		Name:  "refreshMetadata",
		Usage: "refresh metadata",
		Action: func(c *cli.Context) error {
			script.RefreshMetadata()
			return nil
		},
	},
	{
		Name:  "plugin",
		Usage: "plugin sub commands",
		Before: func(c *cli.Context) error {
			srv := service.New()
			_, cancel := context.WithCancel(context.Background())
			c.App.After = func(*cli.Context) error {
				cancel()
				srv.Close()
				return nil
			}
			return nil
		},
		Subcommands: pluginCommands(),
		After:       func(context *cli.Context) error { return nil },
	},
}

func pluginCommands() []cli.Command {
	var cmds []cli.Command
	for name, plugin := range plugins.RegisteredPlugins {
		if !plugin.Enable() {
			continue
		}
		cmds = append(cmds, cli.Command{
			Name: name,
			Before: func(c *cli.Context) error {
				return nil
			},
			Subcommands: plugin.Commands(),
		})
	}
	return cmds
}
