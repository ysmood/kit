package main

import (
	"github.com/ysmood/gokit/pkg/run"
)

var covPath *string

func main() {
	app := run.TasksNew("godev", "dev tool for common go project")

	covPath = app.Flag("cov-path", "path for coverage output").Default("coverage.txt").String()

	run.Tasks().App(app).Add(
		run.Task("test", "run go unit test").Init(cmdTest),
		run.Task("build", "build [and deploy] specified dirs").Init(cmdBuild),
		run.Task("cov", "view html coverage report").Run(cov),
	).Do()
}

func cmdTest(cmd run.TaskCmd) func() {
	cmd.Default()

	match := cmd.Arg("match", "match test name").String()
	path := cmd.Flag("path", "the base dir of path").Short('p').Default("./...").String()

	return func() {
		test(*path, *match, true)
	}
}

func cmdBuild(cmd run.TaskCmd) func() {
	deployTag := cmd.Flag("deploy", "release to github with tag").Short('d').Bool()
	patterns := cmd.Flag("pattern", "folders to build").Short('p').Default(".").Strings()
	ver := cmd.Flag("version", "the name of the tag").Short('v').String()

	return func() {
		lint()
		test("./...", "", false)
		build(*patterns, *deployTag, *ver)
	}
}

func cov() {
	run.Exec("go", "tool", "cover", "-html="+*covPath).MustDo()
}

func lint() {
	run.MustGoTool("golang.org/x/lint/golint")
	run.Exec("golint", "-set_exit_status", "./...").MustDo()
}

func test(path, match string, dev bool) {
	conf := []string{
		"go",
		"test",
		"-coverprofile=" + *covPath,
		"-covermode=count",
		"-run", match,
		path,
	}

	if dev {
		run.MustGoTool("github.com/kyoh86/richgo")
		conf[0] = "richgo"
		run.Guard(conf...).MustDo()
		return
	}

	run.Exec(conf...).MustDo()
}
