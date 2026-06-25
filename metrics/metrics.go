package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	PacketCount      *prometheus.CounterVec
	BytesTransmitted *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		PacketCount: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "packet_count",
				Help: "Number of packets seen",
			},
			[]string{"direction", "network", "transport", "application"},
		),
		BytesTransmitted: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "bytes_transmitted",
				Help: "Number of bytes transmitted",
			},
			[]string{"direction", "network", "transport", "application"},
		),
	}
	return m
}
