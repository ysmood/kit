package gokit_test

import (
	"testing"

	. "github.com/ysmood/gokit"
)

func TestLintWholeProject(_ *testing.T) {
	EnsureGoTool("github.com/golangci/golangci-lint/cmd/golangci-lint")
	Exec("golangci-lint", "run").MustDo()
}

func ExampleExec() {
	Exec("echo", "ok").MustDo()

	str := Exec("echo", "ok").MustString()
	Log(str)
}

func ExampleReq() {
	val := Req("http://test.com").Post().Query(
		"search", "keyword",
		"even", []string{"array", "is", "supported"},
	).MustGJSON("json.path.value")

	Log(val)
}

func ExampleWalk() {
	Log(Walk("**/*.go", "**/*.md", WalkGitIgnore).MustList())
}

func ExampleGuard() {
	Guard("go", "run", "./server").ExecCtx(
		Exec().Prefix("server | @yellow"),
	).MustDo()
}
