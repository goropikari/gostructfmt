// Command gostructfmt formats struct literals in Go source.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/goropikari/gostructfmt/internal/cli"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	command := cli.NewCommand(stdin, stdout, stderr)
	command.SetArgs(args)

	if err := command.Execute(); err != nil {
		code := 1
		if exitErr, ok := err.(*cli.ExitError); ok {
			code = exitErr.Code
		}

		fmt.Fprintln(stderr, "gostructfmt:", err)

		return code
	}

	return 0
}
