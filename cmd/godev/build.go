package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/blang/semver"
	"github.com/mholt/archiver"
	gos "github.com/ysmood/gokit/pkg/os"
	"github.com/ysmood/gokit/pkg/run"
	"github.com/ysmood/gokit/pkg/utils"
)

type buildTask struct {
	dir  string
	os   string
	name string
	bin  string
	out  string
	zip  string
}

func build(patterns []string, deployTag bool, version string) {
	_ = gos.Remove("dist")

	bTasks := genBuildTasks(patterns)
	tasks := []func(){}
	for _, task := range bTasks {
		func(ctx *buildTask) {
			tasks = append(tasks, func() { ctx.build() })
		}(task)
	}
	utils.All(tasks...)

	if deployTag {
		deploy(bTasks, version)
	}
}

func deploy(bTasks []*buildTask, tag string) {
	if tag == "" {
		for _, t := range bTasks {
			if t.os == runtime.GOOS {
				ver, err := run.Exec(t.out, "--version").String()
				if err == nil {
					tag = strings.TrimSpace(ver)
				}
			}
		}
	}

	if _, err := semver.ParseTolerant(tag); err != nil {
		panic("invalid semver flag: --version " + tag + " (" + err.Error() + ")")
	}

	run.Exec("git", "tag", tag).MustDo()
	run.Exec("git", "push", "origin", tag).MustDo()

	_, err := exec.LookPath("hub")
	if err != nil {
		panic("please install hub.github.com first")
	}

	args := []string{"hub", "release", "create", "-m", tag}
	for _, t := range bTasks {
		args = append(args, "-a", t.zip)
	}
	args = append(args, tag)

	run.Exec(args...).Raw().MustDo()
}

func (ctx *buildTask) build() {
	gos.Log("building:", ctx.dir, "->", ctx.out)

	env := []string{
		"GOOS=" + ctx.os,
		"GOARCH=amd64",
	}

	utils.E(run.Exec(
		"go", "build",
		"-ldflags=-w -s",
		"-o", ctx.out,
		ctx.dir,
	).Cmd(&exec.Cmd{
		Env: append(os.Environ(), env...),
	}).Do())

	if ctx.os == "linux" {
		compressGz(ctx.out, ctx.zip, ctx.bin)
	} else {
		compressZip(ctx.out, ctx.zip, ctx.bin)
	}

	gos.Log("build done:", ctx.out)
}

func genBuildTasks(patterns []string) []*buildTask {
	list := gos.Walk(patterns...).MustList()

	tasks := []*buildTask{}
	for _, dir := range list {
		name := filepath.Base(dir)
		for _, os := range []string{"darwin", "linux", "windows"} {
			bin := name
			if os == "windows" {
				bin += ".exe"
			}

			out := "dist/" + name + "-" + os
			if os == "darwin" {
				out = "dist/" + name + "-mac"
			}

			zip := out + ".zip"
			if os == "linux" {
				zip = out + ".tar.gz"
			}

			tasks = append(tasks, &buildTask{
				dir:  dir,
				os:   os,
				name: name,
				bin:  bin,
				out:  out,
				zip:  zip,
			})
		}
	}

	return tasks
}

func compressZip(from, to, name string) {
	file, err := os.Open(from)
	utils.E(err)
	fileInfo, err := file.Stat()
	utils.E(err)

	tar := archiver.NewZip()
	oFile, err := os.Create(to)
	utils.E(err)
	utils.E(tar.Create(oFile))

	utils.E(tar.Write(archiver.File{
		FileInfo: archiver.FileInfo{
			FileInfo:   fileInfo,
			CustomName: name,
		},
		ReadCloser: file,
	}))

	tar.Close()
}

func compressGz(from, to, name string) {
	file, err := os.Open(from)
	utils.E(err)
	fileInfo, err := file.Stat()
	utils.E(err)

	tar := archiver.NewTarGz()
	oFile, err := os.Create(to)
	utils.E(err)
	utils.E(tar.Create(oFile))

	utils.E(tar.Write(archiver.File{
		FileInfo: archiver.FileInfo{
			FileInfo:   fileInfo,
			CustomName: name,
		},
		ReadCloser: file,
	}))

	tar.Close()
}