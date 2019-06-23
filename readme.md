# Gokit

[![GoDoc](https://godoc.org/github.com/ysmood/gokit?status.svg)](http://godoc.org/github.com/ysmood/gokit)
[![Build Status](https://travis-ci.org/ysmood/gokit.svg?branch=master)](https://travis-ci.org/ysmood/gokit)
[![Build status](https://ci.appveyor.com/api/projects/status/b8mkds289asy6q5s/branch/master?svg=true)](https://ci.appveyor.com/project/ysmood/gokit/branch/master)
[![codecov](https://codecov.io/gh/ysmood/gokit/branch/master/graph/badge.svg)](https://codecov.io/gh/ysmood/gokit)
[![goreport](https://goreportcard.com/badge/github.com/ysmood/gokit)](https://goreportcard.com/report/github.com/ysmood/gokit)

This project is a collection of often used io related methods with sane defaults to make coding less verbose.

## Modules

### os

Covers most used os related functions that are missing from the stdlib.

Such as smart glob:

```go
package main

import . "github.com/ysmood/gokit"

func main() {
	Log(Walk("**/*.go", "**/*.md", WalkGitIgnore).MustList())
}

```

A better `Exec` alternatives for the stdlib one.

```go
package main

import . "github.com/ysmood/gokit"

func main() {
	Exec("echo", "ok").MustDo()

	str := Exec("echo", "ok").MustString()
	Log(str)
}

```

### http

The http request lib from stdlib is pretty verbose to use. The `gokit.Req` is a much better
alternative to use with it's fluent api design. You will reduce a lot of your code without sacrificing performance.
It covers all the functions of the Go's stdlib one, no api is hidden from the origin http lib.

```go
package main

import . "github.com/ysmood/gokit"

func main() {
	val := Req("http://test.com").Post().Query(
		"search", "keyword",
		"even", []string{"array", "is", "supported"},
	).MustJSON("json.path.value")

	Log(val)
}

```

## CLI tool

Goto the [release page](https://github.com/ysmood/gokit/releases) download the exectuable for your OS.
Or install with single line command below.

### godev

A general dev tool for go project to lint, test, build cross platform executable, etc.

Install: `curl -L https://git.io/fjaxx | repo=ysmood/gokit bin=godev sh`

```bash
usage: godev [<flags>] <command> [<args> ...]

dev tool for common go project

Flags:
  --help                     Show context-sensitive help (also try --help-long
                             and --help-man).
  --version                  Show application version.
  --cov-path="coverage.txt"  path for coverage output

Commands:
  help [<command>...]
    Show help.

  test* [<flags>] [<match>]
    run go unit test

  lint
    lint project with golint and golangci-lint

  build [<flags>]
    build [and deploy] specified dirs

  cov
    view html coverage report



```

### guard

Install: `curl -L https://git.io/fjaxx | repo=ysmood/gokit bin=guard sh`

```bash
usage: guard [<flags>]

run and guard a command, kill and rerun it when watched files are modified

  Examples:

   # follow the "--" is the command and its arguments you want to execute
   # kill and restart the web server when a file changes
   guard -- node server.js

   # use ! prefix to ignore pattern, the below means watch all files but not those in tmp dir
   guard -w '**' -w '!tmp/**' -- echo changed

   # the special !g pattern will read the gitignore files and ignore patterns in them
   # the below is the default patterns guard will use
   guard -w '**' -w '!g' -- echo changed

   # support go template
   guard -- echo {{op}} {{path}}

   # watch and sync current dir to remote dir with rsync
   guard -n -- rsync {{path}} root@host:/home/me/app/{{path}}

   # the patterns must be quoted
   guard -w '*.go' -w 'lib/**/*.go' -- go run main.go

   # the output will be prefix with red 'my-app | '
   guard -p 'my-app | @red' -- python test.py
   
   # use "---" as separator to guard multiple commands
   guard -w 'a/*' -- ls a --- -w 'b/*' -- ls b


Flags:
      --help             Show context-sensitive help (also try --help-long and
                         --help-man).
  -w, --watch=WATCH ...  the pattern to watch, can set multiple patterns
  -d, --dir=DIR          base dir path
  -p, --prefix="auto"    prefix for command output
  -c, --clear-screen     clear screen before each run
  -n, --no-init-run      don't execute the cmd on startup
      --poll=300ms       poll interval
      --debounce=300ms   suppress the frequency of the event
      --raw              when you need to interact with the subprocess
      --version          Show application version.


```

You can also use it as a lib

```go
package main

import . "github.com/ysmood/gokit"

func main() {
	Guard("go", "run", "./server").ExecCtx(
		Exec().Prefix("server | @yellow"),
	).MustDo()
}

```

## Development

To write testable code, I try to isolate all error related dependencies.

### Build Project

Under project root

```bash
go get ./cmd/...
gokit-dev build
```

Binaries will be built into dist folder.
