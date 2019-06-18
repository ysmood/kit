package main

import (
	"os"
	"os/exec"
	"path"

	"github.com/blang/semver"
	"github.com/mholt/archiver"
	gos "github.com/ysmood/gokit/pkg/os"
	"github.com/ysmood/gokit/pkg/run"
	"github.com/ysmood/gokit/pkg/utils"
)

func build(pattern []string, deployTag bool, version string) {
	list := gos.Walk(pattern...).MustList()

	_ = gos.Remove("dist")

	tasks := []func(){}
	for _, dir := range list {
		name := path.Base(dir)

		for _, osName := range []string{"darwin", "linux", "windows"} {
			tasks = append(tasks, func(n, dir, osn string) func() {
				return func() { buildForOS(n, dir, osn) }
			}(name, dir, osName))
		}
	}
	utils.All(tasks...)

	if deployTag {
		deploy(version)
	}
}

func deploy(tag string) {
	if _, err := semver.ParseTolerant(tag); err != nil {
		panic("invalid semver flag: --version " + tag + " (" + err.Error() + ")")
	}

	files := gos.Walk("dist/*").MustList()

	run.Exec("git", "tag", tag).MustDo()
	run.Exec("git", "push", "origin", tag).MustDo()

	_, err := exec.LookPath("hub")
	if err != nil {
		panic("please install hub.github.com first")
	}

	args := []string{"hub", "release", "create", "-m", tag}
	for _, f := range files {
		args = append(args, "-a", f)
	}
	args = append(args, tag)

	run.Exec(args...).Raw().MustDo()
}

func buildForOS(name, dir, osName string) {
	gos.Log("building:", name, osName)

	env := []string{
		"GOOS=" + osName,
		"GOARCH=amd64",
	}

	oPath := "dist/" + name + "-" + osName

	if osName == "darwin" {
		oPath = "dist/" + name + "-mac"
	}

	utils.E(run.Exec(
		"go", "build",
		"-ldflags=-w -s",
		"-o", oPath,
		dir,
	).Cmd(&exec.Cmd{
		Env: append(os.Environ(), env...),
	}).Do())

	if osName == "linux" {
		compressGz(oPath, oPath+".tar.gz", name+extByOS(osName))
	} else {
		compressZip(oPath, oPath+".zip", name+extByOS(osName))
	}

	_ = os.Remove(oPath)

	gos.Log("build done:", name, osName)
}

func extByOS(osName string) string {
	if osName == "windows" {
		return ".exe"
	}
	return ""
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
