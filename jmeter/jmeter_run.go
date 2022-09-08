package jmeter

import "runtime"

var JR JmeterRunner

func New() JmeterRunner {
	if JR != nil {
		return JR
	}
	if runtime.GOOS == "windows" {
		JR = &JmeterWinsRun{CanRunJmeterScript: true}
	} else {
		JR = &JmeterUnixRun{CanRunJmeterScript: true}
	}
	return JR
}
