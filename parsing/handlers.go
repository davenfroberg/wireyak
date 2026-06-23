package parsing

import (
	"github.com/davenfroberg/wireyak/metrics"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func parseV4Layer(l gopacket.Layer, m *metrics.Metrics) {
	_, ok := l.(*layers.IPv4)
	if ok {
		// srcIp := ip.SrcIP.String()
		// dstIp := ip.DstIP.String()
		m.PacketCount.WithLabelValues("ipv4").Inc()
	}
}

func parseV6Layer(l gopacket.Layer, m *metrics.Metrics) {
	_, ok := l.(*layers.IPv6)
	if ok {
		// srcIp := v6.SrcIP.String()
		// dstIp := v6.DstIP.String()
		m.PacketCount.WithLabelValues("ipv6").Inc()
	}
}

func parseDnsLayer(l gopacket.Layer, m *metrics.Metrics) {
	_, ok := l.(*layers.DNS)
	if ok {
		m.PacketCount.WithLabelValues("dns").Inc()
	}
}
