package gokit_test

import (
	"testing"

	. "github.com/ysmood/gokit"
)

func TestLintWholeProject(_ *testing.T) {
	// lint first
	EnsureGoTool("github.com/golangci/golangci-lint/cmd/golangci-lint")
	E(Exec("golangci-lint", "run").Do())
}

func ExampleExec() {
	E(Exec("echo", "ok").Do())
}

func ExampleReq() {
	val := Req("http://test.com").Post().Query(
		"search", "keyword",
		"even", []string{"array", "is", "supported"},
	).GJSON("json.path.value")

	Log(val)
}

func ExampleWalk() {
	Log(Walk("**/*.go", "**/*.md", WalkGitIgnore).List())
}

func ExampleGuard() {
	_ = Guard("go", "run", "./server").ExecCtx(
		Exec().Prefix("server | @yellow"),
	).Do()
}
