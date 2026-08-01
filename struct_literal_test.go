package gostructfmt_test

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/goropikari/gostructfmt"
	"github.com/stretchr/testify/require"
)

func TestFormatStructLiterals(t *testing.T) {
	t.Run("expands named struct literal", func(t *testing.T) {
		// Arrange
		source := []byte("package example\ntype User struct { Name string; Age int }\nfunc build() User { return User{Name: \"Alice\", Age: 20} }\n")
		file, err := gostructfmt.Parse("example.go", source)
		require.NoError(t, err)

		// Act
		err = gostructfmt.FormatStructLiterals(file)
		output, printErr := file.Print()

		// Assert
		require.NoError(t, err)
		require.NoError(t, printErr)
		require.Equal(t, "package example\n\ntype User struct {\n\tName string\n\tAge  int\n}\n\nfunc build() User {\n\treturn User{\n\t\tName: \"Alice\",\n\t\tAge:  20,\n\t}\n}\n", string(output))
	})

	t.Run("expands anonymous qualified and generic struct literals", func(t *testing.T) {
		// Arrange
		source := []byte("package example\nimport \"other\"\ntype Box[T any] struct { Value T }\nfunc build() { _ = struct { Name string }{Name: \"x\"}; _ = other.User{Name: \"x\"}; _ = Box[string]{Value: \"x\"} }\n")
		file, err := gostructfmt.Parse("example.go", source)
		require.NoError(t, err)

		// Act
		err = gostructfmt.FormatStructLiterals(file)
		output, printErr := file.Print()

		// Assert
		require.NoError(t, err)
		require.NoError(t, printErr)
		require.Contains(t, string(output), "struct{ Name string }{\n\t\tName: \"x\",\n\t}")
		require.Contains(t, string(output), "other.User{\n\t\tName: \"x\",\n\t}")
		require.Contains(t, string(output), "Box[string]{\n\t\tValue: \"x\",\n\t}")
	})

	t.Run("keeps empty and non-struct literals unchanged", func(t *testing.T) {
		// Arrange
		source := []byte("package example\ntype User struct { Name string }\ntype IDs []int\ntype Table map[string]int\ntype Alias User\nvar _ = User{}\nvar _ = IDs{1, 2}\nvar _ = Table{\"x\": 1}\nvar _ = Alias{Name: \"Alice\"}\nvar _ = []int{1, 2}\nvar _ = [2]int{1, 2}\nvar _ = map[string]int{\"x\": 1}\n")
		file, err := gostructfmt.Parse("example.go", source)
		require.NoError(t, err)

		// Act
		err = gostructfmt.FormatStructLiterals(file)
		output, printErr := file.Print()

		// Assert
		require.NoError(t, err)
		require.NoError(t, printErr)
		require.Contains(t, string(output), "var _ = User{}")
		require.Contains(t, string(output), "var _ = IDs{1, 2}")
		require.Contains(t, string(output), "var _ = Table{\"x\": 1}")
		require.Contains(t, string(output), "var _ = Alias{\n\tName: \"Alice\",\n}")
		require.Contains(t, string(output), "var _ = []int{1, 2}")
		require.Contains(t, string(output), "var _ = [2]int{1, 2}")
		require.Contains(t, string(output), "var _ = map[string]int{\"x\": 1}")
	})

	t.Run("preserves direct AST edits before formatting", func(t *testing.T) {
		// Arrange
		file, err := gostructfmt.Parse("example.go", []byte("package example\nvar _ = User{Name: \"Alice\"}\n"))
		require.NoError(t, err)

		file.AST().Name.Name = "renamed"

		// Act
		err = gostructfmt.FormatStructLiterals(file)
		output, printErr := file.Print()

		// Assert
		require.NoError(t, err)
		require.NoError(t, printErr)
		require.Contains(t, string(output), "package renamed")
		require.Contains(t, string(output), "User{\n")
	})

	t.Run("expands nested literals and preserves comments", func(t *testing.T) {
		// Arrange
		source := []byte("package example\ntype Inner struct { Name string }\ntype Outer struct { Inner Inner }\nfunc build() Outer { return Outer{ // outer\nInner: Inner{Name: \"x\"}, // inner\n} }\n")
		file, err := gostructfmt.Parse("example.go", source)
		require.NoError(t, err)

		// Act
		err = gostructfmt.FormatStructLiterals(file)
		output, printErr := file.Print()

		// Assert
		require.NoError(t, err)
		require.NoError(t, printErr)
		require.Equal(t, "package example\n\ntype Inner struct{ Name string }\ntype Outer struct{ Inner Inner }\n\nfunc build() Outer {\n\treturn Outer{ // outer\n\t\tInner: Inner{\n\t\t\tName: \"x\",\n\t\t}, // inner\n\t}\n}\n", string(output))
	})

	t.Run("keeps already multiline literals valid and reparsable", func(t *testing.T) {
		// Arrange
		source := []byte("package example\ntype User struct { Name string }\nvar _ = User{\nName: \"Alice\",\n}\n")
		file, err := gostructfmt.Parse("example.go", source)
		require.NoError(t, err)

		// Act
		err = gostructfmt.FormatStructLiterals(file)
		output, printErr := file.Print()
		_, parseErr := parser.ParseFile(token.NewFileSet(), "formatted.go", output, parser.ParseComments)

		// Assert
		require.NoError(t, err)
		require.NoError(t, printErr)
		require.NoError(t, parseErr)
		require.Equal(t, "package example\n\ntype User struct{ Name string }\n\nvar _ = User{\n\tName: \"Alice\",\n}\n", string(output))
	})
}

func TestTransformStructLiterals(t *testing.T) {
	t.Run("delegates to the struct literal formatter", func(t *testing.T) {
		// Arrange
		file, err := gostructfmt.Parse("example.go", []byte("package example\ntype User struct { Name string }\nvar _ = User{Name: \"Alice\"}\n"))
		require.NoError(t, err)

		// Act
		err = gostructfmt.TransformStructLiterals(file)
		output, printErr := file.Print()

		// Assert
		require.NoError(t, err)
		require.NoError(t, printErr)
		require.Contains(t, string(output), "User{\n")
	})
}
