# Gokit

[![GoDoc](https://godoc.org/github.com/ysmood/kit?status.svg)](http://godoc.org/github.com/ysmood/kit)
[![Build Status](https://travis-ci.org/ysmood/kit.svg?branch=master)](https://travis-ci.org/ysmood/kit)
[![Build status](https://ci.appveyor.com/api/projects/status/im4xdodkpfd5vvwg/branch/master?svg=true)](https://ci.appveyor.com/project/ysmood/kit/branch/master)
[![codecov](https://codecov.io/gh/ysmood/kit/branch/master/graph/badge.svg)](https://codecov.io/gh/ysmood/kit)
[![goreport](https://goreportcard.com/badge/github.com/ysmood/kit)](https://goreportcard.com/report/github.com/ysmood/kit)

This project is a collection of often used io related methods with sane defaults to make coding less verbose.

## Modules

### os

Covers most used os related functions that are missing from the stdlib.

Such as the smart glob that handles git ignore properly:

```go
package main

import "github.com/ysmood/kit"

func main() {
    kit.Log(kit.Walk("**/*.go", "**/*.md", kit.WalkGitIgnore).MustList())
}

```

A better `Exec` alternative:

```go
package main

import "github.com/ysmood/kit"

func main() {
    kit.Exec("echo", "ok").MustDo()

    str := kit.Exec("echo", "ok").MustString()

    kit.Log(str)
}

```

### http

The http lib from stdlib is pretty verbose to use. The `kit.Req` is a much better
alternative to use with it's fluent api design. It helps to reduce the code without sacrificing performance and
flexiblity. No api is hidden from the origin http lib.

```go
package main

import "github.com/ysmood/kit"

func main() {
    val := kit.Req("http://test.com").Post().Query(
        "search", "keyword",
        "even", []string{"array", "is", "supported"},
    ).MustJSON().Get("json.path.value").String()

    kit.Log(val)
}

```

```go
package main

import "github.com/ysmood/kit"

func main() {
    server := kit.MustServer(":8080")
    server.Engine.GET("/", func(ctx kit.GinContext) {
        ctx.String(200, "ok")
    })
    server.MustDo()
}

```

## CLI tool

Goto the [release page](https://github.com/ysmood/kit/releases) to download the exectuable for your OS.
Or install with single line command below.

### godev

A general dev tool for go project to lint, test, build, deploy cross platform executables.

`godev` will release your project to your own github release page.
Project released with `godev` can be installed via [this shell script](https://github.com/ysmood/github-install).
`godev` itself is an example of how to use it.

Install `godev`: `curl -L https://git.io/fjaxx | repo=ysmood/kit bin=godev sh`

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

  build [<flags>] [<pattern>...]
    build [and deploy] specified dirs

  cov
    view html coverage report



```

### guard

Install `guard`: `curl -L https://git.io/fjaxx | repo=ysmood/kit bin=guard sh`

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

import "github.com/ysmood/kit"

func main() {
    kit.Guard("go", "run", "./server").ExecCtx(
        kit.Exec().Prefix("server | @yellow"),
    ).MustDo()
}

```

## Development

To write testable code, I try to isolate all error related dependencies.

### Build Project

Under project root

```bash
go get ./cmd/...
kit-dev build
```

Binaries will be built into dist folder.
