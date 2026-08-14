// Package formatter contains the implementation used by the gotreesj CLI.
package formatter

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

// File is a parsed Go source file and the position information needed to print
// it again. The AST and FileSet must be kept together because AST positions
// refer to files in the FileSet.
type File struct {
	ast      *ast.File
	fileSet  *token.FileSet
	filename string
	source   []byte
	// sourceValid is false after an arbitrary AST mutation because the
	// original source can no longer be used as the transformation base.
	sourceValid bool
}

// AST returns the parsed syntax tree for this file. The returned tree may be
// edited by a formatter; its positions remain associated with FileSet. Calling
// AST marks the source snapshot stale so a subsequent source-based formatter
// first synchronizes it from the current tree.
func (f *File) AST() *ast.File {
	if f == nil {
		return nil
	}

	f.sourceValid = false

	return f.ast
}

// FileSet returns the position set associated with the parsed syntax tree.
// It is exposed for AST transformations that need source positions, but the
// FileSet itself cannot be replaced on File.
func (f *File) FileSet() *token.FileSet {
	if f == nil {
		return nil
	}

	return f.fileSet
}

// Filename returns the filename used while parsing this file.
func (f *File) Filename() string {
	if f == nil {
		return ""
	}

	return f.filename
}

// MutateAST applies edit to the parsed AST while keeping the File/FileSet
// association intact. It is intended for formatter phases that rewrite AST
// nodes before calling Print.
func (f *File) MutateAST(edit func(*ast.File) error) error {
	if f == nil {
		return fmt.Errorf("mutate Go source: nil file")
	}

	if f.ast == nil {
		return fmt.Errorf("mutate Go source %q: nil AST", f.filename)
	}

	if edit == nil {
		return fmt.Errorf("mutate Go source %q: nil edit function", f.filename)
	}

	err := edit(f.ast)
	if err == nil {
		f.sourceValid = false
	}

	return err
}

// Parse parses src as a Go source file while retaining comments. It returns no
// partial file when parsing fails, so callers cannot accidentally print a
// partially parsed or otherwise unsafe result.
func Parse(filename string, src []byte) (*File, error) {
	if filename == "" {
		filename = "<input>"
	}

	fileSet := token.NewFileSet()

	fileAST, err := parser.ParseFile(fileSet, filename, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse Go source %q: %w", filename, err)
	}

	return &File{
		ast:         fileAST,
		fileSet:     fileSet,
		filename:    filename,
		source:      append([]byte(nil), src...),
		sourceValid: true,
	}, nil
}

// Print renders a parsed file using go/format. The result is gofmt-compatible
// and includes comments retained during Parse.
func (f *File) Print() ([]byte, error) {
	if f == nil {
		return nil, fmt.Errorf("print Go source: nil file")
	}

	if f.ast == nil {
		return nil, fmt.Errorf("print Go source %q: nil AST", f.filename)
	}

	if f.fileSet == nil {
		return nil, fmt.Errorf("print Go source %q: nil file set", f.filename)
	}

	var output bytes.Buffer

	err := format.Node(&output, f.fileSet, f.ast)
	if err != nil {
		return nil, fmt.Errorf("print Go source %q: %w", f.filename, err)
	}

	return output.Bytes(), nil
}

// ParseAndPrint parses src with comments and prints it in gofmt-compatible
// form. On parse or print failure it returns no output.
func ParseAndPrint(filename string, src []byte) ([]byte, error) {
	file, err := Parse(filename, src)
	if err != nil {
		return nil, err
	}

	return file.Print()
}
