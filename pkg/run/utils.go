package run

import (
	"os"
	"os/exec"
	os_path "path"
	"strings"

	gos "github.com/ysmood/kit/pkg/os"
	"github.com/ysmood/kit/pkg/utils"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var goPathCache string

// GoPath gets the current GOPATH properly
func GoPath() string {
	if goPathCache != "" {
		return goPathCache
	}
	path, _ := exec.Command("go", "env", "GOPATH").CombinedOutput()
	goPathCache = strings.TrimSpace(string(path))
	return goPathCache
}

var goBinCache string

// GoBin gets the current GOBIN properly
func GoBin() string {
	if goBinCache != "" {
		return goBinCache
	}
	path, _ := exec.Command("go", "env", "GOBIN").CombinedOutput()
	goBinCache = strings.TrimSpace(string(path))

	if goBinCache == "" {
		goBinCache = os_path.Join(GoPath(), "bin")
	}
	return goBinCache
}

// LookPath finds the executable from PATH and GOBIN
func LookPath(name string) string {
	path, err := exec.LookPath(name)
	if err == nil {
		return path
	}

	path = os_path.Join(GoBin(), name)
	if gos.FileExists(path) {
		return path
	}

	return name
}

// MustGoTool try to find executable bin under GOPATH, if not exists run go get to download it
func MustGoTool(path string) {
	p := os_path.Join(GoBin(), os_path.Base(path)+gos.ExecutableExt())
	if !gos.Exists(p) {
		utils.Log("go get", path)
		Exec("go", "get", path).Dir(gos.HomeDir()).MustDo()
	}
}

// TasksContext ...
type TasksContext struct {
	app   *kingpin.Application
	tasks []*TaskContext
}

// Tasks a simple wrapper for kingpin to make it easier to use
func Tasks() *TasksContext {
	return &TasksContext{
		tasks: []*TaskContext{},
	}
}

// TasksNew ...
var TasksNew = kingpin.New

// App ...
func (ctx *TasksContext) App(app *kingpin.Application) *TasksContext {
	ctx.app = app
	return ctx
}

// Add ...
func (ctx *TasksContext) Add(tasks ...*TaskContext) *TasksContext {
	ctx.tasks = append(ctx.tasks, tasks...)
	return ctx
}

// Do ...
func (ctx *TasksContext) Do() {
	if ctx.app == nil {
		ctx.app = kingpin.New("", "")
	}

	for _, task := range ctx.tasks {
		cmd := ctx.app.Command(task.name, task.help)
		if task.run == nil {
			task.run = task.init(cmd)
		}
	}

	target := kingpin.MustParse(ctx.app.Parse(os.Args[1:]))

	for _, task := range ctx.tasks {
		if task.name == target {
			task.run()
		}
	}
}

// TaskCmd ...
type TaskCmd = *kingpin.CmdClause

// TaskContext ...
type TaskContext struct {
	name string
	help string
	run  func()
	init func(TaskCmd) func()
}

// Task ...
func Task(name, help string) *TaskContext {
	return &TaskContext{
		name: name,
		help: help,
	}
}

// Run ...
func (ctx *TaskContext) Run(f func()) *TaskContext {
	ctx.run = f
	return ctx
}

// Init ...
func (ctx *TaskContext) Init(f func(TaskCmd) func()) *TaskContext {
	ctx.init = f
	return ctx
}
