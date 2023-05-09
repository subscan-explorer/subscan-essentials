package observer

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itering/subscan/internal/service"
	"golang.org/x/exp/slog"
)

var (
	srv  *service.Service
	stop = make(chan struct{}, 2)
)

func Run(dt string) {
	srv = service.New()
	defer srv.Close()
	for {
		switch dt {
		case "substrate":

			go srv.Subscribe(stop)
			slog.Debug("Connected to substrate node")
		default:
			log.Fatalf("no such daemon component: %s", dt)
		}
		enableTermSignalHandler()
		if _, ok := <-stop; ok {
			time.Sleep(3 * time.Second)
			break
		}
	}
}

func enableTermSignalHandler() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		slog.Info("Received signal, exiting...", "signal", <-sigs)
		stop <- struct{}{}
	}()
}
