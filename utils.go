package gokit

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"go/build"
	"os"
	"sync"
	"time"

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
