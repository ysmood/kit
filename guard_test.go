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
	d := 10 * time.Millisecond

	go g.Guard([]string{"echo", "ok", "{{path}}"}, []string{p}, g.GuardOptions{
		Stop:     stop,
		Debounce: &d,
	})

	go func() {
		for {
			time.Sleep(300 * time.Millisecond)
			g.OutputFile(p+"/f", "ok", nil)
			g.OutputFile(p+"/b", "ok", nil)
			time.Sleep(300 * time.Millisecond)
			g.Mkdir(p+"/d", nil)
		}
	}()

	time.Sleep(1 * time.Second)

	stop <- g.Nil{}
}
