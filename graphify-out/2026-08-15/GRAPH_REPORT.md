# Graph Report - gostructfmt (2026-08-14)

## Corpus Check

- 22 files · ~10,019 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary

- 134 nodes · 240 edges · 11 communities (9 shown, 2 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 26 edges (avg confidence: 0.82)
- Token cost: 0 input · 0 output

## Graph Freshness

- Built from commit: `0cd2a803`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)

- struct_literal.go
- Parse
- command.go
- Plugin
- Testing Guidelines
- run
- gostructfmt
- Lint and Format Job
- github.com/goropikari/gotreesj
- Factors and evidence
- README.md

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
10. `elidedCompositeLiteralTypes()` - 8 edges

## Surprising Connections (you probably didn't know these)

- `Repository Testing Policy` --semantically_similar_to--> `Testing Guidelines` [INFERRED] [semantically similar]
  AGENTS.md → TESTING.md
- `Test Framework Guidelines` --semantically_similar_to--> `Testing Guidelines` [INFERRED] [semantically similar]
  docs/testing-guidelines.md → TESTING.md
- `AAA Test Structure` --semantically_similar_to--> `Arrange Act Assert Pattern` [INFERRED] [semantically similar]
  docs/testing-guidelines.md → TESTING.md
- `run()` --calls--> `NewCommand()` [INFERRED]
  cmd/gotreesj/main.go → internal/cli/command.go
- `run()` --calls--> `Parse()` [INFERRED]
  plugin/gotreesj/gotreesj.go → internal/formatter/parse.go

## Import Cycles

- None detected.

## Hyperedges (group relationships)

- **Repository Test Quality Practices** — testing_aaa_pattern, testing_testify_assertions, testing_deterministic_tests, docs_testing_guidelines_t_run, agents_testing_policy [INFERRED 0.95]

## Communities (11 total, 2 thin omitted)

### Community 0 - "struct_literal.go"

Cohesion: 0.17
Nodes (30): CompositeLit, Expr, File, sourceEdit, addElidedArrayElementTypes(), addElidedCompositeLiteralType(), addElidedMapEntryTypes(), addLiteralEdits() (+22 more)

### Community 1 - "Parse"

Cohesion: 0.16
Nodes (12): FileSet, File, Parse(), ParseAndPrint(), T, TestFileAccessorsAndMutation(), TestFilePrint(), TestParse() (+4 more)

### Community 2 - "command.go"

Cohesion: 0.21
Nodes (20): ExitError, Command, FileMode, LineRange, atomicWriteFile(), changedGoFiles(), execute(), expandCurrentDirectory() (+12 more)

### Community 3 - "Plugin"

Cohesion: 0.17
Nodes (9): Analyzer, Plugin, LinterPlugin, Pass, init(), New(), run(), T (+1 more)

### Community 4 - "Testing Guidelines"

Cohesion: 0.20
Nodes (10): Repository Testing Policy, AAA Test Structure, Test Framework Guidelines, Subtest Scenarios, Arrange Act Assert Pattern, Deterministic Tests, Race Detection, Test Pyramid (+2 more)

### Community 5 - "run"

Cohesion: 0.31
Nodes (7): Reader, Writer, main(), run(), T, runGitCommand(), TestRun()

### Community 6 - "gostructfmt"

Cohesion: 0.33
Nodes (6): gostructfmt golangci-lint Module Plugin, Atomic File Replacement, Command Line Interface, Git Diff Formatting, golangci-lint Plugin, gostructfmt

### Community 7 - "Lint and Format Job"

Cohesion: 0.33
Nodes (6): Development Checks, Pull Request Requirements, Pull Request Template, CI Workflow, Lint and Format Job, Standard golangci-lint Configuration

### Community 9 - "Factors and evidence"

Cohesion: 0.15
Nodes (12): Deterministic and idempotent formatting, Factors and evidence, Functional correctness, Goals, Handoff, Interface compatibility, Interpretation rules, Maintainable responsibility boundaries (+4 more)

## Knowledge Gaps

- **21 isolated node(s):** `github.com/goropikari/gotreesj`, `Repository quality state`, `Product boundary`, `Goals`, `Functional correctness` (+16 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions

_Questions this graph is uniquely positioned to answer:_

- **Why does `Parse()` connect `Parse` to `struct_literal.go`, `command.go`, `Plugin`?**
  _High betweenness centrality (0.203) - this node is a cross-community bridge._
- **Why does `formatSource()` connect `command.go` to `struct_literal.go`, `Parse`?**
  _High betweenness centrality (0.137) - this node is a cross-community bridge._
- **Why does `LineRange` connect `command.go` to `struct_literal.go`?**
  _High betweenness centrality (0.092) - this node is a cross-community bridge._
- **Are the 11 inferred relationships involving `Parse()` (e.g. with `.syncSourceFromAST()` and `formatSource()`) actually correct?**
  _`Parse()` has 11 INFERRED edges - model-reasoned connections that need verification._
- **Are the 3 inferred relationships involving `FormatStructLiterals()` (e.g. with `formatSource()` and `Parse()`) actually correct?**
  _`FormatStructLiterals()` has 3 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/goropikari/gotreesj`, `Repository quality state`, `Product boundary` to the rest of the system?**
  _21 weakly-connected nodes found - possible documentation gaps or missing edges._
