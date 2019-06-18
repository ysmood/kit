package run

import (
	"os"
	os_path "path"

	gos "github.com/ysmood/gokit/pkg/os"
	"github.com/ysmood/gokit/pkg/utils"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// MustGoTool try to find executable bin under GOPATH, if not exists run go get to download it
func MustGoTool(path string) {
	if !gos.Exists(gos.GoPath() + "/bin/" + os_path.Base(path)) {
		utils.E(Exec("go", "get", path).Dir(gos.HomeDir()).Do())
	}
}

// TasksContext ...
type TasksContext struct {
	app   *kingpin.Application
	tasks map[string]*TaskContext
}

// Tasks a simple wrapper for kingpin to make it easier to use
func Tasks() *TasksContext {
	return &TasksContext{
		tasks: map[string]*TaskContext{},
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
	for _, task := range tasks {
		ctx.tasks[task.name] = task
	}
	return ctx
}

// Do ...
func (ctx *TasksContext) Do() {
	if ctx.app == nil {
		ctx.app = kingpin.New("", "")
	}

	for name, task := range ctx.tasks {
		cmd := ctx.app.Command(name, task.help)
		if task.run == nil {
			task.run = task.init(cmd)
		}
	}

	name := kingpin.MustParse(ctx.app.Parse(os.Args[1:]))

	ctx.tasks[name].run()
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
