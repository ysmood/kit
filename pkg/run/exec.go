package run

import (
	"io"
	"os/exec"
	"strings"

	"github.com/ysmood/gokit/pkg/os"
	gos "github.com/ysmood/gokit/pkg/os"
	"github.com/ysmood/gokit/pkg/utils"
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

func (ctx *ExecContext) do() {
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
}

// Do ...
func (ctx *ExecContext) Do() error {
	ctx.do()

	return run(formatPrefix(ctx.prefix), ctx.isRaw, ctx.cmd)
}

// MustDo ...
func (ctx *ExecContext) MustDo() {
	utils.E(ctx.Do())
}

// String ...
func (ctx *ExecContext) String() (string, error) {
	ctx.do()

	b, err := ctx.cmd.CombinedOutput()

	return string(b), err
}

// MustString ...
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

func pipeToStdoutWithPrefix(prefix string, reader io.Reader) error {
	const size = 32 * 1024
	buf := make([]byte, size)
	prefixBuf := []byte(prefix)
	bufOut := make([]byte, size+len(prefixBuf))

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			_, _ = gos.Stdout.Write(buf[:n])
			return err
		}

		bufOutIndex := 0
		bufOutIndex += copy(bufOut[bufOutIndex:], prefixBuf)
		for _, r := range string(buf[:n]) {
			if r == '\n' {
				bufOutIndex += copy(bufOut[bufOutIndex:], prefixBuf)
			}
			bufOutIndex += copy(bufOut[bufOutIndex:], []byte(string(r)))
		}
		utils.E(gos.Stdout.Write(bufOut[:bufOutIndex]))
	}

	return nil
}
