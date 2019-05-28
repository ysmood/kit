package main

import (
	kit "github.com/ysmood/gokit"
)

func main() {
	kit.Guard("echo", "ok").MustDo()
}
