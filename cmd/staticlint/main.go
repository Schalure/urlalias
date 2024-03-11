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
// 1. appends.Analyzer - Package appends defines an Analyzer that detects if there is only one variable in append.
// 2. assign.Analyzer - Package assign defines an Analyzer that detects useless assignments.
// 3. bools.Analyzer - Package bools defines an Analyzer that detects common mistakes involving boolean operators.
// 4. buildtag.Analyzer - Package buildtag defines an Analyzer that checks build tags.
// 5. copylock.Analyzer - Package copylock defines an Analyzer that checks for locks erroneously passed by value.
// 6. defers.Analyzer - Package defers defines an Analyzer that checks for common mistakes in defer statements.
// 7. fieldalignment.Analyzer - Package fieldalignment defines an Analyzer that detects structs that would use less memory if their fields were sorted.
// 8. nilfunc.Analyzer - Package nilfunc defines an Analyzer that checks for useless comparisons against nil.
// 9. printf.Analyzer - Package printf defines an Analyzer that checks consistency of Printf format strings and arguments
// 10. shadow.Analyzer - Package shadow defines an Analyzer that checks for shadowed variables.
// 11. shift.Analyzer - Package shift defines an Analyzer that checks for shifts that exceed the width of an integer
// 12. slog.Analyzer - Package slog defines an Analyzer that checks for mismatched key-value pairs in log/slog calls
// 13. sortslice.Analyzer - Package sortslice defines an Analyzer that checks for calls to sort.Slice that do not use a slice type as first argument
// 14. structtag.Analyzer - Package structtag defines an Analyzer that checks struct field tags are well formed
// 15. tests.Analyzer - Package tests defines an Analyzer that checks for common mistaken usages of tests and examples
//
// 16. MainExitAnalyzer - MainExitAnalyzer checks if there is a direct call to os.Exit in the main package
//
// 17. statickint uses all staticcheck analizers SA and ST groups
//
// For more inforation about analyzers see:
// https://pkg.go.dev/golang.org/x/tools/go/analysis#section-sourcefiles
// https://staticcheck.dev/docs/checks/
package main

import (
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
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
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
		slog.Analyzer,
		sortslice.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,

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

	//	search main package
	var file *ast.File
	for _, file = range pass.Files {
		if file.Name.Name != "main" {
			continue
		}
		break
	}

	//	search main function
	var mainBody ast.Node
	ast.Inspect(file, func(node ast.Node) bool {
		if f, ok := node.(*ast.FuncDecl); ok {
			if f.Name.Name != "main" {
				return true
			}
			mainBody = f.Body
			return false
		}
		return true
	})

	if mainBody == nil {
		return nil, nil
	}

	//	search os.Exit in main function
	ast.Inspect(mainBody, func(node ast.Node) bool {

		if c, ok := node.(*ast.CallExpr); ok {
			s, ok := c.Fun.(*ast.SelectorExpr)
			if ! ok {
				return true
			}

			iX, ok := s.X.(*ast.Ident)
			if ! ok {
				return true
			}

			if iX.Name == "os" && s.Sel.Name == "Exit" {
				pass.Reportf(iX.Pos(), "Direct call to os.Exit from the main function\n")
			}
		}	
		return true				
	})
	return nil, nil
}
