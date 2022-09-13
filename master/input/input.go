package input

import (
	"Frieza/constant"
	"Frieza/jmeter"
	"Frieza/master"
	"Frieza/master/server"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"log"
	"strings"
)

func CmdInput(jmeterRunner jmeter.JmeterRunner) {
	// master命令交互入口
	slaves, err := server.NewSlaves()
	if err == master.ErrNotFound {
		log.Println(err)
	} else if err != nil {
		panic(err)
	}

	//jmeterRunner := jmeter.New()
	// 重复执行,检查input内容
	for {
		Cmd := ""
		prompt := &survey.Input{
			Message: "Frieze>>>",
		}

		err := survey.AskOne(prompt, &Cmd)
		if err == terminal.InterruptErr {
			// 接收到ctrl+c后master退出，并通知所有slave，发送ServerQuit指令
			slavesSlice := slaves.GetAllSlavesSlice()
			for _, slave := range slavesSlice {
				slave <- constant.ServerQuit
			}
			return
		}
		if strings.HasPrefix(Cmd, "jmeter ") {
			// jmeterRunner 命令行，目前只判断是否是“jmeterRunner ” 开头，未进行更详细的命令校验。

			// 去掉result tree
			jmx, err := jmeter.GetJmxFilename(Cmd)
			if err != nil {
				log.Println("jmeter script not found in (", Cmd, ")")
				continue
			}
			jmeter.DisableJmeterResultTrees(jmx)
			// run
			log.Printf("run jmeterRunner shell >>> %s\n", Cmd)
			if jmeterRunner.CanRun() {
				go jmeterRunner.ExecJmeterScript(Cmd, slaves.GetAllCanRunIP())
			} else {
				log.Println("Unable to execute, wait for the end of the operation")
			}
		}
		if strings.HasPrefix(Cmd, "log") {
			// 打印jmeter-server.log到控制台，此处并未进行log筛选操作
			// TODO 读取jmeter-server.log 并且逐行分析进行数据处理操作
			jmeterRunner.OutputJmeterLog()
		}
		if strings.HasPrefix(Cmd, "stop") {
			// stop：停止所有status是running状态的slave，使其重启进入Idle状态。
			// 目的：结束slave的脚本运行
			if !jmeterRunner.CanRun() {
				// stop master jmeterRunner
				jmeterRunner.StopJmeterThread()
				// call slaves restart
				slavesSlice := slaves.GetAllSlavesSlice()

				for _, slave := range slavesSlice {
					slave <- constant.StopJmeterRunCmd
				}
				// 命令发送成功后，master直接将slaves标记为Idle状态，这里需要确保slaves一定重启成功。
				// 目前无需进行确认，因为重启操作失败会导致slave自动退出程序，退出程序时会自动通知master删除对应的slave
				// 当程序命令未通知到slave会如何？
				// slaves断开会通知slave，且每隔几秒自动触发状态更新。
				jmeterRunner.SetRunStatus(true)
			} else {
				log.Println("jmeterRunner is not run,status is not Running")
			}
		}
		if Cmd == "slave" {
			fmt.Println(slaves)
		}
	}

}
