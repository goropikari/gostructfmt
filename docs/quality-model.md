# Quality Model

This document defines how quality is assessed for `gostructfmt`. It describes
the properties to observe and the evidence that supports an assessment; it does
not prescribe thresholds, required checks, severities, or release decisions.

## Product boundary

The assessed product is the Go formatter CLI, its formatting library, and the
golangci-lint module plugin. The boundary includes source read from standard
input, explicitly named files, recursive file discovery, Git-diff selection,
atomic writes, and diagnostic output. Build tooling and CI are instruments for
evidence, rather than product-quality factors themselves.

## Goals

| ID | Goal                                                                                                                                                 | Questions                                                                                                                 | Metrics                                                                                           |
| -- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| G1 | For formatter users, preserve valid Go source while making populated struct literals deterministic and gofmt-compatible in CLI and library contexts. | Does formatting preserve parseable Go and only alter supported literals? Is equivalent input formatted consistently?      | Contract-test outcomes; parse/format success and error outcomes; regression findings.             |
| G2 | For users modifying files, prevent partial or unintended filesystem changes when input, selection, or writing fails.                                 | Are only intended Go files selected? Are writes atomic and withheld after a failure?                                      | Integration-test outcomes for file discovery, diff selection, and write failure; review findings. |
| G3 | For plugin users and maintainers, keep CLI and golangci-lint integrations compatible with their documented contracts.                                | Do flags, exit status, diagnostics, and analyzer output remain usable by callers?                                         | End-to-end and integration-test outcomes; API/compatibility review findings.                      |
| G4 | For maintainers, keep changes understandable and safely verifiable.                                                                                  | Are responsibilities and error paths locally understandable? Do tests cover observable behavior and remain deterministic? | Semantic review findings; test results; test-design review findings.                              |

Metrics in this table estimate source and integration risk. They do not by
themselves prove runtime reliability, user satisfaction, or delivery readiness.

## Factors and evidence

### Functional correctness

```text
reliability
  -> Go-syntax preservation and scoped struct-literal transformation
  -> documented input/output and error contracts hold
  -> Go tests and focused CLI integration tests
  -> passing results, or a finding that names the input, observed output, and impact
```

| Field              | Definition                                                                                                                                                |
| ------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Factor             | Valid source remains valid and the formatter changes only populated struct literals in its supported scope.                                               |
| Measures           | Do representative valid and invalid inputs have the documented output, diagnostics, and exit status? Are empty, slice, array, and map literals unchanged? |
| Instruments        | `go test ./...`; formatter and CLI contract tests; `blackbox-risk-based-test` when a change introduces new user-visible behavior.                         |
| Evidence           | Test command output and, for failures, the exact input scenario, observed result, source location, and user impact.                                       |
| Limitations        | Example tests cannot prove correctness for every Go syntax form; fuzzing or parser differential testing would be separate evidence.                       |
| Applicability      | Formatter library and CLI behavior.                                                                                                                       |
| Default confidence | High for covered contracts; medium for untested syntax combinations.                                                                                      |

### Deterministic and idempotent formatting

```text
reliability
  -> deterministic formatting for the same source and selection
  -> repeated formatting produces the same bytes after the first successful run
  -> focused repeat-run tests and regression tests
  -> first and repeated outputs, including any diagnostic difference
```

| Field              | Definition                                                                                                                       |
| ------------------ | -------------------------------------------------------------------------------------------------------------------------------- |
| Factor             | Formatting is stable across repeated invocations with the same input and options.                                                |
| Measures           | Does a successful second run leave output unchanged? Do equivalent supported inputs produce consistent formatting?               |
| Instruments        | Unit and CLI tests; `go test ./...`; `property-based-test-review` when input variation makes example coverage insufficient.      |
| Evidence           | First and second output or file contents, plus test output.                                                                      |
| Limitations        | Does not establish deterministic behavior across all Go toolchain versions or operating systems without a representative matrix. |
| Applicability      | Formatter library, standard-input flow, and file-writing flow.                                                                   |
| Default confidence | Medium.                                                                                                                          |

### Safe file selection and recovery

```text
reliability
  -> controlled file targeting and recoverable writes
  -> paths outside the documented selection are skipped; failed operations do not leave partial output
  -> integration tests, source inspection, and runtime failure observation
  -> selected paths, before/after contents, errors, and cleanup results
```

| Field              | Definition                                                                                                                                     |
| ------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| Factor             | Recursive, explicit-file, and Git-diff operations target only intended files; file replacement is atomic after successful formatting.          |
| Measures           | Are hidden, vendor, testdata, and underscore-prefixed paths skipped as documented? Does a parse or write failure preserve prior file contents? |
| Instruments        | CLI integration tests using `t.TempDir()`; `go test ./...`; `data-integrity-review` for changes to write or selection behavior.                |
| Evidence           | Test fixture tree, selected-file list where observable, before/after content, returned error, and temporary-file cleanup result.               |
| Limitations        | Filesystem fault injection and platform-specific atomic-replace semantics may need dedicated environment coverage.                             |
| Applicability      | `-w`, recursive paths, and `--diff`.                                                                                                           |
| Default confidence | Medium.                                                                                                                                        |

### Interface compatibility

```text
maintainability
  -> explicit and stable caller-facing contracts
  -> documented CLI flags, standard streams, exit behavior, and plugin diagnostics remain compatible
  -> contract tests and API compatibility review
  -> command output/exit status and review findings tied to a caller or configuration
```

| Field              | Definition                                                                                                                         |
| ------------------ | ---------------------------------------------------------------------------------------------------------------------------------- |
| Factor             | Users can continue to invoke the CLI and configure the linter plugin according to the documented interface.                        |
| Measures           | Do standard input/output, `-w`, multiple-file restrictions, `--diff`, diagnostics, and plugin configuration retain their contract? |
| Instruments        | CLI and plugin tests; `api-compatibility-review` for changed flags, diagnostics, configuration, or exported APIs.                  |
| Evidence           | Test output, configuration fixture result, and review finding with the affected caller contract.                                   |
| Limitations        | Undocumented downstream integrations are not discoverable solely from repository evidence.                                         |
| Applicability      | `cmd/gostructfmt`, `plugin/gostructfmt`, README configuration examples.                                                            |
| Default confidence | Medium.                                                                                                                            |

### Maintainable responsibility boundaries

```text
maintainability
  -> clear separation of parsing, transformation, CLI orchestration, and plugin analysis
  -> a change has a localized reason and dependencies flow through intentional boundaries
  -> architecture and readability review
  -> finding with file, symbol, coupled concerns, and predicted change impact
```

| Field              | Definition                                                                                                                            |
| ------------------ | ------------------------------------------------------------------------------------------------------------------------------------- |
| Factor             | Parsing, source transformation, file operations, CLI flow, and lint analysis remain independently understandable.                     |
| Measures           | Are unrelated reasons to change coupled in one symbol? Do package dependencies preserve the intended formatter/CLI/plugin boundaries? |
| Instruments        | `architecture-review`; `readability-review`; compilation and tests as supporting, not conclusive, evidence.                           |
| Evidence           | Review findings naming the file, symbol, behavior, dependency direction, and concrete maintenance impact.                             |
| Limitations        | Requires reviewer judgment and does not prove runtime behavior.                                                                       |
| Applicability      | Production Go code.                                                                                                                   |
| Default confidence | Medium.                                                                                                                               |

### Verifiable behavior

```text
maintainability
  -> deterministic tests that assert observable contracts
  -> changed behavior has focused success, failure, and boundary coverage
  -> Go tests and test-quality review
  -> test results plus findings about missing behaviors or brittle assertions
```

| Field              | Definition                                                                                                                                                            |
| ------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Factor             | Tests make user-visible behavior, error handling, and filesystem safety regressions detectable without depending on execution order or hidden implementation details. |
| Measures           | Do tests follow the repository's testify, AAA, and `t.Run` conventions? Do they cover relevant normal, failure, and boundary cases?                                   |
| Instruments        | `go test ./...`; `test-quality-review`; `blackbox-risk-based-test` for acceptance coverage; `whitebox-risk-based-test` for high-risk internal paths.                  |
| Evidence           | Test output and review findings that identify the missing or brittle behavior, test location, and residual risk.                                                      |
| Limitations        | Passing tests demonstrate only the cases exercised; test count and coverage percentages are not substitutes for risk coverage.                                        |
| Applicability      | All tests; priority is behavior changed by the current work.                                                                                                          |
| Default confidence | High for executed scenarios; medium for omitted risks.                                                                                                                |

## Interpretation rules

- Record unavailable evidence as **not measured**, never as a pass.
- Keep raw leaf evidence (commands, scenarios, findings, affected symbols) so a
  quality conclusion can be traced to an improvement action.
- Treat a critical safety or compatibility finding as visible on its own; do
  not let an average score conceal it.
- A clean formatter, lint, or test run supports only the factors it observes.
  It does not replace semantic architecture, compatibility, or safety review.
- Revisit the model when product boundaries change, such as adding new input
  sources, a new write mode, a public API, or another plugin integration.

## Handoff

An enforceable repository policy belongs in a separate quality baseline. That
baseline should select applicable factors from this model, assign required or
critical status, set thresholds only with recorded rationale, and register the
commands and review instruments that produce evidence.
