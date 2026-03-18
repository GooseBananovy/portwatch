package linux

import "os"

type LinuxProvider struct {
	procBase string
	netProcBase string
}

func NewLinuxProvider() *LinuxProvider {
	procBase := os.Getenv("PROC_BASE")
	if len(procBase) == 0 {
		procBase = "/proc"
	}

	pid1ProcBase := os.Getenv("PID1_PROC_BASE")
	if len(pid1ProcBase) == 0 {
		pid1ProcBase = procBase
	}

	return &LinuxProvider{
		procBase: procBase,
		netProcBase: pid1ProcBase,
	}
}
