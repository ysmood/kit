// +build !windows

package os

import (
	"os"
	"syscall"
)

// SendSigInt ...
func SendSigInt(pid int) error {
	p, _ := os.FindProcess(pid)
	return p.Signal(syscall.SIGINT)
}

// ExecutableExt ...
func ExecutableExt() string {
	return ""
}
