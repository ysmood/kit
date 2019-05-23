package example

import (
	g "github.com/ysmood/gokit"
)

func Example() {
	g.Guard("go", "run", "./server").ExecCtx(
		g.Exec().Prefix("server | @yellow"),
	).Do()
}
