package observer

import (
	"context"
	"fmt"
	"github.com/itering/subscan/internal/script"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/util"
	"github.com/robfig/cron/v3"
)

var (
	srv  *service.Service
	stop = make(chan struct{}, 2)
)

func Run(dt string) {
	srv = service.New()
	defer srv.Close()
	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)
	go enableTermSignalHandler(cancel)
	switch dt {
	case "subscribe":
		wg.Add(1)
		go func() {
			defer wg.Done()
			srv.Subscribe(ctx)
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			RunCron()
		}()
	case "worker":
		wg.Add(1)
		go func() {
			defer wg.Done()
			Consumption()
		}()
	default:
		panic(fmt.Sprintf("no such daemon component: %s", dt))
	}
	time.Sleep(3 * time.Second)
	wg.Wait()

}

func enableTermSignalHandler(cancel func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	util.Logger().Info(fmt.Sprintf("Received signal %s, exiting...\n", <-sigs))
	cancel()
	close(stop)
}

func RunCron() {
	// or use cron.DefaultLogger
	c := cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))
	if _, err := c.AddFunc("@every 3m", func() {
		script.RefreshMetadata()
	}); err != nil {
		util.Logger().Error(fmt.Sprintf("Failed to register cron job: %v", err))
		os.Exit(1)
	}
	c.Start()
	<-stop
	<-c.Stop().Done()
	util.Logger().Info("Cron stopped")
}
