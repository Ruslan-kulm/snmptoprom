package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

func WebServer(rawTraps <-chan IntStatuTrap) {

	curIntStatu := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "AO",
			Subsystem: "network",
			Name:      "interface_status",
			Help:      "Interface status oh network devices in AO",
		},
		[]string{
			"hostname",
			"interface",
		},
	)
	prometheus.MustRegister(curIntStatu)

	sumIntStatu := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "AO",
			Subsystem: "network",
			Name:      "interface_status_changes",
			Help:      "Interface status changes oh network devices in AO",
		},
		[]string{
			"hostname",
			"interface",
		},
	)
	prometheus.MustRegister(sumIntStatu)

	for trap := range rawTraps {
		status := 0.0
		if strings.Contains(trap.InterfaceStatus, "up") {
			status = 1.0
		}

		curIntStatu.WithLabelValues(trap.DeviceName, trap.InterfaceName).Set(status)
		sumIntStatu.WithLabelValues(trap.DeviceName, trap.InterfaceName).Inc()
	}

}
