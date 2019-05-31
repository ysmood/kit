package utils_test

import (
	"encoding/base64"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	kit "github.com/ysmood/gokit"
)

type T = testing.T

func TestNoop(t *T) {
	kit.Noop()
}

func TestAll(t *T) {
	kit.All(func() {
		fmt.Println("one")
	}, func() {
		fmt.Println("two")
	})
}

func TestE(t *T) {
	defer func() {
		r := kit.ErrArg(recover())

		assert.EqualError(t, r, "exec: \"exitexit\": executable file not found in $PATH")
	}()

	kit.E(kit.Exec("exitexit").Do())
}

func TestE1(t *T) {
	defer func() {
		r := kit.ErrArg(recover())

		assert.EqualError(t, r, "err")
	}()

	kit.E1("ok", nil)
	kit.E1("ok", errors.New("err"))
}

func TestTry(t *T) {
	err := kit.Try(func() {
		panic("err")
	})

	assert.Equal(t, "err", err)
}

func TestJSON(t *T) {
	a := kit.JSON("10")
	b := kit.JSON([]byte("10"))

	assert.Equal(t, a.Int(), b.Int())
}

func TestGenerateRandomString(t *T) {
	v := kit.GenerateRandomString(10)
	raw, _ := base64.URLEncoding.DecodeString(v)
	assert.Len(t, raw, 10)
}

func TestSTemplate(t *T) {
	out := kit.S(
		"{{.a}} {{.b}} {{.c.A}}",
		"a", "<value>",
		"b", 10,
		"c", struct{ A string }{"ok"},
	)
	assert.Equal(t, "<value> 10 ok", out)
}
