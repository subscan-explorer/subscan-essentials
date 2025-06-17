package prom

import (
	"fmt"
	"net/http"
	"time"

	"github.com/itering/subscan/util"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var MetricsServer *http.Server

func init() {
	port := util.GetEnv("SUBSCAN_PROM_PORT", "8082")
	MetricsServer = &http.Server{Addr: fmt.Sprintf(":%s", port)}
}

func New() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/readiness", readinessProbe)
	time.AfterFunc(time.Second*5, func() {
		ready = true
	})

	if err := MetricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		util.Logger().Error(fmt.Errorf("metrics server failed to listen: %v", err))
	}
}

func healthz(rsp http.ResponseWriter, _ *http.Request) {
	rsp.WriteHeader(http.StatusOK)
}

var ready bool

func readinessProbe(rsp http.ResponseWriter, _ *http.Request) {
	if ready {
		rsp.WriteHeader(http.StatusOK)
	} else {
		rsp.WriteHeader(http.StatusBadRequest)
	}
}
