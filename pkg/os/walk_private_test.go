package os

import (
	"os"
	"testing"
)

func TestGetGitSubmodulesEmpty(t *testing.T) {
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", oldPath)
	l := getGitSubmodules("/")
	if l != nil {
		panic(l)
	}
}
