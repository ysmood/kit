package run_test

import (
	"testing"
	"time"

	kit "github.com/ysmood/gokit"
)

func TestGuardDefaults(t *testing.T) {
	p := "tmp/" + kit.RandString(10)

	_ = kit.OutputFile(p+"/f", "ok", nil)

	guard := kit.Guard("echo", "ok", "{{.path}}").Dir("..").Patterns("*/*.go")
	go guard.MustDo()

	time.Sleep(300 * time.Millisecond)

	guard.Stop()
}

func TestGuardErr(t *testing.T) {
	guard := kit.Guard("exitexit").NoInitRun()
	go guard.MustDo()

	time.Sleep(100 * time.Millisecond)

	guard.Stop()
}

func TestGuard(t *testing.T) {
	p := "tmp/" + kit.RandString(10)

	_ = kit.OutputFile(p+"/f", "ok", nil)

	d := 0 * time.Millisecond
	i := 1 * time.Millisecond

	guard := kit.Guard("echo", "ok", "{{.path}}").
		ExecCtx(kit.Exec()).
		Dir("").
		Patterns(p + "/**").
		Debounce(&d).
		NoInitRun().
		ClearScreen().
		Interval(&i)

	go guard.MustDo()

	go func() {
		time.Sleep(50 * time.Millisecond)
		_ = kit.OutputFile(p+"/f", "ok", nil)
		time.Sleep(50 * time.Millisecond)
		_ = kit.OutputFile(p+"/g", "ok", nil)
		time.Sleep(50 * time.Millisecond)
		_ = kit.Mkdir(p+"/d", nil)
	}()

	time.Sleep(300 * time.Millisecond)

	guard.Stop()
}

func TestGuardDebounce(t *testing.T) {
	p := "tmp/" + kit.RandString(10)

	_ = kit.OutputFile(p+"/f", "ok", nil)

	i := 1 * time.Millisecond

	guard := kit.Guard("echo", "ok", "{{.path}}").Patterns(p + "/**").Interval(&i)
	go guard.MustDo()

	go func() {
		time.Sleep(50 * time.Millisecond)
		_ = kit.OutputFile(p+"/f", "a", nil)
		_ = kit.OutputFile(p+"/b", "b", nil)
	}()

	time.Sleep(200 * time.Millisecond)

	guard.Stop()
}

func TestGuardWatchErr(t *testing.T) {
	p := "tmp/" + kit.RandString(10)

	_ = kit.OutputFile(p+"/f", "ok", nil)

	i := 1 * time.Millisecond

	guard := kit.Guard("echo", "ok").Patterns(p + "/**").Interval(&i)
	go guard.MustDo()

	go func() {
		time.Sleep(50 * time.Millisecond)
		_ = kit.Remove(p+"/f", "a")
	}()

	time.Sleep(200 * time.Millisecond)

	guard.Stop()
}

func TestGuardRunErr(t *testing.T) {
	guard := kit.Guard("exitexit").Patterns("a")
	go guard.MustDo()

	time.Sleep(200 * time.Millisecond)

	guard.Stop()
}
