// Package gostructfmt provides a golangci-lint module plugin for struct
// literal formatting checks.
package gostructfmt

import (
	"bytes"
	"fmt"
	"go/ast"

	"github.com/golangci/plugin-module-register/register"
	"github.com/goropikari/gostructfmt/internal/formatter"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("gostructfmt", New)
}

// Plugin implements the golangci-lint module plugin contract.
type Plugin struct{}

// New creates the gostructfmt linter plugin.
func New(settings any) (register.LinterPlugin, error) {
	return &Plugin{}, nil
}

// BuildAnalyzers returns the struct literal formatting analyzer.
func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "gostructfmt",
			Doc:  "check populated Go struct literals for gostructfmt formatting",
			Run:  run,
		},
	}, nil
}

// GetLoadMode reports that the analyzer only needs syntax information.
func (p *Plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.File(file.Pos()).Name()

		source, err := pass.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", filename, err)
		}

		parsed, err := formatter.Parse(filename, source)
		if err != nil {
			return nil, err
		}

		ast.Inspect(file, func(node ast.Node) bool {
			literal, ok := node.(*ast.CompositeLit)
			if !ok || len(literal.Elts) == 0 {
				return true
			}

			start := pass.Fset.Position(literal.Pos()).Offset

			endOffset := pass.Fset.Position(literal.End()).Offset
			if bytes.Contains(source[start:endOffset], []byte{'\n'}) {
				return true
			}

			position := pass.Fset.Position(literal.Pos())
			end := pass.Fset.Position(literal.End())

			formatted, formatErr := formatter.FormatStructLiteralsInLines(parsed, []formatter.LineRange{{
				Start: position.Line,
				End:   end.Line,
			}})
			if formatErr != nil {
				return true
			}

			if string(formatted) != string(source) {
				pass.Report(analysis.Diagnostic{
					Pos:      literal.Pos(),
					End:      literal.End(),
					Category: "gostructfmt",
					Message:  "struct literal is not formatted by gostructfmt",
				})
			}

			return true
		})
	}

	return nil, nil
}
