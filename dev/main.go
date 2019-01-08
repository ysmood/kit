package main

import (
	"os"

	g "github.com/ysmood/gokit"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app      = kingpin.New("dev", "dev tool for gokit")
	cmdTest  = app.Command("test", "run test").Default()
	cmdLab   = app.Command("lab", "run lab")
	cmdBuild = app.Command("build", "cross build project")
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case cmdLab.FullCommand():
		g.Guard([]string{"go", "run", "./dev/lab"}, nil, nil)

	case cmdTest.FullCommand():
		g.Guard([]string{"go", "test", "-v", "./..."}, nil, nil)

	case cmdBuild.FullCommand():
		build()
	}
}
