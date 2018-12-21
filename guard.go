package gokit

import (
	"os/exec"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
)

// GuardDefaultPatterns match all, then ignore all gitignore rules and hidden files
var GuardDefaultPatterns = []string{".", "**/*", WalkGitIgnore, WalkHidden}

// GuardOptions ...
type GuardOptions struct {
	Interval time.Duration // default 300ms
	Close    chan Nil
	ExecOpts *ExecOptions
}

// Guard run and guard a command, kill and rerun it if watched files are modified.
// Because it's based on polling, so it's cross-platform and file system
func Guard(args, patterns []string, opts *GuardOptions) error {
	prefix := C("[guard]", "cyan")

	if opts == nil {
		opts = &GuardOptions{}
	}
	if opts.ExecOpts == nil {
		opts.ExecOpts = &ExecOptions{}
	}
	opts.ExecOpts.NoWait = true

	if patterns == nil || len(patterns) == 0 {
		patterns = GuardDefaultPatterns
	}

	var cmd *exec.Cmd
	wait := make(chan struct{})

	run := func() {
		Log(prefix, "run command", C(args, "green"))

		var err error
		cmd, err = Exec(args, opts.ExecOpts)

		if err != nil {
			Err(prefix, C(err, "red"))
		}

		err = cmd.Wait()
		if err != nil {
			Err(prefix, C(err, "red"))
		}
		Log(prefix, "command done", C(args, "green"))
		wait <- struct{}{}
	}

	go run()

	w := watcher.New()

	watchFiles := func(dir string) error {
		list, err := Glob(patterns, &WalkOptions{Dir: dir})

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

	m, err := newMatcher(opts.ExecOpts.Dir, patterns)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-w.Event:
				if event.Op == watcher.Create {
					continue
				}

				matched, _, err := m.match(event.Path, event.IsDir())
				if err != nil {
					Err(err)
				}

				if !matched {
					continue
				}

				Log(prefix, event)

				if event.Op == watcher.Create {
					if event.IsDir() {
						if err := watchFiles(event.Path); err != nil {
							Err(err)
						}
					} else {
						w.Add(event.Path)
					}
				}

				err = KillTree(cmd.Process.Pid)

				if err != nil {
					Err(prefix, err)
				}

				<-wait

				go run()

			case err := <-w.Error:
				Err(prefix, err)

			case <-w.Closed:
				return
			}
		}
	}()

	go func() {
		if opts.Close == nil {
			return
		}

		<-opts.Close
		w.Close()
	}()

	interval := opts.Interval
	if opts.Interval == 0 {
		interval = time.Millisecond * 300
	}

	return w.Start(interval)
}
