package main

import (
	"net/http"

	"github.com/davenfroberg/wireyak/metrics"
	"github.com/davenfroberg/wireyak/parsing"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	go http.ListenAndServe(":2112", nil) // open metrics endpoint for Prometheus to poll

	// start packet capture
	if handle, err := pcap.OpenLive("en0", 3200, true, pcap.BlockForever); err != nil {
		panic(err)
	} else {
		parser := parsing.NewParser(metrics)
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			parser.ParsePacket(packet)
		}
	}
}
