package main

import (
	"fmt"
	"hash/fnv"
	"os"
	"strings"

	g "github.com/ysmood/gokit"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New(
		"guard",
		`run and guard a command, kill and rerun it when watched files are modified

		Examples:

		 guard -- ls -al # follow the "--" is the command and its arguments you want to execute
		 guard -w '*.go' -w 'lib/**/*.go' -- go run main.go # the pattern must be quoted
		 guard -p 'my-app | @red' -- python test.py # the output will be prefix with red 'my-app | '
		`,
	)
	patterns := app.Flag("watch", "the pattern to watch, can set multiple patterns, syntax https://github.com/bmatcuk/doublestar#patterns").Short('w').Strings()
	dir := app.Flag("dir", "base dir path").Short('d').String()
	prefix := app.Flag("prefix", "prefix for command output").Short('p').String()

	app.Version("0.0.1")

	args, cmdArgs := getArgs()

	_, err := app.Parse(args)

	if err != nil {
		fmt.Println("for help run: guard --help")
		panic(err)
	}

	if cmdArgs == nil {
		panic("empty command")
	}

	g.Guard(cmdArgs, *patterns, &g.ExecOptions{
		Dir:    *dir,
		Prefix: genPrefix(*prefix, cmdArgs),
	})
}

func getArgs() (args []string, cmdArgs []string) {
	i := indexOf(os.Args, "--")

	if i == -1 {
		return os.Args[1:], nil
	}

	return os.Args[1:i], os.Args[i+1:]
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
