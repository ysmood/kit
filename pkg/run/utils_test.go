package run_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/ysmood/kit"

	"github.com/stretchr/testify/assert"
)

func TestEnsureGoTool(t *testing.T) {
	_ = os.Remove(path.Join(kit.GoBin(), "golint"))
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

func TestGoPath(t *testing.T) {
	old := os.Getenv("GOPATH")
	kit.E(os.Setenv("GOPATH", ""))
	defer func() { kit.E(os.Setenv("GOPATH", old)) }()

	s := kit.GoPath()

	assert.True(t, kit.Exists(s))
}

func TestGoBin(t *testing.T) {
	assert.Contains(t, kit.GoBin(), "/bin")
	assert.Contains(t, kit.GoBin(), "/bin")
}

func TestLookPath(t *testing.T) {
	PATH := os.Getenv("PATH")
	os.Setenv("PATH", strings.ReplaceAll(PATH, kit.GoBin(), ""))
	defer os.Setenv("PATH", PATH)

	p := "golint"
	assert.NotEqual(t, p, kit.LookPath(p))
}
