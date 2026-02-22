# Contributing to triagegeist

Thank you for your interest in contributing. This document describes the project's workflow, quality standards, and legal expectations. The project follows Google-style open source practices where applicable and is licensed under the European Union Public Licence v. 1.2 (EUPL-1.2).

---

## Table of contents

1. [Licence and legal](#licence-and-legal)
2. [Code of conduct](#code-of-conduct)
3. [Getting started](#getting-started)
4. [Development workflow](#development-workflow)
5. [Code and documentation standards](#code-and-documentation-standards)
6. [Testing and benchmarking](#testing-and-benchmarking)
7. [Pull request process](#pull-request-process)
8. [Scope and priorities](#scope-and-priorities)

---

## Licence and legal

By contributing to triagegeist, you agree that your contributions will be licensed under the **European Union Public Licence v. 1.2 (EUPL-1.2)**. You represent that you have the right to license your work under these terms. The project does not require a separate Contributor License Agreement (CLA); the act of submitting a pull request constitutes acceptance of these terms.

| Document | Purpose |
|----------|---------|
| [LICENSE](LICENSE) | Full EUPL-1.2 text |
| [GOVERNANCE.md](GOVERNANCE.md) | Maintainers and decision-making |

---

## Code of conduct

This project adheres to a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold it. Please report unacceptable behaviour as described in that document.

---

## Getting started

### Prerequisites

| Tool | Minimum version | Notes |
|------|-----------------|-------|
| Go  | 1.22           | See [go.mod](go.mod) |
| Git | 2.x            | For clone, branch, push |

### Clone and build

```bash
git clone https://github.com/olaflaitinen/triagegeist.git
cd triagegeist
go mod download
go build ./...
go test ./...
```

### Repository layout (high level)

| Path | Purpose |
|------|---------|
| `*.go` (root) | Core package `triagegeist`: engine, params, level, doc, params_validate |
| `score/*.go` | Subpackage `score`: vitals, acuity formula, normalisation |
| `norm/*.go` | Reference ranges, deviation, normalisation helpers |
| `metrics/*.go` | Confusion matrix, sensitivity, specificity, kappa, AUC, calibration |
| `stats/*.go` | Mean, CI, percentiles, level distribution, RMSE, MAE |
| `validate/*.go` | Vitals and params validation, clamping |
| `export/*.go` | Result struct, JSON/CSV export, level report, summary |
| `docs/` | ARCHITECTURE.md, BENCHMARKS.md, COMPARISON.md, README.md |
| `examples/` | basic/, advanced/ example programs |
| `assets/` | Logo (SVG, 8000x2000, no background) |
| `.github/` | Issue and pull request templates |

---

## Development workflow

1. **Open or find an issue**  
   Check [open issues](https://github.com/olaflaitinen/triagegeist/issues). If your change is non-trivial, open an issue first to align with maintainers.

2. **Fork and branch**  
   Fork the repository and create a branch from `main`. Use a short, descriptive branch name (e.g. `fix/validate-thresholds`, `docs/architecture`).

3. **Implement and test**  
   Make your changes, add or update tests, and ensure all tests and linters pass (see below).

4. **Commit**  
   Write clear commit messages. Prefer present tense and one logical change per commit.

5. **Push and open a pull request**  
   Push your branch and open a PR against `main`. Fill in the PR template and reference any related issues.

---

## Code and documentation standards

### Style and formatting

| Rule | How to enforce |
|------|-----------------|
| Formatting | `gofmt -s -w .` or equivalent; the project uses standard Go formatting |
| Imports | Group: standard library, then third-party, then project imports; use `goimports` if available |
| Line length | Prefer readability; avoid excessively long lines |
| Naming | Follow [Effective Go](https://go.dev/doc/effective_go): mixedCaps, short names for scope, no redundant type in name |

### Package design

| Principle | Application in triagegeist |
|-----------|----------------------------|
| Small, focused packages | Root package: API and types; `score`: formula and vitals only |
| Minimal dependencies | No external dependencies in core; avoid new deps unless justified |
| Exported API stability | Avoid breaking changes to exported names/signatures without a major version or clear deprecation |

### Documentation

- Every exported symbol (type, function, method, constant) must have a doc comment starting with the symbol name.
- Use full sentences and, where helpful, tables or mathematical notation (e.g. in `doc.go`).
- Keep comments accurate when behaviour or parameters change.

### Mathematical notation (in comments and docs)

When referring to formulas in code or package docs, use consistent notation:

| Concept | Notation example |
|---------|-------------------|
| Normalised score | \( s \in [0,1] \) |
| Level | \( L \in \{1,\ldots,5\} \) |
| Weights | \( w_i \), \( \alpha \) |
| Thresholds | \( T_1, T_2, T_3, T_4 \) |

---

## Testing and benchmarking

### Tests

- **Unit tests**: Cover new and changed behaviour. Place tests in `*_test.go` in the same package (or `triagegeist_test` for examples).
- **Run all tests**: `go test ./... -count=1`
- **Run tests with verbose output**: `go test ./... -v -count=1`
- **Coverage** (optional): `go test -cover ./...` or `go test -coverprofile=coverage.out ./...`

**Packages with tests** (all must pass before a PR is merged):

| Package | Test file | What is tested |
|---------|-----------|----------------|
| triagegeist | engine_test.go, example_test.go | Engine Acuity/Level/ScoreAndLevel, FromScore, Params.Validate, benchmarks, examples |
| score | score_test.go | VitalComponent, Acuity, Normalize |
| norm | norm_test.go | DefaultRanges, Deviation, NormalizeLinear, ClampToRange, At/Set, CriticalBounds, WeightedDeviationSum |
| metrics | metrics_test.go | ConfusionMatrix, TP/FP/FN/TN, Sensitivity/Specificity, BinaryCM, AUC, CalibrationError, WeightedKappa |
| stats | stats_test.go | Mean, Variance, StdDev, CI95, Median, Percentile, LevelDistribution, ComputeScoreStats, ExactAgreement, RMSE |
| validate | validate_test.go | Vitals report, ClampVitals, ResourceCount, Params report, AtLeastOneVital |
| export | export_test.go | FromVitalsScoreLevel, ToCSVRow, ToJSON, LevelReport, ComputeSummary, ResultToVitals, WriteCSV |

Examples (examples/basic, examples/advanced) have no `*_test.go` but must run successfully: `go run ./examples/basic` and `go run ./examples/advanced`.

### Benchmarks

- Benchmark code lives in `*_test.go` with functions of the form `BenchmarkXxx(b *testing.B)`.
- Run benchmarks: `go test -bench=. -benchmem ./...`
- Do not commit large binary or generated data; keep benchmarks reproducible and fast enough for CI.

### CI

The project expects that `go build ./...` and `go test ./...` succeed on the supported Go version. PRs should maintain or improve test coverage and not regress benchmarks without justification.

---

## Pull request process

1. **Target branch**  
   All PRs target `main` (or the default branch configured in the repository).

2. **Checks**  
   Before requesting review, ensure:
   - `go build ./...` succeeds
   - `go test ./...` passes
   - New/changed code is documented and formatted

3. **Review**  
   At least one maintainer (see [GOVERNANCE.md](GOVERNANCE.md)) will review. Address feedback by pushing new commits to the same branch.

4. **Merge**  
   Maintainers merge when the PR is approved and CI (if configured) is green. The project may use squash or merge commits depending on repository settings.

5. **Changelog**  
   For user-visible changes, add an entry to [CHANGELOG.md](CHANGELOG.md) under an appropriate version (or "Unreleased").

---

## Scope and priorities

The project focuses on:

- **Parametric triage and acuity scoring** in emergency medicine (formula-based, configurable).
- **Performance and portability**: pure Go, minimal allocations, no mandatory cgo or heavy runtimes.
- **Clarity and maintainability**: clear APIs, good documentation, and alignment with Go and open source best practices.

Out of scope for the core library (unless explicitly agreed via issue/PR):

- Replacement of clinical judgment or institutional protocols.
- Binding to proprietary or non-redistributable triage algorithms.
- Large binary assets or heavy ML frameworks as default dependencies.

If you are unsure whether a change fits the project scope, open an issue before investing in a large PR.

---

## Summary checklist for contributors

- [ ] Read [LICENSE](LICENSE) and [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)
- [ ] Fork, branch from `main`, implement changes
- [ ] Add/update tests; run `go test ./...`
- [ ] Format with `gofmt`; ensure `go build ./...` passes
- [ ] Update documentation for new or changed behaviour
- [ ] Update [CHANGELOG.md](CHANGELOG.md) for user-visible changes
- [ ] Open a PR with a clear description and reference to any issue
