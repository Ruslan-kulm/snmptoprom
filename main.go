package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
)

type IntStatuTrap struct {
	DeviceName      string
	IpAddr          net.IP
	InterfaceName   string
	InterfaceStatus string
}

type Config struct {
	TrapListener struct {
		Port string `envconfig:"SNMP_PORT"`
	}
	TrapHandler struct {
		Community string `envconfig:"SNMP_COMMUNITY"`
	}
	PromExporter struct {
		PromPort string `envconfig:"PROM_PORT"`
	}
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	fmt.Println(cfg)
	if err != nil {
		log.Panicf("error get env variables: %s", err)
	}

	rawTraps := make(chan IntStatuTrap, 1000)
	handledTraps := make(chan IntStatuTrap, 1000)

	trapListener := onNewTrap(rawTraps)
	go CreateTrapListener(trapListener, cfg)
	go TrapHandler(rawTraps, handledTraps, cfg)
	go WebServer(handledTraps)

	http.Handle("/metrics", promhttp.Handler())
	http_err := http.ListenAndServe(":"+cfg.PromExporter.PromPort, nil)
	if http_err != nil {
		log.Panicf("error in listen: %s", http_err)
	}
}
