package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"пше ыефегы
)

type IntStatuTrap struct {
	DeviceName      string
	IpAddr          net.IP
	InterfaceName   string
	InterfaceStatus string
}

func main() {
	cfg := GetConfig()

	rawTraps := make(chan IntStatuTrap, 1000)
	handledTraps := make(chan IntStatuTrap, 1000)

	trapListener := onNewTrap(rawTraps)
	go CreateTrapListener(trapListener, *cfg)
	go TrapHandler(rawTraps, handledTraps, *cfg)
	go WebServer(handledTraps)

	http.Handle("/metrics", promhttp.Handler())
	http_err := http.ListenAndServe(":"+cfg.PromExporter.PromPort, nil)
	if http_err != nil {
		log.Panicf("error in listen: %s", http_err)
	}
}
