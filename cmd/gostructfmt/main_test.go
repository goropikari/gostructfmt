package main

import (
	"bytes"
	"os"
	"os/exec"
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
		require.Contains(t, errors.String(), "Format populated Go struct literals")
		require.Contains(t, errors.String(), "gostructfmt -w ./...")
		require.Contains(t, errors.String(), "gostructfmt --diff -w")
		require.Contains(t, errors.String(), "changed lines")
		require.Contains(t, errors.String(), "-w, --write")
		require.Contains(t, errors.String(), "--diff")
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

	t.Run("formats only changed Go files in diff mode", func(t *testing.T) {
		// Arrange
		directory := t.TempDir()
		t.Chdir(directory)
		require.NoError(t, os.WriteFile("changed.go", []byte("package example\ntype User struct { Name string }\nvar _ = User{Name: \"old\"}\n"), 0o600))
		require.NoError(t, os.WriteFile("unchanged.go", []byte("package example\ntype Other struct { Name string }\nvar _ = Other{Name: \"untouched\"}\n"), 0o600))
		runGitCommand(t, directory, "init")
		runGitCommand(t, directory, "add", "changed.go", "unchanged.go")
		runGitCommand(t, directory, "-c", "user.name=gostructfmt", "-c", "user.email=gostructfmt@example.invalid", "commit", "-m", "initial")
		require.NoError(t, os.WriteFile("changed.go", []byte("package example\ntype User struct { Name string }\nvar _ = User{Name: \"old\"}\nvar _ = User{Name: \"new\"}\n"), 0o600))

		var (
			output bytes.Buffer
			errors bytes.Buffer
		)

		// Act
		status := run([]string{"--diff", "-w"}, bytes.NewBuffer(nil), &output, &errors)
		changed, changedErr := os.ReadFile("changed.go")
		unchanged, unchangedErr := os.ReadFile("unchanged.go")

		// Assert
		require.Equal(t, 0, status)
		require.NoError(t, changedErr)
		require.NoError(t, unchangedErr)
		require.Empty(t, output.String())
		require.Empty(t, errors.String())
		require.Contains(t, string(changed), `User{Name: "old"}`)
		require.Contains(t, string(changed), "User{\n")
		require.Contains(t, string(unchanged), `Other{Name: "untouched"}`)
	})

	t.Run("does nothing when diff is empty", func(t *testing.T) {
		// Arrange
		directory := t.TempDir()
		t.Chdir(directory)
		runGitCommand(t, directory, "init")

		var (
			output bytes.Buffer
			errors bytes.Buffer
		)

		// Act
		status := run([]string{"--diff"}, bytes.NewBufferString("not Go source"), &output, &errors)

		// Assert
		require.Equal(t, 0, status)
		require.Empty(t, output.String())
		require.Empty(t, errors.String())
	})
}

func runGitCommand(t *testing.T, directory string, args ...string) {
	t.Helper()

	command := exec.Command("git", args...)
	command.Dir = directory
	require.NoError(t, command.Run())
}
