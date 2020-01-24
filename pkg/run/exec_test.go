package run_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/kit"
)

func TestMain(m *testing.M) {
	_ = kit.Remove("tmp/**")
	os.Exit(m.Run())
}

func TestExec(t *testing.T) {
	kit.Exec("go", "version").MustDo()
}

func TestExecCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := kit.Exec("go", "version").Context(ctx).Do()

	assert.Error(t, err)
}

func TestExecPrefix(t *testing.T) {
	err := kit.Exec("go", "version").Prefix("[app] ").Do()
	assert.Nil(t, err)
}

func TestExecMustString(t *testing.T) {
	assert.Regexp(t, "go version", kit.Exec("go", "version").MustString())
}

func TestExecMustStringErr(t *testing.T) {
	assert.Panics(t, func() {
		kit.Exec(kit.RandString(16)).MustString()
	})
}

func TestExecEnv(t *testing.T) {
	s := "tmp/" + kit.RandString(10)
	kit.E(kit.Mkdir(s, nil))
	assert.Regexp(t, s, kit.Exec("go", "env").Env("GOTMPDIR="+s).Env("A=B").MustString())
}

func TestExecPrefixColor(t *testing.T) {
	err := kit.Exec("go", "version").Args([]string{"go", "version"}).Prefix("[app] @green").Do()
	assert.Nil(t, err)
}

func TestExecErr(t *testing.T) {
	err := kit.Exec("exitexit")
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

func TestOverrideGoBin(t *testing.T) {
	err := kit.Exec("go", "version").NewEnv("GOBIN=test").Do()
	assert.Nil(t, err)
}
