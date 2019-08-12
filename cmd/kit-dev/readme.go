package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	gos "github.com/ysmood/kit/pkg/os"
	"github.com/ysmood/kit/pkg/run"
	"github.com/ysmood/kit/pkg/utils"
)

func genReadme() {
	fexmaple := utils.E1(gos.ReadString("kit_test.go")).(string)

	fset := token.NewFileSet()
	fast := utils.E1(parser.ParseFile(fset, "kit_test.go", fexmaple, parser.ParseComments)).(*ast.File)

	guardHelp := run.Exec("go", "run", "./cmd/guard", "--help").MustString()
	godevHelp := run.Exec("go", "run", "./cmd/godev", "--help").MustString()

	list := []interface{}{
		"GuardHelp", guardHelp,
		"GodevHelp", godevHelp,
	}

	ast.Inspect(fast, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			name := x.Name.Name
			code := fexmaple[x.Body.Pos()-1 : x.Body.End()]

			if strings.HasPrefix(name, "Example") {
				list = append(list, name, formatExample(code))
			}
		}
		return true
	})

	f := utils.E1(gos.ReadString("readme.tpl.md")).(string)

	utils.E(gos.OutputFile("readme.md", utils.S(f, list...), nil))
}

func formatExample(code string) string {
	return utils.S(
		strings.Join(
			[]string{
				"```go",
				"package main",
				"",
				"import . \"github.com/ysmood/kit\"",
				"",
				"func main() {{.code}}",
				"```",
			},
			"\n",
		),
		"code", code,
	)
}
