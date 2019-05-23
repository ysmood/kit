package main

import (
	. "github.com/ysmood/gokit/pkg/guard"
)

func main() {
	_ = Guard("echo", "ok").Do()
}
