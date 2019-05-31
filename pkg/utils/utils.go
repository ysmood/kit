package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"sync"
	"text/template"

	"github.com/tidwall/gjson"
)

// Nil used to create empty channel
type Nil struct{}

// Noop swallow all args and do nothing
func Noop(_ ...interface{}) {}

// ErrArg get the last arg as error
func ErrArg(args ...interface{}) error {
	return args[len(args)-1].(error)
}

// E the last arg is error, panic it
func E(args ...interface{}) []interface{} {
	err, ok := args[len(args)-1].(error)
	if ok {
		panic(err)
	}
	return args
}

// E1 if the second arg is error panic it, or return the first arg
func E1(arg interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return arg
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
func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}

// GenerateRandomString ...
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
