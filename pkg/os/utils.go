package os

import (
	"os"
	"os/signal"
	"time"
)

// WaitSignal block until get specified os signals
func WaitSignal(signals ...os.Signal) {
	c := make(chan os.Signal, 1)
	if len(signals) == 0 {
		signals = append(signals, os.Interrupt)
	}
	signal.Notify(c, signals...)
	<-c
	signal.Stop(c)
	close(c)
}

// Retry retry function after a duration for several times, if success returns nil
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
