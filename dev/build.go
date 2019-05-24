package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/mholt/archiver"
	"github.com/ysmood/gokit"
	. "github.com/ysmood/gokit/pkg/exec"
	. "github.com/ysmood/gokit/pkg/os"
	. "github.com/ysmood/gokit/pkg/utils"
)

func build(deployTag *bool) {
	list, err := Walk("cmd/*").List()

	if err != nil {
		panic(err)
	}

	_ = Remove("dist/**")

	tasks := []func(){}
	for _, name := range list {
		name = path.Base(name)

		for _, osName := range []string{"darwin", "linux", "windows"} {
			tasks = append(tasks, func(n, osn string) func() {
				return func() { buildForOS(n, osn) }
			}(name, osName))
		}
	}
	All(tasks...)

	if *deployTag {
		deploy(gokit.Version)
	}
}

func deploy(tag string) {
	files, err := Walk("dist/*").List()
	E(err)

	_ = Exec("hub", "release", "delete", tag).Raw().Do()

	args := []string{"hub", "release", "create", "-m", tag}
	for _, f := range files {
		args = append(args, "-a", f)
	}
	args = append(args, tag)

	E(Exec(args...).Do())
}

func buildForOS(name, osName string) {
	Log("building:", name, osName)

	f := fmt.Sprint

	env := []string{
		f("GOOS=", osName),
		"GOARCH=amd64",
	}

	oPath := f("dist/", name, "-", osName)

	if osName == "darwin" {
		oPath = f("dist/", name, "-mac")
	}

	E(Exec(
		"go", "build",
		"-ldflags=-w -s",
		"-o", oPath,
		f("./cmd/", name),
	).Cmd(&exec.Cmd{
		Env: append(os.Environ(), env...),
	}).Do())

	compress(oPath, f(oPath, ".zip"), name+extByOS(osName))

	_ = Remove(oPath)

	Log("build done:", name, osName)
}

func extByOS(osName string) string {
	if osName == "windows" {
		return ".exe"
	}
	return ""
}

func compress(from, to, name string) {
	file, err := os.Open(from)
	E(err)
	fileInfo, err := file.Stat()
	E(err)

	tar := archiver.NewZip()
	tar.CompressionLevel = 9
	oFile, err := os.Create(to)
	E(err)
	E(tar.Create(oFile))

	E(tar.Write(archiver.File{
		FileInfo: archiver.FileInfo{
			FileInfo:   fileInfo,
			CustomName: name,
		},
		ReadCloser: file,
	}))

	tar.Close()
}
