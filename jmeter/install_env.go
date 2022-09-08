package jmeter

import (
	"Frieza/utils/ip"
	"Frieza/utils/properties"
	"Frieza/utils/shell"
	"fmt"
	"log"
	"strings"
)

//system.properties 最后一行添加 防止多网卡情况发生
//java.rmi.server.hostname=<IP addr>

func SetServerHostname() error {
	host := ip.GetOutBoundIp()
	p := properties.NewProperties()
	filename := fmt.Sprintf("%s.properties", PropertiesName.User)
	err := p.LoadFromFile(filename)
	if err != nil {
		log.Println("p.LoadFromFile err :", err)
		return err
	}
	if value, ok := p.Property(RmiServerHostnameKey); !ok || value != host {
		p.SetProperty(RmiServerHostnameKey, host)
	}
	err = p.StoreToFile(filename)

	if err != nil {
		log.Println("set java.rmi.server.hostname is failed")
		return err

	}
	log.Printf("set  java.rmi.server.hostname = %s", host)
	return nil

}
func IsSetServerHostname() bool {
	host := ip.GetOutBoundIp()
	p := properties.NewProperties()
	err := p.LoadFromFile(fmt.Sprintf("%s.properties", PropertiesName.User))
	if err != nil {
		log.Println("p.LoadFromFile err :", err)
		return false
	}
	if value, ok := p.Property(RmiServerHostnameKey); !ok || strings.TrimSpace(value) != host {
		return false
	}
	return true
}

func SetJmeterJVM() {
	// jmeter
	//: "${HEAP:="-Xms1g -Xmx1g -XX:MaxMetaspaceSize=256m"}"
}

// Server failed to start: java.rmi.RemoteException: Cannot start. localhost.localdomain is a loopback address.
// An error occurred: Cannot start. localhost.localdomain is a loopback address

// 解决办法:
// 设置jmeter环境变量（server.rmi.ssl.disable=true）
// 读取user.properties中是否有server.rmi.ssl.disable=true
// 没有则写入server.rmi.ssl.disable=true

func SetServerSSL() error {
	p := properties.NewProperties()
	filename := fmt.Sprintf("%s.properties", PropertiesName.User)
	err := p.LoadFromFile(filename)
	if err != nil {
		log.Println("p.LoadFromFile err :", err)
		log.Println("set server.rmi.ssl.disable is failed: ", err)
		return err
	}
	if value, ok := p.Property(RmiServerDisableKey); !ok || value == "false" {
		p.SetProperty(RmiServerDisableKey, "true")
	}
	err = p.StoreToFile(filename)
	if err != nil {
		return err
	}
	log.Println("set server.rmi.ssl.disable is success")
	return nil

}

func IsOpenServerSSL() bool {
	p := properties.NewProperties()
	err := p.LoadFromFile(fmt.Sprintf("%s.properties", PropertiesName.User))
	if err != nil {
		log.Println("p.LoadFromFile err :", err)
		log.Println("set server.rmi.ssl.disable is failed: ", err)
		return false
	}
	if value, ok := p.Property(RmiServerDisableKey); !ok || value != "true" {
		return false
	}
	return true
}

func IsFindJavaEnv() bool {
	cmd := CheckJDKVersion
	out, err := shell.ExecShell(cmd)
	if err == nil && strings.TrimSpace(out) == "0" { // " 0" || "0"
		return true
	} else {
		log.Println("not found java env:", out, err)
		return false
	}
}

func InstallJDK() error {
	// TODO 待验证
	var cmd string
	out, err := shell.ExecShell(Yum)
	if err == nil && strings.TrimSpace(out) == "0" {
		cmd = InstallJDKCmdCentOS
	}
	out, err = shell.ExecShell(AptGet)
	if err == nil && strings.TrimSpace(out) == "0" {
		cmd = InstallJDKCmdUbuntu
	}

	log.Printf("try install jdk >>> %s", cmd)
	_, err = shell.ExecShell(cmd)
	if err != nil {
		log.Println("install jdk err :", err)
		return err
	}

	return nil
}
