package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/mholt/archiver"
	g "github.com/ysmood/gokit"
)

func main() {
	list, err := g.Glob([]string{"cmd/*"}, nil)

	if err != nil {
		panic(err)
	}

	g.Remove("dist/**")

	for _, name := range list {
		name = path.Base(name)
		build(name, "darwin")
		build(name, "linux")
		build(name, "windows")
	}
}

func build(name, osName string) {
	g.Log("build", name, osName)

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
}

func extByOS(osName string) string {
	if osName == "windows" {
		return ".exe"
	}
	return ""
}
