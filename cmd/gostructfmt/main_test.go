package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("formats standard input to standard output", func(t *testing.T) {
		// Arrange
		input := bytes.NewBufferString("package example\ntype User struct { Name string }\nvar user = User{Name: \"Alice\"}\n")

		var (
			output bytes.Buffer
			errors bytes.Buffer
		)

		// Act
		status := run(nil, input, &output, &errors)

		// Assert
		require.Equal(t, 0, status)
		require.Empty(t, errors.String())
		require.Contains(t, output.String(), "User{\n")
	})

	t.Run("reports invalid standard input and emits no output", func(t *testing.T) {
		// Arrange
		input := bytes.NewBufferString("package example\nfunc broken(\n")

		var (
			output bytes.Buffer
			errors bytes.Buffer
		)

		// Act
		status := run(nil, input, &output, &errors)

		// Assert
		require.Equal(t, 1, status)
		require.Empty(t, output.String())
		require.Contains(t, errors.String(), "<stdin>")
	})

	t.Run("prints help without treating it as a formatting error", func(t *testing.T) {
		// Arrange
		var (
			input  bytes.Buffer
			output bytes.Buffer
			errors bytes.Buffer
		)

		// Act
		status := run([]string{"-h"}, &input, &output, &errors)

		// Assert
		require.Equal(t, 0, status)
		require.Contains(t, errors.String(), "Usage:")
	})

	t.Run("writes one file atomically when requested", func(t *testing.T) {
		// Arrange
		directory := t.TempDir()
		filename := filepath.Join(directory, "input.go")
		require.NoError(t, os.WriteFile(filename, []byte("package example\ntype User struct { Name string }\nvar _ = User{Name: \"Alice\"}\n"), 0o600))

		var (
			output bytes.Buffer
			errors bytes.Buffer
		)

		// Act
		status := run([]string{"-w", filename}, bytes.NewBuffer(nil), &output, &errors)
		formatted, readErr := os.ReadFile(filename)

		// Assert
		require.Equal(t, 0, status)
		require.NoError(t, readErr)
		require.Empty(t, output.String())
		require.Empty(t, errors.String())
		require.Contains(t, string(formatted), "User{\n")
	})

	t.Run("rejects write mode without a file", func(t *testing.T) {
		// Arrange
		var (
			output bytes.Buffer
			errors bytes.Buffer
		)

		// Act
		status := run([]string{"-w"}, bytes.NewBufferString("package example\n"), &output, &errors)

		// Assert
		require.Equal(t, 2, status)
		require.Empty(t, output.String())
		require.Contains(t, errors.String(), "requires at least one file")
	})

	t.Run("rejects multiple files without write mode", func(t *testing.T) {
		// Arrange
		var (
			output bytes.Buffer
			errors bytes.Buffer
		)

		// Act
		status := run([]string{"one.go", "two.go"}, bytes.NewBuffer(nil), &output, &errors)

		// Assert
		require.Equal(t, 2, status)
		require.Empty(t, output.String())
		require.Contains(t, errors.String(), "multiple files require -w")
	})

	t.Run("formats every Go file below the current directory for ./...", func(t *testing.T) {
		// Arrange
		directory := t.TempDir()
		require.NoError(t, os.Mkdir(filepath.Join(directory, "nested"), 0o750))
		require.NoError(t, os.WriteFile(filepath.Join(directory, "main.go"), []byte("package example\ntype User struct { Name string }\nvar _ = User{Name: \"Alice\"}\n"), 0o600))
		nestedFilename := filepath.Join(directory, "nested", "nested.go")
		require.NoError(t, os.WriteFile(nestedFilename, []byte("package nested\ntype User struct { Name string }\nvar _ = User{Name: \"Bob\"}\n"), 0o600))
		t.Chdir(directory)

		var (
			output bytes.Buffer
			errors bytes.Buffer
		)

		// Act
		status := run([]string{"-w", "./..."}, bytes.NewBuffer(nil), &output, &errors)
		mainSource, mainErr := os.ReadFile(filepath.Join(directory, "main.go"))
		nestedSource, nestedErr := os.ReadFile(nestedFilename)

		// Assert
		require.Equal(t, 0, status)
		require.NoError(t, mainErr)
		require.NoError(t, nestedErr)
		require.Empty(t, output.String())
		require.Empty(t, errors.String())
		require.Contains(t, string(mainSource), "User{\n")
		require.Contains(t, string(nestedSource), "User{\n")
	})
}
