package gokit_test

import (
	"testing"

	g "github.com/ysmood/gokit"
)

func TestLog(t *testing.T) {
	g.Log("ok")
	g.Err("err")
	g.Dump(10)
	g.ClearScreen()
}
