package main

import (
	"fmt"
	"github.com/itering/subscan/internal/daemon"
	"github.com/itering/subscan/internal/daemon/script"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/internal/service"
	"os"

	"github.com/itering/subscan/internal/util"
	"github.com/urfave/cli"
)

var (
	Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "Start a daemon",
			Action: func(c *cli.Context) error {
				call(c, "start")
				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "Stop a daemon",
			Action: func(c *cli.Context) error {
				call(c, "stop")
				return nil
			},
		},
		{
			Name:  "status",
			Usage: "Get a daemon status",
			Action: func(c *cli.Context) error {
				call(c, "status")
				return nil
			},
		},
		{
			Name:  "migrate",
			Usage: "migrate database",
			Action: func(c *cli.Context) error {
				srv := service.New()
				_ = os.Setenv("TASK_MOD", "true")
				srv.Migration()
				fmt.Println("migrate success")
				return nil
			},
		},
		{
			Name:  "install",
			Usage: "install with new network",
			Action: func(c *cli.Context) error {
				_ = os.Setenv("TASK_MOD", "true")
				script.Install()
				return nil
			},
		},
	}
)

func call(c *cli.Context, signal string) {
	dt := c.Args().Get(0)
	if util.StringInSlice(dt, dao.DaemonAction) {
		daemons.Run(dt, signal)
	} else {
		fmt.Println("no such daemon: " + dt)
	}
}
