package ip

import (
	"fmt"
	"net"
	"strings"
)

func GetLocalhostIP() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
			}
		}
	}
}

func GetOutBoundIp() string {
	// 该方法可直接获取到对外的IP
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		panic("获取本机对外IP异常")
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(localAddr.String(), ":")[0]
	return ip
}
