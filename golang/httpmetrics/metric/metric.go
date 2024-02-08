package metric

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (m *metrics) WrapHandler(path string, handler http.HandlerFunc) http.HandlerFunc {
	label := prometheus.Labels{
		"handler": path,
	}

	httpRequestTotal := m.httpRequestsTotal.MustCurryWith(label)
	httpDuration := m.httpDuration.MustCurryWith(label)

	instrumentChain := promhttp.InstrumentHandlerCounter(
		httpRequestTotal, promhttp.InstrumentHandlerDuration(
			httpDuration, handler,
		),
	)

	return instrumentChain.ServeHTTP
}

type metrics struct {
	httpRequestsTotal *prometheus.CounterVec
	httpDuration      *prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"code", "method", "handler"},
		),
		httpDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "The HTTP request latencies in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"code", "method", "handler"},
		),
	}

	reg.MustRegister(m.httpRequestsTotal, m.httpDuration)
	return m
}
