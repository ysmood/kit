package run

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/ysmood/kit/pkg/utils"
)

// ExecContext ...
type ExecContext struct {
	context context.Context
	cmd     *exec.Cmd
	dir     string

	// Prefix prefix has a special syntax, the string after "@" can specify the color
	// of the prefix and will be removed from the output
	prefix string

	isRaw bool // Set the terminal to raw mode

	args []string
	env  []string
}

// Exec executes os command and auto pipe stdout and stdin
func Exec(args ...string) *ExecContext {
	return &ExecContext{
		args: args,
	}
}

// Context sets the context of the Command
func (ctx *ExecContext) Context(c context.Context) *ExecContext {
	ctx.context = c
	return ctx
}

// Args sets arguments
func (ctx *ExecContext) Args(args []string) *ExecContext {
	ctx.args = args
	return ctx
}

// Env appends the current env with strings, each string should be something like "key=value"
func (ctx *ExecContext) Env(env ...string) *ExecContext {
	if ctx.env == nil {
		ctx.env = append(os.Environ(), env...)
	} else {
		ctx.env = append(ctx.env, env...)
	}
	return ctx
}

// NewEnv overrides the parrent Env with the env passed in
func (ctx *ExecContext) NewEnv(env ...string) *ExecContext {
	ctx.env = env
	return ctx
}

// Dir sets the working dir to execute
func (ctx *ExecContext) Dir(dir string) *ExecContext {
	ctx.dir = dir
	return ctx
}

// Prefix sets the prefix string of the stdout and stderr
func (ctx *ExecContext) Prefix(p string) *ExecContext {
	ctx.prefix = p
	return ctx
}

// Raw enables the raw stdin mode
func (ctx *ExecContext) Raw() *ExecContext {
	ctx.isRaw = true
	return ctx
}

// GetCmd gets the exec.Cmd to execute
func (ctx *ExecContext) GetCmd() *exec.Cmd {
	if ctx.cmd != nil {
		return ctx.cmd
	}

	if len(ctx.args) == 0 {
		return nil
	}

	if ctx.context == nil {
		ctx.context = context.Background()
	}

	cmd := exec.CommandContext(ctx.context, LookPath(ctx.args[0]), ctx.args[1:]...)

	if ctx.cmd == nil {
		ctx.cmd = cmd
	}
	if ctx.dir != "" {
		ctx.cmd.Dir = ctx.dir
	}
	if ctx.env != nil {
		ctx.cmd.Env = ctx.env
	}

	ctx.cmd.Path = cmd.Path
	ctx.cmd.Args = cmd.Args

	return ctx.cmd
}

// Do the exec.Cmd
func (ctx *ExecContext) Do() error {
	cmd := ctx.GetCmd()

	return run(formatPrefix(ctx.prefix), ctx.isRaw, cmd)
}

// MustDo ...
func (ctx *ExecContext) MustDo() {
	utils.E(ctx.Do())
}

// String ...
func (ctx *ExecContext) String() (string, error) {
	cmd := ctx.GetCmd()

	b, err := cmd.CombinedOutput()

	return string(b), err
}

// MustString ...
func (ctx *ExecContext) MustString() string {
	out, err := ctx.String()
	if err != nil {
		utils.Err(out)
		panic(err)
	}
	return out
}

func formatPrefix(prefix string) string {
	i := strings.LastIndex(prefix, "@")
	if i == -1 {
		return prefix
	}

	color := prefix[i+1:]

	return utils.C(prefix[:i], color)
}

func pipeToStdoutWithPrefix(prefix string, reader io.Reader) {
	const size = 32 * 1024
	buf := make([]byte, size)
	prefixBuf := []byte(prefix)
	bufOut := make([]byte, size+len(prefixBuf))

	bufOutIndex := 0
	newline := true
	for {
		n, rerr := reader.Read(buf)

		for _, r := range string(buf[:n]) {
			if newline {
				bufOutIndex += copy(bufOut[bufOutIndex:], prefixBuf)
				newline = false
			}
			if r == '\n' {
				newline = true
			}
			bufOutIndex += copy(bufOut[bufOutIndex:], []byte(string(r)))
		}
		_, _ = utils.Stdout.Write(bufOut[:bufOutIndex])
		bufOutIndex = 0

		if rerr != nil {
			if rerr == io.EOF {
				break
			}
			return
		}
	}
}
