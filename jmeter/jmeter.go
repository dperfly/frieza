package jmeter

import (
	"context"
)

var CanRunJmeterScript = true

type JmeterRunner interface {
	IsStopped() bool
	CanRun() bool
	SetRunStatus(canRun bool)
	ExecJmeterScript(jmeterCmd string, hosts []string)
	OutputJmeterLog()
	StopJmeterThread()
	StartSlave(ctx context.Context, output chan string, masterHost string)
	CanRunJmeter() bool
}
