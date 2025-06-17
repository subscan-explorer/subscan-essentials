package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(
		// block
		subBlockStatusGauge, SubBlockFillError,
		// worker
		WorkerProcessCost,
	)
}
