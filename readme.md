# Gokit

[![Build Status](https://travis-ci.org/ysmood/gokit.svg?branch=master)](https://travis-ci.org/ysmood/gokit)
[![codecov](https://codecov.io/gh/ysmood/gokit/branch/master/graph/badge.svg)](https://codecov.io/gh/ysmood/gokit)

Some of the io related methods that are often used.

This library won't have the best performance, but it will have sane defaults to make coding less verbose.

This project an example to use Go only for self hosted automation.

## Install CLI Tools

Goto the [release page](https://github.com/ysmood/gokit/releases) download the binary for your OS.

## Modules

### os

Covers most used os related functions that are missing from the stdlib.

Such as smart glob:

```go
package main

import (
    . "github.com/ysmood/gokit"
)

func main() {
    list, _ := Walk("**/*.go", "**/*.md", WalkGitIgnore).List()
    Log(list)
}
```

### process

A better `Exec` alternatives for the stdlib one.

### req

The http request lib from stdlib is pretty verbose to use. The `gokit.Req` is a much better
alternative to use with it's fluent api design. You will reduce a lot of your code without sacrificing performance.
It covers all the functions of the Go's stdlib one, no api is hidden from the origin http lib.

```go
package main

import (
    . "github.com/ysmood/gokit"
    // "github.com/ysmood/gokit/pkg/http" you don't have to import all sub modules
)

func main() {
    val := Req("http://test.com").Post().Query(
        "search", "keyword",
        "even": ["array", "is", "supported"]
    ).GJSON("json.path.value")

    Log(val)
}
```

## Tools

### guard

```bash
$ guard --help
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
     guard -- echo {{.op}} {{.path}}

     # watch and sync current dir to remote dir with rsync
     guard -n -- rsync {{.path}} root@host:/home/me/app/{{.path}}

     # the patterns must be quoted
     guard -w '*.go' -w 'lib/**/*.go' -- go run main.go

     # the output will be prefix with red 'my-app | '
     guard -p 'my-app | @red' -- python test.py

     # use "---" as separator to guard multiple commands
     guard -w 'a/*' -- ls a --- -w 'b/*' -- ls b


Flags:
            --help             Show context-sensitive help (also try --help-long and --help-man).
    -w, --watch=WATCH ...  the pattern to watch, can set multiple patterns
    -d, --dir=DIR          base dir path
    -p, --prefix="auto"    prefix for command output
    -n, --no-init-run      don't execute the cmd on startup
            --poll=300ms       poll interval
            --debounce=300ms   suppress the frequency of the event
            --version          Show application version.
```

## Development

### Build Project

```bash
go run ./dev --help
go run ./dev build
```

Binaries will be built into dist folder.