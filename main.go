package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/davenfroberg/wireyak/metrics"
	"github.com/davenfroberg/wireyak/parsing"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/v3/process"
)

const INTERFACE = "en0"

func main() {
	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	go http.ListenAndServe(":2112", nil) // open metrics endpoint for Prometheus to poll

	myMac, err := getMacAddr()

	go startProcessSniffing(metrics)

	if err != nil {
		panic(err)
	} else {
		startNetworkSniffing(myMac, metrics)
	}
}

func startProcessSniffing(metrics *metrics.Metrics) {
	parser := parsing.NewProcessParser(metrics)

	// poll every 3 seconds indefinitely for new all running processes
	for true {
		processes, err := process.Processes()
		if err != nil {
			log.Fatal(err)
		}
		parser.ParseProcesses(processes)
		time.Sleep(3 * time.Second)
	}
}

func startNetworkSniffing(myMac net.HardwareAddr, metrics *metrics.Metrics) {
	if handle, err := pcap.OpenLive(INTERFACE, 3200, true, pcap.BlockForever); err != nil {
		panic(err)
	} else {
		parser := parsing.NewPacketParser(myMac, metrics)
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			parser.ParsePacket(packet)
		}
	}
}

func getMacAddr() (net.HardwareAddr, error) {
	inf, err := net.InterfaceByName(INTERFACE)
	if err != nil {
		return nil, err
	}
	if inf.HardwareAddr == nil {
		return nil, errors.New("No MAC address associated with interface")
	}
	return inf.HardwareAddr, nil
}
