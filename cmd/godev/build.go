package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/blang/semver/v4"
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

func build(patterns []string, dist string, deploy bool, version string, isZip bool, osList []string) {
	_ = gos.Remove(dist)

	bTasks := genBuildTasks(patterns, dist, osList)
	tasks := []func(){}
	for _, task := range bTasks {
		func(ctx *buildTask) {
			tasks = append(tasks, func() { ctx.build(isZip) })
		}(task)
	}
	utils.All(tasks...)()

	if deploy {
		deployToGithub(bTasks, version)
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

	gos.RetryPanic(5, 3*time.Second, func() {
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
		"-trimpath",
		"-ldflags=-w -s",
		"-o", ctx.out,
		ctx.dir,
	).Env(env...).Do())

	if isZip {
		compress(ctx.out, ctx.zip, ctx.bin)
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

func compress(from, to, name string) {
	fi, err := os.Stat(from)
	utils.E(err)
	src, err := os.Open(from)
	utils.E(err)
	utils.E(gos.OutputFile(to, "", nil))
	dst, err := os.OpenFile(to, os.O_RDWR, 0664)

	var compressor io.Writer
	var close func()
	if filepath.Ext(to) == ".zip" {
		compressor, close = compressZip(fi, dst)
	} else {
		compressor, close = compressGz(fi, dst)
	}

	utils.E(io.Copy(compressor, src))
	close()
}

func compressZip(fi os.FileInfo, dst io.Writer) (io.Writer, func()) {
	zw := zip.NewWriter(dst)
	h, err := zip.FileInfoHeader(fi)
	utils.E(err)
	w, err := zw.CreateHeader(h)
	utils.E(err)
	return w, func() {
		utils.E(zw.Close())
	}
}

func compressGz(fi os.FileInfo, dst io.Writer) (io.Writer, func()) {
	gw := gzip.NewWriter(dst)
	tw := tar.NewWriter(gw)

	h, err := tar.FileInfoHeader(fi, "")
	utils.E(err)
	utils.E(tw.WriteHeader(h))

	return tw, func() {
		utils.E(tw.Close())
		utils.E(gw.Close())
	}
}
