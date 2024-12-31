// Command staticlint is a multichecker for performing static analysis.
//
// It consists of several types of analyzers:
// - standard analysis passes (e.g., `appends`, `assign`, `shadow`).
// - Staticcheck analyzers (SA checks).
// - Stylecheck analyzers (ST checks).
// - Custom analyzers exitanalyser.Analyzer.
//
// Usage:
//
//	First compile code to binary.
//	Then run the compiled binary against Go code to detect issues and problems.
package main

import (
	"fmt"

	"github.com/madcarpet/metrics/cmd/staticlint/exitanalyser"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"honnef.co/go/tools/analysis/facts/nilness"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"honnef.co/go/tools/unused"
)

func main() {
	// listAnalyzers defines custom and Go's standard analyzers.
	listAnalyzers := []*analysis.Analyzer{
		appends.Analyzer,         // Detects suspicious append patterns.
		assign.Analyzer,          // Checks for redundant variable assignments.
		atomic.Analyzer,          // Identifies incorrect use of sync/atomic.
		bools.Analyzer,           // Warns about suspicious boolean logic.
		copylock.Analyzer,        // Reports copying of locks.
		defers.Analyzer,          // Detects deferred calls in loops.
		httpresponse.Analyzer,    // Identifies unclosed HTTP response bodies.
		printf.Analyzer,          // Warns about incorrect use of fmt.Printf-like functions.
		shadow.Analyzer,          // Detects variable shadowing.
		shift.Analyzer,           // Identifies shifts that exceed the bit width of the variable.
		unmarshal.Analyzer,       // Reports issues in unmarshaling JSON/XML data.
		unreachable.Analyzer,     // Detects unreachable code.
		exitanalyser.Analyzer,    // Custom analyzer for handling program exits.
		nilness.Analysis,         // Detects potential nil pointer dereferences.
		unused.Analyzer.Analyzer, // Reports unused variables, functions, etc. (U1000).
	}
	mychecks := make([]*analysis.Analyzer, 0, len(staticcheck.Analyzers)+len(stylecheck.Analyzers)+len(listAnalyzers))
	// Add all analyzers from Staticcheck.
	for _, c := range staticcheck.Analyzers {
		mychecks = append(mychecks, c.Analyzer)

	}
	// Add all analyzers from Stylecheck.
	for _, c := range stylecheck.Analyzers {
		mychecks = append(mychecks, c.Analyzer)
	}
	// Add the custom list of analyzers.
	mychecks = append(mychecks, listAnalyzers...)
	// Print to stdout total number of analyzers.
	fmt.Println("Analyzers count:", len(mychecks))
	// Run the multichecker.
	multichecker.Main(
		mychecks...,
	)
}
