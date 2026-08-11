package internals

import (
	"encoding/binary"
)

func MakeFrame(header *EthernetHeader, arp *EthernetPayload) []byte {

	frame := make([]byte, 42)

	// Ethernet Header
	copy(frame[0:6], header.DestMAC)
	copy(frame[6:12], header.SourceMAC)
	binary.BigEndian.PutUint16(frame[12:14], header.EtherType)

	// ARP Header
	binary.BigEndian.PutUint16(frame[14:16], arp.HardwareType)
	binary.BigEndian.PutUint16(frame[16:18], arp.ProtocolType)

	frame[18] = arp.HardWareSize
	frame[19] = arp.ProtocolSize

	binary.BigEndian.PutUint16(frame[20:22], arp.OptCode)

	// ARP Sender
	copy(frame[22:28], arp.SenderMAC)
	copy(frame[28:32], arp.SenderIP.To4())

	// ARP Target
	copy(frame[32:38], arp.TargetMAC)
	copy(frame[38:42], arp.TargetIP.To4())

	return frame
}
