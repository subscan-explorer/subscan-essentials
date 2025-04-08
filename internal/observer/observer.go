package observer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/util"
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
