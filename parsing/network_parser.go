package parsing

import (
	"bytes"
	"net"
	"strings"

	"github.com/davenfroberg/wireyak/metrics"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type PacketParser struct {
	Metrics  *metrics.Metrics
	LocalMAC net.HardwareAddr
}

func NewPacketParser(mac net.HardwareAddr, m *metrics.Metrics) *PacketParser {
	return &PacketParser{
		Metrics:  m,
		LocalMAC: mac,
	}
}

func (p *PacketParser) ParsePacket(packet gopacket.Packet) {
	packetTags := getPacketTags(packet, p.LocalMAC)
	p.Metrics.IncPacket(packetTags)
	p.Metrics.AddBytesTransmitted(packetTags, packet.Metadata().Length)
}

func getPacketTags(packet gopacket.Packet, localMAC net.HardwareAddr) metrics.PacketTags {
	// tag 1: Direction
	directionTag := "rx"
	if isOutgoing(packet, localMAC) {
		directionTag = "tx"
	}

	// tag 2: network layer
	networkTag := "none"
	if networkLayer := packet.NetworkLayer(); networkLayer != nil {
		networkTag = strings.ToLower(networkLayer.LayerType().String())
	}

	// tag 3/4: transport layer + application layer
	transportTag := "none"
	applicationTag := "none"
	if transportLayer := packet.TransportLayer(); transportLayer != nil {
		transportTag = strings.ToLower(transportLayer.LayerType().String())
		applicationTag = getApplicationLabel(transportLayer)
	}
	if applicationTag == "https" && transportTag == "udp" {
		applicationTag = "quic"
	}

	return metrics.PacketTags{
		Direction:   directionTag,
		Network:     networkTag,
		Transport:   transportTag,
		Application: applicationTag,
	}
}

func getApplicationLabel(transportLayer gopacket.Layer) string {
	// dns, http, https, anything else?
	srcPort := uint16(0)
	dstPort := uint16(0)

	if tcp, ok := transportLayer.(*layers.TCP); ok {
		srcPort = uint16(tcp.SrcPort)
		dstPort = uint16(tcp.DstPort)
	} else if udp, ok := transportLayer.(*layers.UDP); ok {
		srcPort = uint16(udp.SrcPort)
		dstPort = uint16(udp.DstPort)
	} else {
		return "unknown"
	}

	if srcPort == 53 || dstPort == 53 {
		return "dns"
	} else if srcPort == 80 || dstPort == 80 {
		return "http"
	} else if srcPort == 443 || dstPort == 443 {
		return "https"
	}

	return "none"
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
