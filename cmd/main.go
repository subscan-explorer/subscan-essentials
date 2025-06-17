package main

import (
	"context"
	"fmt"
	"github.com/itering/subscan/internal/server/prom"
	"github.com/itering/subscan/util/mq"
	redisUtil "github.com/itering/subscan/util/redis"
	"os"
	"runtime"

	"github.com/go-kratos/kratos/v2"
	"github.com/itering/subscan/configs"
	"github.com/itering/subscan/internal/server/http"
	"github.com/itering/subscan/internal/service"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/urfave/cli"
)

func main() {
	defer func() {
		websocket.Close()
	}()
	if err := setupApp().Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Name = "SUBSCAN"
	app.Usage = "SUBSCAN Backend Service, use -h get help"
	app.Version = "2.0"
	app.Description = "SubScan Backend Service, substrate blockchain explorer"
	app.Flags = []cli.Flag{cli.StringFlag{Name: "conf", Value: "../configs"}}
	app.Before = func(context *cli.Context) error {
		configs.Init()
		redisUtil.Init()
		go prom.New()
		runtime.GOMAXPROCS(runtime.NumCPU())
		mq.New()
		return nil
	}
	app.Commands = commands
	app.Action = func(*cli.Context) error { run(); return nil }
	app.After = func(ctx *cli.Context) error {
		_ = prom.MetricsServer.Shutdown(context.Background())
		return nil
	}
	return app
}

func run() {
	svc := service.New()
	httpSrv := http.NewHTTPServer(configs.Boot.Server, svc)
	defer func() {
		// Micro services
		svc.Close()
	}()

	app := kratos.New(kratos.Metadata(map[string]string{}), kratos.Server(httpSrv))
	if err := app.Run(); err != nil {
		panic(err)
	}
}
