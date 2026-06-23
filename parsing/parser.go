package parsing

import (
	"bytes"
	"net"

	"github.com/davenfroberg/wireyak/metrics"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type LayerHandler func(l gopacket.Layer, m *metrics.Metrics, isOutgoing bool)

type Parser struct {
	Metrics  *metrics.Metrics
	Handlers map[gopacket.LayerType]LayerHandler
	LocalMAC net.HardwareAddr
}

func NewParser(mac net.HardwareAddr, m *metrics.Metrics) *Parser {
	return &Parser{
		Metrics: m,
		Handlers: map[gopacket.LayerType]LayerHandler{
			layers.LayerTypeIPv4: parseV4Layer,
			layers.LayerTypeIPv6: parseV6Layer,
			layers.LayerTypeDNS:  parseDnsLayer,
		},
		LocalMAC: mac,
	}
}

func (p *Parser) ParsePacket(packet gopacket.Packet) {
	isOutgoing := isOutgoing(packet, p.LocalMAC)

	for _, layer := range packet.Layers() {
		handler, ok := p.Handlers[layer.LayerType()]
		if ok {
			handler(layer, p.Metrics, isOutgoing)
		}
	}
}

func isOutgoing(packet gopacket.Packet, localMAC net.HardwareAddr) bool {
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethLayer != nil {
		eth, ok := ethLayer.(*layers.Ethernet)
		if ok {
			return bytes.Equal(eth.SrcMAC, localMAC)
		}
	}
	return false
}
