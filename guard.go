package gokit

import (
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/hoisie/mustache"
	"github.com/radovskyb/watcher"
)

// GuardContext ...
type GuardContext struct {
	Interval  *time.Duration // default 300ms
	Stop      chan Nil       // send signal to it to stop the watcher
	ExecOpts  ExecOptions
	Debounce  *time.Duration // default 300ms, suppress the frequency of the event
	NoInitRun bool

	args     []string
	patterns []string
	cmd      *exec.Cmd
	prefix   string
	count    int
	wait     chan Nil
}

// Guard run and guard a command, kill and rerun it if watched files are modified.
// Because it's based on polling, so it's cross-platform and file system.
// The args supports mustach template, variables {{path}}, {{op}} are available.
// The default patterns are GuardDefaultPatterns
func Guard(args ...string) *GuardContext {
	return &GuardContext{
		args:   args,
		prefix: C("[guard]", "cyan"),
		count:  0,
		wait:   make(chan Nil),
	}
}

// GuardDefaultPatterns match all, then ignore all gitignore rules and all submodules
func GuardDefaultPatterns() []string {
	return []string{"**", WalkGitIgnore}
}

// Patterns set patterns
func (ctx *GuardContext) Patterns(patterns ...string) *GuardContext {
	ctx.patterns = patterns
	return ctx
}

// Context set context
func (ctx *GuardContext) Context(src GuardContext) *GuardContext {
	ctx.Interval = src.Interval
	ctx.Stop = src.Stop
	ctx.ExecOpts = src.ExecOpts
	ctx.Debounce = src.Debounce
	ctx.NoInitRun = src.NoInitRun
	return ctx
}

// Do run
func (ctx *GuardContext) Do() error {
	onStart := ctx.ExecOpts.OnStart
	ctx.ExecOpts.OnStart = func(opts *ExecOptions) {
		if onStart != nil {
			onStart(opts)
		}

		ctx.count++
		Log(ctx.prefix, "run", ctx.count, C(opts.Cmd.Args, "green"))
		ctx.cmd = opts.Cmd
	}

	if ctx.patterns == nil || len(ctx.patterns) == 0 {
		ctx.patterns = GuardDefaultPatterns()
	}

	unescapeArgs := func(args []string, e *watcher.Event) []string {
		if e == nil {
			e = &watcher.Event{}
		}

		newArgs := []string{}
		for _, arg := range args {
			dir, err := filepath.Abs(ctx.ExecOpts.Dir)
			if err != nil {
				Err(err)
			}

			p, err := filepath.Abs(e.Path)
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

	run := func(e *watcher.Event) {
		eArgs := unescapeArgs(ctx.args, e)

		err := Exec(eArgs, ctx.ExecOpts)
		errMsg := ""
		if err != nil {
			errMsg = C(err, "red")
		}
		Log(ctx.prefix, "done", ctx.count, C(ctx.args, "green"), errMsg)

		ctx.wait <- struct{}{}
	}

	w := watcher.New()
	matcher, err := NewMatcher(ctx.ExecOpts.Dir, ctx.patterns)
	if err != nil {
		return err
	}

	watchFiles := func(dir string) error {
		list, err := Glob(ctx.patterns, &WalkOptions{Dir: dir, Matcher: matcher})

		if err != nil {
			return err
		}

		dict := map[string]Nil{}

		for _, p := range list {
			dict[p] = Nil{}
		}

		for _, p := range list {
			dir := filepath.Dir(p)
			_, has := dict[dir]

			if !has {
				dict[dir] = Nil{}
				w.Add(dir)
			}
			w.Add(p)
		}

		var watched string
		if len(list) > 10 {
			watched = strings.Join(append(list[0:10], "..."), " ")
		} else {
			watched = strings.Join(list, " ")
		}

		Log(ctx.prefix, "watched", len(list), "files:", C(watched, "green"))

		return nil
	}

	if err := watchFiles(ctx.ExecOpts.Dir); err != nil {
		return err
	}

	go func() {
		debounce := ctx.Debounce
		var lastRun time.Time
		if debounce == nil {
			t := time.Millisecond * 300
			debounce = &t
		}

		for {
			select {
			case e := <-w.Event:
				matched, _, err := matcher.match(e.Path, e.IsDir())
				if err != nil {
					Err(err)
				}

				if !matched {
					continue
				}

				if time.Since(lastRun) < *debounce {
					lastRun = time.Now()
					continue
				}
				lastRun = time.Now()

				Log(ctx.prefix, e)

				if e.Op == watcher.Create {
					if e.IsDir() {
						if err := watchFiles(e.Path); err != nil {
							Err(err)
						}
					} else {
						w.Add(e.Path)
					}
				}

				if ctx.cmd != nil {
					KillTree(ctx.cmd.Process.Pid)

					<-ctx.wait
				}

				go run(&e)

			case err := <-w.Error:
				Log(ctx.prefix, err)

			case <-w.Closed:
				return
			}
		}
	}()

	go func() {
		if ctx.Stop == nil {
			return
		}

		<-ctx.Stop
		w.Close()
	}()

	if !ctx.NoInitRun {
		go run(nil)
	}

	interval := ctx.Interval
	if interval == nil {
		t := time.Millisecond * 300
		interval = &t
	}

	return w.Start(*interval)
}
