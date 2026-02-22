# Architecture

This document describes the high-level architecture of triagegeist: package layout, data flow, design decisions, and extension points. It is intended for contributors, integrators, and researchers who need to understand or extend the library. Nothing is omitted; every package and its role are described.

---

## Overview

triagegeist is a **library** (not a daemon or service). It exposes a pure-Go API for computing a continuous acuity score and a discrete triage level from vital signs and an expected resource count. All core logic is deterministic and parameterised; there are no network calls or external model files by default. The library is organised into a root package and several subpackages, each with a single responsibility.

---

## Package structure (complete)

| Package | Path | Responsibility | Depends on |
|---------|------|----------------|------------|
| **triagegeist** | Root `*.go` | Public API: Engine, Params, Level, FromScore; batch evaluation; presets; validation bridge | score, validate |
| **score** | `score/*.go` | Acuity formula: Vitals struct, deviation, VitalComponent, ResourceComponent, Acuity, AcuityWithNorms; default norms and weights | (none) |
| **norm** | `norm/*.go` | Reference ranges (Ranges), Deviation, NormalizeLinear, ClampToRange, CriticalBounds, WeightedDeviationSum; DefaultRanges, PediatricRanges | (none) |
| **metrics** | `metrics/*.go` | ConfusionMatrix, TP/FP/FN/TN, Sensitivity, Specificity, PPV, NPV, F1, CohenKappa, BinaryCM, AUC, CalibrationError, WeightedKappa | (none) |
| **stats** | `stats/*.go` | Mean, Variance, StdDev, SE, CI95, Median, Percentile, LevelDistribution, ScoreStats, LevelStats, RMSE, MAE, ExactAgreement, WithinLevel | (none) |
| **validate** | `validate/*.go` | Vitals validation (Vitals, ClampVitals, VitalsValid), ResourceCount, Params validation (ParamsLike, Params, ParamsValid), AtLeastOneVital | score |
| **export** | `export/*.go` | Result struct, FromVitalsScoreLevel, ToJSON, CSVHeader, ToCSVRow, WriteCSV, Batch, LevelReport, ReportRow, ComputeSummary, ReadResultJSON, ResultToVitals | score |

**Dependency rule**: No cycles. The root package may import score and validate; score and norm and metrics and stats have no internal project imports; validate imports only score; export imports only score.

---

## Directory layout (full)

```
triagegeist/
├── .github/
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.md
│   │   └── feature_request.md
│   └── PULL_REQUEST_TEMPLATE.md
├── assets/
│   └── README.md
├── docs/
│   ├── README.md
│   ├── ARCHITECTURE.md   (this file)
│   ├── BENCHMARKS.md
│   └── COMPARISON.md
├── examples/
│   ├── README.md
│   ├── basic/
│   │   └── main.go
│   └── advanced/
│       └── main.go
├── norm/
│   ├── norm.go
│   └── norm_test.go
├── score/
│   ├── score.go
│   └── score_test.go
├── metrics/
│   ├── metrics.go
│   └── metrics_test.go
├── stats/
│   ├── stats.go
│   └── stats_test.go
├── validate/
│   ├── validate.go
│   └── validate_test.go
├── export/
│   ├── export.go
│   └── export_test.go
├── doc.go
├── params.go
├── params_validate.go
├── level.go
├── engine.go
├── engine_test.go
├── example_test.go
├── go.mod
├── LICENSE
├── NOTICE
├── CHANGELOG.md
├── CONTRIBUTING.md
├── SECURITY.md
├── CODE_OF_CONDUCT.md
├── GOVERNANCE.md
└── README.md
```

---

## Data flow (step-by-step)

1. **Input**
   - Caller provides: `score.Vitals` (HR, RR, SBP, DBP, Temp, SpO2, GCS; 0 or missing means not present) and `resourceCount` (non-negative integer).
   - Engine holds a copy of `Params` (VitalWeights, MaxResources, ResourceWeight, T1, T2, T3, T4).

2. **Optional validation**
   - Caller may use `validate.Vitals(v)` to obtain a report; `validate.ClampVitals(v)` to clamp out-of-range values; `validate.ResourceCount(count, maxResources)` to clamp resource count.
   - Root package offers `ValidateParamsExternal(p)` which delegates to `validate.Params` with a ParamsLike struct.

3. **Vital component**
   - For each present vital (e.g. HR > 0), compute deviation \( d_i = \min(1, |x_i - \mu_i| / \sigma_i) \) using either score package default norms or custom norms via `score.VitalComponentWithNorms`.
   - Weighted sum over present vitals, normalised by the sum of weights of present vitals, yields \( V \in [0,1] \).

4. **Resource component**
   - \( R = \alpha \cdot \min(1, \texttt{resourceCount} / \texttt{maxResources}) \).

5. **Raw and normalised score**
   - \( \text{raw} = V + R \), \( s = \text{raw} / (\sum_i w_i + \alpha) \), then \( s \) is clamped to \( [0,1] \).

6. **Level**
   - Compare \( s \) to thresholds \( T_1 > T_2 > T_3 > T_4 \): Level 1 if \( s \geq T_1 \), Level 2 if \( T_2 \leq s < T_1 \), etc., Level 5 if \( s < T_4 \).

7. **Output**
   - `Engine.Acuity` returns \( s \); `Engine.Level` returns \( L \); `Engine.ScoreAndLevel` returns both. Batch methods (`BatchScoreAndLevel`, `BatchEvaluate`) repeat this for slices.

8. **Optional export and metrics**
   - Use `export.FromVitalsScoreLevel` to build a `Result` for JSON/CSV; `export.WriteCSV` for batch CSV; `metrics.NewConfusionMatrix(pred, ref)` when reference levels are available; `stats.ComputeScoreStats`, `stats.ComputeLevelStats` for aggregates.

No persistent state is modified inside the library. Engine is safe for concurrent use provided Params is not mutated during calls.

---

## Design decisions

| Decision | Rationale |
|----------|------------|
| **Pure Go, no cgo** | Portability, easy cross-compilation, no C toolchain dependency. |
| **No default external dependencies** | Keeps the core small; avoids supply-chain and versioning friction. |
| **Parametric formula only** | Reproducibility, auditability, low latency. ML can be plugged via wrappers if needed. |
| **Vitals as struct** | Predictable memory layout, no allocations in hot path, explicit typing. |
| **Separate score package** | Formula is independent and testable without the full API. |
| **norm package** | Reference ranges and deviation helpers reusable outside score (e.g. validation, custom formulae). |
| **metrics and stats** | Accuracy and descriptive statistics are separate so callers can choose what to use. |
| **validate package** | Input validation and clamping are explicit; root package can stay thin and delegate. |
| **export package** | Serialisation (JSON, CSV) and reporting (LevelReport, Summary) in one place. |

---

## Extension points

- **Custom parameters**: Set `Params` (weights, thresholds, maxResources, resourceWeight) and pass to `NewEngine`. Use `PresetStrict`, `PresetLenient`, `PresetResearch` or build from `DefaultParams()` and override.
- **Custom norms**: Use `score.VitalComponentWithNorms` and `score.AcuityWithNorms` with a `[7][2]float64` norms array, or use `norm.Ranges` and `norm.WeightedDeviationSum` for custom aggregation.
- **External predictors**: Implement a type that takes vitals (and optionally resource count) and returns a score; then use `FromScore(score, params)` to map to level. The library does not depend on any external model runtime.
- **Validation**: Use `validate` before calling the engine; use `ValidateParamsExternal` in the root package to check Params with the same logic as `validate.Params`.
- **Benchmarks**: See [BENCHMARKS.md](BENCHMARKS.md). Add new benchmarks in the appropriate `*_test.go` and document in that file.

---

## Testing strategy

| Package | What is tested |
|---------|----------------|
| **triagegeist** | Engine Acuity/Level/ScoreAndLevel, FromScore boundaries, Params.Validate, batch helpers, example tests. |
| **score** | VitalComponent, Acuity, Normalize, default behaviour. |
| **norm** | DefaultRanges, Deviation, NormalizeLinear, ClampToRange, At/Set, CriticalBounds, WeightedDeviationSum, Valid. |
| **metrics** | NewConfusionMatrix, TP/FP/FN/TN, Sensitivity/Specificity, perfect agreement, BinaryCM, AUC, CalibrationError, WeightedKappa. |
| **stats** | Mean, Variance, StdDev, CI95, Median, Percentile, LevelDistribution, ComputeScoreStats, ExactAgreement, RMSE. |
| **validate** | Vitals report, ClampVitals, ResourceCount, Params report, AtLeastOneVital. |
| **export** | FromVitalsScoreLevel, ToCSVRow, ToJSON, LevelReport, ComputeSummary, ResultToVitals, WriteCSV. |

Run all tests: `go test ./... -count=1`. Examples (examples/basic, examples/advanced) do not contain `*_test.go` but can be run with `go run ./examples/basic` and `go run ./examples/advanced`.

---

## Conventions

- **Exported names**: Use clear, consistent names; avoid abbreviations except standard ones (HR, RR, SBP, DBP, GCS, etc.).
- **Errors**: The core API does not return errors; invalid inputs (e.g. out-of-range vitals) are handled by the validate package or by clamping in the caller. Functions that do I/O (e.g. export.WriteCSV) return errors.
- **Immutability**: Params and Ranges are copied by value; Engine holds a copy of Params and does not mutate it. Callers should not mutate Params while an Engine is in use from multiple goroutines.
- **Documentation**: Every exported symbol has a doc comment. Package doc includes tables and formula references where helpful.
