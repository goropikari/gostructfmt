package formatter

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"sort"
)

// FormatStructLiterals expands populated struct composite literals into the
// multi-line form understood by gofmt. It leaves empty literals and
// slice/array/map literals unchanged, and preserves comments by applying
// edits to the original source before parsing it again.
//
// Formatting reparses the result, so callers must reacquire AST and FileSet
// pointers after this function returns. Qualified and indexed type
// expressions are treated as struct literals because this package does not
// load external type information; callers requiring semantic type filtering
// should perform that filtering before invoking this function.
func FormatStructLiterals(file *File) error {
	if file == nil {
		return fmt.Errorf("format struct literals: nil file")
	}

	if !file.sourceValid {
		if err := file.syncSourceFromAST(); err != nil {
			return fmt.Errorf("format struct literals: synchronize source: %w", err)
		}
	}

	structTypes := declaredStructTypes(file.ast)

	fileInfo := file.fileSet.File(file.ast.Pos())
	if fileInfo == nil {
		return fmt.Errorf("format struct literals: source file unavailable")
	}

	edits := make([]sourceEdit, 0)

	ast.Inspect(file.ast, func(node ast.Node) bool {
		literal, ok := node.(*ast.CompositeLit)
		if !ok || len(literal.Elts) == 0 || !isStructLiteralType(literal.Type, structTypes) {
			return true
		}

		addLiteralEdits(&edits, literal, fileInfo, file.source)

		return true
	})

	if len(edits) == 0 {
		return nil
	}

	updatedSource := applySourceEdits(file.source, edits)

	updatedFile, err := Parse(file.filename, updatedSource)
	if err != nil {
		return fmt.Errorf("format struct literals: %w", err)
	}

	file.ast = updatedFile.ast
	file.fileSet = updatedFile.fileSet
	file.source = updatedFile.source
	file.sourceValid = true

	return nil
}

func (f *File) syncSourceFromAST() error {
	output, err := f.Print()
	if err != nil {
		return err
	}

	updated, err := Parse(f.filename, output)
	if err != nil {
		return err
	}

	f.ast = updated.ast
	f.fileSet = updated.fileSet
	f.source = updated.source
	f.sourceValid = true

	return nil
}

// TransformStructLiterals is an alias for FormatStructLiterals for callers
// that describe formatter phases as AST transformations.
func TransformStructLiterals(file *File) error {
	return FormatStructLiterals(file)
}

type sourceEdit struct {
	offset int
	end    int
	text   string
}

// LineRange identifies an inclusive, one-based range of source lines.
type LineRange struct {
	Start int
	End   int
}

// FormatStructLiteralsInLines formats only struct literals that overlap one
// of the supplied source ranges. Source outside those literals is preserved
// byte-for-byte.
func FormatStructLiteralsInLines(file *File, ranges []LineRange) ([]byte, error) {
	if file == nil {
		return nil, fmt.Errorf("format struct literals in lines: nil file")
	}

	if !file.sourceValid {
		if err := file.syncSourceFromAST(); err != nil {
			return nil, fmt.Errorf("format struct literals in lines: synchronize source: %w", err)
		}
	}

	updatedSource, err := formatSelectedLiteralSource(file, ranges)
	if err != nil {
		return nil, err
	}

	if _, err := Parse(file.filename, updatedSource); err != nil {
		return nil, fmt.Errorf("format struct literals in lines: %w", err)
	}

	return updatedSource, nil
}

func formatSelectedLiteralSource(file *File, ranges []LineRange) ([]byte, error) {
	fileInfo := file.fileSet.File(file.ast.Pos())
	if fileInfo == nil {
		return nil, fmt.Errorf("format struct literals in lines: source file unavailable")
	}

	literals := selectedLiterals(file.ast, fileInfo, ranges)

	edits := make([]sourceEdit, 0, len(literals))
	for _, literal := range literals {
		if hasSelectedAncestor(literal, literals, fileInfo) {
			continue
		}

		start := fileInfo.Offset(literal.Pos())
		end := fileInfo.Offset(literal.End())

		formatted, err := formatLiteralSource(file.source[start:end])
		if err != nil {
			return nil, fmt.Errorf("format struct literals in lines: %w", err)
		}

		edits = append(edits, sourceEdit{offset: start, end: end, text: string(formatted)})
	}

	return applySourceEdits(file.source, edits), nil
}

func selectedLiterals(fileAST *ast.File, fileInfo *token.File, ranges []LineRange) []*ast.CompositeLit {
	structTypes := declaredStructTypes(fileAST)
	literals := make([]*ast.CompositeLit, 0)

	ast.Inspect(fileAST, func(node ast.Node) bool {
		literal, ok := node.(*ast.CompositeLit)
		if !ok || len(literal.Elts) == 0 || !isStructLiteralType(literal.Type, structTypes) {
			return true
		}

		if literalOverlapsLines(literal, fileInfo, ranges) {
			literals = append(literals, literal)
		}

		return true
	})

	return literals
}

func literalOverlapsLines(literal *ast.CompositeLit, file *token.File, ranges []LineRange) bool {
	start := file.Line(literal.Pos())

	end := file.Line(literal.End())
	for _, lineRange := range ranges {
		if lineRange.Start <= end && start <= lineRange.End {
			return true
		}
	}

	return false
}

func hasSelectedAncestor(literal *ast.CompositeLit, selected []*ast.CompositeLit, file *token.File) bool {
	start := file.Offset(literal.Pos())

	end := file.Offset(literal.End())
	for _, candidate := range selected {
		candidateStart := file.Offset(candidate.Pos())

		candidateEnd := file.Offset(candidate.End())
		if candidateStart < start && end < candidateEnd {
			return true
		}
	}

	return false
}

func formatLiteralSource(source []byte) ([]byte, error) {
	wrapped := append([]byte("package example\nvar _ = "), source...)
	wrapped = append(wrapped, '\n')

	file, err := Parse("<literal>", wrapped)
	if err != nil {
		return nil, err
	}

	if err := FormatStructLiterals(file); err != nil {
		return nil, err
	}

	formatted, err := file.Print()
	if err != nil {
		return nil, err
	}

	const marker = "var _ = "

	_, after, ok := bytes.Cut(formatted, []byte(marker))
	if !ok {
		return nil, fmt.Errorf("formatted literal marker unavailable")
	}

	literal := bytes.TrimSpace(after)

	return append([]byte(nil), literal...), nil
}

func addLiteralEdits(edits *[]sourceEdit, literal *ast.CompositeLit, file *token.File, source []byte) {
	lbrace := file.Offset(literal.Lbrace)
	rbrace := file.Offset(literal.Rbrace)
	multiline := bytes.Contains(source[lbrace:rbrace], []byte{'\n'})

	if !multiline {
		*edits = append(*edits, sourceEdit{
			offset: lbrace + 1,
			text:   "\n",
		})
		for _, element := range literal.Elts[1:] {
			*edits = append(*edits, sourceEdit{
				offset: file.Offset(element.Pos()),
				text:   "\n",
			})
		}

		*edits = append(*edits, sourceEdit{
			offset: rbrace,
			text:   "\n",
		})
	}

	for index, element := range literal.Elts {
		end := rbrace
		if index+1 < len(literal.Elts) {
			end = file.Offset(literal.Elts[index+1].Pos())
		}

		start := file.Offset(element.End())
		if !hasCommaBefore(source, start, end) {
			*edits = append(*edits, sourceEdit{
				offset: start,
				text:   ",",
			})
		}
	}
}

func hasCommaBefore(source []byte, start, end int) bool {
	for start < end {
		start = skipWhitespace(source, start, end)
		if start >= end {
			return false
		}

		if source[start] == ',' {
			return true
		}

		next, ok := skipComment(source, start, end)
		if !ok {
			return false
		}

		start = next
	}

	return false
}

func skipWhitespace(source []byte, start, end int) int {
	for start < end {
		switch source[start] {
		case ' ', '\t', '\r', '\n':
			start++
		default:
			return start
		}
	}

	return start
}

func skipComment(source []byte, start, end int) (int, bool) {
	if start+1 >= end || source[start] != '/' {
		return start, false
	}

	if source[start+1] == '/' {
		return skipLineComment(source, start+2, end), true
	}

	if source[start+1] != '*' {
		return start, false
	}

	return skipBlockComment(source, start+2, end), true
}

func skipLineComment(source []byte, start, end int) int {
	for start < end && source[start] != '\n' {
		start++
	}

	return start
}

func skipBlockComment(source []byte, start, end int) int {
	for start+1 < end && (source[start] != '*' || source[start+1] != '/') {
		start++
	}

	if start+1 < end {
		start += 2
	}

	return start
}

func applySourceEdits(source []byte, edits []sourceEdit) []byte {
	sort.SliceStable(edits, func(i, j int) bool {
		if edits[i].offset == edits[j].offset {
			return edits[i].text == "," && edits[j].text != ","
		}

		return edits[i].offset < edits[j].offset
	})

	var output bytes.Buffer

	previous := 0
	for _, edit := range edits {
		if edit.offset < previous || edit.offset > len(source) {
			continue
		}

		output.Write(source[previous:edit.offset])
		output.WriteString(edit.text)

		previous = max(edit.end, edit.offset)
	}

	output.Write(source[previous:])

	return output.Bytes()
}

func declaredStructTypes(file *ast.File) map[string]bool {
	expressions := make(map[string]ast.Expr)

	for _, declaration := range file.Decls {
		group, ok := declaration.(*ast.GenDecl)
		if !ok || group.Tok != token.TYPE {
			continue
		}

		for _, specification := range group.Specs {
			typeSpec, ok := specification.(*ast.TypeSpec)
			if !ok {
				continue
			}

			expressions[typeSpec.Name.Name] = typeSpec.Type
		}
	}

	result := make(map[string]bool, len(expressions))
	for name := range expressions {
		result[name] = resolveStructType(name, expressions, make(map[string]bool))
	}

	return result
}

func resolveStructType(name string, expressions map[string]ast.Expr, visiting map[string]bool) bool {
	if visiting[name] {
		return false
	}

	expression, ok := expressions[name]
	if !ok {
		return false
	}

	visiting[name] = true
	defer delete(visiting, name)

	return isStructType(expression, expressions, visiting)
}

func isStructLiteralType(expression ast.Expr, declared map[string]bool) bool {
	switch typeExpression := expression.(type) {
	case *ast.StructType:
		return true
	case *ast.Ident:
		structType, known := declared[typeExpression.Name]
		return !known || structType
	case *ast.SelectorExpr, *ast.IndexExpr, *ast.IndexListExpr:
		return true
	case *ast.ParenExpr:
		return isStructLiteralType(typeExpression.X, declared)
	default:
		return false
	}
}

func isStructType(expression ast.Expr, expressions map[string]ast.Expr, visiting map[string]bool) bool {
	switch typeExpression := expression.(type) {
	case *ast.StructType:
		return true
	case *ast.Ident:
		return resolveStructType(typeExpression.Name, expressions, visiting)
	case *ast.ParenExpr:
		return isStructType(typeExpression.X, expressions, visiting)
	default:
		return false
	}
}
