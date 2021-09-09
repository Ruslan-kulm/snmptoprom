package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const appName = "snmpToProm"
const appVer = "1.0.0"

func Init() *Context {
	ctx := new(Context)
	cfg := GetConfig()
	ctx.config = cfg

	logger := GetLogger(ctx)
	ctx.logger = logger
	ctx.logger.Info("Logger initialised")

	ctx.logger.WithFields(log.Fields{
		"SNMP_PORT":      cfg.TrapListener.Port,
		"SNMP_COMMUNITY": cfg.TrapHandler.Community,
		"PROM_PORT":      cfg.PromExporter.PromPort,
	}).Info("Config initialised")

	return ctx
}

func main() {
	ctx := Init()

	rawTraps := make(chan IntStatuTrap, 1000)
	handledTraps := make(chan IntStatuTrap, 1000)

	trapListener := onNewTrap(rawTraps, ctx)
	go CreateTrapListener(trapListener, ctx)
	go TrapHandler(rawTraps, handledTraps, ctx)
	go WebServer(handledTraps, ctx)

	http.Handle("/metrics", promhttp.Handler())
	ctx.logger.Info("Listening %s", ":"+ctx.config.PromExporter.PromPort)
	err := http.ListenAndServe(":"+ctx.config.PromExporter.PromPort, nil)
	if err != nil {
		ctx.logger.WithFields(log.Fields{
			"error": err,
		}).Error("error on listen %s", ":"+ctx.config.PromExporter.PromPort)
		os.Exit(1)
	}
}
