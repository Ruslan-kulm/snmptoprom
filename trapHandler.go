package main

import (
	g "github.com/gosnmp/gosnmp"
	"log"
	"time"
)

const hostNameOID = "1.3.6.1.4.1.9.2.1.3.0"

func TrapHandler(rawTraps <-chan IntStatuTrap, handledTraps chan<- IntStatuTrap, cfg Config) {
	for trap := range rawTraps {
		params := &g.GoSNMP{
			Target:    trap.IpAddr.String(),
			Port:      161,
			Community: cfg.TrapHandler.Community,
			Version:   g.Version2c,
			Timeout:   time.Duration(10) * time.Second,
		}

		conn_err := params.Connect()
		if conn_err != nil {
			log.Printf("Cann't conntct to %s, erro: %v", trap.IpAddr.String(), conn_err)
		}
		defer func() {
			close_err := params.Conn.Close()
			if close_err != nil {
				log.Printf("Can't close connection to %s, error: %v", trap.IpAddr.String(), close_err)
			}
		}()

		hostName := ""
		oids := []string{hostNameOID}
		result, get_err := params.Get(oids)
		if get_err != nil {
			log.Printf("Can't fetch hostname for %s, use ip addres instead error: %v", trap.IpAddr.String(), get_err)
			hostName = trap.IpAddr.String()

		} else {
			hostName = string(result.Variables[0].Value.([]byte))
		}
		trap.DeviceName = hostName
		handledTraps <- trap

	}
}
