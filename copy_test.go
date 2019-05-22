package gokit_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	g "github.com/ysmood/gokit"
)

func TestCopy(t *testing.T) {
	str, _ := g.GenerateRandomString(10)
	p := fmt.Sprintf("fixtures/deep/path/%s", str)

	g.OutputFile(p+"/a/b", "ok", nil)
	g.OutputFile(p+"/a/c/c", "ok", nil)

	g.Copy(p+"/a", p+"/d")

	assert.True(t, g.Exists(p+"/d/b"))
	assert.True(t, g.Exists(p+"/d/c/c"))
}
