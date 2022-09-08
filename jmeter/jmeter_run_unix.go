package jmeter

import (
	"Frieza/utils/shell"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type JmeterUnixRun struct {
	CanRunJmeterScript bool
}

func (r *JmeterUnixRun) IsStopped() bool {
	//检查jmeter是否处于停止状态
	cmd := CheckJmeterIsRunning
	if out, err := shell.ExecShell(cmd); err == nil {
		if s, err := strconv.Atoi(out); err == nil && s >= 1 {
			// wc -l >= 1 则说明jmeter正在运行，返回false
			return false
		}
	}
	return true
}

func (r *JmeterUnixRun) CanRun() bool {
	return r.CanRunJmeterScript
}

func (r *JmeterUnixRun) SetRunStatus(canRun bool) {
	r.CanRunJmeterScript = canRun
}

func (r *JmeterUnixRun) ExecJmeterScript(jmeterCmd string, hosts []string) {
	// 当执行此命令时会限制再次执行jmeter命令
	// 当程序执行结束后再次释放标记为可以再次执行jmeter脚本
	r.SetRunStatus(false)
	defer r.SetRunStatus(true)
	//server jmeter分布式命令执行
	// jmeter -n -t XX.JMX -l xxx.jtl -R xx.xx.xx.xx,xx.xx.xx.xx
	strings.TrimSpace(jmeterCmd)

	for i, h := range hosts {
		hosts[i] = strings.TrimSpace(h)
	}
	host := strings.Join(hosts, ",")

	cmd := ""
	if host != "" {
		cmd = fmt.Sprintf("./%s -R %s", jmeterCmd, host)
	} else {
		cmd = fmt.Sprintf("./%s", jmeterCmd)
	}
	log.Println("run cmd :", cmd)
	_, err := shell.ExecShell(cmd)

	// 需要特殊处理一下，stop会杀掉jmeter进程，所以执行的时候会报错：exit status 137
	if err != nil {
		if err.Error() == "exit status 137" {
			return
		}
		log.Println("ExecJmeterScript err: ", err)
	}

}
func (r *JmeterUnixRun) OutputJmeterLog() {
	//读取jmeter.log 过滤掉没用的内容
	b, err := ioutil.ReadFile("jmeter.log")
	if err != nil {
		log.Println("OutputJmeterLog err: ", err)
	}
	log.Println(string(b))
}

func (r *JmeterUnixRun) StopJmeterThread() {
	// 快速停止jmeter-server,快速停止jmeter-slave
	log.Println("kill all jmeter starting...")
	cmd := KillJmeterCmd
	log.Println("run cmd >>> ", cmd)
	_, err := shell.ExecShell(cmd)
	if err != nil {
		if err.Error() == "exit status 123" {
			log.Println("[Warning] jmeter not run ...")
		}
		log.Fatal("kill all jmeter error , cause : ", err)
		return
	}
	log.Println("kill all jmeter success")
}

func (r *JmeterUnixRun) StartSlave(ctx context.Context, output chan string, masterHost string) {
	// 1.启动jmeter-server 并且返回管道将输出内容交给管道
	// TODO 需要优化，管道由外部创建，然后带进来，就不用再次创建了。
	var cmd string
	cmd = RunJmeterServerCmd
	strings.TrimSpace(masterHost)
	cmd = fmt.Sprintf(cmd, masterHost)
	log.Println("run cmd >>> ", cmd)
	go shell.CommandAndOutputChan(ctx, output, cmd)
}

func (r *JmeterUnixRun) CanRunJmeter() bool {
	// 首次运行前检测jmeter一系列参数是否配置
	// 检测系统防火墙是否已经关闭
	if r.CanRun() {
		return true
	}
	// check java jdk
	if !IsFindJavaEnv() {
		InstallJDK()
	}
	// check env : server.rmi.ssl.disable
	if !IsOpenServerSSL() {
		SetServerSSL()
	}
	// check env : java.rmi.server.hostname
	//if !IsSetServerHostname() {
	//	SetServerHostname()
	//}
	// 检查jmeter是否已经启动
	if r.IsStopped() {
		return false
	}
	// 检查jmeter -version是否可以正常显示版本
	if out, err := shell.ExecShell(CheckJmeterVersionCmd); err == nil {
		if s, err := strconv.Atoi(out); err == nil && s == 0 {
			return true
		}
	}
	return false
}

//func GetJmeterSlaveStatus() int {
//	// 检查jmeter是否已经启动
//	if out, err := shell.ExecShell(CheckJmeterIsRunning); err == nil {
//		if s, err := strconv.Atoi(out); err == nil && s >= 1 {
//			// 判断jmeter-server.log文件最后一行内容, 包含started,说明正在运行中 Running
//
//		}
//	}
//	return 0
//}

//func GetDirPath() string {
//	if ex, err := os.Executable(); err == nil {
//		path, _ := filepath.EvalSymlinks(filepath.Dir(ex))
//		return path
//	}
//	return "."
//
//}
