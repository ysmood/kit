package main

import (
	. "github.com/ysmood/gokit"
)

func main() {
	_ = Guard("echo", "ok").Do()
}
