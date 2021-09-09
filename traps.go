package main

import (
	"github.com/google/uuid"
	"net"
)

type IntStatuTrap struct {
	Id              uuid.UUID
	DeviceName      string
	IpAddr          net.IP
	InterfaceName   string
	InterfaceStatus string
}
