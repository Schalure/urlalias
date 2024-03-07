// statickint is used to statically analyze the code of your application
//
// How to use:
// In windows: from command line to execut command from root directory of your solution
//
//	.\<path to staticlint.exe> ./...
//
// In linux: from terminal to execut command from root directory of your solution
//
//	.\<path to staticlint> ./...
//
// Description of the analyzers used:
// 1. Appends analyzer - Package appends defines an Analyzer that detects if there is only one variable in append.
// 2. Assign analyzer - Package assign defines an Analyzer that detects useless assignments.
// 3. Bools analyzer - Package bools defines an Analyzer that detects common mistakes involving boolean operators.
// 4. Buildtag analyzer - Package buildtag defines an Analyzer that checks build tags.
// 5. Copylock analyzer - Package copylock defines an Analyzer that checks for locks erroneously passed by value.
// 6. Defers analyzer - Package defers defines an Analyzer that checks for common mistakes in defer statements.
// 7. Fieldalignment analyzer - Package fieldalignment defines an Analyzer that detects structs that would use less memory if their fields were sorted.
// 8. Nilfunc analyzer - Package nilfunc defines an Analyzer that checks for useless comparisons against nil.
//
// For more inforation about analyzers see: https://pkg.go.dev/golang.org/x/tools/go/analysis#section-sourcefiles
package main

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

// main function of statickint
func main() {

	checks := []*analysis.Analyzer{
		//	standart analyzers
		appends.Analyzer,
		assign.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		copylock.Analyzer,
		defers.Analyzer,
		fieldalignment.Analyzer,
		nilfunc.Analyzer,
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

// MainExitAnalyzer checks if there is a direct call to os.Exit in the main package
var MainExitAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "checks if there is a direct call to os.Exit in the main package",
	Run:  runOsExitCheck,
}

// runOsExitCheck is run function of MainExitAnalyzer
func runOsExitCheck(pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			if c, ok := node.(*ast.CallExpr); ok {
				if s, ok := c.Fun.(*ast.SelectorExpr); ok {
					var packName string
					pack, ok := s.X.(*ast.Ident)
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
