package main

import (
	"github.com/ysmood/gokit/pkg/run"
	"github.com/ysmood/gokit/pkg/utils"
)

func main() {
	app := run.TasksNew("dev", "dev tool for gokit")
	app.Version(utils.Version)

	run.Tasks().App(app).Add(
		run.Task("build", "").Init(cmdBuild),
		run.Task("readme", "build readme").Run(genReadme),
		run.Task("export", "export all submodules under gokit namespace").Run(export),
	).Do()
}

func cmdBuild(cmd run.TaskCmd) func() {
	deploy := cmd.Flag("deploy", "release to github with tag").Short('d').Bool()

	args := []string{
		"godev", "build",
		"--strict",
		"cmd/*", "!cmd/gokit-dev",
	}

	return func() {
		if *deploy {
			args = append(args, "-d")
		}

		export()
		genReadme()
		run.Exec(args...).MustDo()
	}
}
