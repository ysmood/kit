package os_test

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/kit"
)

type T = testing.T

func TestExecutableExt(t *T) {
	var expected string
	if runtime.GOOS == "windows" {
		expected = ".exe"
	} else {
		expected = ""
	}
	assert.Equal(t, expected, kit.ExecutableExt())
}

func TestWaitSignal(t *T) {
	if runtime.GOOS == "windows" {
		// TODO: seems like the SIGINT will force the exit code to 2 event tests all pass.
		// Not sure why this happens, for now I have to skip this test for windows.
		return
	}

	go func() {
		time.Sleep(time.Second)
		kit.E(kit.SendSigInt(os.Getpid()))
	}()
	kit.WaitSignal()
}

func TestRetry(t *T) {
	count := 0
	errs := kit.RetryPanic(3, time.Nanosecond, func() {
		count = count + 1
	})

	assert.Equal(t, true, errs == nil)
	assert.Equal(t, 1, count)
}

func TestRetryHalf(t *T) {
	count := 0
	errs := kit.RetryPanic(5, time.Nanosecond, func() {
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
	errs := kit.RetryPanic(3, time.Nanosecond, func() {
		count = count + 1
		panic(count)
	})

	assert.Equal(t, []interface{}{1, 2, 3}, errs)
	assert.Equal(t, 3, count)
}

func TestEscape(t *T) {
	expected := "/?*"

	if runtime.GOOS == "windows" {
		expected = "／？＊"
	}

	assert.Equal(t, expected, kit.Escape("/?*"))
}
