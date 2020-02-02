package utils_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/kit"
	"github.com/ysmood/kit/pkg/utils"
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

		assert.EqualError(t, r, "err")
	}()

	kit.E(func() error {
		return errors.New("err")
	}())
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

	assert.Equal(t, `{"a":1}`, kit.MustToJSON(map[string]int{"a": 1}))
}

func TestGenerateRandomString(t *T) {
	v := kit.RandString(10)
	raw, _ := hex.DecodeString(v)
	assert.Len(t, raw, 10)
}

func TestSTemplate(t *T) {
	out := kit.S(
		"{{.a}} {{.b}} {{.c.A}} {{d}}",
		"a", "<value>",
		"b", 10,
		"c", struct{ A string }{"ok"},
		"d", func() string {
			return "ok"
		},
	)
	assert.Equal(t, "<value> 10 ok ok", out)
}

func TestObservable(t *testing.T) {
	o := utils.Observable{}

	go func() {
		for i := 0; i < 10; i++ {
			o.Publish(i)
		}
	}()

	e := o.Until(func(e utils.Event) bool {
		return e.(int) == 5
	})

	assert.Equal(t, 5, e)
}
