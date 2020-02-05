package utils

import (
	"runtime"
	"testing"
)

func TestWindowsLog(t *testing.T) {
	goos = "windows"
	defer func() { goos = runtime.GOOS }()

	_ = ClearScreen()

	stdout()
	stderr()
}
