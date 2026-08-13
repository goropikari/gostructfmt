# Graph Report - gostructfmt (2026-08-14)

## Corpus Check

- 24 files · ~10,699 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary

- 146 nodes · 251 edges · 12 communities (10 shown, 2 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 26 edges (avg confidence: 0.82)
- Token cost: 0 input · 0 output

## Graph Freshness

- Built from commit: `283bb9bc`
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
- github.com/goropikari/gostructfmt
- Factors and evidence
- Implementation Plan
- README.md

## God Nodes (most connected - your core abstractions)

1. `Parse()` - 14 edges
2. `FormatStructLiterals()` - 13 edges
3. `File` - 10 edges
4. `LineRange` - 10 edges
5. `Implementation Plan` - 10 edges
6. `selectedLiterals()` - 9 edges
7. `execute()` - 8 edges
8. `File` - 8 edges
9. `FormatStructLiteralsInLines()` - 8 edges
10. `formatSelectedLiteralSource()` - 8 edges

## Surprising Connections (you probably didn't know these)

- `Repository Testing Policy` --semantically_similar_to--> `Testing Guidelines` [INFERRED] [semantically similar]
  AGENTS.md → TESTING.md
- `Test Framework Guidelines` --semantically_similar_to--> `Testing Guidelines` [INFERRED] [semantically similar]
  docs/testing-guidelines.md → TESTING.md
- `AAA Test Structure` --semantically_similar_to--> `Arrange Act Assert Pattern` [INFERRED] [semantically similar]
  docs/testing-guidelines.md → TESTING.md
- `run()` --calls--> `NewCommand()` [INFERRED]
  cmd/gostructfmt/main.go → internal/cli/command.go
- `run()` --calls--> `Parse()` [INFERRED]
  plugin/gostructfmt/gostructfmt.go → internal/formatter/parse.go

## Import Cycles

- None detected.

## Hyperedges (group relationships)

- **Repository Test Quality Practices** — testing_aaa_pattern, testing_testify_assertions, testing_deterministic_tests, docs_testing_guidelines_t_run, agents_testing_policy [INFERRED 0.95]
- **Repository Quality Gate** — agents_development_checks, github_workflows_ci_lint_and_format, golangci_standard_linters, golangci_gostructfmt_gostructfmt_linter [INFERRED 0.85]

## Communities (12 total, 2 thin omitted)

### Community 0 - "struct_literal.go"

Cohesion: 0.18
Nodes (28): CompositeLit, Expr, File, sourceEdit, addElidedArrayElementTypes(), addElidedCompositeLiteralType(), addElidedMapEntryTypes(), addLiteralEdits() (+20 more)

### Community 1 - "Parse"

Cohesion: 0.15
Nodes (13): FileSet, File, Parse(), ParseAndPrint(), T, TestFileAccessorsAndMutation(), TestFilePrint(), TestParse() (+5 more)

### Community 2 - "command.go"

Cohesion: 0.20
Nodes (21): ExitError, Command, FileMode, LineRange, atomicWriteFile(), changedGoFiles(), execute(), expandCurrentDirectory() (+13 more)

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

Cohesion: 0.29
Nodes (7): gostructfmt golangci-lint Module Plugin, gostructfmt Linter Configuration, Atomic File Replacement, Command Line Interface, Git Diff Formatting, golangci-lint Plugin, gostructfmt

### Community 7 - "Lint and Format Job"

Cohesion: 0.33
Nodes (6): Development Checks, Pull Request Requirements, Pull Request Template, CI Workflow, Lint and Format Job, Standard golangci-lint Configuration

### Community 9 - "Factors and evidence"

Cohesion: 0.15
Nodes (12): Deterministic and idempotent formatting, Factors and evidence, Functional correctness, Goals, Handoff, Interface compatibility, Interpretation rules, Maintainable responsibility boundaries (+4 more)

### Community 10 - "Implementation Plan"

Cohesion: 0.18
Nodes (10): Acceptance and validation, Change size and staging, Change summary, Completion evidence, Decision gates, Implementation Plan, Implementation steps, Required route decision (+2 more)

## Knowledge Gaps

- **30 isolated node(s):** `github.com/goropikari/gostructfmt`, `Change summary`, `Risk and decision analysis`, `Required route decision`, `Change size and staging` (+25 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions

_Questions this graph is uniquely positioned to answer:_

- **Why does `Parse()` connect `Parse` to `struct_literal.go`, `command.go`, `Plugin`?**
  _High betweenness centrality (0.170) - this node is a cross-community bridge._
- **Why does `formatSource()` connect `command.go` to `struct_literal.go`, `Parse`?**
  _High betweenness centrality (0.116) - this node is a cross-community bridge._
- **Why does `LineRange` connect `command.go` to `struct_literal.go`?**
  _High betweenness centrality (0.077) - this node is a cross-community bridge._
- **Are the 11 inferred relationships involving `Parse()` (e.g. with `.syncSourceFromAST()` and `formatSource()`) actually correct?**
  _`Parse()` has 11 INFERRED edges - model-reasoned connections that need verification._
- **Are the 3 inferred relationships involving `FormatStructLiterals()` (e.g. with `formatSource()` and `Parse()`) actually correct?**
  _`FormatStructLiterals()` has 3 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/goropikari/gostructfmt`, `Change summary`, `Risk and decision analysis` to the rest of the system?**
  _30 weakly-connected nodes found - possible documentation gaps or missing edges._
