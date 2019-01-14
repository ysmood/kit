package gokit

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	ps "github.com/mitchellh/go-ps"
)

// ExecOptions ...
type ExecOptions struct {
	Cmd *exec.Cmd
	Dir string

	// Prefix prefix has a special syntax, the string after "@" can specify the color
	// of the prefix and will be removed from the output
	Prefix string

	IsRaw bool
}

// Exec execute os command and auto pipe stdout and stdin
func Exec(args []string, opts *ExecOptions) error {
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
	}
	if opts.Dir != "" {
		opts.Cmd.Dir = opts.Dir
	}

	opts.Cmd.Path = cmd.Path
	opts.Cmd.Args = cmd.Args

	return run(formatPrefix(opts.Prefix), opts.IsRaw, opts.Cmd)
}

// KillTree kill process and all its children process
func KillTree(pid int) error {
	ids, err := childrenProcessIDs(pid)
	ids = append(ids, pid)

	if err != nil {
		return err
	}

	for _, id := range ids {
		p, err := os.FindProcess(id)

		if err != nil {
			return err
		}

		err = p.Signal(syscall.SIGINT)

		if err != nil {
			return err
		}
	}

	return nil
}

func formatPrefix(prefix string) string {
	i := strings.LastIndex(prefix, "@")
	if i == -1 {
		return prefix
	}

	color := prefix[i+1:]

	return C(prefix[:i], color)
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
