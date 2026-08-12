package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/termtrix/internals"
	"golang.org/x/sys/unix"
)

type WorkerP struct {
	SourceIP              net.IP
	SourceMAC             net.HardwareAddr
	LocalNetWorkInterface string
	HardwareType          uint16
	ProtocolType          uint16
	HardWareSize          uint8
	ProtocolSize          uint8

	OptCode uint16

	SenderMAC net.HardwareAddr
	SenderIP  net.IP

	TargetMAC net.HardwareAddr

	FD int
}

func main() {

	var BroadcastMAC = net.HardwareAddr{
		0xff,
		0xff,
		0xff,
		0xff,
		0xff,
		0xff,
	}
	clientINF := internals.NetworkInterface{}
	INFC := internals.SourceIP_MAC(&clientINF)

	_, network, _ := net.ParseCIDR("192.168.1.22/24") // change this to your network range

	targets := hosts(network)

	ETH := &internals.EthernetHeader{
		DestMAC:   BroadcastMAC,
		SourceMAC: INFC.SourceMAC,
		EtherType: 0x0806,
	}

	var TagetMac = net.HardwareAddr{
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
	}

	fd, err := unix.Socket(
		unix.AF_PACKET,
		unix.SOCK_RAW,
		int(htons(unix.ETH_P_ARP)),
	)

	if err != nil {
		log.Fatal("Failed to connect socket", err)
		return
	}

	WorkerPayload := &WorkerP{
		SourceIP:     INFC.SourceIP,
		SourceMAC:    INFC.SourceMAC,
		HardwareType: 0x0001,
		ProtocolType: 0x0800,
		HardWareSize: 6,
		ProtocolSize: 4,
		OptCode:      1,
		SenderMAC:    INFC.SourceMAC,
		SenderIP:     INFC.SourceIP,
		TargetMAC:    TagetMac,
		FD:           fd,
	}

	// 1. Start receiver FIRST
	go func() {
		buf := make([]byte, 2048)

		for {
			n, _, err := unix.Recvfrom(fd, buf, 0)
			if err != nil {
				fmt.Println("Recv error:", err)
				return
			}

			if n < 42 {
				continue
			}

			opcode := binary.BigEndian.Uint16(buf[20:22])

			if opcode != 2 {
				continue
			}

			senderIP := net.IPv4(
				buf[28],
				buf[29],
				buf[30],
				buf[31],
			)

			// senderMAC := FormatMAC(buf[22:28])

			fmt.Println(senderIP, "-->", net.HardwareAddr(buf[22:28]))
		}
	}()

	jobs := make(chan net.IP)

	var wg sync.WaitGroup

	workers := 10

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go worker(i, ETH, WorkerPayload, jobs, &wg)
	}

	for _, target := range targets {
		jobs <- target
	}

	close(jobs)

	wg.Wait()

	time.Sleep(100 * time.Second)

}

func htons(v uint16) uint16 {
	return (v<<8)&0xff00 | v>>8
}

func FormatMAC(b []byte) string {
	if len(b) < 6 {
		return ""
	}
	return fmt.Sprintf("%02x %02x %02x %02x %02x %02x",
		b[0], b[1], b[2], b[3], b[4], b[5])
}

func hosts(network *net.IPNet) []net.IP {
	var result []net.IP

	ip := network.IP.To4()
	// mask := network.Mask

	for i := 0; i < 256; i++ {
		target := make(net.IP, 4)
		copy(target, ip)

		target[3] = byte(i)

		if network.Contains(target) {
			result = append(result, target)
		}
	}

	return result
}

func worker(id int, ETH_ *internals.EthernetHeader, ETH *WorkerP, jobs chan net.IP, wg *sync.WaitGroup) {
	defer wg.Done()

	INFACE, _ := net.InterfaceByName("ens33")

	addr := &unix.SockaddrLinklayer{
		Ifindex: INFACE.Index,
	}
	for target := range jobs {

		arpPacket := &internals.EthernetPayload{
			HardwareType: 0x0001,
			ProtocolType: 0x0800,
			HardWareSize: 6,
			ProtocolSize: 4,
			OptCode:      1,
			SenderMAC:    ETH.SourceMAC,
			SenderIP:     ETH.SenderIP,
			TargetMAC:    ETH.TargetMAC,
			TargetIP:     target,
		}

		frame := internals.MakeFrame(ETH_, arpPacket)

		err := unix.Sendto(
			ETH.FD,
			frame,
			0,
			addr,
		)

		if err != nil {
			fmt.Println("Sendto failed:", err)
			continue
		}
	}
}
