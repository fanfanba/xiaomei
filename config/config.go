package config

import (
	"net"
	"strings"
	"time"
)

var Data = parseConfigData()
var TimeZone = time.FixedZone(Data.TimeZoneName, Data.TimeZoneOffset)

func init() {
	setupMailer()
}

func CurrentAppServer() ServerConfig {
	ifcAddrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, server := range Data.DeployServers {
		if server.AppAddr != `` {
			for _, ifcAddr := range ifcAddrs {
				if strings.HasPrefix(ifcAddr.String(), server.Addr+`/`) {
					return server
				}
			}
		}
	}
	return ServerConfig{}
}
