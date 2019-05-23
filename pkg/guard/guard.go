package guard

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/hoisie/mustache"
	"github.com/radovskyb/watcher"
	. "github.com/ysmood/gokit/pkg/exec"
	. "github.com/ysmood/gokit/pkg/os"
	. "github.com/ysmood/gokit/pkg/utils"
)

type GuardContext struct {
	args     []string
	patterns []string
	dir      string

	clearScreen bool
	interval    *time.Duration // default 300ms
	execCtx     *ExecContext
	debounce    *time.Duration // default 300ms
	noInitRun   bool

	prefix  string
	count   int
	wait    chan Nil
	watcher *watcher.Watcher
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

// Dir set dir
func (ctx *GuardContext) Dir(d string) *GuardContext {
	ctx.dir = d
	return ctx
}

// Patterns set patterns
func (ctx *GuardContext) Patterns(patterns ...string) *GuardContext {
	ctx.patterns = patterns
	return ctx
}

// NoInitRun don't execute the cmd on startup
func (ctx *GuardContext) NoInitRun() *GuardContext {
	ctx.noInitRun = true
	return ctx
}

// ClearScreen clear screen before each run
func (ctx *GuardContext) ClearScreen() *GuardContext {
	ctx.clearScreen = true
	return ctx
}

// Interval poll interval
func (ctx *GuardContext) Interval(interval *time.Duration) *GuardContext {
	ctx.interval = interval
	return ctx
}

// Debounce suppress the frequency of the event
func (ctx *GuardContext) Debounce(debounce *time.Duration) *GuardContext {
	ctx.debounce = debounce
	return ctx
}

func (ctx *GuardContext) ExecCtx(c *ExecContext) *GuardContext {
	ctx.execCtx = c
	return ctx
}

// Stop stop watching
func (ctx *GuardContext) Stop() {
	ctx.watcher.Close()
}

// Do run
func (ctx *GuardContext) Do() error {
	if ctx.patterns == nil || len(ctx.patterns) == 0 {
		ctx.patterns = GuardDefaultPatterns()
	}

	if ctx.execCtx == nil {
		ctx.execCtx = Exec()
	}

	logErr := func(err error) {
		if err != nil {
			Err(err)
		}
	}

	// unescape the {{path}} {{op}} placeholders
	unescapeArgs := func(args []string, e *watcher.Event) []string {
		if e == nil {
			e = &watcher.Event{}
		}

		newArgs := []string{}
		for _, arg := range args {
			dir, err := filepath.Abs(ctx.dir)
			logErr(err)

			p, err := filepath.Abs(e.Path)
			logErr(err)

			p, err = filepath.Rel(dir, p)
			logErr(err)

			newArgs = append(
				newArgs,
				mustache.Render(arg, map[string]string{"path": p, "op": e.Op.String()}),
			)
		}
		return newArgs
	}

	var execCtxClone ExecContext
	run := func(e *watcher.Event) {
		if ctx.clearScreen {
			_ = ClearScreen()
		}

		ctx.count++
		Log(ctx.prefix, "run", ctx.count, C(ctx.args, "green"))

		execCtxClone = *ctx.execCtx
		err := execCtxClone.Dir(ctx.dir).Args(unescapeArgs(ctx.args, e)).Do()

		errMsg := ""
		if err != nil {
			errMsg = C(err, "red")
		}
		Log(ctx.prefix, "done", ctx.count, C(ctx.args, "green"), errMsg)

		ctx.wait <- Nil{}
	}

	ctx.watcher = watcher.New()
	matcher, err := NewMatcher(ctx.dir, ctx.patterns)
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
				_ = ctx.watcher.Add(dir)
			}
			_ = ctx.watcher.Add(p)
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

	if err := watchFiles(ctx.dir); err != nil {
		return err
	}

	go func() {
		debounce := ctx.debounce
		var lastRun time.Time
		if debounce == nil {
			t := time.Millisecond * 300
			debounce = &t
		}

		for {
			select {
			case e := <-ctx.watcher.Event:
				matched, _, err := matcher.Match(e.Path, e.IsDir())
				logErr(err)

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
						err := watchFiles(e.Path)
						logErr(err)
					} else {
						_ = ctx.watcher.Add(e.Path)
					}
				}

				if execCtxClone.GetCmd() != nil {
					_ = KillTree(execCtxClone.GetCmd().Process.Pid)

					<-ctx.wait
				}

				go run(&e)

			case err := <-ctx.watcher.Error:
				Log(ctx.prefix, err)

			case <-ctx.watcher.Closed:
				return
			}
		}
	}()

	if !ctx.noInitRun {
		go run(nil)
	}

	interval := ctx.interval
	if interval == nil {
		t := time.Millisecond * 300
		interval = &t
	}

	return ctx.watcher.Start(*interval)
}
