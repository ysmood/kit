package gokit_test

import (
	"fmt"
	"testing"
	"time"

	g "github.com/ysmood/gokit"
)

func TestGuard(t *testing.T) {
	str, _ := g.GenerateRandomString(10)
	p := fmt.Sprintf("fixtures/%s", str)

	g.OutputFile(p+"/f", "ok", nil)

	stop := make(chan g.Nil)
	d := 0 * time.Millisecond

	go g.Guard("echo", "ok", "{{path}}").Patterns(p + "/**").Context(g.GuardContext{
		Stop:     stop,
		Debounce: &d,
	}).Do()

	go func() {
		time.Sleep(100 * time.Millisecond)
		g.OutputFile(p+"/f", "ok", nil)
		time.Sleep(100 * time.Millisecond)
		g.Mkdir(p+"/d", nil)
	}()

	time.Sleep(1 * time.Second)

	stop <- g.Nil{}
}
