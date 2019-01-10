package main

import (
	g "github.com/ysmood/gokit"
)

func main() {
	g.Exec([]string{"ls", "-aG"}, &g.ExecOptions{
		Prefix: "test | ",
	})
}
