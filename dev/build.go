package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/mholt/archiver"
	g "github.com/ysmood/gokit"
)

func build(deployTag *string) {
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
				return func() { buildForOS(n, osn) }
			}(name, osName))
		}
	}
	g.All(tasks...)

	if deployTag != nil {
		deploy(*deployTag)
	}
}

func deploy(tag string) {
	files, err := g.Glob([]string{"dist/*"}, nil)
	g.E(err)

	g.Exec([]string{"hub", "release", "delete", tag}, &g.ExecOptions{IsRaw: true})

	args := []string{"hub", "release", "create", "-m", tag}
	for _, f := range files {
		args = append(args, "-a", f)
	}
	args = append(args, tag)

	g.E(g.Exec(args, nil))
}

func buildForOS(name, osName string) {
	g.Log("building:", name, osName)

	f := fmt.Sprint

	env := []string{
		f("GOOS=", osName),
		"GOARCH=amd64",
	}

	oPath := f("dist/", name, "-", osName)

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

	compress(oPath, f(oPath, ".tar.gz"), name+extByOS(osName))

	g.Remove(oPath)

	g.Log("build done:", name, osName)
}

func extByOS(osName string) string {
	if osName == "windows" {
		return ".exe"
	}
	return ""
}

func compress(from, to, name string) {
	file, err := os.Open(from)
	g.E(err)
	fileInfo, err := file.Stat()
	g.E(err)

	tar := archiver.NewTarGz()
	tar.CompressionLevel = 9
	oFile, err := os.Create(to)
	tar.Create(oFile)

	tar.Write(archiver.File{
		FileInfo: archiver.FileInfo{
			FileInfo:   fileInfo,
			CustomName: name,
		},
		ReadCloser: file,
	})

	g.E(err)
	tar.Close()
}