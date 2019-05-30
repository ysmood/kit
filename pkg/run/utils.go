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

// TaskCmd ...
type TaskCmd = *kingpin.CmdClause

// Task ...
type Task struct {
	Help string
	Init func(TaskCmd) func()
	Task func()
}

// Tasks ...
type Tasks = map[string]Task

// TaskNew ...
var TaskNew = kingpin.New

// TaskRun a simple wrapper for kingpin to make it easier to use
// The app arg can be nil
func TaskRun(app *kingpin.Application, tasks Tasks) {
	if app == nil {
		app = kingpin.New("", "")
	}

	callbacks := map[string]func(){}
	for name, task := range tasks {
		cmd := app.Command(name, task.Help)
		if task.Init == nil {
			callbacks[name] = task.Task
		} else {
			callbacks[name] = task.Init(cmd)
		}
	}

	name := kingpin.MustParse(app.Parse(os.Args[1:]))

	callbacks[name]()
}
