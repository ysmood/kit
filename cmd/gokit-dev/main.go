package main

import (
	"fmt"
	"os"

	gos "github.com/ysmood/gokit/pkg/os"
	"github.com/ysmood/gokit/pkg/run"
	"github.com/ysmood/gokit/pkg/utils"
)

const covPath = "coverage.txt"

func main() {
	os.Chdir(gos.ThisDirPath() + "/../..")

	run.Tasks().App(run.TasksNew("dev", "dev tool for gokit")).Add(
		run.Task("test", "").Init(func(cmd run.TaskCmd) func() {
			cmd.Default()

			match := cmd.Arg("match", "match test name").String()

			return func() {
				test(*match, true)
			}
		}),
		run.Task("lab", "run temp random experimental code").Run(func() {
			run.Guard("go", "run", "./dev/lab").MustDo()
		}),
		run.Task("build", "").Init(func(cmd run.TaskCmd) func() {
			deployTag := cmd.Flag("deploy", "release to github with tag").Short('d').Bool()
			return func() {
				export()
				lint()
				test("", false)
				genReadme()
				build(*deployTag)
			}
		}),
		run.Task("readme", "build readme").Run(genReadme),
		run.Task("export", "export all submodules under gokit namespace").Run(export),
		run.Task("cov", "view html coverage report").Run(func() {
			run.Exec("go", "tool", "cover", "-html="+covPath).MustDo()
		}),
	).Do()
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
		"-coverprofile=" + covPath,
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
