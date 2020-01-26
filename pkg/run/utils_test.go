package run_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ysmood/kit"

	"github.com/stretchr/testify/assert"
)

func TestEnsureGoTool(t *testing.T) {
	_ = os.Remove(filepath.Join(kit.GoBin(), "kit-gotool-test"))
	kit.MustGoTool("github.com/ysmood/kit/pkg/run/fixtures/kit-gotool-test")
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
	assert.Contains(t, kit.GoBin(), "bin")
	assert.Contains(t, kit.GoBin(), "bin") // use cache
}

func TestLookPath(t *testing.T) {
	kit.MustGoTool("github.com/ysmood/kit/pkg/run/fixtures/kit-gotool-test")

	PATH := os.Getenv("PATH")
	kit.E(os.Setenv("PATH", strings.ReplaceAll(PATH, kit.GoBin(), "")))
	defer func() { kit.E(os.Setenv("PATH", PATH)) }()

	p := "kit-gotool-test"
	assert.NotEqual(t, p, kit.LookPath(p))
}
