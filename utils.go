package gokit

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"go/build"
	"os"
	"sync"

	"github.com/mgutz/ansi"
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
