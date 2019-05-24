package utils_test

import (
	"encoding/base64"
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
		r := recover()

		assert.Equal(t, "exec: \"exitexit\": executable file not found in $PATH", r.(error).Error())
	}()

	E(Exec("exitexit").Do())
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
	out := S("{{.key}}", "key", "value")

	assert.Equal(t, "value", out)
}
