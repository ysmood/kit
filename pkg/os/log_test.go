package os_test

import (
	"testing"

	. "github.com/ysmood/gokit/pkg/os"
	. "github.com/ysmood/gokit/pkg/utils"
)

func TestLog(t *testing.T) {
	Log("ok")
	Err("err")
	Dump(10)
	E(ClearScreen())
}
