package gostructfmt_test

import (
	"testing"

	"github.com/goropikari/gostructfmt"
	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	t.Run("formats struct literals and is deterministic", func(t *testing.T) {
		// Arrange
		source := []byte("package example\n\ntype User struct { Name string; Age int }\n\nvar user = User{Name: \"Alice\", Age: 20}\n")

		// Act
		first, firstErr := gostructfmt.Format("example.go", source)
		second, secondErr := gostructfmt.Format("example.go", source)

		// Assert
		require.NoError(t, firstErr)
		require.NoError(t, secondErr)
		require.Equal(t, first, second)
		require.Equal(t, "package example\n\ntype User struct {\n\tName string\n\tAge  int\n}\n\nvar user = User{\n\tName: \"Alice\",\n\tAge:  20,\n}\n", string(first))

		third, thirdErr := gostructfmt.Format("example.go", first)
		require.NoError(t, thirdErr)
		require.Equal(t, first, third)
	})

	t.Run("returns no output for invalid source", func(t *testing.T) {
		// Arrange
		source := []byte("package example\nfunc broken(\n")

		// Act
		output, err := gostructfmt.Format("broken.go", source)

		// Assert
		require.Error(t, err)
		require.Nil(t, output)
		require.Contains(t, err.Error(), "broken.go")
	})
}
