package main

import (
	"github.com/google/uuid"
	g "github.com/gosnmp/gosnmp"
	"github.com/sirupsen/logrus"
	"net"
	"os"
)

const interfaceStatus = ".1.3.6.1.6.3.1.1.5"

func TrapListen(th g.TrapHandlerFunc, ctx *Context) {
	ctx.logger.Debug("TrapListener is preparing")
	tl := g.NewTrapListener()
	tl.OnNewTrap = th
	tl.Params = g.Default
	ctx.logger.Debug("TrapListener is prepared. Listening", "0.0.0.0:"+ctx.config.TrapListener.Port)
	err := tl.Listen("0.0.0.0:" + ctx.config.TrapListener.Port)
	if err != nil {
		ctx.logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("error on listen ", "0.0.0.0:"+ctx.config.TrapListener.Port)
		os.Exit(1)
	}
}

func onNewTrap(ctx *Context) func(packet *g.SnmpPacket, addr *net.UDPAddr) {
	return func(packet *g.SnmpPacket, addr *net.UDPAddr) {
		id := uuid.New()
		ctx.logger.WithFields(logrus.Fields{
			"id":       id,
			"trapBody": packet,
		}).Debug("A new trap received")

		switch packet.Enterprise {
		case interfaceStatus:
			interfaceName := string(packet.Variables[1].Value.([]byte))
			interfaceStatus := string(packet.Variables[3].Value.([]byte))
			t := IntStatuTrap{
				Id:              id,
				DeviceName:      "",
				IpAddr:          addr.IP,
				InterfaceName:   interfaceName,
				InterfaceStatus: interfaceStatus,
			}
			ctx.logger.WithFields(logrus.Fields{
				"id":        t.Id,
				"ip":        t.IpAddr,
				"interface": t.InterfaceName,
				"status":    t.InterfaceStatus,
			}).Debug("Interface Status trap initialized")

			go IntStatuTrapHandler(t, ctx)

		default:
			ctx.logger.Debug("Unsupported trap ", packet.Enterprise)
		}
	}
}
