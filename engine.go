// Copyright (c) triagegeist authors: Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov.
// Licensed under the EUPL.

package triagegeist

import (
	"github.com/olaflaitinen/triagegeist/score"
)

// Engine evaluates acuity and level from vitals and resource count using
// a fixed parameter set. Safe for concurrent use; no mutable state.
//
// All methods that take (vitals, resourceCount) use the engine's Params
// for weights, thresholds, and maxResources. The engine does not modify
// Params during evaluation.
//
//	| Method              | Returns                   | Use case                 |
//	|---------------------|---------------------------|--------------------------|
//	| Acuity              | score in [0, 1]           | Continuous outcome       |
//	| Level               | Level 1..5                | Discrete triage level    |
//	| ScoreAndLevel       | (acuity, level)           | Single evaluation         |
//	| BatchScoreAndLevel  | (acuities, levels)        | Batch evaluation          |
//	| BatchAcuity         | []float64                 | Batch acuity only         |
//	| BatchLevel          | []Level                   | Batch level only          |
//	| Evaluate            | EvaluateResult            | Single with struct        |
//	| BatchEvaluate       | []EvaluateResult          | Batch with struct         |
type Engine struct {
	P Params
}

// NewEngine returns an engine with the given parameters. Use DefaultParams()
// if no custom calibration is needed.
func NewEngine(p Params) *Engine {
	return &Engine{P: p}
}

// Acuity returns the normalized acuity score in [0, 1] for the given vitals
// and resource count, using the engine's parameters.
func (e *Engine) Acuity(v score.Vitals, resourceCount int) float64 {
	return score.Acuity(v, resourceCount, e.P.MaxResources, e.P.VitalWeights, e.P.ResourceWeight)
}

// Level returns the discrete triage level (1 to 5) for the given vitals and
// resource count.
func (e *Engine) Level(v score.Vitals, resourceCount int) Level {
	s := e.Acuity(v, resourceCount)
	return FromScore(s, e.P)
}

// ScoreAndLevel returns both the normalized acuity score and the level.
func (e *Engine) ScoreAndLevel(v score.Vitals, resourceCount int) (acuity float64, level Level) {
	acuity = e.Acuity(v, resourceCount)
	level = FromScore(acuity, e.P)
	return acuity, level
}

// BatchScoreAndLevel evaluates acuity and level for each (vitals, resourceCount) pair.
// vitals and resourceCounts must have the same length. Returns two slices of that length.
func (e *Engine) BatchScoreAndLevel(vitals []score.Vitals, resourceCounts []int) (acuities []float64, levels []Level) {
	n := len(vitals)
	if len(resourceCounts) != n {
		return nil, nil
	}
	acuities = make([]float64, n)
	levels = make([]Level, n)
	for i := 0; i < n; i++ {
		acuities[i], levels[i] = e.ScoreAndLevel(vitals[i], resourceCounts[i])
	}
	return acuities, levels
}

// BatchAcuity returns acuity scores for each (vitals, resourceCount) pair.
func (e *Engine) BatchAcuity(vitals []score.Vitals, resourceCounts []int) []float64 {
	n := len(vitals)
	if len(resourceCounts) != n {
		return nil
	}
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = e.Acuity(vitals[i], resourceCounts[i])
	}
	return out
}

// BatchLevel returns levels for each (vitals, resourceCount) pair.
func (e *Engine) BatchLevel(vitals []score.Vitals, resourceCounts []int) []Level {
	n := len(vitals)
	if len(resourceCounts) != n {
		return nil
	}
	out := make([]Level, n)
	for i := 0; i < n; i++ {
		out[i] = e.Level(vitals[i], resourceCounts[i])
	}
	return out
}

// Params returns a copy of the engine's parameters.
func (e *Engine) Params() Params {
	return e.P.Clone()
}

// WithParams returns a new Engine with the given params. The receiver is unchanged.
func (e *Engine) WithParams(p Params) *Engine {
	return NewEngine(p)
}

// ScoreAndLevelWithResourceClamp evaluates ScoreAndLevel after clamping resourceCount
// to [0, MaxResources]. Use when input may exceed the parameter cap.
func (e *Engine) ScoreAndLevelWithResourceClamp(v score.Vitals, resourceCount int) (acuity float64, level Level) {
	rc := resourceCount
	if rc < 0 {
		rc = 0
	}
	if rc > e.P.MaxResources {
		rc = e.P.MaxResources
	}
	return e.ScoreAndLevel(v, rc)
}

// EvaluateResult holds acuity, level, and optional metadata for one evaluation.
type EvaluateResult struct {
	Acuity float64
	Level  Level
}

// Evaluate returns a single EvaluateResult.
func (e *Engine) Evaluate(v score.Vitals, resourceCount int) EvaluateResult {
	a, l := e.ScoreAndLevel(v, resourceCount)
	return EvaluateResult{Acuity: a, Level: l}
}

// BatchEvaluate returns a slice of EvaluateResult for each (vitals, resourceCount) pair.
func (e *Engine) BatchEvaluate(vitals []score.Vitals, resourceCounts []int) []EvaluateResult {
	n := len(vitals)
	if len(resourceCounts) != n {
		return nil
	}
	out := make([]EvaluateResult, n)
	for i := 0; i < n; i++ {
		out[i] = e.Evaluate(vitals[i], resourceCounts[i])
	}
	return out
}

// CountByLevel returns the number of evaluations in results that have the given level.
func CountByLevel(results []EvaluateResult, level Level) int {
	var c int
	for _, r := range results {
		if r.Level == level {
			c++
		}
	}
	return c
}

// MeanAcuity returns the mean acuity over results.
func MeanAcuity(results []EvaluateResult) float64 {
	if len(results) == 0 {
		return 0
	}
	var sum float64
	for _, r := range results {
		sum += r.Acuity
	}
	return sum / float64(len(results))
}

// LevelDistributionFromResults returns counts per level (index 1..5) for the given results.
func LevelDistributionFromResults(results []EvaluateResult) [6]int {
	var c [6]int
	for _, r := range results {
		if r.Level >= 1 && r.Level <= 5 {
			c[r.Level.Int()]++
		}
	}
	return c
}

// MinAcuity returns the minimum acuity in results. Returns 0 if empty.
func MinAcuity(results []EvaluateResult) float64 {
	if len(results) == 0 {
		return 0
	}
	min := results[0].Acuity
	for _, r := range results[1:] {
		if r.Acuity < min {
			min = r.Acuity
		}
	}
	return min
}

// MaxAcuity returns the maximum acuity in results. Returns 0 if empty.
func MaxAcuity(results []EvaluateResult) float64 {
	if len(results) == 0 {
		return 0
	}
	max := results[0].Acuity
	for _, r := range results[1:] {
		if r.Acuity > max {
			max = r.Acuity
		}
	}
	return max
}

// FilterByLevel returns the subset of results whose level equals the given level.
func FilterByLevel(results []EvaluateResult, level Level) []EvaluateResult {
	var out []EvaluateResult
	for _, r := range results {
		if r.Level == level {
			out = append(out, r)
		}
	}
	return out
}

// FilterHighAcuity returns results with level 1 or 2.
func FilterHighAcuity(results []EvaluateResult) []EvaluateResult {
	var out []EvaluateResult
	for _, r := range results {
		if r.Level.IsHighAcuity() {
			out = append(out, r)
		}
	}
	return out
}

// FilterLowAcuity returns results with level 4 or 5.
func FilterLowAcuity(results []EvaluateResult) []EvaluateResult {
	var out []EvaluateResult
	for _, r := range results {
		if r.Level.IsLowAcuity() {
			out = append(out, r)
		}
	}
	return out
}

// AcuityStats holds min, max, mean for a slice of EvaluateResult.
type AcuityStats struct {
	Min  float64
	Max  float64
	Mean float64
	N    int
}

// AcuityStatsFromResults computes AcuityStats for the given results.
func AcuityStatsFromResults(results []EvaluateResult) AcuityStats {
	var s AcuityStats
	s.N = len(results)
	if s.N == 0 {
		return s
	}
	s.Min = MinAcuity(results)
	s.Max = MaxAcuity(results)
	s.Mean = MeanAcuity(results)
	return s
}

// NewDefaultEngine returns an engine with DefaultParams().
func NewDefaultEngine() *Engine {
	return NewEngine(DefaultParams())
}

// NewStrictEngine returns an engine with PresetStrict().
func NewStrictEngine() *Engine {
	return NewEngine(PresetStrict())
}

// NewLenientEngine returns an engine with PresetLenient().
func NewLenientEngine() *Engine {
	return NewEngine(PresetLenient())
}

// NewResearchEngine returns an engine with PresetResearch().
func NewResearchEngine() *Engine {
	return NewEngine(PresetResearch())
}

