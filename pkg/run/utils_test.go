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

	kit.TaskRun(nil, kit.Tasks{
		"act": kit.Task{
			Help: "description",
			Init: func(cmd kit.TaskCmd) func() {
				flag := cmd.Flag("ok", "info").Bool()

				return func() {
					outA = *flag
				}
			},
		},
		"fn": kit.Task{Task: func() {
			outB = true
		}},
	})

	assert.True(t, outA)
	assert.False(t, outB)
}
