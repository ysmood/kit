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

// Escape file name based on the os. Half-width illegal char will be replaced with its full-width version.
func Escape(name string) string {
	return name
}
