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

{{.ExampleWalk}}

A better `Exec` alternatives for the stdlib one.

{{.ExampleExec}}

### http

The http request lib from stdlib is pretty verbose to use. The `gokit.Req` is a much better
alternative to use with it's fluent api design. You will reduce a lot of your code without sacrificing performance.
It covers all the functions of the Go's stdlib one, no api is hidden from the origin http lib.

{{.ExampleReq}}

{{.ExampleServer}}

## CLI tool

Goto the [release page](https://github.com/ysmood/gokit/releases) download the exectuable for your OS.
Or install with single line command below.

### godev

A general dev tool for go project to lint, test, build cross platform executable, etc.

`godev` will release your project to your own github release page.
Project released with `godev` can be installed via [this one line command](https://github.com/ysmood/github-install).
`godev` itself is an example of how to use it.

Install `godev`: `curl -L https://git.io/fjaxx | repo=ysmood/gokit bin=godev sh`

```bash
{{.GodevHelp}}
```

### guard

Install `guard`: `curl -L https://git.io/fjaxx | repo=ysmood/gokit bin=guard sh`

```bash
{{.GuardHelp}}
```

You can also use it as a lib

{{.ExampleGuard}}

## Development

To write testable code, I try to isolate all error related dependencies.

### Build Project

Under project root

```bash
go get ./cmd/...
gokit-dev build
```

Binaries will be built into dist folder.
