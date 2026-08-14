package gotreesj_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/goropikari/gotreesj/plugin/gotreesj"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
)

func TestRun(t *testing.T) {
	t.Run("reports a populated single-line struct literal", func(t *testing.T) {
		// Arrange
		source := []byte("package example\ntype User struct { Name string }\nvar _ = User{Name: \"Alice\"}\n")

		// Act
		reports := runAnalyzer(t, source)

		// Assert
		require.Len(t, reports, 1)
		require.Equal(t, "gotreesj", reports[0].Category)
	})

	t.Run("reports a long function call with single-line arguments", func(t *testing.T) {
		// Arrange
		source := []byte("package example\nvar _ = configureUser(\"Alice\", \"Administrator\", \"Tokyo\", \"engineering\", true, 42, \"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\")\n")

		// Act
		reports := runAnalyzer(t, source)

		// Assert
		require.Len(t, reports, 1)
		require.Equal(t, "function call arguments are not formatted by gotreesj", reports[0].Message)
	})

	t.Run("accepts a formatted multi-line struct literal", func(t *testing.T) {
		// Arrange
		source := []byte("package example\ntype User struct { Name string }\nvar _ = User{\n\tName: \"Alice\",\n}\n")

		// Act
		reports := runAnalyzer(t, source)

		// Assert
		require.Empty(t, reports)
	})
}

func runAnalyzer(t *testing.T, source []byte) []analysis.Diagnostic {
	t.Helper()

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

	plugin, pluginErr := gotreesj.New(nil)
	require.NoError(t, pluginErr)

	analyzers, analyzerErr := plugin.BuildAnalyzers()
	require.NoError(t, analyzerErr)

	_, err = analyzers[0].Run(pass)
	require.NoError(t, err)

	return reports
}
