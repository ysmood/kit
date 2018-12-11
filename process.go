package gokit

import (
	"io"
	"os"
	"os/exec"
	"regexp"
	"syscall"

	"github.com/mitchellh/go-ps"
)

// ExecOptions ...
type ExecOptions struct {
	Cmd    *exec.Cmd
	Prefix string
	NoWait bool
}

type prefixWriter struct {
	prefix string
	output io.Writer
}

var prefixReg = regexp.MustCompile(`.*\n`)

func (pw *prefixWriter) Write(p []byte) (int, error) {
	_, err := pw.output.Write(
		prefixReg.ReplaceAllFunc(p, func(m []byte) []byte {
			return append([]byte(pw.prefix), m...)
		}),
	)

	// TODO: I don't why we have to return the same size, even the size has changed,
	// or it will throw an error
	return len(p), err
}

// Exec execute os command and auto pipe stdout and stdin
func Exec(args []string, opts *ExecOptions) (*exec.Cmd, error) {
	cmd := exec.Command(args[0], args[1:]...)

	if opts == nil {
		opts = &ExecOptions{}
	} else {
		clone := *opts
		opts = &clone
	}

	if opts.Cmd == nil {
		opts.Cmd = cmd
	} else {
		clone := *opts.Cmd
		opts.Cmd = &clone
		opts.Cmd.Path = cmd.Path
		opts.Cmd.Args = cmd.Args
	}

	if opts.Cmd.Stdout == nil {
		opts.Cmd.Stdout = &prefixWriter{opts.Prefix, os.Stdout}
	}
	if opts.Cmd.Stdin == nil {
		opts.Cmd.Stderr = &prefixWriter{opts.Prefix, os.Stderr}
	}
	if opts.Cmd.Stdin == nil {
		opts.Cmd.Stdin = os.Stdin
	}

	if opts.NoWait {
		return opts.Cmd, opts.Cmd.Start()
	}
	return opts.Cmd, opts.Cmd.Run()
}

// KillTree kill process and all its children process
func KillTree(pid int) error {
	ids, err := childrenProcessIDs(pid)

	if err != nil {
		return err
	}

	for _, id := range ids {
		err = syscall.Kill(id, syscall.SIGINT)

		if err != nil {
			return err
		}
	}

	return nil
}

func childrenProcessIDs(id int) ([]int, error) {
	list, err := ps.Processes()

	if err != nil {
		return nil, err
	}

	children := []int{}
	for _, p := range list {
		if p.PPid() == id {
			children = append(children, p.Pid())
		}
	}

	return children, nil
}
