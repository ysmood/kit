package gokit_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kataras/iris/core/errors"
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

	g.E(g.Exec("exitexit"))
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

func TestParamsAssign(t *testing.T) {
	type test_type struct {
		str string
	}
	type test_typep struct {
		str string
	}

	var str string
	var i int
	var tt test_type
	var ttp *test_typep

	err := g.Params(
		[]interface{}{test_type{"ok"}, "ok", &test_typep{"yes"}, 10},
		&str,
		&i,
		&tt,
		&ttp,
	)

	assert.Nil(t, err)
	assert.Equal(t, "ok", str)
	assert.Equal(t, 10, i)
	assert.Equal(t, "ok", tt.str)
	assert.Equal(t, "yes", ttp.str)
}

func TestParamsSameType(t *testing.T) {
	var a int
	var b int

	err := g.Params(
		[]interface{}{1, 2},
		&a,
		&b,
	)

	assert.Nil(t, err)
	assert.Equal(t, 1, a)
	assert.Equal(t, 2, b)
}

func TestParamsRest(t *testing.T) {
	rest := []int{}

	err := g.Params(
		[]interface{}{1, 2},
		g.ParamsRest{&rest},
	)

	assert.Nil(t, err)
	assert.Equal(t, 1, rest[0])
	assert.Equal(t, 2, rest[1])
}

func TestParamsRestErr(t *testing.T) {
	rest := []int{}

	err := g.Params(
		[]interface{}{1, "err"},
		g.ParamsRest{&rest},
	)

	assert.EqualError(t, err, "rest params type error: (string) (len=3) \"err\"\n")
}

func TestParamsError(t *testing.T) {
	var str string

	err := g.Params(
		[]interface{}{10},
		&str,
	)

	assert.EqualError(t, err, "params type error: (int) 10\n")
}

func TestParamsFunc(t *T) {
	type test_type struct {
		S string
	}

	var v *test_type
	var n int

	err := g.Params(
		[]interface{}{&test_type{"ok"}, 10},
		func(tt *test_type) {
			v = tt
			v.S += "ok"
		},
		func(v int) {
			n = v + 1
		},
	)

	assert.Nil(t, err)
	assert.Equal(t, "okok", v.S)
	assert.Equal(t, 11, n)
}

func TestParamsFuncErr(t *T) {

	err := g.Params(
		[]interface{}{10},
		func(v int) error {
			return errors.New("err")
		},
	)

	assert.EqualError(t, err, "err")
}
