package main

import (
	"github.com/itering/subscan/internal/observer"
	"github.com/itering/subscan/internal/script"
	"github.com/urfave/cli"
)

var commands = []cli.Command{
	{
		Name:  "start",
		Usage: "Start one worker, E.g. subscribe",
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
		Action: func(c *cli.Context) error {
			script.CheckCompleteness()
			return nil
		},
	},
}
