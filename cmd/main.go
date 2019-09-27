package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"runtime"
	"subscan-end/internal/jobs"
	"subscan-end/internal/server/http"
	"subscan-end/internal/service"
	"subscan-end/libs/substrate"
	"subscan-end/sub"
	"syscall"
	"time"
)

func main() {
	defer log.Close()
	defer substrate.CloseWsConnection()
	if err := setupApp().Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Subscan"
	app.Usage = "Subscan End"
	app.Version = "1.0"
	app.Action = func() error { runServe(); return nil }
	app.Commands = sub.Commands
	app.Flags = []cli.Flag{cli.StringFlag{Name: "conf"}}
	app.Before = func(context *cli.Context) error {
		flag.Parse()
		if err := paladin.Init(); err != nil {
			panic(err)
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
	log.Info("subscan-end run ......")
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
			_ = httpSrv.Shutdown(ctx)
			log.Info("subscan-end exit")
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
