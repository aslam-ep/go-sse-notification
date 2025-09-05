package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ActiveClients = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "sse_active_clients",
			Help: "Current number of SSE clients",
		},
	)

	MessageCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sse_message_sent_total",
			Help: "Total number of message sent to clients",
		},
		[]string{"user"},
	)

	DroppedMessages = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "sse_dropped_message_total",
			Help: "Number of dropped message due to slow clients",
		},
	)
)
