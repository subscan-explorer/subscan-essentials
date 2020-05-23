package daemons

import (
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan/daemons/script"
	"github.com/itering/subscan/daemons/tasks"
	"github.com/robfig/cron/v3"
)

type cronLog struct{}

func (l cronLog) Print(v ...interface{}) {
	log.Error(strings.Repeat("%v ", len(v)), v...)
}

func (l cronLog) Info(msg string, v ...interface{}) {
	log.Info("%v", msg)
	log.Error(strings.Repeat("%v ", len(v)), v...)
}

func (l cronLog) Error(err error, msg string, v ...interface{}) {
	log.Info("err: %v,msg: %v", err, msg)
	log.Error(strings.Repeat("%v ", len(v)), v...)
}

func DoCron(stop chan struct{}) {
	c := cron.New(cron.WithChain(
		cron.Recover(cronLog{}), // or use cron.DefaultLogger
	))

	if _, err := c.AddFunc("* * * * *", func() {
		tasks.RefreshMetadata(srv)
	}); err != nil {
		panic(err)
	}

	// every hour
	if _, err := c.AddFunc("@hourly", func() {
		script.RefreshAccountInfo(srv)
	}); err != nil {
		panic(err)
	}

	if _, err := c.AddFunc("@hourly", func() {
		script.RepairCodecError("finalized = 0")
	}); err != nil {
		panic(err)
	}

	c.Start()
	fmt.Println("cron start")
	<-stop
	c.Stop()
}
