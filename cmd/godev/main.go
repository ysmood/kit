package main

import (
	"github.com/ysmood/gokit/pkg/run"
	"github.com/ysmood/gokit/pkg/utils"
)

var covPath *string

func main() {
	app := run.TasksNew("godev", "dev tool for common go project")
	app.Version(utils.Version)

	covPath = app.Flag("cov-path", "path for coverage output").Default("coverage.txt").String()

	run.Tasks().App(app).Add(
		run.Task("test", "run go unit test").Init(cmdTest),
		run.Task("lint", "lint project with golint and golangci-lint").Run(lint),
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
	patterns := cmd.Arg("pattern", "folders to build").Default(".").Strings()
	dir := cmd.Flag("dir", "the output dir").Default("dist").String()
	deployTag := cmd.Flag("deploy", "release to github with tag").Short('d').Bool()
	ver := cmd.Flag("deploy-version", "the name of the tag").Short('v').String()
	noZip := cmd.Flag("no-zip", "don't generate zip file").Short('n').Bool()
	osList := cmd.Flag("os", "os to build, by default mac, linux and windows will be built").Strings()

	return func() {
		lint()
		test("./...", "", false)
		build(*patterns, *dir, *deployTag, *ver, !*noZip, *osList)
	}
}

func cov() {
	run.Exec("go", "tool", "cover", "-html="+*covPath).MustDo()
}

func lint() {
	run.MustGoTool("golang.org/x/lint/golint")
	run.Exec("golint", "-set_exit_status", "./...").MustDo()

	run.MustGoTool("github.com/golangci/golangci-lint/cmd/golangci-lint")
	run.Exec("golangci-lint", "run").MustDo()
}

func test(path, match string, dev bool) {
	conf := []string{
		"go",
		"test",
		"-coverprofile=" + *covPath,
		"-count=1", // prevent the go test cache
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
