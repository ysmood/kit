package utils_test

import (
	"testing"

	"github.com/ysmood/kit"
)

func TestLog(t *testing.T) {
	kit.Log("ok")
	kit.Err("err")
	kit.Dump(10)
	kit.Sdump("ok")
	kit.E(kit.ClearScreen())
}
