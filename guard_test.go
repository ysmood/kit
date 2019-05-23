package gokit_test

import (
	"fmt"
	"testing"
	"time"

	g "github.com/ysmood/gokit"
)

func TestGuardDefaults(t *testing.T) {
	str, _ := g.GenerateRandomString(10)
	p := fmt.Sprintf("fixtures/%s", str)

	g.OutputFile(p+"/f", "ok", nil)

	guard := g.Guard("echo", "ok", "{{path}}")
	go guard.Do()

	time.Sleep(300 * time.Millisecond)

	guard.Stop()
}

func TestGuard(t *testing.T) {
	str, _ := g.GenerateRandomString(10)
	p := fmt.Sprintf("fixtures/%s", str)

	g.OutputFile(p+"/f", "ok", nil)

	d := 0 * time.Millisecond
	i := 1 * time.Millisecond

	guard := g.Guard("echo", "ok", "{{path}}").
		ExecCtx(g.Exec()).
		Dir("").
		Patterns(p + "/**").
		Debounce(&d).
		NoInitRun().
		ClearScreen().
		Interval(&i)
	go guard.Do()

	go func() {
		time.Sleep(50 * time.Millisecond)
		g.OutputFile(p+"/f", "ok", nil)
		time.Sleep(50 * time.Millisecond)
		g.OutputFile(p+"/g", "ok", nil)
		time.Sleep(50 * time.Millisecond)
		g.Mkdir(p+"/d", nil)
	}()

	time.Sleep(300 * time.Millisecond)

	guard.Stop()
}

func TestGuardDebounce(t *testing.T) {
	str, _ := g.GenerateRandomString(10)
	p := fmt.Sprintf("fixtures/%s", str)

	g.OutputFile(p+"/f", "ok", nil)

	i := 1 * time.Millisecond

	guard := g.Guard("echo", "ok", "{{path}}").Patterns(p + "/**").Interval(&i)
	go guard.Do()

	go func() {
		time.Sleep(50 * time.Millisecond)
		g.OutputFile(p+"/f", "a", nil)
		g.OutputFile(p+"/b", "b", nil)
	}()

	time.Sleep(200 * time.Millisecond)

	guard.Stop()
}
