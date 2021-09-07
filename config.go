package main

import (
	"os"
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
}

func GetConfig() *Config {
	snmpPort := os.Getenv("SNMP_PORT")
	community := os.Getenv("SNMP_COMMUNITY")
	promPort := os.Getenv("PROM_PORT")

	if snmpPort == "" {
		snmpPort = "162"
	}
	if community == "" {
		community = "public"
	}
	if promPort == "" {
		promPort = "8080"
	}

	cfg := new(Config)
	cfg.TrapListener.Port = snmpPort
	cfg.TrapHandler.Community = community
	cfg.PromExporter.PromPort = promPort

	return cfg
}
