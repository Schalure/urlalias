package main

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {

	checks := []*analysis.Analyzer {
		//	standart analyzers
		printf.Analyzer,
        shadow.Analyzer,
        shift.Analyzer,
        structtag.Analyzer,	

		//	custom analyzers
		MainExitAnalyzer,
	}

	//	staticcheck analyzers
	for _, a := range staticcheck.Analyzers {
		if strings.Contains(a.Analyzer.Name, "SA") || strings.Contains(a.Analyzer.Name, "ST") {
			checks = append(checks, a.Analyzer)
		}
	}

	multichecker.Main(checks...)
}

var MainExitAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc: "checks if there is a direct call to os.Exit in the main package",
	Run: runOsExitCheck,
}

func runOsExitCheck (pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			if c, ok := node.(*ast.CallExpr); ok {
				if s, ok := c.Fun.(*ast.SelectorExpr); ok {
					var packName string
					pack, ok := s.X.(*ast.Ident); 
					if ok {
						packName = pack.Name
					}
					if packName == "os" && s.Sel.Name == "Exit" {
						fmt.Printf("%v: Direct call to os.Exit from the main package\n", pass.Fset.Position(pack.Pos()).String())
					}
				}
			}
			return true
		})
	}
	return nil, nil
}