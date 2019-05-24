package os_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/ysmood/gokit/pkg/os"
	. "github.com/ysmood/gokit/pkg/utils"
)

func TestCopy(t *testing.T) {
	p := "tmp/" + GenerateRandomString(10)

	_ = OutputFile(p+"/a/b", "ok", nil)
	_ = OutputFile(p+"/a/c/c", "ok", nil)

	_ = Copy(p+"/a", p+"/d")

	assert.True(t, Exists(p+"/d/b"))
	assert.True(t, Exists(p+"/d/c/c"))
}

func TestCopyFile(t *testing.T) {
	p := "tmp/" + GenerateRandomString(10)

	_ = OutputFile(p+"/a", "ok", nil)

	_ = Copy(p+"/a", p+"/b")

	assert.True(t, Exists(p+"/b"))
}

func TestCopyErr(t *testing.T) {
	p := "tmp/" + GenerateRandomString(10)

	assert.EqualError(t, Copy(p, p), fmt.Sprintf("stat %s: no such file or directory", p))
}

func TestCopyDirErr(t *testing.T) {
	p := "tmp/" + GenerateRandomString(10)

	_ = OutputFile(p+"/a", "ok", nil)

	err := Copy(p+"/a", "/")

	assert.EqualError(t, err, "open /: is a directory")
}
