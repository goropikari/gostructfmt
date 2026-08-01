package gostructfmt_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/goropikari/gostructfmt/plugin/gostructfmt"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
)

func TestRun(t *testing.T) {
	t.Run("reports a populated single-line struct literal", func(t *testing.T) {
		// Arrange
		source := []byte("package example\ntype User struct { Name string }\nvar _ = User{Name: \"Alice\"}\n")
		fileSet := token.NewFileSet()
		file, err := parser.ParseFile(fileSet, "example.go", source, parser.ParseComments)
		require.NoError(t, err)

		reports := make([]analysis.Diagnostic, 0)
		pass := &analysis.Pass{
			Fset:  fileSet,
			Files: []*ast.File{file},
			ReadFile: func(string) ([]byte, error) {
				return source, nil
			},
			Report: func(diagnostic analysis.Diagnostic) {
				reports = append(reports, diagnostic)
			},
		}

		// Act
		plugin, pluginErr := gostructfmt.New(nil)
		require.NoError(t, pluginErr)

		analyzers, analyzerErr := plugin.BuildAnalyzers()
		require.NoError(t, analyzerErr)

		_, err = analyzers[0].Run(pass)

		// Assert
		require.NoError(t, err)
		require.Len(t, reports, 1)
		require.Equal(t, "gostructfmt", reports[0].Category)
	})

	t.Run("accepts a formatted multi-line struct literal", func(t *testing.T) {
		// Arrange
		source := []byte("package example\ntype User struct { Name string }\nvar _ = User{\n\tName: \"Alice\",\n}\n")
		fileSet := token.NewFileSet()
		file, err := parser.ParseFile(fileSet, "example.go", source, parser.ParseComments)
		require.NoError(t, err)

		reports := make([]analysis.Diagnostic, 0)
		pass := &analysis.Pass{
			Fset:  fileSet,
			Files: []*ast.File{file},
			ReadFile: func(string) ([]byte, error) {
				return source, nil
			},
			Report: func(diagnostic analysis.Diagnostic) {
				reports = append(reports, diagnostic)
			},
		}

		// Act
		plugin, pluginErr := gostructfmt.New(nil)
		require.NoError(t, pluginErr)

		analyzers, analyzerErr := plugin.BuildAnalyzers()
		require.NoError(t, analyzerErr)

		_, err = analyzers[0].Run(pass)

		// Assert
		require.NoError(t, err)
		require.Empty(t, reports)
	})
}
