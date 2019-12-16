package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"sync"
	"text/template"

	"github.com/tidwall/gjson"
)

// JSONResult shortcut for gjson.Result
type JSONResult = *gjson.Result

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

// MustToJSONBytes encode data to json bytes
func MustToJSONBytes(data interface{}) []byte {
	bytes, err := json.Marshal(data)
	E(err)
	return bytes
}

// MustToJSON encode data to json string
func MustToJSON(data interface{}) string {
	return string(MustToJSONBytes(data))
}

// JSON parse json for easily access the value from json path
func JSON(data interface{}) JSONResult {
	var res gjson.Result
	switch v := data.(type) {
	case string:
		res = gjson.Parse(v)
	case []byte:
		res = gjson.ParseBytes(v)
	}

	return &res
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

// RandBytes generate random bytes with specified byte length
func RandBytes(len int) []byte {
	b := make([]byte, len)
	_, _ = rand.Read(b)
	return b
}

// RandString generate random string with specified string length
func RandString(len int) string {
	b := RandBytes(len)
	return hex.EncodeToString(b)
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

	dict := map[string]interface{}{}
	fnDict := template.FuncMap{}

	l := len(params)
	for i := 0; i < l-1; i += 2 {
		k := params[i].(string)
		v := params[i+1]
		if reflect.TypeOf(v).Kind() == reflect.Func {
			fnDict[k] = v
		} else {
			dict[k] = v
		}
	}

	t := template.Must(template.New("").Funcs(fnDict).Parse(tpl))
	E(t.Execute(&out, dict))

	return out.String()
}
