package main

import (
	"fmt"

	"github.com/ysmood/gokit/pkg/run"
	"github.com/ysmood/gokit/pkg/utils"
)

var covPath *string

var tasks = run.Tasks{
	"test": run.Task{Init: func(cmd run.TaskCmd) func() {
		cmd.Default()

		match := cmd.Arg("match", "match test name").String()

		return func() {
			test(*match, true)
		}
	}},
	"lab": run.Task{Help: "run temp random experimental code", Task: func() {
		run.Guard("go", "run", "./dev/lab").MustDo()
	}},
	"build": run.Task{Init: func(cmd run.TaskCmd) func() {

		deployTag := cmd.Flag("deploy", "release to github with tag").Short('d').Bool()
		return func() {
			export()
			lint()
			test("", false)
			genReadme()
			build(*deployTag)
		}
	}},
	"readme": run.Task{Task: genReadme, Help: "build readme"},
	"export": run.Task{Task: export, Help: "export all submodules under gokit namespace"},
	"cov": run.Task{Help: "view html coverage report", Task: func() {
		run.Exec("go", "tool", "cover", "-html="+*covPath).MustDo()
	}},
}

func main() {
	app := run.TaskNew("dev", "dev tool for gokit")
	covPath = app.Flag("cov-path", "coverage output file path").Default("coverage.txt").String()

	run.TaskRun(app, tasks)
}

func lint() {
	run.MustGoTool("golang.org/x/lint/golint")
	run.Exec("golint", "-set_exit_status", "./...").MustDo()
}

func test(match string, dev bool) {
	utils.Noop(Build)

	conf := []string{
		"go",
		"test",
		"-coverprofile=" + *covPath,
		"-covermode=count",
		"-run", match,
		"./...",
	}

	if dev {
		run.MustGoTool("github.com/kyoh86/richgo")
		conf[0] = "richgo"
		run.Guard(conf...).MustDo()
		return
	}

	run.Exec(conf...).MustDo()
}

// BuildArgs ...
type BuildArgs struct {
	Deploy bool `default:"true" desc:"release to github with tag (install hub.github.com first)"`
}

// Build ...
func Build(b *BuildArgs) {
	if b.Deploy {
		fmt.Println("OK")
	}
}
