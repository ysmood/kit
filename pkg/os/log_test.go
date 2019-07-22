package os_test

import (
	"testing"

	kit "github.com/ysmood/gokit"
)

func TestLog(t *testing.T) {
	kit.Log("ok")
	kit.Err("err")
	kit.Dump(10)
	kit.Sdump("ok")
	kit.E(kit.ClearScreen())
}
