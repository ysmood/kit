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

```bash
{{.GuardHelp}}
```

## Development

### Build Project

```bash
go run ./dev --help
go run ./dev build
```

Binaries will be built into dist folder.