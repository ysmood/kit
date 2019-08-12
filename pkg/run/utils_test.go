package run_test

import (
	"os"
	"path"
	"testing"

	"github.com/ysmood/kit"

	"github.com/stretchr/testify/assert"
)

func TestEnsureGoTool(t *testing.T) {
	_ = os.Remove(path.Join(kit.GoPath(), "bin", "golint"))
	kit.MustGoTool("golang.org/x/lint/golint")
}

func TestRunTask(t *testing.T) {
	outA := false
	outB := false

	old := os.Args
	os.Args = []string{"test", "act", "--ok"}
	defer func() { os.Args = old }()

	kit.Tasks().App(nil).Add(
		kit.Task("act", "").Init(func(cmd kit.TaskCmd) func() {
			flag := cmd.Flag("ok", "info").Bool()

			return func() {
				outA = *flag
			}
		}),
		kit.Task("fn", "").Run(func() {
			outB = true
		}),
	).Do()

	assert.True(t, outA)
	assert.False(t, outB)
}
