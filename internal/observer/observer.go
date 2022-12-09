package observer

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itering/substrate-api-rpc/pkg/recws"

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
	for {
		switch dt {
		case "substrate":
			subscribeConn := &recws.RecConn{KeepAliveTimeout: 10 * time.Second, WriteTimeout: time.Second * 5, ReadTimeout: 10 * time.Second}
			subscribeConn.Dial(util.WSEndPoint, nil)
			go srv.Subscribe(subscribeConn, stop)
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
		log.Printf("Received signal %s, exiting...\n", <-sigs)
		stop <- struct{}{}
	}()
}
