package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"os"
	"os/signal"
	"sync"
	"text/template"
	"time"

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

func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}

func GenerateRandomString(s int) string {
	b := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b)
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

// S Template render, the params is key-value pairs
func S(tpl string, params ...interface{}) string {
	var out bytes.Buffer

	t := template.Must(template.New("").Parse(tpl))

	dict := map[string]interface{}{}
	l := len(params)
	for i := 0; i < l-1; i += 2 {
		dict[params[i].(string)] = params[i+1]
	}

	E(t.Execute(&out, dict))

	return out.String()
}

func WaitSignal(sig os.Signal) {
	c := make(chan os.Signal, 1)
	if sig == nil {
		sig = os.Interrupt
	}
	signal.Notify(c, sig)
	<-c
	close(c)
}
