package kit

func ExampleExec() {
	Exec("echo", "ok").MustDo()

	str := Exec("echo", "ok").MustString()
	Log(str)
}

func ExampleReq() {
	val := Req("http://test.com").Post().Query(
		"search", "keyword",
		"even", []string{"array", "is", "supported"},
	).MustJSON("json.path.value")

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
