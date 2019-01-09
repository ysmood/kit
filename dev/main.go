package main

import (
	"os"

	g "github.com/ysmood/gokit"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app       = kingpin.New("dev", "dev tool for gokit")
	cmdTest   = app.Command("test", "run test").Default()
	cmdLab    = app.Command("lab", "run lab")
	cmdBuild  = app.Command("build", "cross build project")
	deployTag = cmdBuild.Flag("deploy", "release to github with tag (install hub.github.com first)").Short('d').String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case cmdLab.FullCommand():
		lab()

	case cmdTest.FullCommand():
		test()

	case cmdBuild.FullCommand():
		g.E(g.Exec([]string{"go", "test", "./..."}, nil))
		build(deployTag)
	}
}

func lab() {
	g.Guard([]string{"go", "run", "./dev/lab"}, nil, nil)
}

func test() {
	g.Guard([]string{"go", "test", "-v", "./..."}, nil, nil)
}
