package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/mholt/archiver"
	gos "github.com/ysmood/gokit/pkg/os"
	"github.com/ysmood/gokit/pkg/run"
	"github.com/ysmood/gokit/pkg/utils"
)

func build(deployTag bool) {
	list := gos.Walk("cmd/*").MustList()

	_ = gos.Remove("dist/**")

	tasks := []func(){}
	for _, name := range list {
		name = path.Base(name)

		for _, osName := range []string{"darwin", "linux", "windows"} {
			tasks = append(tasks, func(n, osn string) func() {
				return func() { buildForOS(n, osn) }
			}(name, osName))
		}
	}
	utils.All(tasks...)

	if deployTag {
		deploy(utils.Version)
	}
}

func deploy(tag string) {
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

func buildForOS(name, osName string) {
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
		"./cmd/"+name,
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
	tar.CompressionLevel = 9
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
	tar.CompressionLevel = 9
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

func genReadme() {
	fexmaple := utils.E1(gos.ReadStringFile("kit_test.go")).(string)

	fset := token.NewFileSet()
	fast := utils.E1(parser.ParseFile(fset, "kit_test.go", fexmaple, parser.ParseComments)).(*ast.File)

	guardHelp := run.Exec("go", "run", "./cmd/guard", "--help").MustString()

	list := []interface{}{
		"GuardHelp", guardHelp,
	}

	ast.Inspect(fast, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			name := x.Name.Name
			code := fexmaple[x.Body.Pos()-1 : x.Body.End()]

			if strings.HasPrefix(name, "Example") {
				list = append(list, name, formatExample(code))
			}
		}
		return true
	})

	f := utils.E1(gos.ReadStringFile("readme.tpl.md")).(string)

	utils.E(gos.OutputFile("readme.md", utils.S(f, list...), nil))
}

func formatExample(code string) string {
	return utils.S(
		strings.Join(
			[]string{
				"```go",
				"package main",
				"",
				"import . \"github.com/ysmood/gokit\"",
				"",
				"func main() {{.code}}",
				"```",
			},
			"\n",
		),
		"code", code,
	)
}
