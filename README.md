# gotreesj

`gotreesj` formats populated Go struct literals and long function calls as
deterministic, gofmt-compatible multi-line syntax. Empty and non-struct
literals are left unchanged; this includes slice, array, and map containers,
though populated struct literals within their elements, keys, or values are
formatted. Function calls longer than 120 characters are split one argument
per line.

## Command line

Install the formatter with:

```sh
go install github.com/goropikari/gotreesj/cmd/gotreesj@latest
```

It reads Go source from standard input and writes formatted source to standard
output, which makes it suitable for formatter and lint pipelines:

```sh
gotreesj < input.go > output.go
```

Files may be passed directly. Without `-w`, the result is written to standard
output for one file; with `-w`, each file is atomically replaced only after
successful parsing and formatting:

```sh
gotreesj file.go
gotreesj -w file.go
```

To format all Go files recursively under the current directory, use:

```sh
gotreesj -w ./...
```

To format only struct literals overlapping changed lines in the working tree's
Git diff, use:

```sh
gotreesj --diff -w
```

Hidden directories, `vendor`, `testdata`, and underscore-prefixed directories
are skipped.

`-w` requires at least one file. Multiple files also require `-w`; this avoids
concatenating multiple package declarations into one invalid standard-output
source file.
Invalid Go source is reported on standard error and produces a non-zero exit
status. No partial formatted source is emitted for an invalid input.

## golangci-lint

The repository includes a golangci-lint module plugin. Build a custom
golangci-lint binary and run it with the bundled plugin configuration:

```sh
golangci-lint custom
./custom-gcl run --config .golangci-gotreesj.yml ./...
```

The plugin reports populated struct literals and long function calls that need
formatting. Apply the changes with `gotreesj --diff -w`.
