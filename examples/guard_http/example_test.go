package example

import (
	g "github.com/ysmood/gokit"
)

func Example() {
	g.Guard("go", "run", "./server").Context(
		g.GuardContext{
			ExecOpts: g.ExecOptions{
				Prefix: "server | @yellow",
			},
		},
	).Do()
}
