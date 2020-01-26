package os

import (
	"os"
	"testing"

	"github.com/ysmood/kit/pkg/utils"
)

func TestGetGitSubmodulesEmpty(t *testing.T) {
	oldPath := os.Getenv("PATH")
	utils.E(os.Setenv("PATH", ""))
	defer func() { utils.E(os.Setenv("PATH", oldPath)) }()
	l := getGitSubmodules("/")
	if l != nil {
		panic(l)
	}
}
