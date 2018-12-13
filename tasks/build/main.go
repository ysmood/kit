package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	g "github.com/ysmood/gokit"
)

func main() {
	list, err := g.Glob([]string{"cmd/*"}, nil)

	if err != nil {
		panic(err)
	}

	for _, name := range list {
		name = path.Base(name)
		build(name, "darwin")
		build(name, "linux")
		build(name, "windows")
	}
}

func build(name, osName string) {
	env := []string{
		fmt.Sprint("GOOS=", osName),
		"GOARCH=amd64",
	}

	g.Exec([]string{
		"go", "build",
		"-o", fmt.Sprint("dist/", fmt.Sprint(name, "-", osName)),
		fmt.Sprint("./cmd/", name),
	}, &g.ExecOptions{
		Cmd: &exec.Cmd{
			Env: append(os.Environ(), env...),
		},
	})
}
