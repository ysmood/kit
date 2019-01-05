package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/mholt/archiver"
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
		g.Log("OK")

	case cmdTest.FullCommand():
		g.Guard([]string{"go", "test", "./..."}, nil, nil)

	case cmdBuild.FullCommand():
		list, err := g.Glob([]string{"cmd/*"}, nil)

		if err != nil {
			panic(err)
		}

		g.Remove("dist/**")

		tasks := []func(){}
		for _, name := range list {
			name = path.Base(name)

			for _, osName := range []string{"darwin", "linux", "windows"} {
				tasks = append(tasks, func(n, osn string) func() {
					return func() { build(n, osn) }
				}(name, osName))
			}
		}
		g.All(tasks...)
	}
}

func build(name, osName string) {
	g.Log("building:", name, osName)

	f := fmt.Sprint

	env := []string{
		f("GOOS=", osName),
		"GOARCH=amd64",
	}

	oPath := f("dist/", name, "-", osName, extByOS(osName))

	g.Exec([]string{
		"go", "build",
		"-ldflags=-w -s",
		"-o", oPath,
		f("./cmd/", name),
	}, &g.ExecOptions{
		Cmd: &exec.Cmd{
			Env: append(os.Environ(), env...),
		},
	})

	archiver.Archive([]string{oPath}, f(oPath, ".tar.gz"))

	g.Remove(oPath)

	g.Log("build done:", name, osName)
}

func extByOS(osName string) string {
	if osName == "windows" {
		return ".exe"
	}
	return ""
}
