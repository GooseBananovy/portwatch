package linux

import "os"

type LinuxProvider struct {
	procBase string
}

func NewLinuxProvider() *LinuxProvider {
	procBase := os.Getenv("PROC_BASE")
	if len(procBase) == 0 {
		procBase = "/proc"
	}
	return &LinuxProvider{
		procBase: procBase,
	}
}
