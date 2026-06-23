package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	PacketTotals *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		PacketTotals: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "packet_count",
				Help: "Number of packets seen",
			},
			[]string{"protocol"},
		),
	}
	return m
}
