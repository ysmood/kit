package exec_test

import (
	"testing"

	. "github.com/ysmood/gokit/pkg/exec"
)

func TestEnsureGoTool(t *testing.T) {
	EnsureGoTool("github.com/ysmood/gokit")
}
