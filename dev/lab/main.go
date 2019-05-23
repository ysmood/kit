package main

import (
	. "github.com/ysmood/gokit"
)

func main() {
	Guard("echo", "ok").Do()
}
