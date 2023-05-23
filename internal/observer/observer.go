package observer

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/pkg/recws"
	"golang.org/x/exp/slog"
)

var (
	srv  *service.Service
	stop = make(chan struct{}, 2)
)

func Run(dt string) {
	srv = service.New(stop)
	srv.Run()
	defer srv.Close()
	for {
		switch dt {
		case "substrate":
			conn := &recws.RecConn{KeepAliveTimeout: 5 * time.Second, WriteTimeout: time.Second * 5, ReadTimeout: 60 * time.Second}
			conn.Dial(util.WSEndPoint, nil)
			go srv.Subscribe(conn, stop)
			slog.Debug("Connected to substrate node")
		default:
			log.Fatalf("no such daemon component: %s", dt)
		}
		enableTermSignalHandler()
		if _, ok := <-stop; !ok {
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
		close(stop)
	}()
}
