# Examples

This directory contains example programs that demonstrate triagegeist usage for education, research, and integration. Each example is self-contained and can be run without modifying the library. Nothing is omitted below; every step and requirement is stated explicitly.

---

## Structure

| Directory | Purpose |
|-----------|---------|
| [basic/](basic/) | Single evaluation: build params and engine, validate vitals, compute acuity and level, build an export result. Use this to understand the minimal workflow. |
| [advanced/](advanced/) | Batch evaluation: validate and clamp inputs, run batch scoring, compute descriptive statistics (mean, CI95, level distribution), confusion matrix and metrics (accuracy, kappa, sensitivity/specificity, binary high/low acuity), level report, export summary, and write CSV to stdout. Use this for research or auditing pipelines. |

---

## Requirements

- **Go**: 1.22 or later. Check with `go version`.
- **Module**: The examples are part of the triagegeist module. From the repository root, run `go mod tidy` once so that the module graph is resolved. No separate `go get` is needed when running from the repo.
- **Platform**: Any platform supported by Go (Linux, Windows, macOS, etc.). The examples do not use OS-specific code.

---

## Running

**From the repository root (recommended):**

```bash
go run ./examples/basic
go run ./examples/advanced
```

**From within an example directory:**

```bash
cd examples/basic
go run .
```

```bash
cd examples/advanced
go run .
```

**Expected output (summary):**

- **basic**: Prints acuity (e.g. 0.7465), level (e.g. 2 Emergent), and wait time (e.g. 15 min). No CSV or file output.
- **advanced**: Prints acuity statistics (\(N\), mean \( \bar{x} \), std \( \sigma \), \(95\%\ \mathrm{CI}\), min, max, percentiles), level distribution (counts and percentages per level \( L \in \{1,\ldots,5\} \)), agreement metrics (overall accuracy, Cohen's \( \kappa \), macro sensitivity/specificity), binary metrics (sensitivity, specificity, PPV, NPV, F1), weighted kappa, exact and within-one-level agreement, export summary, level report, and then a CSV table of all results to stdout.

---

## Learning path

1. **basic**: Run one evaluation with default params. Note how vitals and resource count are passed, how validation is used (Vitals report, ClampVitals, ResourceCount), and how the result is passed to export.FromVitalsScoreLevel. This is the minimal integration pattern.
2. **advanced**: Run batch evaluations with synthetic data. See how to use Engine.BatchScoreAndLevel, stats (e.g. \( \bar{x} \), \( \mathrm{CI}_{95\%} \), level distribution), metrics (confusion matrix, \( \kappa \), sensitivity/specificity), export.LevelReport, export.ComputeSummary, and export.WriteCSV. This mirrors a research or audit workflow.

---

## Disclaimer

Examples use synthetic or illustrative data. They are not a substitute for clinical protocols or regulatory compliance. Validate and calibrate parameters for your setting. The library is for research and decision support only; see the root [README](../README.md) and [LICENSE](../LICENSE) for full terms.
