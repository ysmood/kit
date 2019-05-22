package gokit

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"go/build"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/kataras/iris/core/errors"
	"github.com/mgutz/ansi"
	"github.com/tidwall/gjson"
)

// Nil used to create empty channel
type Nil struct{}

// E if find an error in args, panic it
func E(args ...interface{}) []interface{} {
	for _, arg := range args {
		err, ok := arg.(error)

		if ok {
			panic(err)
		}
	}
	return args
}

// JSON parse json for easily access the value from json path
func JSON(data interface{}) (res gjson.Result) {
	switch v := data.(type) {
	case string:
		res = gjson.Parse(v)
	case []byte:
		res = gjson.ParseBytes(v)
	}

	return res
}

// All run all actions concurrently
func All(actions ...func()) {
	wg := &sync.WaitGroup{}

	wg.Add(len(actions))

	runner := func(action func()) {
		defer wg.Done()
		action()
	}

	for _, action := range actions {
		go runner(action)
	}

	wg.Wait()
}

// GenerateRandomBytes ...
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString ...
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// C color terminal string
func C(str interface{}, color string) string {
	return ansi.Color(fmt.Sprint(str), color)
}

// GoPath get the current GOPATH properly
func GoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

// Try try fn with recover, return the panic as value
func Try(fn func()) (err interface{}) {
	defer func() {
		err = recover()
	}()

	fn()

	return err
}

// Retry retry function after a duration for several times
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

// ParamsRest ...
type ParamsRest struct {
	Params interface{}
}

// Params auto assign params by their types
func Params(values []interface{}, typedParams ...interface{}) error {
	dict, list, rest := paramsMap(typedParams)

	for _, value := range values {
		v := reflect.ValueOf(value)

		typeKey := v.Type()

		arr := dict[typeKey]
		var p *reflect.Value
		if arr != nil {
			p, dict[typeKey] = dict[typeKey][0], dict[typeKey][1:]
		}

		if p == nil {
			t := v.Type()
			for _, pp := range list {
				if pp.Kind() == reflect.Ptr && pp.Elem().Kind() == reflect.Interface {
					if t.Implements(pp.Type().Elem()) {
						p = pp
					}
				} else if pp.Kind() == reflect.Func && pp.Type().In(0).Kind() == reflect.Interface {
					if t.Implements(pp.Type().In(0)) {
						p = pp
					}
				}
			}
			if p == nil {
				if rest.IsValid() {
					if rest.Elem().Type().Elem() == v.Type() {
						rest.Elem().Set(reflect.Append(rest.Elem(), v))
						continue
					} else {
						return errors.New("rest params type error: " + spew.Sdump(value))
					}
				}
				return errors.New("params type error: " + spew.Sdump(value))
			}
		}

		if p.Kind() == reflect.Func {
			err := p.Call([]reflect.Value{v})
			if err != nil {
				e := err[0]
				if !e.IsNil() {
					return e.Interface().(error)
				}
			}
		} else {
			p.Elem().Set(v)
		}
	}
	return nil
}

func paramsMap(typedParams []interface{}) (map[reflect.Type][]*reflect.Value, []*reflect.Value, reflect.Value) {
	dict := map[reflect.Type][]*reflect.Value{}
	list := []*reflect.Value{}
	var rest reflect.Value

	for _, param := range typedParams {
		r, ok := param.(ParamsRest)
		if ok {
			rest = reflect.ValueOf(r.Params)
			continue
		}

		v := reflect.ValueOf(param)

		var key reflect.Type
		if v.Kind() == reflect.Func {
			key = v.Type().In(0)
		} else {
			key = v.Elem().Type()
		}
		list = append(list, &v)
		if dict[key] == nil {
			dict[key] = []*reflect.Value{&v}
		} else {
			dict[key] = append(dict[key], &v)
		}
	}
	return dict, list, rest
}
