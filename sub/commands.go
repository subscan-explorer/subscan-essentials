package sub

import (
	"fmt"
	"github.com/urfave/cli"
	"subscan-end/daemons"
	"subscan-end/internal/dao"
	"subscan-end/utiles"
)

var (
	Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "Start a daemon",
			Action: func(c *cli.Context) error {
				dt := c.Args().Get(0)
				if utiles.StringInSlice(dt, dao.DaemonAction) {
					daemons.Run(dt, "start")
				} else {
					fmt.Println("no such daemon: " + dt)
				}
				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "Stop a daemon",
			Action: func(c *cli.Context) error {
				dt := c.Args().Get(0)
				if utiles.StringInSlice(dt, dao.DaemonAction) {
					daemons.Run(dt, "stop")
				} else {
					fmt.Println("no such daemon: " + dt)
				}
				return nil
			},
		},
		{
			Name:  "status",
			Usage: "Get a daemon status",
			Action: func(c *cli.Context) error {
				dt := c.Args().Get(0)
				if utiles.StringInSlice(dt, dao.DaemonAction) {
					daemons.Run(dt, "status")
				} else {
					fmt.Println("no such daemon: " + dt)
				}
				return nil
			},
		},
		{
			Name:  "repairBlock",
			Usage: "Repair All Block",
			Action: func(c *cli.Context) error {
				daemons.RepairBlock()
				return nil
			},
		},
		{
			Name:  "repairBlockData",
			Usage: "Repair All Block Data",
			Action: func(c *cli.Context) error {
				daemons.RepairBlockData()
				return nil
			},
		},
		{
			Name:  "RepairValidateInfo",
			Usage: "Repair Validate Info",
			Action: func(c *cli.Context) error {
				daemons.RepairValidateInfo()
				return nil
			},
		},
	}
)
