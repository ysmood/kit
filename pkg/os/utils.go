package os

import (
	"os"
	"os/signal"
	"time"
)

// WaitSignal ...
func WaitSignal(sig os.Signal) {
	c := make(chan os.Signal, 1)
	if sig == nil {
		sig = os.Interrupt
	}
	signal.Notify(c, sig)
	<-c
	signal.Stop(c)
	close(c)
}

// Retry retry function after a duration for several times
func Retry(times int, wait time.Duration, fn func()) (errs []interface{}) {
	var try func(int)

	try = func(countdown int) {
		defer func() {
			if r := recover(); r != nil {
				errs = append(errs, r)
				if countdown <= 1 {
					return
				}
				time.Sleep(wait)
				try(countdown - 1)
			}
		}()

		fn()
	}

	try(times)

	if len(errs) < times {
		return nil
	}
	return errs
}
