# Comparison with alternatives

This document positions triagegeist against other approaches to emergency triage and acuity scoring: other programming languages, other Go libraries, and generic ML runtimes. The goal is to help integrators choose the right tool and to clarify where triagegeist excels. Every comparison dimension is covered.

---

## Comparison criteria

| Criterion | Description |
|-----------|-------------|
| **Performance** | Latency per evaluation (ns/op), allocations per call, throughput per core. |
| **Portability** | Pure Go vs cgo; dependency footprint; cross-compilation to various OS/arch. |
| **Determinism and auditability** | Formula-based vs black-box; reproducibility across runs and machines. |
| **Flexibility** | Configurable parameters (weights, thresholds, norms); extensibility (e.g. custom models). |
| **Documentation and standards** | API docs, mathematical specification, licence, governance, contribution process. |
| **Domain focus** | Purpose-built for ED triage/acuity vs general-purpose ML or clinical tools. |
| **Validation and metrics** | Built-in input validation, accuracy metrics (sensitivity, specificity, kappa), export (JSON/CSV). |

---

## triagegeist vs general-purpose ML runtimes (Python, TensorFlow, ONNX, etc.)

| Aspect | triagegeist | Typical ML stack (e.g. Python + TF/ONNX) |
|--------|-------------|------------------------------------------|
| **Latency** | Sub-microsecond per call; no interpreter or graph load | Millisecond to tens of ms per call; cold start and model load overhead |
| **Allocations** | Zero in hot path (stack-allocated structs) | Allocations for tensors, Python objects, or runtime buffers |
| **Deployment** | Single static binary; no Python or runtime install | Requires Python interpreter, libraries, often GPU/CUDA for large models |
| **Determinism** | Fully deterministic, formula-based | Depends on model and runtime; floating-point and hardware can vary |
| **Auditability** | Explicit weights and thresholds; no black box | Model weights and architecture may be opaque |
| **Use case** | Parametric triage/acuity; rule-like, interpretable | When you need learned models (e.g. neural nets) and accept higher cost |

**When to use triagegeist**: When you need fast, deterministic, auditable acuity and level from vitals and resources in Go, without ML dependencies or heavy infrastructure.

**When to use ML runtimes**: When you have trained models (e.g. deep learning) that must run in production and you accept the operational and latency cost.

---

## triagegeist vs other Go libraries

As of this writing, there are few dedicated, widely adopted Go libraries for emergency triage or acuity scoring. Many systems use:

- **Ad-hoc logic in application code**: Hard to reuse, test, or document. triagegeist provides a single, documented, tested package and subpackages (norm, metrics, stats, validate, export).
- **Wrappers around C or external services**: Introduce cgo or network latency and failure modes. triagegeist is pure Go.
- **Generic numeric or stats libraries**: Not tailored to triage (vitals, levels, thresholds). triagegeist is domain-specific and includes validation, metrics, and export.

triagegeist aims to be the **reference Go library** for parametric ED triage and acuity: small API, clear formula, no mandatory external deps, strong documentation, and full test coverage across all packages.

---

## triagegeist vs reference triage systems (ESI, MTS, etc.)

| Aspect | triagegeist | ESI / MTS (conceptually) |
|--------|-------------|---------------------------|
| **Basis** | Parametric formula (weights + norms + thresholds) | Proprietary algorithms; flowcharts and decision trees |
| **Implementation** | Open source (EUPL-1.2), full formula in code | Often proprietary or licence-restricted; not fully replicable in code |
| **Calibration** | All parameters configurable; defaults are generic | Tied to specific system; may require licence to use officially |
| **Output** | Continuous score \( s \in [0,1] \) plus discrete level \( L \in \{1,2,3,4,5\} \) | Typically discrete level only |

triagegeist does **not** implement ESI or MTS verbatim (which would require their licence and exact logic). It provides a **parametric, auditable alternative** that can be calibrated toward similar behaviour where legally and clinically appropriate.

---

## Summary table

| Dimension | triagegeist | Heavy ML (e.g. Python/TF) | Ad-hoc Go | Proprietary triage (ESI/MTS) |
|-----------|-------------|----------------------------|-----------|-----------------------------|
| Speed | Very high | Lower | Depends | N/A (often external) |
| Allocations | Zero (hot path) | Higher | Depends | N/A |
| Determinism | Full | Model-dependent | Depends | Yes (but closed) |
| Configurability | Full (params, norms) | Model-dependent | Depends | Limited |
| Licence | EUPL-1.2, open | Varies | Varies | Often restricted |
| Domain | ED triage/acuity | General ML | Varies | ED triage |
| Validation | validate package | Usually external | Ad-hoc | Built into product |
| Metrics | metrics package | Varies | Ad-hoc | Varies |
| Export | export package (JSON/CSV) | Varies | Ad-hoc | Varies |

---

## Conclusion

Use **triagegeist** when you need:

- Fast, low-allocation acuity and level computation in Go.
- A single, well-documented, formula-based implementation with full test coverage.
- No external ML runtime or proprietary triage algorithm dependency.
- Full parameter control and reproducibility for research or deployment.
- Built-in validation, metrics, and export for integration and auditing.

Consider other tools when you need official ESI/MTS certification, trained neural models, or integration with existing non-Go triage services that cannot be replaced.
