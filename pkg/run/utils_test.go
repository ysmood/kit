package run_test

import (
	"testing"

	kit "github.com/ysmood/gokit"
)

func TestEnsureGoTool(t *testing.T) {
	kit.MustGoTool("github.com/ysmood/gokit")
}
