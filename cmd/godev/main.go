package main

import (
	"fmt"
	gos "os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ysmood/kit/pkg/os"
	"github.com/ysmood/kit/pkg/run"
	"github.com/ysmood/kit/pkg/utils"
)

var covPath *string

func main() {
	defer func() {
		err := recover()
		if err != nil {
			_, _ = fmt.Fprintln(utils.Stderr, err)
			gos.Exit(1)
		}
	}()

	app := run.TasksNew("godev", "dev tool for common go project")
	app.Version(utils.Version)

	covPath = app.Flag("cov-path", "path for coverage output").Default("coverage.txt").String()

	run.Tasks().App(app).Add(
		run.Task("test", "run go unit test").Init(cmdTest),
		run.Task("lint", "lint project with golint and golangci-lint").Init(cmdLint),
		run.Task("build", "build [and deploy] specified dirs").Init(cmdBuild),
		run.Task("cov", "view html coverage report").Run(cov),
	).Do()
}

func cmdTest(cmd run.TaskCmd) func() {
	cmd.Default()

	match := cmd.Arg("match", "match test name").String()
	path := cmd.Flag("path", "the base dir of path").Short('p').Default("./...").Strings()
	dev := cmd.Flag("dev", "run as dev mode").Short('d').Bool()
	fast := cmd.Flag("failfast", "fail on first error").Short('f').Bool()
	short := cmd.Flag("short", "run as short mode").Short('s').Bool()
	race := cmd.Flag("race", "enable race detector").Short('r').Bool()
	min := cmd.Flag(
		"min", "if total coverage percentage is lower than the minimum exit with non-zero",
	).Short('m').Default("0.0").Float64()
	lint := cmd.Flag("lint", "lint before test").Short('l').Bool()
	verbose := cmd.Flag("verbose", "enable verbose").Short('v').Bool()

	return func() {
		test(*path, *match, *min, *lint, *dev, *fast, *short, *race, *verbose)
	}
}

func cmdBuild(cmd run.TaskCmd) func() {
	patterns := cmd.Arg("pattern", "folders to build").Default(".").Strings()
	dir := cmd.Flag("dir", "the output dir").Default("dist").String()
	deploy := cmd.Flag("deploy", "release to github with tag").Short('d').Bool()
	ver := cmd.Flag("deploy-version", "the name of the tag").Short('v').String()
	noZip := cmd.Flag("no-zip", "don't generate zip file").Short('n').Bool()
	osList := cmd.Flag("os", "os to build, by default mac, linux and windows will be built").Strings()
	strict := cmd.Flag("strict", "strictly lint and test before build").Bool()

	return func() {
		if *strict {
			test([]string{"./..."}, "", 100, true, false, false, false, false, false)
		}
		build(*patterns, *dir, *deploy, *ver, !*noZip, *osList)
	}
}

func cov() {
	run.Exec("go", "tool", "cover", "-html="+*covPath).MustDo()
}

func cmdLint(cmd run.TaskCmd) func() {
	path := cmd.Arg("path", "match test name").Default("./...").Strings()

	return func() {
		lint(*path)
	}
}

func lint(path []string) {
	run.Exec("go", "mod", "tidy").MustDo()
	checkGitClean()

	run.MustGoTool("golang.org/x/lint/golint")
	fmt.Println(fmt.Sprintf("[lint] golang.org/x/lint/golint %v", path))
	args := append([]string{"golint", "-set_exit_status"}, path...)
	run.Exec(args...).MustDo()

	run.MustGoTool("github.com/kisielk/errcheck")
	fmt.Println(fmt.Sprintf("[lint] github.com/kisielk/errcheck %v", path))
	args = append([]string{"errcheck"}, path...)
	run.Exec(args...).MustDo()

	cyclePath := path
	if path[0] == "./..." {
		cyclePath = []string{"."}
	}
	run.MustGoTool("github.com/fzipp/gocyclo")
	fmt.Println(fmt.Sprintf("[lint] github.com/fzipp/gocyclo %v", cyclePath))
	args = append([]string{"gocyclo", "-over", "15"}, cyclePath...)
	run.Exec(args...).MustDo()

	fmtPath := path
	if path[0] == "./..." {
		fmtPath = []string{"."}
	}
	fmt.Println(fmt.Sprintf("[lint] gofmt -s -l -w %v", fmtPath))
	// gofmt doesn't return non-zero when fails, we have to check return manually
	args = append([]string{"gofmt", "-s", "-l", "-w"}, fmtPath...)
	out := strings.TrimSpace(run.Exec(args...).MustString())
	if out != "" {
		panic("\"gofmt -s\" check failed:\n" + out)
	}
}

func test(path []string, match string, min float64, isLint, dev, fast, short, race, verbose bool) {
	if isLint {
		lint(path)
	}

	conf := []string{
		"go",
		"test",
		"-coverprofile=" + *covPath,
		"-count=1", // prevent the go test cache
		"-run", match,
	}

	conf = append(conf, path...)

	if dev {
		run.MustGoTool("github.com/kyoh86/richgo")
		conf[0] = "richgo"
	}

	if fast {
		conf = append(conf, "-failfast")
	}

	if short {
		conf = append(conf, "-short")
	}

	if race {
		conf = append(conf, "-race")
	} else {
		conf = append(conf, "-covermode=count")
	}

	if verbose {
		conf = append(conf, "-v")
	}

	run.Exec(conf...).Raw().MustDo()

	checkCoverage(min)
}

func checkCoverage(min float64) {
	if s, _ := os.ReadString(*covPath); s == "" || s == "mode: count\n" {
		return
	}

	out := run.Exec("go", "tool", "cover", "-func="+*covPath).MustString()
	totalStr := regexp.MustCompile(`(\d+.\d+)%\n\z`).FindStringSubmatch(out)[1]
	total, _ := strconv.ParseFloat(totalStr, 64)
	if total < min {
		panic(fmt.Sprintf("[lint] Test coverage %.1f%% must >= %.1f%%", total, min))
	}
	fmt.Printf("Test coverage: %.1f%%\n", total)
}

func checkGitClean() {
	out := run.Exec("git", "status", "--porcelain").MustString()
	if out != "" {
		panic("[lint] Changes of \"go generate\", \"lint auto fix\", etc are not committed:\n" + out)
	}
}
