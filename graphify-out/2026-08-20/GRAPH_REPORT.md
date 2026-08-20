# Graph Report - . (2026-08-18)

## Corpus Check

- Corpus is ~10,488 words - fits in a single context window. You may not need a graph.

## Summary

- 136 nodes · 262 edges · 12 communities (10 shown, 2 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 28 edges (avg confidence: 0.82)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)

- Repository Quality Rules
- CLI Command Pipeline
- AST Parsing and Tests
- Linter Plugin Integration
- Struct Literal Formatting
- Composite Literal Analysis
- CLI Entry Point
- Source Edit Application
- Plugin Diagnostics
- Struct Type Resolution
- Quality Gate Metadata
- Go Module Root

## God Nodes (most connected - your core abstractions)

1. `Parse()` - 14 edges
2. `FormatStructLiterals()` - 14 edges
3. `File` - 12 edges
4. `LineRange` - 11 edges
5. `formatSelectedLiteralSource()` - 9 edges
6. `addLongCallEdits()` - 9 edges
7. `selectedLiterals()` - 9 edges
8. `execute()` - 8 edges
9. `File` - 8 edges
10. `FormatStructLiteralsInLines()` - 8 edges

## Surprising Connections (you probably didn't know these)

- `run()` --calls--> `NewCommand()` [INFERRED]
  cmd/gotreesj/main.go → internal/cli/command.go
- `run()` --calls--> `Parse()` [INFERRED]
  plugin/gotreesj/gotreesj.go → internal/formatter/parse.go
- `TestFileAccessorsAndMutation()` --calls--> `New()` [INFERRED]
  internal/formatter/parse_test.go → plugin/gotreesj/gotreesj.go
- `reportFormattingIssue()` --calls--> `FormatStructLiteralsInLines()` [INFERRED]
  plugin/gotreesj/gotreesj.go → internal/formatter/struct_literal.go
- `CI workflow` --conceptually_related_to--> `Interface compatibility` [INFERRED]
  .github/workflows/ci.yml → docs/quality-model.md

## Import Cycles

- None detected.

## Hyperedges (group relationships)

- **Repository quality and verification practices** — agents_required_checks, testing_behavior_contracts, quality_verifiable_behavior, github_workflows_ci [INFERRED 0.85]
- **Formatter behavior and integration contracts** — readme_struct_literal_formatting, readme_long_function_call_formatting, readme_atomic_writes, readme_golangci_plugin [EXTRACTED 1.00]

## Communities (12 total, 2 thin omitted)

### Community 0 - "Repository Quality Rules"

Cohesion: 0.08
Nodes (27): Agent and development instructions, Change size guidance, Required development checks, Custom golangci-lint configuration, Quality model, Repository testing guidelines, Pull request template, CI workflow (+19 more)

### Community 1 - "CLI Command Pipeline"

Cohesion: 0.21
Nodes (20): ExitError, Command, FileMode, LineRange, atomicWriteFile(), changedGoFiles(), execute(), expandCurrentDirectory() (+12 more)

### Community 2 - "AST Parsing and Tests"

Cohesion: 0.16
Nodes (12): FileSet, File, Parse(), ParseAndPrint(), T, TestFileAccessorsAndMutation(), TestFilePrint(), TestParse() (+4 more)

### Community 3 - "Linter Plugin Integration"

Cohesion: 0.24
Nodes (9): Analyzer, Plugin, Pass, formattingMessage(), File, Node, init(), reportFormattingIssue() (+1 more)

### Community 4 - "Struct Literal Formatting"

Cohesion: 0.29
Nodes (10): CallExpr, addLongCallEdits(), formatLiteralSource(), formatSelectedLiteralSource(), FormatStructLiteralsInLines(), File, Node, hasSelectedAncestor() (+2 more)

### Community 5 - "Composite Literal Analysis"

Cohesion: 0.42
Nodes (10): CompositeLit, Expr, addElidedArrayElementTypes(), addElidedCompositeLiteralType(), addElidedMapEntryTypes(), elidedCompositeLiteralTypes(), isStructLiteralType(), isStructLiteralTypeOrElided() (+2 more)

### Community 6 - "CLI Entry Point"

Cohesion: 0.31
Nodes (7): Reader, Writer, main(), run(), T, runGitCommand(), TestRun()

### Community 7 - "Source Edit Application"

Cohesion: 0.42
Nodes (8): sourceEdit, addLiteralEdits(), applySourceEdits(), hasCommaBefore(), skipBlockComment(), skipComment(), skipLineComment(), skipWhitespace()

### Community 8 - "Plugin Diagnostics"

Cohesion: 0.38
Nodes (6): Diagnostic, LinterPlugin, New(), T, runAnalyzer(), TestRun()

### Community 9 - "Struct Type Resolution"

Cohesion: 0.33
Nodes (6): declaredStructTypes(), FormatStructLiterals(), File, isStructType(), resolveStructType(), TransformStructLiterals()

## Knowledge Gaps

- **8 isolated node(s):** `github.com/goropikari/gotreesj`, `Custom golangci-lint configuration`, `Pull request template`, `gotreesj golangci-lint configuration`, `Repository quality state documentation` (+3 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions

_Questions this graph is uniquely positioned to answer:_

- **Why does `Parse()` connect `AST Parsing and Tests` to `CLI Command Pipeline`, `Linter Plugin Integration`, `Struct Literal Formatting`, `Struct Type Resolution`?**
  _High betweenness centrality (0.235) - this node is a cross-community bridge._
- **Why does `formatSource()` connect `CLI Command Pipeline` to `Struct Type Resolution`, `AST Parsing and Tests`, `Struct Literal Formatting`?**
  _High betweenness centrality (0.150) - this node is a cross-community bridge._
- **Why does `FormatStructLiteralsInLines()` connect `Struct Literal Formatting` to `CLI Command Pipeline`, `AST Parsing and Tests`, `Linter Plugin Integration`, `Source Edit Application`?**
  _High betweenness centrality (0.109) - this node is a cross-community bridge._
- **Are the 11 inferred relationships involving `Parse()` (e.g. with `.syncSourceFromAST()` and `formatSource()`) actually correct?**
  _`Parse()` has 11 INFERRED edges - model-reasoned connections that need verification._
- **Are the 3 inferred relationships involving `FormatStructLiterals()` (e.g. with `formatSource()` and `Parse()`) actually correct?**
  _`FormatStructLiterals()` has 3 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/goropikari/gotreesj`, `Custom golangci-lint configuration`, `Pull request template` to the rest of the system?**
  _8 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Repository Quality Rules` be split into smaller, more focused modules?**
  _Cohesion score 0.08262108262108261 - nodes in this community are weakly interconnected._
