package constant

const (
	Starting = iota // 启动中 Created remote object之前
	Idle            // 空闲状态
	Running         // 运行中
	Stopped         // 停止状态
	Failed          // 启动失败
)

const (
	PingBackCmd      = iota + 1 //心跳检测，检测slave的状态，定时反馈
	StopJmeterRunCmd            //停止脚本运行命令
	ServerQuit                  // master退出
	RunSlaveCmd                 //运行命令
	RestartSlaveCmd             //重启命令
)
