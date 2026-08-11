package internals

import (
	"net"

	"github.com/rs/zerolog/log"
)

func SourceIP_MAC(ClienfINF *NetworkInterface) *NetworkInterface {
	INFC, err := net.Interfaces()

	if err != nil {
		log.Error().Msg(err.Error())

	}

	for _, IFC := range INFC {
		ipAddr, _ := IFC.Addrs()

		if IFC.Name == "ens33" {
			for _, clientIP := range ipAddr {
				isV4 := clientIP.(*net.IPNet).IP.To4()
				if isV4 != nil {
					ClienfINF.SourceIP = isV4
					ClienfINF.SourceMAC = IFC.HardwareAddr
					ClienfINF.LocalNetWorkInterface = IFC.Name
				}
			}

		}

	}
	return ClienfINF
}
