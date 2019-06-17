package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	gos "github.com/ysmood/gokit/pkg/os"
	"github.com/ysmood/gokit/pkg/run"
	"github.com/ysmood/gokit/pkg/utils"
)

func genReadme() {
	if !gos.FileExists("readme.tpl.md") {
		return
	}

	fexmaple := utils.E1(gos.ReadString("kit_test.go")).(string)

	fset := token.NewFileSet()
	fast := utils.E1(parser.ParseFile(fset, "kit_test.go", fexmaple, parser.ParseComments)).(*ast.File)

	guardHelp := run.Exec("go", "run", "./cmd/guard", "--help").MustString()

	list := []interface{}{
		"GuardHelp", guardHelp,
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
				"import . \"github.com/ysmood/gokit\"",
				"",
				"func main() {{.code}}",
				"```",
			},
			"\n",
		),
		"code", code,
	)
}
