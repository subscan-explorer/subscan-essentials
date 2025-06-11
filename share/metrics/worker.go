package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	WorkerProcessCost = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "subscan",
		Subsystem: "worker",
		Name:      "process_duration_seconds",
		Buckets:   []float64{0.05, 0.1, 0.2, 0.5, 1, 2, 5, 10, 20, 30, 60},
	}, []string{"queue", "class"})
)
