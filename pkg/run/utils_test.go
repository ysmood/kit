package run_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	kit "github.com/ysmood/gokit"
)

func TestEnsureGoTool(t *testing.T) {
	kit.MustGoTool("github.com/ysmood/gokit")
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
