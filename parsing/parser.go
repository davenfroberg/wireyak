package parsing

import (
	"github.com/davenfroberg/wireyak/metrics"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type LayerHandler func(l gopacket.Layer, m *metrics.Metrics)

type Parser struct {
	Metrics  *metrics.Metrics
	Handlers map[gopacket.LayerType]LayerHandler
}

func NewParser(m *metrics.Metrics) *Parser {
	return &Parser{
		Metrics: m,
		Handlers: map[gopacket.LayerType]LayerHandler{
			layers.LayerTypeIPv4: parseV4Layer,
			layers.LayerTypeIPv6: parseV6Layer,
			layers.LayerTypeDNS:  parseDnsLayer,
		},
	}
}

func (p *Parser) ParsePacket(packet gopacket.Packet) {
	for _, layer := range packet.Layers() {
		handler, ok := p.Handlers[layer.LayerType()]
		if ok {
			handler(layer, p.Metrics)
		}
	}
}
