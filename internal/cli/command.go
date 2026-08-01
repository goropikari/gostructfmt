package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goropikari/gostructfmt"
	"github.com/spf13/cobra"
)

// ExitError carries the process status required by the CLI wrapper.
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string { return e.Err.Error() }

// NewCommand creates the gostructfmt command with injectable streams.
func NewCommand(stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	writeFiles := false

	command := &cobra.Command{
		Use:           "gostructfmt [files...]",
		Short:         "format Go struct literals",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(command *cobra.Command, args []string) error {
			return execute(command, args, stdin, stdout, stderr, writeFiles)
		},
	}
	command.SetOut(stderr)
	command.SetErr(stderr)
	command.Flags().BoolVarP(&writeFiles, "write", "w", false, "write result to source files instead of standard output")

	return command
}

func execute(command *cobra.Command, args []string, stdin io.Reader, stdout, stderr io.Writer, writeFiles bool) error {
	filenames, err := expandPackagePatterns(args)
	if err != nil {
		return &ExitError{Code: 2, Err: err}
	}

	if len(filenames) == 0 {
		if writeFiles {
			return &ExitError{Code: 2, Err: fmt.Errorf("-w requires at least one file")}
		}

		return formatReader("<stdin>", stdin, stdout, stderr)
	}

	if len(filenames) > 1 && !writeFiles {
		return &ExitError{Code: 2, Err: fmt.Errorf("multiple files require -w")}
	}

	status := 0

	for _, filename := range filenames {
		if err := formatFile(filename, writeFiles, stdout); err != nil {
			fmt.Fprintf(stderr, "gostructfmt: %s: %v\n", filepath.Clean(filename), err)

			status = 1
		}
	}

	if status != 0 {
		return &ExitError{Code: status, Err: fmt.Errorf("formatting failed")}
	}

	return nil
}

func formatReader(filename string, input io.Reader, output, stderr io.Writer) error {
	source, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("gostructfmt: %s: read: %w", filename, err)
	}

	formatted, err := formatSource(filename, source)
	if err != nil {
		return fmt.Errorf("gostructfmt: %s: %w", filename, err)
	}

	if _, err := output.Write(formatted); err != nil {
		return fmt.Errorf("gostructfmt: %s: write: %w", filename, err)
	}

	return nil
}

func formatFile(filename string, writeFiles bool, stdout io.Writer) error {
	source, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	formatted, err := formatSource(filename, source)
	if err != nil {
		return err
	}

	if !writeFiles {
		_, err = stdout.Write(formatted)
		return err
	}

	info, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("stat: %w", err)
	}

	if err := atomicWriteFile(filename, formatted, info.Mode().Perm()); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func formatSource(filename string, source []byte) ([]byte, error) {
	file, err := gostructfmt.Parse(filename, source)
	if err != nil {
		return nil, err
	}

	if err := gostructfmt.FormatStructLiterals(file); err != nil {
		return nil, err
	}

	return file.Print()
}

func expandPackagePatterns(paths []string) ([]string, error) {
	expanded := make([]string, 0, len(paths))
	for _, path := range paths {
		if path != "./..." {
			expanded = append(expanded, path)
			continue
		}

		files, err := expandCurrentDirectory()
		if err != nil {
			return nil, fmt.Errorf("expand ./...: %w", err)
		}

		expanded = append(expanded, files...)
	}

	sort.Strings(expanded)

	return expanded, nil
}

func expandCurrentDirectory() ([]string, error) {
	files := make([]string, 0)
	err := filepath.WalkDir(".", func(current string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() {
			return skipPackageDirectory(current, entry.Name())
		}

		if filepath.Ext(current) == ".go" {
			files = append(files, current)
		}

		return nil
	})

	return files, err
}

func skipPackageDirectory(path, name string) error {
	if path == "." {
		return nil
	}

	if name == "vendor" || name == "testdata" || strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
		return filepath.SkipDir
	}

	return nil
}

func atomicWriteFile(filename string, contents []byte, mode os.FileMode) error {
	temporary, err := os.CreateTemp(filepath.Dir(filename), ".gostructfmt-*")
	if err != nil {
		return err
	}

	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)

	if err := temporary.Chmod(mode); err != nil {
		_ = temporary.Close()
		return err
	}

	if _, err := temporary.Write(contents); err != nil {
		_ = temporary.Close()
		return err
	}

	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return err
	}

	if err := temporary.Close(); err != nil {
		return err
	}

	return os.Rename(temporaryName, filename)
}
