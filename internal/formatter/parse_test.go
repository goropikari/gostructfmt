package formatter_test

import (
	"errors"
	"go/ast"
	"go/token"
	"testing"

	"github.com/goropikari/gotreesj/internal/formatter"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Run("parses Go source and retains comments", func(t *testing.T) {
		// Arrange
		source := []byte("package example\n\n// User comment\ntype User struct { Name string }\n")

		// Act
		file, err := formatter.Parse("user.go", source)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, file)
		require.Equal(t, "user.go", file.Filename())
		require.NotNil(t, file.AST())
		require.NotNil(t, file.FileSet())
		require.Len(t, file.AST().Comments, 1)
		require.Equal(t, "User comment\n", file.AST().Comments[0].Text())
	})

	t.Run("returns an error and no file for invalid source", func(t *testing.T) {
		// Arrange
		source := []byte("package example\nfunc broken(\n")

		// Act
		file, err := formatter.Parse("broken.go", source)

		// Assert
		require.Error(t, err)
		require.Nil(t, file)
		require.ErrorContains(t, err, `parse Go source "broken.go"`)
	})

	t.Run("uses a stable default filename", func(t *testing.T) {
		// Arrange
		source := []byte("package example\n")

		// Act
		file, err := formatter.Parse("", source)

		// Assert
		require.NoError(t, err)
		require.Equal(t, "<input>", file.Filename())
	})
}

func TestFilePrint(t *testing.T) {
	t.Run("prints formatted source with comments", func(t *testing.T) {
		// Arrange
		source := []byte("package example\n// comment\nvar value=1\n")
		file, err := formatter.Parse("example.go", source)
		require.NoError(t, err)

		// Act
		output, err := file.Print()

		// Assert
		require.NoError(t, err)
		require.Equal(t, "package example\n\n// comment\nvar value = 1\n", string(output))
	})

	t.Run("prints imports functions and nested struct literals", func(t *testing.T) {
		// Arrange
		source := []byte("package example\nimport \"fmt\"\ntype Config struct { Name string }\nfunc build() Config { return Config{Name: fmt.Sprint(\"x\")} }\n")
		file, err := formatter.Parse("example.go", source)
		require.NoError(t, err)

		// Act
		output, err := file.Print()

		// Assert
		require.NoError(t, err)
		require.Equal(t, "package example\n\nimport \"fmt\"\n\ntype Config struct{ Name string }\n\nfunc build() Config { return Config{Name: fmt.Sprint(\"x\")} }\n", string(output))
	})

	t.Run("prints after a reachable AST mutation using the original file set", func(t *testing.T) {
		// Arrange
		file, err := formatter.Parse("example.go", []byte("package example\n"))
		require.NoError(t, err)
		require.NoError(t, file.MutateAST(func(parsed *ast.File) error {
			parsed.Decls = []ast.Decl{&ast.GenDecl{
				Tok: token.VAR,
			}}

			return nil
		}))

		// Act
		output, err := file.Print()

		// Assert
		require.NoError(t, err)
		require.Equal(t, "package example\n\nvar ()\n", string(output))
	})

	t.Run("rejects nil receiver", func(t *testing.T) {
		// Arrange
		var file *formatter.File

		// Act
		output, err := file.Print()

		// Assert
		require.Error(t, err)
		require.Nil(t, output)
		require.ErrorContains(t, err, "nil file")
	})
}

func TestParseAndPrint(t *testing.T) {
	t.Run("parses and prints imports functions and nested syntax", func(t *testing.T) {
		// Arrange
		source := []byte("package example\nimport \"fmt\"\ntype Config struct { DB struct { Host string } }\nfunc build() Config { return Config{DB: struct { Host string }{Host: fmt.Sprint(\"localhost\")}} }\n")

		// Act
		output, err := formatter.ParseAndPrint("example.go", source)

		// Assert
		require.NoError(t, err)

		expected := "package example\n\nimport \"fmt\"\n\ntype Config struct{ DB struct{ Host string } }\n\nfunc build() Config { return Config{DB: struct{ Host string }{Host: fmt.Sprint(\"localhost\")}} }\n"
		require.Equal(t, expected, string(output))
	})

	t.Run("returns no output for invalid source", func(t *testing.T) {
		// Arrange
		source := []byte("package example\nvar value =\n")

		// Act
		output, err := formatter.ParseAndPrint("broken.go", source)

		// Assert
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("preserves doc block trailing and multiple comments", func(t *testing.T) {
		// Arrange
		source := []byte("package example\n\n// doc comment\n// second doc comment\ntype User struct {\n\tName string /* field block */\n}\n\nfunc build() User {\n\t// before literal\n\treturn User{ // trailing literal\n\t\tName: \"Alice\", // field trailing\n\t}\n}\n")

		// Act
		output, err := formatter.ParseAndPrint("comments.go", source)

		// Assert
		require.NoError(t, err)
		require.Contains(t, string(output), "doc comment")
		require.Contains(t, string(output), "second doc comment")
		require.Contains(t, string(output), "field block")
		require.Contains(t, string(output), "trailing literal")
		require.Contains(t, string(output), "field trailing")
	})
}

func TestFileAccessorsAndMutation(t *testing.T) {
	t.Run("keeps file metadata stable while allowing AST mutation", func(t *testing.T) {
		// Arrange
		file, err := formatter.Parse("example.go", []byte("package example\nvar value = 1\n"))
		require.NoError(t, err)

		originalFileSet := file.FileSet()
		originalFilename := file.Filename()

		// Act
		err = file.MutateAST(func(parsed *ast.File) error {
			parsed.Decls = append(parsed.Decls, &ast.GenDecl{
				Tok: token.VAR,
			})

			return nil
		})

		// Assert
		require.NoError(t, err)
		require.Same(t, originalFileSet, file.FileSet())
		require.Equal(t, originalFilename, file.Filename())
		require.Len(t, file.AST().Decls, 2)
	})

	t.Run("rejects a nil AST edit function", func(t *testing.T) {
		// Arrange
		file, err := formatter.Parse("example.go", []byte("package example\n"))
		require.NoError(t, err)

		// Act
		err = file.MutateAST(nil)

		// Assert
		require.Error(t, err)
		require.ErrorContains(t, err, "nil edit function")
	})

	t.Run("propagates an AST edit error without printing", func(t *testing.T) {
		// Arrange
		file, err := formatter.Parse("example.go", []byte("package example\n"))
		require.NoError(t, err)

		expectedErr := errors.New("edit rejected")

		// Act
		err = file.MutateAST(func(*ast.File) error {
			return expectedErr
		})

		// Assert
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("nil accessors are safe", func(t *testing.T) {
		// Arrange
		var file *formatter.File

		// Act
		astFile := file.AST()
		fileSet := file.FileSet()
		filename := file.Filename()

		// Assert
		require.Nil(t, astFile)
		require.Nil(t, fileSet)
		require.Empty(t, filename)
	})
}
