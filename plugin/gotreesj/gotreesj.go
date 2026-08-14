// Package gotreesj provides a golangci-lint module plugin for struct literal
// and long function call formatting checks.
package gotreesj

import (
	"bytes"
	"fmt"
	"go/ast"

	"github.com/golangci/plugin-module-register/register"
	"github.com/goropikari/gotreesj/internal/formatter"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("gotreesj", New)
}

// Plugin implements the golangci-lint module plugin contract.
type Plugin struct{}

// New creates the gotreesj linter plugin.
func New(settings any) (register.LinterPlugin, error) {
	return &Plugin{}, nil
}

// BuildAnalyzers returns the gotreesj formatting analyzer.
func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "gotreesj",
			Doc:  "check populated Go struct literals and long function calls for gotreesj formatting",
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
			reportFormattingIssue(pass, parsed, source, node)
			return true
		})
	}

	return nil, nil
}

func reportFormattingIssue(pass *analysis.Pass, parsed *formatter.File, source []byte, node ast.Node) {
	message, ok := formattingMessage(node)
	if !ok {
		return
	}

	start := pass.Fset.Position(node.Pos()).Offset

	end := pass.Fset.Position(node.End()).Offset
	if bytes.Contains(source[start:end], []byte{'\n'}) {
		return
	}

	position := pass.Fset.Position(node.Pos())
	last := pass.Fset.Position(node.End())

	formatted, err := formatter.FormatStructLiteralsInLines(parsed, []formatter.LineRange{{
		Start: position.Line,
		End:   last.Line,
	}})
	if err != nil || string(formatted) == string(source) {
		return
	}

	pass.Report(analysis.Diagnostic{
		Pos:      node.Pos(),
		End:      node.End(),
		Category: "gotreesj",
		Message:  message,
	})
}

func formattingMessage(node ast.Node) (string, bool) {
	switch node := node.(type) {
	case *ast.CompositeLit:
		if len(node.Elts) == 0 {
			return "", false
		}

		return "struct literal is not formatted by gotreesj", true
	case *ast.CallExpr:
		if len(node.Args) < 2 {
			return "", false
		}

		return "function call arguments are not formatted by gotreesj", true
	default:
		return "", false
	}
}
