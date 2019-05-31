package os_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	kit "github.com/ysmood/gokit"
)

type T = testing.T

func TestWaitSignal(t *T) {
	go func() {
		time.Sleep(time.Millisecond)
		p, err := os.FindProcess(os.Getpid())
		kit.E(err)

		p.Signal(os.Interrupt)
	}()
	kit.WaitSignal(nil)
}

func TestRetry(t *T) {
	count := 0
	errs := kit.Retry(3, time.Nanosecond, func() {
		count = count + 1
	})

	assert.Equal(t, true, errs == nil)
	assert.Equal(t, 1, count)
}

func TestRetryHalf(t *T) {
	count := 0
	errs := kit.Retry(5, time.Nanosecond, func() {
		count = count + 1

		if count < 3 {
			panic(count)
		}
	})

	assert.Equal(t, true, errs == nil)
	assert.Equal(t, 3, count)
}

func TestRetry3Times(t *T) {
	count := 0
	errs := kit.Retry(3, time.Nanosecond, func() {
		count = count + 1
		panic(count)
	})

	assert.Equal(t, []interface{}{1, 2, 3}, errs)
	assert.Equal(t, 3, count)
}
