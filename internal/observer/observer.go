package observer

import (
	"fmt"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/substrate-api-rpc/pkg/recws"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/util"
	"github.com/sevlyar/go-daemon"
)

var (
	srv *service.Service
)

func Run(dt, signal string) {
	daemon.AddCommand(daemon.StringFlag(&signal, "stop"), syscall.SIGQUIT, termHandler)
	doAction(dt)
}

func doAction(dt string) {
	if !util.StringInSlice(dt, dao.DaemonAction) {
		log.Println("no such daemon")
		return
	}

	logDir := util.GetEnv("LOG_DIR", "../log/")
	pid := fmt.Sprintf("%s%s_pid", logDir, dt)
	logName := fmt.Sprintf("%s%s_log", logDir, dt)

	dc := &daemon.Context{
		PidFileName: pid,
		PidFilePerm: 0644,
		LogFileName: logName,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        nil,
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := dc.Search()
		if err != nil {
			log.Println(dt, "not running")
		} else {
			_ = daemon.SendCommands(d)
		}
		return
	}

	d, err := dc.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer func() {
		err = dc.Release()
		if err != nil {
			log.Println("Error:", err)
		}
	}()

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("daemon started")

	go doRun(dt)

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
	log.Println("daemon terminated")
}

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func doRun(dt string) {
	srv = service.New()
	defer srv.Close()
LOOP:
	for {
		if dt == "substrate" {
			interrupt := make(chan os.Signal, 1)
			subscribeConn := &recws.RecConn{KeepAliveTimeout: 10 * time.Second}
			subscribeConn.Dial(util.WSEndPoint, nil)
			srv.Subscribe(subscribeConn, interrupt)
		} else {
			go heartBeat(dt)
			switch dt {
			default:
				break LOOP
			}
		}
		if _, ok := <-stop; ok {
			break LOOP
		}
	}
	done <- struct{}{}
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

func heartBeat(dt string) {
	for {
		srv.SetHeartBeat(fmt.Sprintf("%s:heartBeat:%s", util.NetworkNode, dt))
		time.Sleep(10 * time.Second)
	}
}
