package daemons

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/internal/util"
	"github.com/sevlyar/go-daemon"
)

var (
	srv *service.Service
)

func Run(dt, signal string) {
	daemon.AddCommand(daemon.StringFlag(&signal, "stop"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(&signal, "status"), syscall.SIGSEGV, statusHandler)
	doAction(dt, signal)
}

func doAction(dt, signal string) {
	pid := fmt.Sprintf("../log/%s_pid", dt)
	logName := fmt.Sprintf("../log/%s_log", dt)
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
			if signal == "status" {
				log.Println(dt, "running", "pid", d.Pid)
			}
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
			Subscribe()
		} else {
			go heartBeat(dt)
			switch dt {
			case "worker":
				RunWorker()
			case "cronWorker":
				go DoCron(stop)
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

func statusHandler(sig os.Signal) error {
	log.Println("configuration status", sig)
	return nil
}

func heartBeat(dt string) {
	for {
		setHeartBeat(dt)
		time.Sleep(time.Duration(10) * time.Second)
	}
}

func setHeartBeat(dt string) {
	srv.SetHeartBeat(fmt.Sprintf("%s:heartBeat:%s", util.NetworkNode, dt))
}
