package fdyn

import (
	"github.com/coredns/coredns/plugin"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Variables declared for monitoring.
var (
	RequestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "fdyn",
		Name:      "request_count_total",
		Help:      "Counter of requests made per upstream.",
	}, []string{"to"})
	RcodeCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "fdyn",
		Name:      "response_rcode_count_total",
		Help:      "Counter of requests made per upstream.",
	}, []string{"rcode", "to"})
	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: plugin.Namespace,
		Subsystem: "fdyn",
		Name:      "request_duration_seconds",
		Buckets:   plugin.TimeBuckets,
		Help:      "Histogram of the time each request took.",
	}, []string{"to"})
	HealthcheckFailureCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "fdyn",
		Name:      "healthcheck_failure_count_total",
		Help:      "Counter of the number of failed healtchecks.",
	}, []string{"to"})
	HealthcheckBrokenCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "fdyn",
		Name:      "healthcheck_broken_count_total",
		Help:      "Counter of the number of complete failures of the healtchecks.",
	})
	SocketGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: plugin.Namespace,
		Subsystem: "fdyn",
		Name:      "sockets_open",
		Help:      "Gauge of open sockets per upstream.",
	}, []string{"to"})
)
