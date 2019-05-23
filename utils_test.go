package gokit_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	g "github.com/ysmood/gokit"
)

type T = testing.T

func TestAll(t *testing.T) {
	g.All(func() {
		fmt.Println("one")
	}, func() {
		fmt.Println("two")
	})
}

func TestE(t *testing.T) {
	defer func() {
		r := recover()

		assert.Equal(t, "exec: \"exitexit\": executable file not found in $PATH", r.(error).Error())
	}()

	g.E(g.Exec("exitexit").Do())
}

func TestRetry(t *testing.T) {
	count := 0
	errs := g.Retry(3, 30*time.Millisecond, func() {
		count = count + 1
	})

	assert.Equal(t, true, errs == nil)
	assert.Equal(t, 1, count)
}

func TestRetryHalf(t *testing.T) {
	count := 0
	errs := g.Retry(5, 30*time.Millisecond, func() {
		count = count + 1

		if count < 3 {
			panic(count)
		}
	})

	assert.Equal(t, true, errs == nil)
	assert.Equal(t, 3, count)
}

func TestRetry3Times(t *testing.T) {
	count := 0
	errs := g.Retry(3, 30*time.Millisecond, func() {
		count = count + 1
		panic(count)
	})

	assert.Equal(t, []interface{}{1, 2, 3}, errs)
	assert.Equal(t, 3, count)
}

func TestGoPath(t *T) {
	s := g.GoPath()

	assert.True(t, g.Exists(s))
}

func TestTry(t *T) {
	err := g.Try(func() {
		panic("err")
	})

	assert.Equal(t, "err", err)
}
