package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/ysmood/kit/pkg/run"
	"github.com/ysmood/kit/pkg/utils"
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
	dev := cmd.Flag("dev", "run as dev mode").Short('d').Bool()
	min := cmd.Flag(
		"min", "if total coverage is lower than the minimum exit with non-zero",
	).Short('m').Default("0.0").Float64()
	lint := cmd.Flag("lint", "lint before test").Short('l').Bool()

	return func() {
		test(*path, *match, *min, *lint, *dev)
	}
}

func cmdBuild(cmd run.TaskCmd) func() {
	patterns := cmd.Arg("pattern", "folders to build").Default(".").Strings()
	dir := cmd.Flag("dir", "the output dir").Default("dist").String()
	deployTag := cmd.Flag("deploy", "release to github with tag").Short('d').Bool()
	ver := cmd.Flag("deploy-version", "the name of the tag").Short('v').String()
	noZip := cmd.Flag("no-zip", "don't generate zip file").Short('n').Bool()
	osList := cmd.Flag("os", "os to build, by default mac, linux and windows will be built").Strings()
	strict := cmd.Flag("strict", "strictly lint and test before build").Bool()

	return func() {
		if *strict {
			test("./...", "", 100, true, false)
		}
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

func test(path, match string, min float64, isLint, dev bool) {
	if isLint {
		lint()
	}

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
	}

	run.Exec(conf...).MustDo()

	checkCoverage(min)
}

func checkCoverage(min float64) {
	out := run.Exec("go", "tool", "cover", "-func="+*covPath).MustString()
	totalStr := regexp.MustCompile(`(\d+.\d+)%\n\z`).FindStringSubmatch(out)[1]
	total, _ := strconv.ParseFloat(totalStr, 64)
	if total < min {
		panic(fmt.Errorf("total coverage %.1f%% is less than minimum %.1f%%", total, min))
	}
	fmt.Printf("Total Cover: %.1f%%\n", total)
}
