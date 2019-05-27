package exec

import (
	os_exec "os/exec"
	"strings"

	"github.com/ysmood/gokit/pkg/os"
	"github.com/ysmood/gokit/pkg/utils"
)

type ExecContext struct {
	cmd *os_exec.Cmd
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

func (ctx *ExecContext) Args(args []string) *ExecContext {
	ctx.args = args
	return ctx
}

func (ctx *ExecContext) Cmd(cmd *os_exec.Cmd) *ExecContext {
	ctx.cmd = cmd
	return ctx
}

func (ctx *ExecContext) GetCmd() *os_exec.Cmd {
	return ctx.cmd
}

func (ctx *ExecContext) Dir(dir string) *ExecContext {
	ctx.dir = dir
	return ctx
}

func (ctx *ExecContext) Prefix(p string) *ExecContext {
	ctx.prefix = p
	return ctx
}

func (ctx *ExecContext) Raw() *ExecContext {
	ctx.isRaw = true
	return ctx
}

func (ctx *ExecContext) do() {
	cmd := os_exec.Command(ctx.args[0], ctx.args[1:]...)

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
}

func (ctx *ExecContext) Do() error {
	ctx.do()

	return run(formatPrefix(ctx.prefix), ctx.isRaw, ctx.cmd)
}

func (ctx *ExecContext) MustDo() {
	utils.E(ctx.Do())
}

func (ctx *ExecContext) String() (string, error) {
	ctx.do()

	b, err := ctx.cmd.CombinedOutput()

	return string(b), err
}

func (ctx *ExecContext) MustString() string {
	return utils.E1(ctx.String()).(string)
}

func formatPrefix(prefix string) string {
	i := strings.LastIndex(prefix, "@")
	if i == -1 {
		return prefix
	}

	color := prefix[i+1:]

	return os.C(prefix[:i], color)
}
