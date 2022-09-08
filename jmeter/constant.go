package jmeter

import (
	"errors"
)

// linux命令
const (
	AptGet                = "apt-get -h &>/etc/null ; echo $?"
	Yum                   = "yum -h &>/etc/null ; echo $?"
	RunJmeterServerCmd    = "./jmeter -Djava.rmi.server.hostname=%s -s -j jmeter-server.log"
	CheckJmeterVersionCmd = "./jmeter -v &>/dev/null;echo $?"
	InstallJDKCmdCentOS   = "yum -y install java-1.8.0-openjdk.x86_64"              // centos
	InstallJDKCmdUbuntu   = "sudo apt-get update; sudo apt-get install default-jdk" //ubuntu  | debian
	KillJmeterCmd         = "ps -ef | grep ApacheJMeter |grep -v grep |awk '{print $2}' | xargs kill -9"
	CheckJmeterIsRunning  = "ps -ef | grep ApacheJMeter |grep -v grep | wc -l"
	CheckJDKVersion       = "java -version &>/dev/null;echo $?"
)

// windows 命令
const (
	RunJmeterServerCmdWins    = ".\\jmeter.bat -Djava.rmi.server.hostname=%s -s -j jmeter-server.log"
	CheckJmeterVersionCmdWins = ".\\jmeter.bat -v > nul && echo %errorlevel%"
	InstallJDKCmdWins         = ""
	KillJmeterCmdWins         = ""
	CheckJmeterIsRunningWins  = ""
	CheckJDKVersionWins       = ""
)

const (
	RmiServerHostnameKey = "java.rmi.server.hostname"
	RmiServerDisableKey  = "server.rmi.ssl.disable"
)

var PropertiesName = struct {
	User   string
	System string
	Jmeter string
}{User: "user", System: "system", Jmeter: "jmeter"}

var (
	ErrRmiDisable = errors.New("set server.rmi.ssl.disable is failed")
)
