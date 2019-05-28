# Gokit

[![Build Status](https://travis-ci.org/ysmood/gokit.svg?branch=master)](https://travis-ci.org/ysmood/gokit)
[![codecov](https://codecov.io/gh/ysmood/gokit/branch/master/graph/badge.svg)](https://codecov.io/gh/ysmood/gokit)
[![goreport](https://goreportcard.com/badge/github.com/ysmood/gokit)](https://goreportcard.com/report/github.com/ysmood/gokit)

This project is a collection of often used io related methods.

- Sane defaults to make coding less verbose.

- Focus on performance.

- Won't produce error, all errors are come from its dependencies.

- Won't use any other lanuage for the development of this project. So it's a good example to show how to use Go only for automation.

- 100% test coverage and goreportcard

## Modules

### os

Covers most used os related functions that are missing from the stdlib.

Such as smart glob:

{{.ExampleWalk}}

### process

A better `Exec` alternatives for the stdlib one.

{{.ExampleExec}}

### req

The http request lib from stdlib is pretty verbose to use. The `gokit.Req` is a much better
alternative to use with it's fluent api design. You will reduce a lot of your code without sacrificing performance.
It covers all the functions of the Go's stdlib one, no api is hidden from the origin http lib.

{{.ExampleReq}}

### guard

{{.ExampleGuard}}

#### CLI tool

## Install CLI Tools

Goto the [release page](https://github.com/ysmood/gokit/releases) download the binary for your OS.

```bash
{{.GuardHelp}}
```

## Development

To write testable code, I try to isolate all error related dependencies.

### Build Project

```bash
go run ./dev --help
go run ./dev build
```

Binaries will be built into dist folder.