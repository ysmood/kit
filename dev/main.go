package main

import (
	"os"

	. "github.com/ysmood/gokit/pkg/exec"
	. "github.com/ysmood/gokit/pkg/guard"

	. "github.com/ysmood/gokit/pkg/utils"
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
		test()

	case cmdBuild.FullCommand():
		E(Exec("go", "test", "-coverprofile="+covPath, "-covermode=count", "./...").Do())
		export()
		build(deployTag)

	case viewCov.FullCommand():
		E(Exec("go", "tool", "cover", "-html="+covPath).Do())

	case cmdExport.FullCommand():
		export()
	}
}

func lab() {
	E(Guard("go", "run", "./dev/lab").Do())
}

func test() {
	EnsureGoTool("github.com/kyoh86/richgo")

	_ = Guard(
		"richgo", "test",
		"-coverprofile="+covPath,
		"-covermode=count",
		"-run", *testMatch,
		"./...",
	).Do()
}
