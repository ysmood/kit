package utils

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

var chPause = make(chan Nil)

// Pause the goroutine forever
func Pause() {
	<-chPause
}

// Sleeper sleeps for sometime, returns the reason to wake, if ctx is done release resource
type Sleeper func(ctx context.Context) error

// ErrMaxSleepCount ...
var ErrMaxSleepCount = errors.New("max sleep count")

// CountSleeper wake when counts to max and return
func CountSleeper(max int) Sleeper {
	count := 0
	return func(ctx context.Context) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if count == max {
			return ErrMaxSleepCount
		}
		count++
		return nil
	}
}

// DefaultBackoff algorithm: A(n) = A(n-1) * random[1.9, 2.1)
func DefaultBackoff(interval time.Duration) time.Duration {
	scale := 2 + (rand.Float64()-0.5)*0.2
	return time.Duration(float64(interval) * scale)
}

// BackoffSleeper returns a sleeper that sleeps in a backoff manner every time get called.
// If algorithm is nil, DefaultBackoff will be used.
// Set interval and maxInterval to the same value to make it a constant interval sleeper.
// If maxInterval is not greater than 0, it will wake immediately.
func BackoffSleeper(init, maxInterval time.Duration, algorithm func(time.Duration) time.Duration) Sleeper {
	if algorithm == nil {
		algorithm = DefaultBackoff
	}

	return func(ctx context.Context) error {
		// wake immediately
		if maxInterval <= 0 {
			return nil
		}

		interval := init
		if init < maxInterval {
			interval = algorithm(init)
		} else {
			interval = maxInterval
		}

		t := time.NewTicker(interval)
		defer t.Stop()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			init = interval
		}

		return nil
	}
}

// MergeSleepers into one, wakes when first sleeper wakes.
// Such as you want to poll to check if wifi is connected, when you poll you also
// want to do do the check whenever wifi driver is enabled, then you can
// merge BackoffSleeper and ChannelSleeper to achieve it.
func MergeSleepers(list ...Sleeper) Sleeper {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithCancel(ctx)
		done := make(chan error, len(list))

		sleep := func(s Sleeper) {
			done <- s(ctx)
			cancel()
		}

		for _, s := range list {
			go sleep(s)
		}

		return <-done
	}
}

// Sleep the goroutine for specified seconds, such as 2.3 seconds
func Sleep(seconds float64) {
	d := time.Duration(seconds * float64(time.Second))
	time.Sleep(d)
}

// Retry fn and sleeper until fn returns true or s returns error
func Retry(ctx context.Context, s Sleeper, fn func() (stop bool, err error)) error {
	for {
		stop, err := fn()
		if stop {
			return err
		}
		err = s(ctx)
		if err != nil {
			return err
		}
	}
}
