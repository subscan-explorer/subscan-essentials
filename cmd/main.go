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
	flag.Parse()
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	log.Init(nil)
	defer afterClose()
	jobs.Init()
	if err := setupApp().Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setupApp() *cli.App {
	return &cli.App{
		Name:  "Subscan",
		Usage: "Subscan End",
		Action: func(c *cli.Context) error {
			if len(c.Args()) == 0 {
				runServe()
			}
			return nil
		},
		Version:  "1.0",
		Commands: sub.Commands,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "config, conf"},
		},
		Before: func(context *cli.Context) error {
			runtime.GOMAXPROCS(runtime.NumCPU())
			return nil
		},
	}
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

func afterClose() {
	_ = log.Close()
	substrate.CloseWsConnection()
}
