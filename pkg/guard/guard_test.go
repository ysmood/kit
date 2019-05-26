package guard_test

import (
	"testing"
	"time"

	. "github.com/ysmood/gokit/pkg/exec"
	. "github.com/ysmood/gokit/pkg/guard"
	. "github.com/ysmood/gokit/pkg/os"
	. "github.com/ysmood/gokit/pkg/utils"
)

func TestGuardDefaults(t *testing.T) {
	p := "tmp/" + GenerateRandomString(10)

	_ = OutputFile(p+"/f", "ok", nil)

	guard := Guard("echo", "ok", "{{.path}}").Dir("..").Patterns("*/*.go")
	go guard.MustDo()

	time.Sleep(300 * time.Millisecond)

	guard.Stop()
}

func TestGuardErr(t *testing.T) {
	guard := Guard("exitexit").NoInitRun()
	go guard.MustDo()

	time.Sleep(100 * time.Millisecond)

	guard.Stop()
}

func TestGuard(t *testing.T) {
	p := "tmp/" + GenerateRandomString(10)

	_ = OutputFile(p+"/f", "ok", nil)

	d := 0 * time.Millisecond
	i := 1 * time.Millisecond

	guard := Guard("echo", "ok", "{{.path}}").
		ExecCtx(Exec()).
		Dir("").
		Patterns(p + "/**").
		Debounce(&d).
		NoInitRun().
		ClearScreen().
		Interval(&i)

	go guard.Do() // nolint:errcheck

	go func() {
		time.Sleep(50 * time.Millisecond)
		_ = OutputFile(p+"/f", "ok", nil)
		time.Sleep(50 * time.Millisecond)
		_ = OutputFile(p+"/g", "ok", nil)
		time.Sleep(50 * time.Millisecond)
		_ = Mkdir(p+"/d", nil)
	}()

	time.Sleep(300 * time.Millisecond)

	guard.Stop()
}

func TestGuardDebounce(t *testing.T) {
	p := "tmp/" + GenerateRandomString(10)

	_ = OutputFile(p+"/f", "ok", nil)

	i := 1 * time.Millisecond

	guard := Guard("echo", "ok", "{{.path}}").Patterns(p + "/**").Interval(&i)
	go guard.Do() // nolint:errcheck

	go func() {
		time.Sleep(50 * time.Millisecond)
		_ = OutputFile(p+"/f", "a", nil)
		_ = OutputFile(p+"/b", "b", nil)
	}()

	time.Sleep(200 * time.Millisecond)

	guard.Stop()
}
