package gokit

import (
	"os/exec"
	"strings"
)

// ExecContext ...
type ExecContext struct {
	cmd *exec.Cmd
	dir string

	// Prefix prefix has a special syntax, the string after "@" can specify the color
	// of the prefix and will be removed from the output
	prefix string

	isRaw bool // Set the terminal to raw mode

	args []string
}

// Exec execute os command and auto pipe stdout and stdin
func Exec(args ...string) *ExecContext {
	return &ExecContext{
		args: args,
	}
}

// Args ...
func (ctx *ExecContext) Args(args []string) *ExecContext {
	ctx.args = args
	return ctx
}

// Cmd ...
func (ctx *ExecContext) Cmd(cmd *exec.Cmd) *ExecContext {
	ctx.cmd = cmd
	return ctx
}

// GetCmd ...
func (ctx *ExecContext) GetCmd() *exec.Cmd {
	return ctx.cmd
}

// Dir ...
func (ctx *ExecContext) Dir(dir string) *ExecContext {
	ctx.dir = dir
	return ctx
}

// Prefix ...
func (ctx *ExecContext) Prefix(p string) *ExecContext {
	ctx.prefix = p
	return ctx
}

// Raw ...
func (ctx *ExecContext) Raw() *ExecContext {
	ctx.isRaw = true
	return ctx
}

// Do ...
func (ctx *ExecContext) Do() error {
	cmd := exec.Command(ctx.args[0], ctx.args[1:]...)

	if ctx.cmd == nil {
		ctx.cmd = cmd
	} else {
		clone := *ctx.cmd
		ctx.cmd = &clone
	}
	if ctx.dir != "" {
		ctx.cmd.Dir = ctx.dir
	}

	ctx.cmd.Path = cmd.Path
	ctx.cmd.Args = cmd.Args

	return run(formatPrefix(ctx.prefix), ctx.isRaw, ctx.cmd)
}

func formatPrefix(prefix string) string {
	i := strings.LastIndex(prefix, "@")
	if i == -1 {
		return prefix
	}

	color := prefix[i+1:]

	return C(prefix[:i], color)
}
