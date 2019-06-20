package run_test

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	kit "github.com/ysmood/gokit"
)

func TestMain(m *testing.M) {
	_ = kit.Remove("tmp/**")
	os.Exit(m.Run())
}

func TestExec(t *testing.T) {
	kit.Exec("go", "version").MustDo()
}

func TestExecPrefix(t *testing.T) {
	err := kit.Exec("go", "version").Prefix("[app] ").Do()
	assert.Nil(t, err)
}

func TestExecMustString(t *testing.T) {
	assert.Regexp(t, "go version", kit.Exec("go", "version").MustString())
}

func TestExecPrefixColor(t *testing.T) {
	err := kit.Exec("go", "version").Args([]string{"go", "version"}).Prefix("[app] @green").Do()
	assert.Nil(t, err)
}

func TestExecErr(t *testing.T) {
	err := kit.Exec("").Cmd(exec.Command("exitexit"))
	assert.Regexp(t, "exec: \"exitexit\": executable file not found in", err.Do().Error())
}

func TestExecRaw(t *testing.T) {
	err := kit.Exec("go", "version").Raw().Do()
	assert.Nil(t, err)
}

func TestExecKillTree(t *testing.T) {
	exe := kit.Exec("go", "run", "./fixtures/sleep")
	go func() { kit.Noop(exe.Do()) }()

	time.Sleep(time.Second)

	err := kit.KillTree(exe.GetCmd().Process.Pid)

	assert.Nil(t, err)
}
