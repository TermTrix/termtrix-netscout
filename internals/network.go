package internals

import "net"

type NetworkInterface struct {
	SourceIP              net.IP
	SourceMAC             net.HardwareAddr
	LocalNetWorkInterface string
}

type EthernetPayload struct {
	HardwareType uint16
	ProtocolType uint16
	HardWareSize uint8
	ProtocolSize uint8

	OptCode uint16

	SenderMAC net.HardwareAddr
	SenderIP  net.IP

	TargetMAC net.HardwareAddr
	TargetIP  net.IP
}

type EthernetHeader struct {
	DestMAC   net.HardwareAddr
	SourceMAC net.HardwareAddr
	EtherType uint16
}
