package main

import (
	g "github.com/ysmood/gokit"
)

func main() {
	g.E(
		g.Guard(
			[]string{"go", "run", "./server"},
			nil,
			&g.GuardOptions{
				ExecOpts: &g.ExecOptions{
					Prefix: "server | @yellow",
				},
			},
		),
	)
}
