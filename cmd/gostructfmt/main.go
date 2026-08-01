// Command gostructfmt formats struct literals in Go source.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/goropikari/gostructfmt"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("gostructfmt", flag.ContinueOnError)
	flags.SetOutput(stderr)

	writeFiles := flags.Bool("w", false, "write result to source files instead of standard output")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}

		return 2
	}

	filenames := flags.Args()
	if len(filenames) == 0 {
		if *writeFiles {
			fmt.Fprintln(stderr, "gostructfmt: -w requires at least one file")
			return 2
		}

		return formatReader("<stdin>", stdin, stdout, stderr)
	}

	if len(filenames) > 1 && !*writeFiles {
		fmt.Fprintln(stderr, "gostructfmt: multiple files require -w")
		return 2
	}

	status := 0

	for _, filename := range filenames {
		if err := formatFile(filename, *writeFiles, stdout); err != nil {
			fmt.Fprintf(stderr, "gostructfmt: %s: %v\n", filepath.Clean(filename), err)

			status = 1
		}
	}

	return status
}

func formatReader(filename string, input io.Reader, output, stderr io.Writer) int {
	source, err := io.ReadAll(input)
	if err != nil {
		fmt.Fprintf(stderr, "gostructfmt: %s: read: %v\n", filename, err)
		return 1
	}

	formatted, err := gostructfmt.Format(filename, source)
	if err != nil {
		fmt.Fprintf(stderr, "gostructfmt: %s: %v\n", filename, err)
		return 1
	}

	if _, err := output.Write(formatted); err != nil {
		fmt.Fprintf(stderr, "gostructfmt: %s: write: %v\n", filename, err)
		return 1
	}

	return 0
}

func formatFile(filename string, writeFiles bool, stdout io.Writer) error {
	source, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	formatted, err := gostructfmt.Format(filename, source)
	if err != nil {
		return err
	}

	if writeFiles {
		info, err := os.Stat(filename)
		if err != nil {
			return fmt.Errorf("stat: %w", err)
		}

		// The path is explicitly supplied by the CLI user and is intentionally
		// the target of -w; no path is derived from untrusted file contents.
		if err := atomicWriteFile(filename, formatted, info.Mode().Perm()); err != nil {
			return fmt.Errorf("write: %w", err)
		}

		return nil
	}

	if _, err := stdout.Write(formatted); err != nil {
		return fmt.Errorf("write: %w", err)
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
