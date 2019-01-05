package main

import (
	"fmt"
	"hash/fnv"
	"os"
	"strings"

	g "github.com/ysmood/gokit"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type options struct {
	patterns *[]string
	dir      *string
	cmd      []string
	prefix   *string
}

func main() {
	optsList := []*options{}
	for _, args := range split(os.Args[1:], "---") {
		optsList = append(optsList, genOptions(args))
	}

	fns := []func(){}
	for _, opts := range optsList {
		fns = append(fns, func(opts *options) func() {
			return func() {
				g.E(g.Guard(opts.cmd, *opts.patterns, &g.GuardOptions{
					ExecOpts: &g.ExecOptions{
						Dir:    *opts.dir,
						Prefix: genPrefix(*opts.prefix, opts.cmd),
					},
				}))
			}
		}(opts))
	}
	g.All(fns...)
}

func genOptions(args []string) *options {
	opts := &options{}

	app := kingpin.New(
		"guard",
		`run and guard a command, kill and rerun it when watched files are modified

		Examples:

		 # follow the "--" is the command and its arguments you want to execute
		 guard -- ls -al

		 # the pattern must be quoted
		 guard -w '*.go' -w 'lib/**/*.go' -- go run main.go

		 # the output will be prefix with red 'my-app | '
		 guard -p 'my-app | @red' -- python test.py
		 
		 # use "---" as separator to watch multiple commands
		 guard -p a/* -- ls a --- -p b/* -- ls b
		`,
	)
	opts.patterns = app.Flag("watch", "the pattern to watch, can set multiple patterns, syntax https://github.com/bmatcuk/doublestar#patterns").Short('w').Strings()
	opts.dir = app.Flag("dir", "base dir path").Short('d').String()
	opts.prefix = app.Flag("prefix", "prefix for command output").Short('p').String()

	app.Version("0.0.2")

	args, cmdArgs := parseArgs(args)

	_, err := app.Parse(args)

	if err != nil {
		fmt.Println("for help run: guard --help")
		panic(err)
	}

	if cmdArgs == nil {
		panic("empty command")
	}

	opts.cmd = cmdArgs

	return opts
}

func parseArgs(args []string) (appArgs []string, cmdArgs []string) {
	i := indexOf(args, "--")

	if i == -1 {
		return args, nil
	}

	return args[:i], args[i+1:]
}

func indexOf(list []string, str string) int {
	for i, elem := range list {
		if elem == str {
			return i
		}
	}

	return -1
}

func genPrefix(prefix string, args []string) string {
	if prefix != "" {
		return prefix
	}

	h := fnv.New32a()
	h.Write([]byte(strings.Join(args, "")))

	return g.C(fmt.Sprint(args[0], " | "), fmt.Sprint(h.Sum32()%256))
}

func split(args []string, sep string) [][]string {
	list := [][]string{}

	tmp := []string{}
	for _, arg := range args {
		if arg == "---" {
			list = append(list, tmp)
			tmp = []string{}
			continue
		}
		tmp = append(tmp, arg)
	}

	return append(list, tmp)
}
