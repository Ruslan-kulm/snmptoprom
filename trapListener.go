package main

import (
	g "github.com/gosnmp/gosnmp"
	"log"
	"net"
)

func CreateTrapListener(th g.TrapHandlerFunc, cfg Config) {
	tl := g.NewTrapListener()
	tl.OnNewTrap = th
	tl.Params = g.Default
	err := tl.Listen("0.0.0.0:" + cfg.TrapListener.Port)
	if err != nil {
		log.Panicf("error in listen: %s", err)
	}
}

func onNewTrap(rawTraps chan<- IntStatuTrap) func(packet *g.SnmpPacket, addr *net.UDPAddr) {
	return func(packet *g.SnmpPacket, addr *net.UDPAddr) {
		switch packet.Enterprise {
		case ".1.3.6.1.6.3.1.1.5":
			interfaceName := string(packet.Variables[1].Value.([]byte))
			interfaceStatus := string(packet.Variables[3].Value.([]byte))

			t := IntStatuTrap{
				DeviceName:      "",
				IpAddr:          addr.IP,
				InterfaceName:   interfaceName,
				InterfaceStatus: interfaceStatus,
			}
			rawTraps <- t
		default:
			log.Printf("got trapdata from %s\n", packet.Enterprise)
		}
	}
}
