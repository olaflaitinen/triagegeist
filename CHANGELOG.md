# Changelog

All notable changes to the triagegeist project are documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/). The project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) where applicable.

---

## [Unreleased]

### Added

- (Placeholder for future changes. Move items here when they are merged but not yet released.)

### Changed

- (None.)

### Deprecated

- (None.)

### Removed

- (None.)

### Fixed

- (None.)

### Security

- (None.)

---

## [1.0.0] (initial release)

### Added

- **Core package `triagegeist`**
  - `Params`: configurable weights, thresholds, and resource cap; `DefaultParams()` and `Validate()`.
  - `Level`: five-level enum (1–5) with `String()`; `FromScore(score, params)` for threshold mapping.
  - `Engine`: `NewEngine(params)`, `Acuity(vitals, resourceCount)`, `Level(...)`, `ScoreAndLevel(...)`.
- **Subpackage `score`**
  - `Vitals`: struct for HR, RR, SBP, DBP, Temp, SpO2, GCS with documented units.
  - `VitalComponent`, `ResourceComponent`, `AcuityRaw`, `Normalize`, `Acuity`: formula implementation and normalisation.
  - Default vital weights and reference ranges (mid, half-width) for all seven vitals.
- **Documentation**
  - Package doc with acuity formula summary and level table.
  - Inline doc comments for all exported symbols.
  - README with mathematical model (LaTeX-style formulas), parameter tables, usage examples.
- **Testing**
  - Unit tests for `Engine`, `FromScore`, `Params.Validate`, and `score` package (VitalComponent, Acuity, Normalize).
  - Example tests for `Engine.ScoreAndLevel` and `FromScore` for pkg.go.dev.
- **Project layout**
  - `LICENSE` (EUPL-1.2), `CONTRIBUTING.md`, `SECURITY.md`, `CODE_OF_CONDUCT.md`, `GOVERNANCE.md`, `CHANGELOG.md`.
  - `assets/` for logo (e.g. 8000×2000 px SVG, no background).
  - `.gitignore` for Go and common IDE/OS artifacts.
  - Documentation in `docs/` (architecture, benchmarks, comparison) for maintainers and users.

### Requirements

- Go 1.22 or later.
- No external dependencies for core and `score` packages.

---

## Version history summary

| Version | Date       | Notes |
|---------|------------|-------|
| 1.0.0   | (TBD)      | Initial public release |
| Unreleased | (ongoing) | Development branch |

---

[Unreleased]: https://github.com/olaflaitinen/triagegeist/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/olaflaitinen/triagegeist/releases/tag/v1.0.0
