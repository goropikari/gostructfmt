package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/goropikari/gotreesj/internal/formatter"
	"github.com/spf13/cobra"
)

// ExitError carries the process status required by the CLI wrapper.
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string { return e.Err.Error() }

// NewCommand creates the gotreesj command with injectable streams.
func NewCommand(stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	writeFiles := false
	diffMode := false

	command := &cobra.Command{
		Use:   "gotreesj [files...]",
		Short: "format Go struct literals",
		Long: "Format populated Go struct literals as deterministic, gofmt-compatible " +
			"multi-line literals. With no files, source is read from standard input. " +
			"Use -w to update files in place; multiple files and ./... require -w. " +
			"Use --diff to format only literals overlapping changed lines in the working " +
			"tree's git diff.",
		Example: `  # Format standard input and print the result
  gotreesj < input.go

  # Format one file and print the result
  gotreesj input.go

  # Format one or more files in place
  gotreesj -w input.go other.go

  # Format all Go files below the current directory in place
  gotreesj -w ./...

  # Format literals overlapping changed lines in the working tree
  gotreesj --diff -w`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(command *cobra.Command, args []string) error {
			return execute(command, args, stdin, stdout, stderr, writeFiles, diffMode)
		},
	}
	command.SetOut(stderr)
	command.SetErr(stderr)
	command.Flags().BoolVarP(&writeFiles, "write", "w", false, "write result to source files instead of standard output")
	command.Flags().BoolVar(&diffMode, "diff", false, "format literals overlapping changed lines in the working tree's git diff")

	return command
}

func execute(command *cobra.Command, args []string, stdin io.Reader, stdout, stderr io.Writer, writeFiles, diffMode bool) error {
	filenames, diffRanges, err := resolveFilenames(args, diffMode)
	if err != nil {
		return &ExitError{
			Code: 2,
			Err:  err,
		}
	}

	if diffMode && len(filenames) == 0 {
		return nil
	}

	if len(filenames) == 0 {
		return formatWithoutFiles(stdin, stdout, stderr, writeFiles)
	}

	if len(filenames) > 1 && !writeFiles {
		return &ExitError{
			Code: 2,
			Err:  fmt.Errorf("multiple files require -w"),
		}
	}

	return formatFiles(filenames, writeFiles, stdout, stderr, diffRanges)
}

func formatWithoutFiles(stdin io.Reader, stdout, stderr io.Writer, writeFiles bool) error {
	if writeFiles {
		return &ExitError{
			Code: 2,
			Err:  fmt.Errorf("-w requires at least one file"),
		}
	}

	return formatReader("<stdin>", stdin, stdout, stderr)
}

func formatFiles(filenames []string, writeFiles bool, stdout, stderr io.Writer, diffRanges map[string][]formatter.LineRange) error {
	status := 0

	for _, filename := range filenames {
		if err := formatFile(filename, writeFiles, stdout, diffRanges[filename]); err != nil {
			fmt.Fprintf(stderr, "gotreesj: %s: %v\n", filepath.Clean(filename), err)

			status = 1
		}
	}

	if status != 0 {
		return &ExitError{
			Code: status,
			Err:  fmt.Errorf("formatting failed"),
		}
	}

	return nil
}

func resolveFilenames(args []string, diffMode bool) ([]string, map[string][]formatter.LineRange, error) {
	if diffMode && len(args) > 0 {
		return nil, nil, fmt.Errorf("--diff cannot be combined with file arguments")
	}

	if diffMode {
		return changedGoFiles()
	}

	filenames, err := expandPackagePatterns(args)

	return filenames, nil, err
}

func formatReader(filename string, input io.Reader, output, stderr io.Writer) error {
	source, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("gotreesj: %s: read: %w", filename, err)
	}

	formatted, err := formatSource(filename, source, nil)
	if err != nil {
		return fmt.Errorf("gotreesj: %s: %w", filename, err)
	}

	if _, err := output.Write(formatted); err != nil {
		return fmt.Errorf("gotreesj: %s: write: %w", filename, err)
	}

	return nil
}

func formatFile(filename string, writeFiles bool, stdout io.Writer, lineRanges []formatter.LineRange) error {
	source, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	formatted, err := formatSource(filename, source, lineRanges)
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

func formatSource(filename string, source []byte, lineRanges []formatter.LineRange) ([]byte, error) {
	file, err := formatter.Parse(filename, source)
	if err != nil {
		return nil, err
	}

	if lineRanges != nil {
		return formatter.FormatStructLiteralsInLines(file, lineRanges)
	}

	if err := formatter.FormatStructLiterals(file); err != nil {
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

func changedGoFiles() ([]string, map[string][]formatter.LineRange, error) {
	command := exec.Command("git", "diff", "--unified=0", "--no-color", "--diff-filter=ACMRTUXB", "--", "*.go")

	output, err := command.Output()
	if err != nil {
		return nil, nil, fmt.Errorf("read git diff: %w", err)
	}

	lineRanges := make(map[string][]formatter.LineRange)
	currentPath := ""

	for line := range strings.SplitSeq(string(output), "\n") {
		if after, ok := strings.CutPrefix(line, "+++ b/"); ok {
			currentPath = filepath.FromSlash(after)
			continue
		}

		if !strings.HasPrefix(line, "@@ ") || currentPath == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		start, count, err := parseDiffRange(strings.TrimPrefix(fields[2], "+"))
		if err != nil {
			return nil, nil, fmt.Errorf("parse git diff range %q: %w", line, err)
		}

		if count > 0 {
			lineRanges[currentPath] = append(lineRanges[currentPath], formatter.LineRange{
				Start: start,
				End:   start + count - 1,
			})
		}
	}

	files := make([]string, 0, len(lineRanges))
	for path := range lineRanges {
		files = append(files, path)
	}

	sort.Strings(files)

	return files, lineRanges, nil
}

func parseDiffRange(value string) (int, int, error) {
	parts := strings.SplitN(value, ",", 2)

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	count := 1
	if len(parts) == 2 {
		count, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, err
		}
	}

	return start, count, nil
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
	temporary, err := os.CreateTemp(filepath.Dir(filename), ".gotreesj-*")
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
