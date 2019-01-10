package main

import (
	g "github.com/ysmood/gokit"
)

func main() {
	g.Guard([]string{"echo", "ok"}, nil, &g.GuardOptions{})
}
