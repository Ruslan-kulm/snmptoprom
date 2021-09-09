package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	dSnmpPort  = "162"
	dCommunity = "public"
	dPromPort  = "8080"
	dLogLevel  = log.DebugLevel
)

type Config struct {
	TrapListener struct {
		Port string
	}
	TrapHandler struct {
		Community string
	}
	PromExporter struct {
		PromPort string
	}
	Logger struct {
		level log.Level
	}
}

func GetConfig() *Config {
	snmpPort := os.Getenv("SNMP_PORT")
	community := os.Getenv("SNMP_COMMUNITY")
	promPort := os.Getenv("PROM_PORT")
	logLevel := os.Getenv("LOG_LEVEL")

	if snmpPort == "" {
		snmpPort = dSnmpPort
	}
	if community == "" {
		community = dCommunity
	}
	if promPort == "" {
		promPort = dPromPort
	}

	var lLevel log.Level
	switch logLevel {
	case "DEBUG":
		lLevel = log.DebugLevel
	case "INFO":
		lLevel = log.InfoLevel
	case "ERROR":
		lLevel = log.ErrorLevel
	default:
		lLevel = dLogLevel
	}

	cfg := new(Config)
	cfg.TrapListener.Port = snmpPort
	cfg.TrapHandler.Community = community
	cfg.PromExporter.PromPort = promPort
	cfg.Logger.level = lLevel

	return cfg
}
