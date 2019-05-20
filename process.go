package gokit

import (
	"os/exec"
	"strings"
)

// ExecOptions ...
type ExecOptions struct {
	Cmd *exec.Cmd
	Dir string

	// Prefix prefix has a special syntax, the string after "@" can specify the color
	// of the prefix and will be removed from the output
	Prefix string

	IsRaw bool // Set the terminal to raw mode

	OnStart func(opts *ExecOptions)
}

// Exec execute os command and auto pipe stdout and stdin
func Exec(params ...interface{}) error {
	var args []string
	var opts ExecOptions

	err := Params(params, ParamsRest{&args}, &args, &opts)
	if err != nil {
		return err
	}

	cmd := exec.Command(args[0], args[1:]...)

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

	if opts.OnStart != nil {
		opts.OnStart(&opts)
	}

	return run(formatPrefix(opts.Prefix), opts.IsRaw, opts.Cmd)
}

func formatPrefix(prefix string) string {
	i := strings.LastIndex(prefix, "@")
	if i == -1 {
		return prefix
	}

	color := prefix[i+1:]

	return C(prefix[:i], color)
}
