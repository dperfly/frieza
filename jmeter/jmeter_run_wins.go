package jmeter

import "context"

type JmeterWinsRun struct {
	CanRunJmeterScript bool
}

func (j *JmeterWinsRun) IsStopped() bool {
	//TODO implement me
	panic("implement me")
}

func (j *JmeterWinsRun) CanRun() bool {
	//TODO implement me
	panic("implement me")
}

func (j *JmeterWinsRun) SetRunStatus(canRun bool) {
	//TODO implement me
	panic("implement me")
}

func (j *JmeterWinsRun) ExecJmeterScript(jmeterCmd string, hosts []string) {
	//TODO implement me
	panic("implement me")
}

func (j *JmeterWinsRun) OutputJmeterLog() {
	//TODO implement me
	panic("implement me")
}

func (j *JmeterWinsRun) StopJmeterThread() {
	//TODO implement me
	panic("implement me")
}

func (j *JmeterWinsRun) StartSlave(ctx context.Context, output chan string, masterHost string) {
	//TODO implement me
	panic("implement me")
}

func (j *JmeterWinsRun) CanRunJmeter() bool {
	//TODO implement me
	panic("implement me")
}
