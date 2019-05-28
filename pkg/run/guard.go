package run

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/ysmood/gokit/pkg/os"
	"github.com/ysmood/gokit/pkg/utils"
)

// GuardContext ...
type GuardContext struct {
	args     []string
	patterns []string
	dir      string

	clearScreen  bool
	interval     *time.Duration // default 300ms
	execCtx      *ExecContext
	execCtxClone ExecContext
	debounce     *time.Duration // default 300ms
	noInitRun    bool

	prefix  string
	count   int
	wait    chan utils.Nil
	watcher *watcher.Watcher
	matcher *os.Matcher
}

// Guard run and guard a command, kill and rerun it if watched files are modified.
// Because it's based on polling, so it's cross-platform and file system.
// The args supports go template, variables {{.path}}, {{.op}} are available.
// The default patterns are GuardDefaultPatterns
func Guard(args ...string) *GuardContext {
	return &GuardContext{
		args:   args,
		prefix: os.C("[guard]", "cyan"),
		count:  0,
		wait:   make(chan utils.Nil),
	}
}

// GuardDefaultPatterns match all, then ignore all gitignore rules and all submodules
func GuardDefaultPatterns() []string {
	return []string{"**", os.WalkGitIgnore}
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

// ExecCtx ...
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

	ctx.watcher = watcher.New()
	matcher, err := os.NewMatcher(ctx.dir, ctx.patterns)
	if err != nil {
		return err
	}
	ctx.matcher = matcher

	ctx.addWatchFiles(ctx.dir)

	go ctx.watch()

	if !ctx.noInitRun {
		go ctx.run(nil)
	}

	interval := ctx.interval
	if interval == nil {
		t := time.Millisecond * 300
		interval = &t
	}

	return ctx.watcher.Start(*interval)
}

// unescape the {{.path}} {{.op}} placeholders
func (ctx *GuardContext) unescapeArgs(args []string, e *watcher.Event) []string {
	if e == nil {
		e = &watcher.Event{}
	}

	newArgs := []string{}
	for _, arg := range args {
		dir, err := filepath.Abs(ctx.dir)
		ctx.logErr(err)

		p, err := filepath.Abs(e.Path)
		ctx.logErr(err)

		p, err = filepath.Rel(dir, p)
		ctx.logErr(err)

		newArgs = append(
			newArgs,
			utils.S(arg, "path", p, "op", e.Op.String()),
		)
	}
	return newArgs
}

func (ctx *GuardContext) logErr(err error) {
	if err != nil {
		os.Log(ctx.prefix, err)
	}
}

func (ctx *GuardContext) run(e *watcher.Event) {
	if ctx.clearScreen {
		_ = os.ClearScreen()
	}

	ctx.count++
	os.Log(ctx.prefix, "run", ctx.count, os.C(ctx.args, "green"))

	ctx.execCtxClone = *ctx.execCtx
	err := ctx.execCtxClone.Dir(ctx.dir).Args(ctx.unescapeArgs(ctx.args, e)).Do()

	errMsg := ""
	if err != nil {
		errMsg = os.C(err, "red")
	}
	os.Log(ctx.prefix, "done", ctx.count, os.C(ctx.args, "green"), errMsg)

	ctx.wait <- utils.Nil{}
}

func (ctx *GuardContext) addWatchFiles(dir string) {
	list, _ := os.Walk().Dir(dir).Matcher(ctx.matcher).List()

	dict := map[string]utils.Nil{}

	for _, p := range list {
		dict[p] = utils.Nil{}
	}

	for _, p := range list {
		dir := filepath.Dir(p)
		_, has := dict[dir]

		if !has {
			dict[dir] = utils.Nil{}
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

	os.Log(ctx.prefix, "watched", len(list), "files:", os.C(watched, "green"))
}

func (ctx *GuardContext) watch() {
	debounce := ctx.debounce
	var lastRun time.Time
	if debounce == nil {
		t := time.Millisecond * 300
		debounce = &t
	}

	for {
		select {
		case e := <-ctx.watcher.Event:
			matched, _, err := ctx.matcher.Match(e.Path, e.IsDir())
			ctx.logErr(err)

			if !matched {
				continue
			}

			if time.Since(lastRun) < *debounce {
				lastRun = time.Now()
				continue
			}
			lastRun = time.Now()

			os.Log(ctx.prefix, e)

			if e.Op == watcher.Create {
				if e.IsDir() {
					ctx.addWatchFiles(e.Path)
					ctx.logErr(err)
				} else {
					_ = ctx.watcher.Add(e.Path)
				}
			}

			if ctx.execCtxClone.GetCmd() != nil {
				_ = KillTree(ctx.execCtxClone.GetCmd().Process.Pid)

				<-ctx.wait
			}

			go ctx.run(&e)

		case err := <-ctx.watcher.Error:
			ctx.logErr(err)

		case <-ctx.watcher.Closed:
			return
		}
	}
}

// MustDo ...
func (ctx *GuardContext) MustDo() {
	utils.E(ctx.Do())
}
