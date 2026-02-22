# Benchmarks

This document describes how triagegeist is benchmarked, what each benchmark measures, how to run and interpret results, and what to expect. No detail is omitted.

---

## How to run

From the repository root.

**All benchmarks, with memory stats:**

```bash
go test -bench=. -benchmem ./...
```

**Only the main package (Engine and root):**

```bash
go test -bench=BenchmarkEngine -benchmem .
```

**Only the score package:**

```bash
go test -bench=BenchmarkScore -benchmem ./score
```

**Longer run for stability (e.g. 3 seconds per benchmark):**

```bash
go test -bench=. -benchmem -benchtime=3s ./...
```

**Exclude benchmarks, run only unit tests:**

```bash
go test ./... -count=1
```

---

## Benchmark definitions (complete)

| Benchmark | Package | What it measures |
|-----------|---------|-------------------|
| `BenchmarkEngine_ScoreAndLevel` | triagegeist | Full path: one `score.Vitals` and one resource count; call `Engine.ScoreAndLevel`; reports ns/op, B/op, allocs/op. |
| `BenchmarkEngine_Acuity` | triagegeist | Same input; call `Engine.Acuity` only (no level mapping). |
| `BenchmarkScore_Acuity` | score | Direct call to `score.Acuity` with default weights and norms, no Engine. |

All use fixed inputs (e.g. `benchVitals`, `benchResources`). No I/O, no network, no file access.

---

## Interpreting results

- **ns/op**: Nanoseconds per operation. Lower is better. Depends on CPU and Go version.
- **B/op**: Bytes allocated per operation. Zero or very low is desired for hot paths.
- **allocs/op**: Number of heap allocations per operation. Zero is ideal for inner loops.

Example output (format may vary by Go version):

```
BenchmarkEngine_ScoreAndLevel-8    xxxxxxx   xxx ns/op   x B/op   x allocs/op
BenchmarkEngine_Acuity-8           xxxxxxx   xxx ns/op   x B/op   x allocs/op
BenchmarkScore_Acuity-8            xxxxxxx   xxx ns/op   x B/op   x allocs/op
```

The `-8` suffix indicates GOMAXPROCS=8. If you see non-zero allocs/op in the main path, check that you are not accidentally passing pointers or slices that escape; the design intends zero allocations for a single evaluation with stack-allocated Vitals and Params.

---

## Expected order of magnitude

- **Latency**: Single evaluation (acuity + level) should be on the order of **hundreds of nanoseconds** per call on modern hardware (single goroutine). Exact numbers depend on CPU, Go version, and inlining.
- **Allocations**: The design aims for **0 allocs/op** in the hot path when Vitals and Params are stack-allocated and not escaped. If allocs appear, they should be documented (e.g. future optional features or interface calls).
- **Throughput**: With 0 allocs and sub-microsecond latency, theoretical throughput is in the **millions of evaluations per second** per core. Real pipelines will add cost for I/O, validation, logging, and serialisation.

---

## Regression policy

- Significant increases in ns/op or allocs/op should be justified in the pull request (e.g. new feature, correctness fix).
- Benchmark results may be summarised in release notes or in this document when major versions are cut.
- If the project adds CI, benchmarks can be run and compared against a baseline; large regressions may fail the build.

---

## Adding new benchmarks

1. Add a function `BenchmarkXxx(b *testing.B)` in the appropriate `*_test.go` file.
2. Use `b.ResetTimer()` after setup if setup is expensive.
3. Run the operation under test inside `for i := 0; i < b.N; i++ { ... }`.
4. Document the new benchmark in this file under "Benchmark definitions".

---

## Comparison with other implementations

See [COMPARISON.md](COMPARISON.md) for triagegeist versus other triage/acuity libraries and ML runtimes. This document focuses only on how to run and interpret triagegeist benchmarks.
