# triagegeist

<div align="center">

<img src="assets/logo.svg" width="320" alt="triagegeist logo" />

</div>

High-performance, parametric AI toolkit for emergency medicine triage and acuity scoring in Go. Designed for minimal latency, zero unnecessary allocations, and state-of-the-art results in its domain without heavy runtime or hardware requirements. Suitable for education, research, and production decision support.

| Attribute   | Value |
|-------------|-------|
| Authors     | Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov |
| License     | European Union Public Licence v. 1.2 (EUPL-1.2) |
| Repository  | [github.com/olaflaitinen/triagegeist](https://github.com/olaflaitinen/triagegeist) |
| Documentation | [pkg.go.dev/github.com/olaflaitinen/triagegeist](https://pkg.go.dev/github.com/olaflaitinen/triagegeist) |
| Go version  | 1.22+ |

Logo: SVG, transparent background, 8000 x 2000 px. Place **assets/logo.svg** (or **assets/triagegeist-logo.svg** and update the `img` src above). See [assets/README.md](assets/README.md).

---

## Table of contents

1. [Overview](#overview)
2. [Repository structure](#repository-structure)
3. [Mathematical model](#mathematical-model)
4. [Subpackages](#subpackages)
5. [Installation and usage](#installation-and-usage)
6. [Examples](#examples)
7. [Metrics and accuracy](#metrics-and-accuracy)
8. [Documentation and standards](#documentation-and-standards)
9. [Contributing and governance](#contributing-and-governance)
10. [Disclaimer](#disclaimer)

---

## Overview

triagegeist computes a **normalised acuity score** \( s \in [0,1] \) and a **discrete triage level** \( L \in \{1,2,3,4,5\} \) from vital signs and expected resource consumption. The implementation is fully parametric and deterministic: no external model files or heavy ML runtimes are required. The design prioritises speed, low memory use, and suitability for embedded or server-side deployment while remaining aligned with established emergency triage concepts.

### Design goals

| Goal           | Approach |
|----------------|----------|
| Speed          | Pure Go, no cgo by default; formula-based scoring; minimal allocations in hot path |
| Determinism    | All thresholds and weights are explicit parameters; reproducible across runs |
| Flexibility    | Pluggable parameters; optional external predictors via interfaces |
| Documentation  | Inline doc comments, tables, LaTeX-style formulas; pkg.go.dev and docs aligned |
| Standards      | Google OSS-style layout; EUPL-1.2; governance and contribution process |
| Education      | Examples, metrics, stats, and export for teaching and research |

### Use cases

| Use case        | How triagegeist helps |
|-----------------|------------------------|
| ED triage support | Single or batch evaluation; configurable thresholds and weights |
| Research         | Metrics (sensitivity, specificity, kappa, AUC); stats (CI, percentiles); export (JSON, CSV) |
| Education        | Clear formulas, examples, and subpackages (norm, validate, export) |
| Auditing         | Export results to CSV/JSON; level reports and summary stats |

---

## Repository structure

Layout follows common Go and GitHub conventions:

```
triagegeist/
├── .github/                    # Issue and PR templates
│   ├── ISSUE_TEMPLATE/
│   └── PULL_REQUEST_TEMPLATE.md
├── assets/                     # Logo (SVG, 8000x2000, no background)
│   └── README.md
├── docs/                       # Design and reference documentation
│   ├── README.md
│   ├── ARCHITECTURE.md
│   ├── BENCHMARKS.md
│   └── COMPARISON.md
├── examples/
│   ├── README.md
│   ├── basic/                  # Single evaluation, validation, export
│   └── advanced/               # Batch, metrics, stats, CSV/JSON
├── norm/                       # Reference ranges, deviation, normalisation
├── score/                      # Acuity formula, Vitals, vital/resource components
├── metrics/                    # Sensitivity, specificity, AUC, kappa, calibration
├── stats/                      # Mean, CI, percentiles, aggregation
├── validate/                   # Input validation, clamping
├── export/                     # JSON, CSV, batch export, level reports
├── doc.go, params.go, level.go, engine.go   # Core API
├── LICENSE, NOTICE, CHANGELOG.md
├── CONTRIBUTING.md, SECURITY.md, CODE_OF_CONDUCT.md, GOVERNANCE.md
└── README.md
```

---

## Mathematical model

### Acuity score

The raw acuity aggregate is a weighted combination of vital-sign deviation and resource count, then normalised to the unit interval.

**Vital component**

For each vital \( i \) with observed value \( x_i \), reference midpoint \( \mu_i \), and half-width \( \sigma_i \):

$$
d_i = \min\left(1,\; \frac{|x_i - \mu_i|}{\sigma_i}\right)
$$

Only **present** vitals (e.g. \( x_i \neq 0 \) where 0 denotes missing) are included. The vital component \( V \) is:

$$
V = \frac{\sum_{i \in \mathcal{I}} w_i \, d_i}{\sum_{i \in \mathcal{I}} w_i}
$$

where \( \mathcal{I} \) is the set of indices with present values and \( w_i \geq 0 \) are configurable weights.

**Resource component**

$$
R = \alpha \cdot \min\left(1,\; \frac{\texttt{resourceCount}}{\texttt{maxResources}}\right)
$$

with \( \alpha \) = `resourceWeight` and `maxResources` the cap on expected resources.

**Normalised score**

$$
\text{raw} = V + R,\qquad
s = \frac{\text{raw}}{\sum_i w_i + \alpha},\qquad
s \in [0,1]
$$

Implementations clamp \( s \) to \( [0,1] \) when necessary.

### Default parameters

| Symbol / name   | Meaning                  | Default value |
|-----------------|--------------------------|---------------|
| \( w_{\text{HR}} \)  | Heart rate weight        | 0.18 |
| \( w_{\text{RR}} \)  | Respiratory rate weight  | 0.22 |
| \( w_{\text{SBP}} \) | Systolic BP weight       | 0.16 |
| \( w_{\text{DBP}} \) | Diastolic BP weight      | 0.10 |
| \( w_{\text{Temp}} \) | Temperature weight     | 0.08 |
| \( w_{\text{SpO2}} \) | Oxygen saturation weight | 0.16 |
| \( w_{\text{GCS}} \)  | Glasgow Coma Scale weight | 0.10 |
| \( \alpha \)     | Resource weight          | 0.25 |
| maxResources    | Cap on resource count    | 6 |

**Reference ranges (mid \( \mu \), half-width \( \sigma \))**

| Vital | \( \mu \) | \( \sigma \) | Unit |
|-------|-----------|--------------|------|
| HR    | 80        | 40           | bpm  |
| RR    | 16        | 10           | /min |
| SBP   | 120       | 40           | mmHg |
| DBP   | 80        | 30           | mmHg |
| Temp  | 37.0      | 2.0          | °C   |
| SpO2  | 98        | 8            | %    |
| GCS   | 15        | 6            | 3–15 |

### Level assignment

Discrete level \( L \) is obtained by thresholding \( s \) with four cutpoints \( T_1 > T_2 > T_3 > T_4 \):

| Level \( L \) | Label          | Condition              | Typical wait (guidance) |
|---------------|----------------|------------------------|--------------------------|
| 1             | Resuscitation  | \( s \geq T_1 \)       | Immediate                |
| 2             | Emergent       | \( T_2 \leq s < T_1 \) | &lt; 15 min            |
| 3             | Urgent         | \( T_3 \leq s < T_2 \) | &lt; 60 min            |
| 4             | Less urgent    | \( T_4 \leq s < T_3 \) | &lt; 120 min           |
| 5             | Non-urgent     | \( s < T_4 \)          | &lt; 240 min           |

Default thresholds: \( T_1 = 0.85 \), \( T_2 = 0.60 \), \( T_3 = 0.35 \), \( T_4 = 0.15 \). All are configurable via [Params](https://pkg.go.dev/github.com/olaflaitinen/triagegeist#Params).

---

## Subpackages

| Package   | Purpose |
|-----------|---------|
| [score](https://pkg.go.dev/github.com/olaflaitinen/triagegeist/score)   | Acuity formula, Vitals struct, vital/resource components, normalisation |
| [norm](https://pkg.go.dev/github.com/olaflaitinen/triagegeist/norm)     | Reference ranges (Ranges), Deviation, normalisation helpers |
| [metrics](https://pkg.go.dev/github.com/olaflaitinen/triagegeist/metrics) | Confusion matrix, sensitivity, specificity, PPV, NPV, F1, Cohen's kappa, AUC, calibration |
| [stats](https://pkg.go.dev/github.com/olaflaitinen/triagegeist/stats)   | Mean, variance, StdDev, SE, CI95, median, percentiles, level distribution, RMSE, MAE |
| [validate](https://pkg.go.dev/github.com/olaflaitinen/triagegeist/validate) | Vitals validation, clamping, Params validation, resource count clamp |
| [export](https://pkg.go.dev/github.com/olaflaitinen/triagegeist/export) | Result struct, JSON/CSV write, batch export, level report, summary |

---

## Installation and usage

Requires **Go 1.22 or later**.

```bash
go get github.com/olaflaitinen/triagegeist
```

**Basic: default parameters**

```go
package main

import (
	"fmt"
	"github.com/olaflaitinen/triagegeist"
	"github.com/olaflaitinen/triagegeist/score"
)

func main() {
	p := triagegeist.DefaultParams()
	eng := triagegeist.NewEngine(p)

	v := score.Vitals{HR: 120, RR: 24, SBP: 90, SpO2: 92}
	resourceCount := 3

	acuity, level := eng.ScoreAndLevel(v, resourceCount)
	fmt.Printf("acuity: %.3f, level: %d (%s)\n", acuity, level, level.String())
}
```

**Custom parameters**

```go
	p := triagegeist.DefaultParams()
	p.T1, p.T2 = 0.90, 0.65
	if !p.Validate() { return }
	eng := triagegeist.NewEngine(p)
```

**Batch evaluation**

```go
	acuities, levels := eng.BatchScoreAndLevel(vitals, resourceCounts)
```

**Validation and export**

```go
	report := validate.Vitals(v)
	if !report.Valid { v = validate.ClampVitals(v) }
	res := export.FromVitalsScoreLevel(v, resourceCount, acuity, level.Int(), level.String())
```

---

## Examples

| Example   | Path              | Description |
|-----------|-------------------|-------------|
| Basic     | [examples/basic](examples/basic)   | Single evaluation, validation, export struct |
| Advanced  | [examples/advanced](examples/advanced) | Batch evaluation, metrics, stats, CSV output |

Run from repository root:

```bash
go run ./examples/basic
go run ./examples/advanced
```

See [examples/README.md](examples/README.md) for a learning path and requirements.

---

## Metrics and accuracy

When reference (ground truth) levels are available, use the **metrics** package:

| Metric        | Use |
|---------------|-----|
| ConfusionMatrix | 5x5 counts; TP, FP, FN, TN per class |
| Sensitivity, Specificity | Per level or binary |
| PPV, NPV, F1  | Positive/negative predictive value, F1 |
| CohenKappa    | Agreement vs chance |
| WeightedKappa | Adjacent-level agreement |
| AUC           | From scores and binary outcomes |
| CalibrationError | Score vs outcome calibration |

Use **stats** for descriptive statistics (mean, CI95, percentiles, level distribution, RMSE, MAE, exact/within-level agreement). See [docs/BENCHMARKS.md](docs/BENCHMARKS.md) and [docs/COMPARISON.md](docs/COMPARISON.md) for performance and comparison with alternatives.

---

## Documentation and standards

| Resource   | Location |
|------------|----------|
| Package API | [pkg.go.dev/github.com/olaflaitinen/triagegeist](https://pkg.go.dev/github.com/olaflaitinen/triagegeist) |
| Architecture | [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) |
| Benchmarks   | [docs/BENCHMARKS.md](docs/BENCHMARKS.md) |
| Comparison   | [docs/COMPARISON.md](docs/COMPARISON.md) |
| Changelog    | [CHANGELOG.md](CHANGELOG.md) |
| Contributing | [CONTRIBUTING.md](CONTRIBUTING.md) |
| Governance   | [GOVERNANCE.md](GOVERNANCE.md) |
| Security     | [SECURITY.md](SECURITY.md) |
| Code of conduct | [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) |

After pushing to GitHub, the module is indexed at **https://pkg.go.dev/github.com/olaflaitinen/triagegeist**. No extra configuration required.

---

## Contributing and governance

Contributions are welcome under EUPL-1.2. Please read [CONTRIBUTING.md](CONTRIBUTING.md) for workflow, code style, and testing. By submitting contributions, you agree to license them under EUPL-1.2. Report security issues privately; see [SECURITY.md](SECURITY.md).

---

## Requirements and compatibility

| Requirement   | Version / note |
|---------------|----------------|
| Go            | 1.22 or later  |
| Dependencies  | None for core and score; norm, metrics, stats, validate, export are part of the module |
| Platforms     | All supported by Go (linux, windows, darwin, etc.) |

---

## Disclaimer

This library is for **research and decision support only**. It is not a substitute for clinical judgment or institutional triage protocols. Use in production clinical systems is at the user's risk. Operators must validate and calibrate parameters for their setting and regulatory context.

---

## Quick reference

| Task              | Package or API |
|-------------------|----------------|
| Single evaluation | `eng.ScoreAndLevel(v, resources)` |
| Batch evaluation  | `eng.BatchScoreAndLevel(vitals, resources)` |
| Custom thresholds | `p.T1, p.T2, p.T3, p.T4` then `NewEngine(p)` |
| Validate vitals   | `validate.Vitals(v)`, `validate.ClampVitals(v)` |
| Validate params   | `p.Validate()` or `ValidateParamsExternal(p)` |
| Metrics           | `metrics.NewConfusionMatrix(pred, ref)`, `cm.Sensitivity(1)`, etc. |
| Statistics        | `stats.Mean(scores)`, `stats.CI95(scores)`, `stats.ComputeLevelStats(levels)` |
| Export            | `export.FromVitalsScoreLevel(...)`, `export.WriteCSV(w, results)` |
