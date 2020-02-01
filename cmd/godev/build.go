package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/mholt/archiver"
	gos "github.com/ysmood/kit/pkg/os"
	"github.com/ysmood/kit/pkg/run"
	"github.com/ysmood/kit/pkg/utils"
)

type buildTask struct {
	dir  string
	os   string
	name string
	bin  string
	out  string
	zip  string
}

func build(patterns []string, dist string, deploy, noGitClean bool, version string, isZip bool, osList []string) {
	if deploy && !noGitClean {
		checkGitClean()
	}

	_ = gos.Remove(dist)

	bTasks := genBuildTasks(patterns, dist, osList)
	tasks := []func(){}
	for _, task := range bTasks {
		func(ctx *buildTask) {
			tasks = append(tasks, func() { ctx.build(isZip) })
		}(task)
	}
	utils.All(tasks...)

	if deploy {
		deployToGithub(bTasks, version)
	}
}

func checkGitClean() {
	out := run.Exec("git", "status", "--porcelain").MustString()
	if out != "" {
		panic("git status must be clean before deployment:\n" + out)
	}
}

func deployToGithub(bTasks []*buildTask, tag string) {
	if tag == "" {
		if len(bTasks) > 0 {
			ver, err := run.Exec("go", "run", bTasks[0].dir, "--version").String()
			if err == nil {
				tag = strings.TrimSpace(ver)
			}
		}
	}

	if _, err := semver.ParseTolerant(tag); err != nil {
		panic("invalid semver flag: --version " + tag + " (" + err.Error() + ")")
	}

	_ = run.Exec("git", "tag", tag).Do()
	_ = run.Exec("git", "push", "origin", tag).Do()

	_, err := exec.LookPath("hub")
	if err != nil {
		panic("please install hub.github.com first")
	}

	gos.Retry(5, 3*time.Second, func() {
		_ = run.Exec("hub", "release", "delete", tag).Raw().Do()

		args := []string{"hub", "release", "create", "-m", tag}
		for _, t := range bTasks {
			args = append(args, "-a", t.zip)
		}
		args = append(args, tag)

		run.Exec(args...).Raw().MustDo()
	})
}

func (ctx *buildTask) build(isZip bool) {
	utils.Log("building:", ctx.dir, "->", ctx.out)

	env := []string{
		"GOOS=" + goos(ctx.os),
		"GOARCH=amd64",
	}

	utils.E(run.Exec(
		"go", "build",
		"-ldflags=-w -s",
		"-o", ctx.out,
		ctx.dir,
	).Env(env...).Do())

	if isZip {
		if ctx.os == "linux" {
			compressGz(ctx.out, ctx.zip, ctx.bin)
		} else {
			compressZip(ctx.out, ctx.zip, ctx.bin)
		}
	}

	utils.Log("build done:", ctx.out)
}

func goos(name string) string {
	if name == "mac" {
		return "darwin"
	}
	return name
}

func genBuildTasks(patterns []string, dist string, osList []string) []*buildTask {
	if osList == nil {
		osList = []string{"mac", "linux", "windows"}
	}

	list := gos.Walk(patterns...).MustList()

	tasks := []*buildTask{}
	for _, dir := range list {
		name := filepath.Base(dir)
		for _, os := range osList {
			bin := name
			if os == "windows" {
				bin += ".exe"
			}

			out := filepath.Join(dist, name+"-"+os)

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

	utils.E(tar.Close())
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

	utils.E(tar.Close())
}
