//go:generate  goreleaser release --snapshot --rm-dist

package main

import (
	"Frieza/jmeter"
	"Frieza/master/input"
	master "Frieza/master/server"
	"Frieza/slave"
	"Frieza/version"
	"context"
	"flag"
	"fmt"
	"log"
	"runtime"
)

func main() {

	var host *string
	var isMaster *bool
	var ver *bool
	host = flag.String("R", "", "master host")
	isMaster = flag.Bool("M", false, "default slave server")
	ver = flag.Bool("v", false, "version")
	flag.Parse()
	fmt.Println(version.Logo)

	if runtime.GOOS == "windows" {
		log.Println("Windows operating system is not currently supported")
		return
	}

	if *ver == true {
		fmt.Println("Frieza-version: ", version.Version)
		return
	}
	log.Println("isMaster", *isMaster)
	jmeterRunner := jmeter.New()
	if *isMaster == false {
		if *host == "" {
			log.Println("[WARNING] please input master host; default localhost")
		} else {
			log.Println("master is : ", fmt.Sprintf("%s:8081", *host))
		}
		slave.StartClient(context.Background(), fmt.Sprintf("%s:8081", *host))
	} else {

		// set jmeter env
		if !jmeter.IsOpenServerSSL() {
			err := jmeter.SetServerSSL()
			if err != nil || !jmeter.IsOpenServerSSL() {
				log.Println(err)
				return
			}
		}
		// set java Server use host
		//if !jmeter.IsSetServerHostname() {
		//	err := jmeter.SetServerHostname()
		//	if err != nil || !jmeter.IsSetServerHostname() {
		//		log.Println(err)
		//		return
		//	}
		//}

		// set java
		if !jmeter.IsFindJavaEnv() {
			err := jmeter.InstallJDK()
			if err != nil || !jmeter.IsFindJavaEnv() {
				log.Println("install jdk failed")
				return
			}
		}

		go master.StartServer(":8081")
		defer jmeterRunner.StopJmeterThread()
		input.CmdInput(jmeterRunner)
	}

}
