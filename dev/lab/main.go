package main

import (
	"strings"

	g "github.com/ysmood/gokit"
)

func main() {
	list, _ := g.Glob(g.GuardDefaultPatterns, nil)

	g.Log(strings.Join(list, "\n"))
}
