package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"time"
)

const MetricsNamespace = "httpserver"

var (
	latencyFn = CreateExecutionTimeMetric(MetricsNamespace, "Time spent.")
)

func Register() {
	err := prometheus.Register(latencyFn)
	if err != nil {
		log.Fatal().Err(err).Msg("Register prometheus fails.")
	}
}

func CreateExecutionTimeMetric(namespace string, help string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "execution_latency_seconds",
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{"step"},
	)
}

type ExecutionTimer struct {
	histo *prometheus.HistogramVec
	start time.Time
}

func NewTimer() *ExecutionTimer {
	return NewExecutionTimer(latencyFn)
}

func NewExecutionTimer(histogramVec *prometheus.HistogramVec) *ExecutionTimer {
	now := time.Now()
	return &ExecutionTimer{
		histo: histogramVec,
		start: now,
	}
}

func (t *ExecutionTimer) ObserveTotal() {
	(*t.histo).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}
