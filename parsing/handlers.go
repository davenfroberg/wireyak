package parsing

import (
	"github.com/davenfroberg/wireyak/metrics"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func parseV4Layer(l gopacket.Layer, m *metrics.Metrics, isOutgoing bool) {
	_, ok := l.(*layers.IPv4)
	if ok {
		direction := "rx"
		if isOutgoing {
			direction = "tx"
		}
		m.PacketCount.WithLabelValues("ipv4", direction).Inc()
	}
}

func parseV6Layer(l gopacket.Layer, m *metrics.Metrics, isOutgoing bool) {
	_, ok := l.(*layers.IPv6)
	if ok {
		direction := "rx"
		if isOutgoing {
			direction = "tx"
		}
		m.PacketCount.WithLabelValues("ipv6", direction).Inc()
	}
}

func parseDnsLayer(l gopacket.Layer, m *metrics.Metrics, isOutgoing bool) {
	_, ok := l.(*layers.DNS)
	if ok {
		direction := "rx"
		if isOutgoing {
			direction = "tx"
		}
		m.PacketCount.WithLabelValues("dns", direction).Inc()
	}
}
