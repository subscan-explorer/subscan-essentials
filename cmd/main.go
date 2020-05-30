package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan/internal/di"
	"github.com/itering/subscan/internal/jobs"
	"github.com/itering/subscan/internal/substrate/websocket"
)

func main() {
	defer func() {
		_ = log.Close()
		websocket.CloseWsConnection()
	}()

	// init configs
	err := flag.Set("conf", "../configs")
	if err != nil {
		panic(err)
	}
	err = paladin.Init()
	if err != nil {
		panic(err)
	}
	jobs.Init()
	log.Init(nil)
	runtime.GOMAXPROCS(runtime.NumCPU())

	// start service
	_, closeFunc, err := di.InitApp()
	if err != nil {
		panic(err)
	}

	// handle signals
	c := make(chan os.Signal, 1)
	log.Info("SubScan End run ......")
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeFunc()
			log.Info("SubScan End exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
