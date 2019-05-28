package main

import (
	"os"

	"github.com/ysmood/gokit/pkg/run"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const covPath = "coverage.txt"

var (
	app       = kingpin.New("dev", "dev tool for gokit")
	cmdTest   = app.Command("test", "run test").Default()
	cmdLab    = app.Command("lab", "run lab")
	cmdBuild  = app.Command("build", "cross build project")
	cmdExport = app.Command("export", "export all submodules under gokit namespace")
	testMatch = cmdTest.Arg("match", "match test name").String()
	deployTag = cmdBuild.Flag("deploy", "release to github with tag (install hub.github.com first)").Short('d').Bool()
	viewCov   = app.Command("cov", "view html coverage report")
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case cmdLab.FullCommand():
		lab()

	case cmdTest.FullCommand():
		test(true)

	case cmdBuild.FullCommand():
		export()
		lint()
		test(false)
		genReadme()
		build(deployTag)

	case viewCov.FullCommand():
		run.Exec("go", "tool", "cover", "-html="+covPath).MustDo()

	case cmdExport.FullCommand():
		export()
	}
}

func lab() {
	run.Guard("go", "run", "./dev/lab").MustDo()
}

func lint() {
	run.MustGoTool("golang.org/x/lint/golint")
	run.Exec("golint", "-set_exit_status", "./...").MustDo()
}

func test(dev bool) {
	conf := []string{
		"go",
		"test",
		"-coverprofile=" + covPath,
		"-covermode=count",
		"-run", *testMatch,
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
