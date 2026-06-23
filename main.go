package main

import (
	"errors"
	"net"
	"net/http"

	"github.com/davenfroberg/wireyak/metrics"
	"github.com/davenfroberg/wireyak/parsing"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const INTERFACE = "en0"

func main() {
	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	go http.ListenAndServe(":2112", nil) // open metrics endpoint for Prometheus to poll

	myMac, err := getMacAddr()

	if err != nil {
		panic(err)
	} else if handle, err := pcap.OpenLive(INTERFACE, 3200, true, pcap.BlockForever); err != nil {
		panic(err)
	} else {
		parser := parsing.NewParser(myMac, metrics)
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			parser.ParsePacket(packet)
		}
	}
}

func getMacAddr() (string, error) {
	inf, err := net.InterfaceByName(INTERFACE)
	if err != nil {
		return "", err
	}
	if inf.HardwareAddr == nil {
		return "", errors.New("No MAC address associated with interface")
	}
	return inf.HardwareAddr.String(), nil
}
