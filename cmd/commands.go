package main

import (
	"github.com/itering/subscan/internal/observer"
	"github.com/itering/subscan/internal/script"
	"github.com/urfave/cli"
)

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
}
