package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan/internal/jobs"
	"github.com/itering/subscan/internal/server/http"
	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/internal/substrate/websocket"
	"github.com/urfave/cli"
)

func main() {
	defer afterClose()
	if err := setupApp().Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Name = "SubScan"
	app.Usage = "SubScan Backend Service, use -h get help"
	app.Version = "1.0"
	app.Action = func(*cli.Context) error { runServe(); return nil }
	// app.Commands = Commands
	app.Description = "SubScan Backend Service, substrate blockchain explorer"
	app.Flags = []cli.Flag{cli.StringFlag{Name: "conf", Value: "../configs"}}
	app.Before = func(context *cli.Context) error {
		if client, err := paladin.NewFile(context.String("conf")); err != nil {
			panic(err)
		} else {
			paladin.DefaultClient = client
		}
		jobs.Init()
		log.Init(nil)
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func runServe() {
	svc := service.New()
	httpSrv := http.New(svc)
	c := make(chan os.Signal, 1)
	log.Info("SubScan End run ......")
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
			_ = httpSrv.Shutdown(ctx)
			log.Info("SubScan End exit")
			svc.Close()
			cancel()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func afterClose() {
	_ = log.Close()
	websocket.CloseWsConnection()
}
