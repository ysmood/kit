package main

import (
	"fmt"
	"os"
	"strings"

	gos "github.com/ysmood/gokit/pkg/os"
	"github.com/ysmood/gokit/pkg/run"
	"github.com/ysmood/gokit/pkg/utils"
	"golang.org/x/tools/go/packages"
)

// Export all the public members of each package under pkg folfer into gokit_exports.go
func export() {
	if !gos.FileExists("kit_test.go") {
		return
	}

	paths := gos.Walk("pkg/*").MustList()

	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedDeps | packages.NeedTypes}
	pkgs, err := packages.Load(cfg, paths...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\n", err)
		os.Exit(1)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	header := ""
	statements := ""

	for _, pkg := range pkgs {
		header += fmt.Sprintf("    \"github.com/ysmood/gokit/pkg/%s\"\n", pkg.Name)

		s := pkg.Types.Scope()
		for _, n := range s.Names() {
			v := s.Lookup(n)
			if v.Exported() {
				statements += "// " + n + " imported\n"
				if strings.HasPrefix(v.String(), "type") {
					statements += fmt.Sprintf("type %s = %s.%s\n", n, pkg.Name, n)
				} else {
					statements += fmt.Sprintf("var %s = %s.%s\n", n, pkg.Name, n)
				}
			}
		}
	}
	header = "package kit\n\nimport(\n" + header + ")\n\n"

	utils.E(gos.OutputFile("kit.go", header+statements, nil))

	run.Exec("gofmt", "-w", "kit.go").MustDo()
}
