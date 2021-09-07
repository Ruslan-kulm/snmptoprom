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
			log.Fatalf("Connect() err: %v", conn_err)
		}
		defer func() {
			close_err := params.Conn.Close()
			log.Fatalf("Get() err: %v", close_err)
		}()

		oids := []string{hostNameOID}
		result, get_err := params.Get(oids) // Get() accepts up to g.MAX_OIDS
		if get_err != nil {
			log.Fatalf("Get() err: %v", get_err)
		}
		hostName := string(result.Variables[0].Value.([]byte))
		trap.DeviceName = hostName
		handledTraps <- trap

	}
}
