package exec_test

import (
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	. "github.com/ysmood/gokit/pkg/exec"
	. "github.com/ysmood/gokit/pkg/utils"
)

func TestExec(t *testing.T) {
	Exec("go", "version").MustDo()
}

func TestExecPrefix(t *testing.T) {
	err := Exec("echo", "test").Prefix("[app] ").Do()
	assert.Nil(t, err)
}

func TestExecMustString(t *testing.T) {
	assert.Equal(t, "ok\n", Exec("echo", "ok").MustString())
}

func TestExecPrefixColor(t *testing.T) {
	err := Exec("echo", "test").Args([]string{"echo", "ok"}).Prefix("[app] @green").Do()
	assert.Nil(t, err)
}

func TestExecErr(t *testing.T) {
	err := Exec("").Cmd(exec.Command("exitexit"))
	assert.EqualError(t, err.Do(), "exec: \"exitexit\": executable file not found in $PATH")
}

func TestExecRaw(t *testing.T) {
	err := Exec("echo", "ok").Raw().Do()
	assert.Nil(t, err)
}

func TestExecKillTree(t *testing.T) {
	exe := Exec("go", "run", "./fixtures/sleep")
	go func() { Noop(exe.Do()) }()

	time.Sleep(100 * time.Millisecond)

	err := KillTree(exe.GetCmd().Process.Pid)

	assert.Nil(t, err)
}
