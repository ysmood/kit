package gokit_test

import (
	"testing"

	. "github.com/ysmood/gokit/pkg/exec"
	. "github.com/ysmood/gokit/pkg/utils"
)

func TestLintWholeProject(_ *testing.T) {
	// lint first
	EnsureGoTool("github.com/golangci/golangci-lint/cmd/golangci-lint")
	E(Exec("golangci-lint", "run").Do())
}
