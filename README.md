# gostructfmt

`gostructfmt` formats populated Go struct literals as deterministic,
gofmt-compatible multi-line literals. Empty literals and slice, array, and map
literals are left unchanged.

## Command line

Install the formatter with:

```sh
go install github.com/goropikari/gostructfmt/cmd/gostructfmt@latest
```

It reads Go source from standard input and writes formatted source to standard
output, which makes it suitable for formatter and lint pipelines:

```sh
gostructfmt < input.go > output.go
```

Files may be passed directly. Without `-w`, the result is written to standard
output for one file; with `-w`, each file is atomically replaced only after
successful parsing and formatting:

```sh
gostructfmt file.go
gostructfmt -w file.go
```

To format all Go files recursively under the current directory, use:

```sh
gostructfmt -w ./...
```

To format only struct literals overlapping changed lines in the working tree's
Git diff, use:

```sh
gostructfmt --diff -w
```

Hidden directories, `vendor`, `testdata`, and underscore-prefixed directories
are skipped.

`-w` requires at least one file. Multiple files also require `-w`; this avoids
concatenating multiple package declarations into one invalid standard-output
source file.
Invalid Go source is reported on standard error and produces a non-zero exit
status. No partial formatted source is emitted for an invalid input.
