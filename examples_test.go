package kit_test

import (
	"github.com/ysmood/kit"
)

func ExampleExec() {
	kit.Exec("echo", "ok").MustDo()

	str := kit.Exec("echo", "ok").MustString()

	kit.Log(str)
}

func ExampleReq() {
	val := kit.Req("http://test.com").Post().Query(
		"search", "keyword",
		"even", []string{"array", "is", "supported"},
	).MustJSON().Get("json.path.value").String()

	kit.Log(val)
}

func ExampleServer() {
	server := kit.MustServer(":8080")
	server.Engine.GET("/", func(ctx kit.GinContext) {
		ctx.String(200, "ok")
	})
	server.MustDo()
}

func ExampleWalk() {
	kit.Log(kit.Walk("**/*.go", "**/*.md", kit.WalkGitIgnore).MustList())
}

func ExampleGuard() {
	kit.Guard("go", "run", "./server").ExecCtx(
		kit.Exec().Prefix("server | @yellow"),
	).MustDo()
}
