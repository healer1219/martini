package utils

import (
	"errors"
	"log"
	"net"
)

func GetLocalIp() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Printf("err get net interfaces, %e", err)
		return "", err
	}
	for _, ifc := range interfaces {
		if ifc.Flags&net.FlagUp == 0 {
			continue
		}
		if ifc.Flags&net.FlagLoopback == 1 {
			continue
		}

		addrs, err := ifc.Addrs()
		if err != nil {
			log.Println(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}

	}
	return "", errors.New("get local ip failed")

}
