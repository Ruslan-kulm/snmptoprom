package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"strings"
)

var curIntStatu *prometheus.GaugeVec
var sumIntStatu *prometheus.GaugeVec

func CreateMetrics() {
	fmt.Printf("CreateMetrics \n")
	curIntStatu = prometheus.NewGaugeVec(
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

	sumIntStatu = prometheus.NewGaugeVec(
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

	fmt.Printf("CreateMetrics Done \n")
}

func UpdateIntStatuMetrics(trap IntStatuTrap, ctx *Context) {
	status := 0.0
	if strings.Contains(trap.InterfaceStatus, "up") {
		status = 1.0
	}

	curIntStatu.WithLabelValues(trap.DeviceName, trap.InterfaceName).Set(status)
	sumIntStatu.WithLabelValues(trap.DeviceName, trap.InterfaceName).Inc()
	ctx.logger.WithFields(logrus.Fields{
		"id": trap.Id,
	}).Debug("Metrics updated")

}
