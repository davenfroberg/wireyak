package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	PacketCount      *prometheus.CounterVec
	BytesTransmitted *prometheus.CounterVec
}

type PacketTags struct {
	Direction   string
	Network     string
	Transport   string
	Application string
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

func (m *Metrics) IncPacket(tags PacketTags) {
	// this order should match the order in metric definition
	m.PacketCount.WithLabelValues(
		tags.Direction,
		tags.Network,
		tags.Transport,
		tags.Application,
	).Inc()
}
