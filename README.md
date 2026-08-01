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
output; with `-w`, each file is replaced only after successful parsing and
formatting:

```sh
gostructfmt file.go
gostructfmt -w file.go
```

Invalid Go source is reported on standard error and produces a non-zero exit
status. No partial formatted source is emitted for an invalid input.

## Go API

```go
formatted, err := gostructfmt.Format("example.go", source)
```

The API applies `FormatStructLiterals` and then prints the result using the Go
formatter. Calling it repeatedly with the same input produces the same output.
