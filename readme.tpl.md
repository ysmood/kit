# Kit

[![GoDoc](https://godoc.org/github.com/ysmood/kit?status.svg)](http://godoc.org/github.com/ysmood/kit)
[![Build Status](https://travis-ci.org/ysmood/kit.svg?branch=master)](https://travis-ci.org/ysmood/kit)
[![Build status](https://ci.appveyor.com/api/projects/status/im4xdodkpfd5vvwg/branch/master?svg=true)](https://ci.appveyor.com/project/ysmood/kit/branch/master)
[![codecov](https://codecov.io/gh/ysmood/kit/branch/master/graph/badge.svg)](https://codecov.io/gh/ysmood/kit)
[![goreport](https://goreportcard.com/badge/github.com/ysmood/kit)](https://goreportcard.com/report/github.com/ysmood/kit)

This project is a collection of often-used io related methods with sane defaults to make coding less verbose.

## Modules

### os

Covers most used os related functions that are missing from the stdlib.

Such as the smart glob that handles git ignore properly:

{{.ExampleWalk}}

A better `Exec` alternative:

{{.ExampleExec}}

### http

The http lib from stdlib is pretty verbose to use. The `kit.Req` is a much better
alternative to use with its fluent api design. It helps to reduce the code without sacrificing performance and
flexibility. No api is hidden from the origin of http lib.

{{.ExampleReq}}

{{.ExampleServer}}

## CLI tool

Goto the [release page](https://github.com/ysmood/kit/releases) to download the executable for your OS.
Or install with single line command below.

### godev

A general dev tool for go project to lint, test, build, deploy cross-platform executable.

`godev` will release your project to your own github release page.
Project released with `godev` can be installed via [this shell script](https://github.com/ysmood/github-install).
`godev` itself is an example of how to use it.

Install `godev`: `curl -L https://git.io/fjaxx | repo=ysmood/kit bin=godev sh`

```bash
{{.GodevHelp}}
```

### guard

Install `guard`: `curl -L https://git.io/fjaxx | repo=ysmood/kit bin=guard sh`

```bash
{{.GuardHelp}}
```

You can also use it as a lib:

{{.ExampleGuard}}

## Development

To write testable code, I try to isolate all error related dependencies.

### Build Project

Under project root

```bash
go run ./cmd/kit-dev build
```

Binaries will be built into dist folder.
