package utils_test

import (
	"encoding/base64"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	. "github.com/ysmood/gokit/pkg/exec"
	. "github.com/ysmood/gokit/pkg/utils"
)

type T = testing.T

func TestAll(t *testing.T) {
	All(func() {
		fmt.Println("one")
	}, func() {
		fmt.Println("two")
	})
}

func TestE(t *testing.T) {
	defer func() {
		r := ErrArg(recover())

		assert.EqualError(t, r, "exec: \"exitexit\": executable file not found in $PATH")
	}()

	E(Exec("exitexit").Do())
}

func TestE1(t *testing.T) {
	defer func() {
		r := ErrArg(recover())

		assert.EqualError(t, r, "err")
	}()

	E1("ok", nil)
	E1("ok", errors.New("err"))
}

func TestRetry(t *testing.T) {
	count := 0
	errs := Retry(3, 30*time.Millisecond, func() {
		count = count + 1
	})

	assert.Equal(t, true, errs == nil)
	assert.Equal(t, 1, count)
}

func TestRetryHalf(t *testing.T) {
	count := 0
	errs := Retry(5, 30*time.Millisecond, func() {
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
	errs := Retry(3, 30*time.Millisecond, func() {
		count = count + 1
		panic(count)
	})

	assert.Equal(t, []interface{}{1, 2, 3}, errs)
	assert.Equal(t, 3, count)
}

func TestTry(t *T) {
	err := Try(func() {
		panic("err")
	})

	assert.Equal(t, "err", err)
}

func TestJSON(t *T) {
	a := JSON("10")
	b := JSON([]byte("10"))

	assert.Equal(t, a.Int(), b.Int())
}

func TestGenerateRandomString(t *T) {
	v := GenerateRandomString(10)
	raw, _ := base64.URLEncoding.DecodeString(v)
	assert.Len(t, raw, 10)
}

func TestSTemplate(t *T) {
	out := S(
		"{{.a}} {{.b}} {{.c.A}}",
		"a", "<value>",
		"b", 10,
		"c", struct{ A string }{"ok"},
	)
	assert.Equal(t, "<value> 10 ok", out)
}

func TestWaitSignal(t *T) {
	go WaitSignal(nil)
}
