# Documentation index

This directory contains design and reference documentation for the triagegeist project. It is intended for contributors, integrators, and researchers who need more detail than the root [README](../README.md) and the API docs on [pkg.go.dev](https://pkg.go.dev/github.com/olaflaitinen/triagegeist). Nothing is omitted; each document is described in full below.

---

## Documents

| Document | Purpose |
|----------|---------|
| [ARCHITECTURE.md](ARCHITECTURE.md) | Package layout (all packages and paths), directory layout, data flow (step-by-step), design decisions, extension points, testing strategy, conventions. Use this to understand how the library is structured and how data moves from input to output. |
| [BENCHMARKS.md](BENCHMARKS.md) | How to run benchmarks (exact commands), benchmark definitions (every benchmark listed), interpreting results (ns/op, B/op, allocs/op), expected order of magnitude, regression policy, adding new benchmarks. Use this to measure and guard performance. |
| [COMPARISON.md](COMPARISON.md) | Comparison criteria, triagegeist vs ML runtimes, vs other Go libraries, vs proprietary triage systems (ESI/MTS), summary table, when to use triagegeist vs alternatives. Use this to decide if triagegeist fits your use case. |

---

## Mathematical notation

Across the docs and the codebase we use consistent notation:

| Symbol | Meaning |
|--------|----------|
| \( s \) | Normalised acuity score in \( [0,1] \). |
| \( L \) | Discrete triage level in \( \{1,2,3,4,5\} \). |
| \( x_i, \mu_i, \sigma_i \) | Vital value, reference midpoint, and half-width for vital \( i \). |
| \( w_i \) | Weight for vital \( i \). |
| \( \alpha \) | Resource weight. |
| \( T_1, T_2, T_3, T_4 \) | Score thresholds for level assignment (\( T_1 > T_2 > T_3 > T_4 \)). |
| \( V, R \) | Vital component and resource component before normalisation. |
| \( d_i \) | Deviation for vital \( i \): \( \min(1, |x_i - \mu_i| / \sigma_i) \). |

**LaTeX support:** All mathematical content in the documentation uses LaTeX so that it renders correctly on GitHub and in viewers that support MathJax or KaTeX.

- **Inline math:** \( \ldots \) e.g. \( s \in [0,1] \), \( L \in \{1,2,3,4,5\} \).
- **Display math:** $$ \ldots $$ for block equations, e.g.
  $$
  d_i = \min\left(1,\ \frac{|x_i - \mu_i|}{\sigma_i}\right).
  $$

If your viewer does not render math, the raw LaTeX will still be readable (e.g. `\( s \in [0,1] \)`).

---

## Where to find what

- **Quick start and API overview**: Root [README](../README.md).
- **Package API (functions, types, methods)**: [pkg.go.dev/github.com/olaflaitinen/triagegeist](https://pkg.go.dev/github.com/olaflaitinen/triagegeist) and subpackages (score, norm, metrics, stats, validate, export).
- **How the library is built and why**: [ARCHITECTURE.md](ARCHITECTURE.md).
- **How to run and interpret benchmarks**: [BENCHMARKS.md](BENCHMARKS.md).
- **How triagegeist compares to other solutions**: [COMPARISON.md](COMPARISON.md).
- **How to contribute**: [CONTRIBUTING.md](../CONTRIBUTING.md).
- **Security and vulnerability reporting**: [SECURITY.md](../SECURITY.md).
- **Governance and maintainers**: [GOVERNANCE.md](../GOVERNANCE.md).
- **Changelog**: [CHANGELOG.md](../CHANGELOG.md).
