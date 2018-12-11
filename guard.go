package gokit

import (
	"time"

	"github.com/radovskyb/watcher"
)

// Guard run and guard a command, kill and rerun it if watched files are modified
func Guard(args, patterns []string, opts *ExecOptions) error {
	if opts == nil {
		opts = &ExecOptions{}
	}
	opts.NoWait = true

	cmd, err := Exec(args, opts)

	if err != nil {
		return err
	}

	w := watcher.New()
	w.SetMaxEvents(1)

	for _, pattern := range patterns {
		list, err := Glob(pattern)

		if err != nil {
			return err
		}

		for _, path := range list {
			w.Add(path)
		}
	}

	go func() {
		for {
			select {
			case event := <-w.Event:
				Log(event)

				err := KillTree(cmd.Process.Pid)

				if err != nil {
					Err(err)
					break
				}

				cmd, err = Exec(args, opts)
				if err != nil {
					Err(err)
					break
				}
			}
		}
	}()

	return w.Start(time.Millisecond * 100)
}
