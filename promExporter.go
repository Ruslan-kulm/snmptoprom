package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

func WebServer(rawTraps <-chan IntStatuTrap) {

	opsQueued := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "AO",
			Subsystem: "network",
			Name:      "interface_status",
			Help:      "Interface status changes oh network devices in AO",
		},
		[]string{
			"hostname",
			"interface",
		},
	)
	prometheus.MustRegister(opsQueued)

	for trap := range rawTraps {
		status := 0.0
		if strings.Contains(trap.InterfaceStatus, "up") {
			status = 1.0
		}

		opsQueued.WithLabelValues(trap.DeviceName, trap.InterfaceName).Set(status)
	}

}
