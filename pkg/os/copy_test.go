package os_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/ysmood/gokit/pkg/os"
	. "github.com/ysmood/gokit/pkg/utils"
)

func TestCopy(t *testing.T) {
	str := GenerateRandomString(10)
	p := fmt.Sprintf("tmp/deep/path/%s", str)

	_ = OutputFile(p+"/a/b", "ok", nil)
	_ = OutputFile(p+"/a/c/c", "ok", nil)

	_ = Copy(p+"/a", p+"/d")

	assert.True(t, Exists(p+"/d/b"))
	assert.True(t, Exists(p+"/d/c/c"))
}
