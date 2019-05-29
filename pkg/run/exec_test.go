package run_test

import (
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	kit "github.com/ysmood/gokit"
)

func TestExec(t *testing.T) {
	kit.Exec("go", "version").MustDo()
}

func TestExecPrefix(t *testing.T) {
	err := kit.Exec("echo", "test").Prefix("[app] ").Do()
	assert.Nil(t, err)
}

func TestExecMustString(t *testing.T) {
	assert.Equal(t, "ok\n", kit.Exec("echo", "ok").MustString())
}

func TestExecPrefixColor(t *testing.T) {
	err := kit.Exec("echo", "test").Args([]string{"echo", "ok"}).Prefix("[app] @green").Do()
	assert.Nil(t, err)
}

func TestExecErr(t *testing.T) {
	err := kit.Exec("").Cmd(exec.Command("exitexit"))
	assert.EqualError(t, err.Do(), "exec: \"exitexit\": executable file not found in $PATH")
}

func TestExecRaw(t *testing.T) {
	err := kit.Exec("echo", "ok").Raw().Do()
	assert.Nil(t, err)
}

func TestExecKillTree(t *testing.T) {
	exe := kit.Exec("go", "run", "./fixtures/sleep")
	go func() { kit.Noop(exe.Do()) }()

	time.Sleep(30 * time.Millisecond)

	err := kit.KillTree(exe.GetCmd().Process.Pid)

	assert.Nil(t, err)
}
