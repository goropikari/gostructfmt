# Graph Report - . (2026-08-14)

## Corpus Check

- Corpus is ~7,031 words - fits in a single context window. You may not need a graph.

## Summary

- 116 nodes · 215 edges · 9 communities (8 shown, 1 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 26 edges (avg confidence: 0.82)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)

- Struct Literal Formatting
- Go AST File Processing
- Command File Operations
- golangci Linter Plugin
- Testing Standards
- CLI Execution Flow
- Product Features and Config
- CI and Pull Requests
- Go Module Dependency

## God Nodes (most connected - your core abstractions)

1. `Parse()` - 14 edges
2. `FormatStructLiterals()` - 13 edges
3. `File` - 10 edges
4. `LineRange` - 10 edges
5. `selectedLiterals()` - 9 edges
6. `execute()` - 8 edges
7. `File` - 8 edges
8. `FormatStructLiteralsInLines()` - 8 edges
9. `formatSelectedLiteralSource()` - 8 edges
10. `formatSource()` - 7 edges

## Surprising Connections (you probably didn't know these)

- `Repository Testing Policy` --semantically_similar_to--> `Testing Guidelines` [INFERRED] [semantically similar]
  AGENTS.md → TESTING.md
- `Test Framework Guidelines` --semantically_similar_to--> `Testing Guidelines` [INFERRED] [semantically similar]
  docs/testing-guidelines.md → TESTING.md
- `AAA Test Structure` --semantically_similar_to--> `Arrange Act Assert Pattern` [INFERRED] [semantically similar]
  docs/testing-guidelines.md → TESTING.md
- `run()` --calls--> `NewCommand()` [INFERRED]
  cmd/gostructfmt/main.go → internal/cli/command.go
- `TestFileAccessorsAndMutation()` --calls--> `New()` [INFERRED]
  internal/formatter/parse_test.go → plugin/gostructfmt/gostructfmt.go

## Import Cycles

- None detected.

## Hyperedges (group relationships)

- **Repository Test Quality Practices** — testing_aaa_pattern, testing_testify_assertions, testing_deterministic_tests, docs_testing_guidelines_t_run, agents_testing_policy [INFERRED 0.95]
- **Repository Quality Gate** — agents_development_checks, github_workflows_ci_lint_and_format, golangci_standard_linters, golangci_gostructfmt_gostructfmt_linter [INFERRED 0.85]

## Communities (9 total, 1 thin omitted)

### Community 0 - "Struct Literal Formatting"

Cohesion: 0.21
Nodes (24): CompositeLit, Expr, File, sourceEdit, addLiteralEdits(), applySourceEdits(), declaredStructTypes(), elidedCompositeLiteralTypes() (+16 more)

### Community 1 - "Go AST File Processing"

Cohesion: 0.13
Nodes (16): FileSet, File, Parse(), ParseAndPrint(), T, TestFileAccessorsAndMutation(), TestFilePrint(), TestParse() (+8 more)

### Community 2 - "Command File Operations"

Cohesion: 0.21
Nodes (20): ExitError, Command, FileMode, LineRange, atomicWriteFile(), changedGoFiles(), execute(), expandCurrentDirectory() (+12 more)

### Community 3 - "golangci Linter Plugin"

Cohesion: 0.20
Nodes (7): Analyzer, Plugin, LinterPlugin, init(), New(), T, TestRun()

### Community 4 - "Testing Standards"

Cohesion: 0.20
Nodes (10): Repository Testing Policy, AAA Test Structure, Test Framework Guidelines, Subtest Scenarios, Arrange Act Assert Pattern, Deterministic Tests, Race Detection, Test Pyramid (+2 more)

### Community 5 - "CLI Execution Flow"

Cohesion: 0.31
Nodes (7): Reader, Writer, main(), run(), T, runGitCommand(), TestRun()

### Community 6 - "Product Features and Config"

Cohesion: 0.29
Nodes (7): gostructfmt golangci-lint Module Plugin, gostructfmt Linter Configuration, Atomic File Replacement, Command Line Interface, Git Diff Formatting, golangci-lint Plugin, gostructfmt

### Community 7 - "CI and Pull Requests"

Cohesion: 0.33
Nodes (6): Development Checks, Pull Request Requirements, Pull Request Template, CI Workflow, Lint and Format Job, Standard golangci-lint Configuration

## Knowledge Gaps

- **10 isolated node(s):** `github.com/goropikari/gostructfmt`, `Git Diff Formatting`, `Testify Assertions`, `Race Detection`, `Test Framework Guidelines` (+5 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions

_Questions this graph is uniquely positioned to answer:_

- **Why does `Parse()` connect `Go AST File Processing` to `Struct Literal Formatting`, `Command File Operations`?**
  _High betweenness centrality (0.260) - this node is a cross-community bridge._
- **Why does `formatSource()` connect `Command File Operations` to `Struct Literal Formatting`, `Go AST File Processing`?**
  _High betweenness centrality (0.180) - this node is a cross-community bridge._
- **Why does `LineRange` connect `Command File Operations` to `Struct Literal Formatting`, `Go AST File Processing`?**
  _High betweenness centrality (0.108) - this node is a cross-community bridge._
- **Are the 11 inferred relationships involving `Parse()` (e.g. with `.syncSourceFromAST()` and `formatSource()`) actually correct?**
  _`Parse()` has 11 INFERRED edges - model-reasoned connections that need verification._
- **Are the 3 inferred relationships involving `FormatStructLiterals()` (e.g. with `formatSource()` and `Parse()`) actually correct?**
  _`FormatStructLiterals()` has 3 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/goropikari/gostructfmt`, `Git Diff Formatting`, `Testify Assertions` to the rest of the system?**
  _10 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Go AST File Processing` be split into smaller, more focused modules?**
  _Cohesion score 0.12666666666666668 - nodes in this community are weakly interconnected._
