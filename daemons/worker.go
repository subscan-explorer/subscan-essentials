package daemons

import (
	"github.com/freehere107/go-workers"
	"subscan-end/internal/jobs"
)

func RunWorker() {
	jobs.Init()
	RegWorker()
	go workers.StatsServer(8080)
	workers.Run()
}

func RegWorker() {
	workers.Process("jobs", myJob, 10)
}

func myJob(message *workers.Msg) {
}
