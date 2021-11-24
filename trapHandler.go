package main

import (
	g "github.com/gosnmp/gosnmp"
	"github.com/sirupsen/logrus"
	"time"
)

const hostNameOID = "1.3.6.1.4.1.9.2.1.3.0"

func IntStatuTrapHandler(trap IntStatuTrap, ctx *Context) {

	ctx.logger.WithFields(logrus.Fields{
		"id": trap.Id,
	}).Debug("Trying to get hostname for ", trap.IpAddr)
	params := &g.GoSNMP{
		Target:    trap.IpAddr.String(),
		Port:      161,
		Community: ctx.config.TrapHandler.Community,
		Version:   g.Version2c,
		Timeout:   time.Duration(10) * time.Second,
	}

	err := params.Connect()
	if err != nil {
		ctx.logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Cann't connect to ", trap.IpAddr.String())
	}

	ctx.logger.WithFields(logrus.Fields{
		"id": trap.Id,
	}).Debug("Conneted to device ", trap.IpAddr)

	hostName := ""
	oids := []string{hostNameOID}
	result, err := params.Get(oids)
	if err != nil || result.Variables[0].Value == nil {
		ctx.logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Can't fetch hostname for %s, use ip address instead", trap.IpAddr.String())
		hostName = trap.IpAddr.String()

	} else {
		ctx.logger.WithFields(logrus.Fields{
			"id":     trap.Id,
			"result": result,
		}).Debug("Got the answer")
		hostName = string(result.Variables[0].Value.([]byte))
		ctx.logger.WithFields(logrus.Fields{
			"id":     trap.Id,
			"result": result,
		}).Debug("Hostname is ", hostName)
	}
	trap.DeviceName = hostName

	go UpdateIntStatuMetrics(trap, ctx)

	err = params.Conn.Close()
	if err != nil {
		ctx.logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Can't close connection to ", trap.IpAddr.String())
	}
}
