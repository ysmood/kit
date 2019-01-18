package gokit

import (
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/hoisie/mustache"
)

// GuardDefaultPatterns match all, then ignore all gitignore rules and all submodules
// "." is for CREAT and REMOVE event
var GuardDefaultPatterns = []string{".", "**", WalkGitIgnore}

// GuardOptions ...
type GuardOptions struct {
	Interval  *time.Duration // default 300ms
	Stop      chan Nil       // send signal to it to stop the watcher
	ExecOpts  *ExecOptions
	Debounce  *time.Duration // default 300ms, suppress the frequency of the event
	NoInitRun bool
}

// Guard run and guard a command, kill and rerun it if watched files are modified.
// Because it's based on polling, so it's cross-platform and file system.
// The args supports mustach template, variables {{path}}, {{op}} are available.
// The default patterns are GuardDefaultPatterns
func Guard(args, patterns []string, opts *GuardOptions) error {
	prefix := C("[guard]", "cyan")

	if opts == nil {
		opts = &GuardOptions{}
	}
	if opts.ExecOpts == nil {
		opts.ExecOpts = &ExecOptions{}
	}

	if patterns == nil || len(patterns) == 0 {
		patterns = GuardDefaultPatterns
	}

	var cmd *exec.Cmd
	wait := make(chan struct{})

	opts.ExecOpts.OnStart = func(opts *ExecOptions) {
		cmd = opts.Cmd
	}

	unescapeArgs := func(args []string, e *fsnotify.Event) []string {
		if e == nil {
			e = &fsnotify.Event{}
		}

		newArgs := []string{}
		for _, arg := range args {
			dir, err := filepath.Abs(opts.ExecOpts.Dir)
			if err != nil {
				Err(err)
			}

			p, err := filepath.Abs(e.Name)
			if err != nil {
				Err(err)
			}

			p, err = filepath.Rel(dir, p)
			if err != nil {
				Err(err)
			}

			newArgs = append(
				newArgs,
				mustache.Render(arg, map[string]string{"path": p, "op": e.Op.String()}),
			)
		}
		return newArgs
	}

	run := func(e *fsnotify.Event) {
		eArgs := unescapeArgs(args, e)
		Log(prefix, "run", C(eArgs, "green"))

		err := Exec(eArgs, opts.ExecOpts)
		errMsg := ""
		if err != nil {
			errMsg = C(err, "red")
		}
		Log(prefix, "done", C(args, "green"), errMsg)

		wait <- struct{}{}
	}

	w, err := NewWatcher()
	if err != nil {
		return err
	}
	matcher, err := NewMatcher(opts.ExecOpts.Dir, patterns)
	if err != nil {
		return err
	}

	watchFiles := func(dir string) error {
		list, err := Glob(patterns, &WalkOptions{Dir: dir, Matcher: matcher})

		if err != nil {
			return err
		}

		for _, p := range list {
			w.Add(p)
		}

		var watched string
		if len(list) > 10 {
			watched = strings.Join(append(list[0:10], "..."), " ")
		} else {
			watched = strings.Join(list, " ")
		}

		Log(prefix, "watched", len(list), "files:", C(watched, "green"))

		return nil
	}

	if err := watchFiles(opts.ExecOpts.Dir); err != nil {
		return err
	}

	go func() {
		debounce := opts.Debounce
		var lastRun time.Time
		if debounce == nil {
			t := time.Second
			debounce = &t
		}

		for {
			select {
			case e, ok := <-w.Events:
				if !ok {
					return
				}

				if time.Since(lastRun) < *debounce {
					lastRun = time.Now()
					continue
				}
				lastRun = time.Now()

				if e.Op&fsnotify.Create == fsnotify.Create {
					continue
				}

				isDir := DirExists(e.Name)

				matched, _, err := matcher.match(e.Name, isDir)
				if err != nil {
					Err(err)
				}

				if !matched {
					continue
				}

				Log(prefix, e)

				if e.Op&fsnotify.Create == fsnotify.Create {
					if isDir {
						if err := watchFiles(e.Name); err != nil {
							Err(err)
						}
					} else {
						w.Add(e.Name)
					}
				}

				if cmd != nil {
					KillTree(cmd.Process.Pid)

					<-wait
				}

				go run(&e)

			case err, ok := <-w.Errors:
				if !ok {
					return
				}

				Log(prefix, err)
			}
		}
	}()

	go func() {
		if opts.Stop == nil {
			return
		}

		<-opts.Stop
		w.Close()
	}()

	if !opts.NoInitRun {
		go run(nil)
	}

	interval := opts.Interval
	if interval == nil {
		t := time.Millisecond * 300
		interval = &t
	}

	return w.Start(*interval)
}
