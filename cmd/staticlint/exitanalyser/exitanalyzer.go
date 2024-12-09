// Package exitanalyser - checks existence of "os.Exit" in main function of main file.
package exitanalyser

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer - define the analyzer.
var Analyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for exit() calls in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	// Iterate over all files in pkg.
	for _, file := range pass.Files {
		// Check if the file name is main.
		if !strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "main.go") {
			continue
		}
		// Inspect nodes to find main function.
		ast.Inspect(file, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "main" {
				return true // Continue inspection.
			}
			// Inspect main function and its nodes.
			ast.Inspect(file, func(n ast.Node) bool {
				// Check every call expression.
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				// Devide call expression to parts, acording schema (os.Exit X - os, Sel - Exit)
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				// Check X part.
				pkg, ok := sel.X.(*ast.Ident)
				if !ok || pkg.Name != "os" {
					return true
				}
				// Check Sel part.
				if sel.Sel.Name != "Exit" {
					return true
				}
				pass.Reportf(sel.Sel.Pos(), "avoid using 'os.Exit' in main function")

				return true
			})

			return false
		})

	}
	return nil, nil
}
