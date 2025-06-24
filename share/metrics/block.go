package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	subBlockStatusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "subscan",
			Subsystem: "substrate",
			Name:      "block_status",
			Help:      "Subscan block status statistics",
		}, []string{"status"},
	)
	SubBlockFillError = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "subscan",
			Subsystem: "substrate",
			Name:      "block_fill_error",
			Help:      "The number of error occurred when exec FillBlockData",
		},
	)
)

func SubBlockGauge(status string, val uint64) {
	subBlockStatusGauge.WithLabelValues(status).Set(float64(val))
}
