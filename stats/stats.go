// Copyright (c) triagegeist authors: Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov.
// Licensed under the EUPL.
//
// Package stats provides statistical helpers for triagegeist: descriptive
// statistics, confidence intervals, and aggregation over score/level outputs.
// All functions are pure and safe for concurrent use.
//
// # Formulas
//
// Mean:   μ = (1/n) Σ x_i
// Var:    σ² = (1/(n-1)) Σ (x_i - μ)²  (sample variance)
// StdDev: σ = sqrt(Var)
// SE:     SE = σ / sqrt(n)
// 95% CI: [μ - 1.96*SE, μ + 1.96*SE]  (normal approximation)
//
// Percentile: linear interpolation between order statistics.
package stats

import (
	"math"
	"sort"
)

// Mean returns the arithmetic mean of x. Returns 0 if len(x)==0.
func Mean(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}
	var sum float64
	for _, v := range x {
		sum += v
	}
	return sum / float64(len(x))
}

// Variance returns the sample variance (unbiased: divisor n-1). Returns 0 if n<2.
func Variance(x []float64) float64 {
	n := float64(len(x))
	if n < 2 {
		return 0
	}
	mu := Mean(x)
	var sum float64
	for _, v := range x {
		d := v - mu
		sum += d * d
	}
	return sum / (n - 1)
}

// StdDev returns the sample standard deviation (sqrt of sample variance).
func StdDev(x []float64) float64 {
	return math.Sqrt(Variance(x))
}

// SE returns the standard error of the mean: StdDev(x) / sqrt(n). Returns 0 if n<2.
func SE(x []float64) float64 {
	n := len(x)
	if n < 2 {
		return 0
	}
	return StdDev(x) / math.Sqrt(float64(n))
}

// CI95 returns the approximate 95% confidence interval for the mean using
// the normal approximation: [mean - 1.96*SE, mean + 1.96*SE].
// If n<2, returns (0, 0).
func CI95(x []float64) (low, high float64) {
	n := len(x)
	if n < 2 {
		return 0, 0
	}
	mu := Mean(x)
	se := SE(x)
	const z = 1.96
	return mu - z*se, mu + z*se
}

// Median returns the median (50th percentile). For even n, returns the
// average of the two middle values. Copies and sorts; does not modify x.
func Median(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}
	cp := make([]float64, len(x))
	copy(cp, x)
	sort.Float64s(cp)
	n := len(cp)
	if n%2 == 1 {
		return cp[n/2]
	}
	return (cp[n/2-1] + cp[n/2]) / 2
}

// Percentile returns the p-th percentile (0 <= p <= 100) using linear
// interpolation between order statistics. Copies and sorts; does not modify x.
func Percentile(x []float64, p float64) float64 {
	if len(x) == 0 || p < 0 || p > 100 {
		return 0
	}
	cp := make([]float64, len(x))
	copy(cp, x)
	sort.Float64s(cp)
	n := float64(len(cp))
	idx := p / 100 * (n - 1)
	i := int(idx)
	if i < 0 {
		i = 0
	}
	if i >= len(cp)-1 {
		return cp[len(cp)-1]
	}
	w := idx - float64(i)
	return cp[i]*(1-w) + cp[i+1]*w
}

// Min returns the minimum value. Returns 0 if len(x)==0.
func Min(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}
	m := x[0]
	for _, v := range x[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

// Max returns the maximum value. Returns 0 if len(x)==0.
func Max(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}
	m := x[0]
	for _, v := range x[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

// Sum returns the sum of x.
func Sum(x []float64) float64 {
	var s float64
	for _, v := range x {
		s += v
	}
	return s
}

// Count returns the number of elements in x that satisfy f.
func Count(x []float64, f func(float64) bool) int {
	var c int
	for _, v := range x {
		if f(v) {
			c++
		}
	}
	return c
}

// CountInt returns the number of elements in x that satisfy f.
func CountInt(x []int, f func(int) bool) int {
	var c int
	for _, v := range x {
		if f(v) {
			c++
		}
	}
	return c
}

// LevelDistribution returns counts of level 1..5 in levels (values 1..5 only).
func LevelDistribution(levels []int) [6]int {
	var out [6]int
	for _, L := range levels {
		if L >= 1 && L <= 5 {
			out[L]++
		}
	}
	return out
}

// ScoreStats holds summary statistics for a slice of acuity scores.
type ScoreStats struct {
	N      int
	Mean   float64
	StdDev float64
	SE     float64
	CI95Lo float64
	CI95Hi float64
	Min    float64
	Max    float64
	P25    float64
	P50    float64
	P75    float64
}

// ComputeScoreStats returns ScoreStats for the given scores (0..1).
func ComputeScoreStats(scores []float64) ScoreStats {
	var s ScoreStats
	s.N = len(scores)
	if s.N == 0 {
		return s
	}
	s.Mean = Mean(scores)
	s.StdDev = StdDev(scores)
	s.SE = SE(scores)
	s.CI95Lo, s.CI95Hi = CI95(scores)
	s.Min = Min(scores)
	s.Max = Max(scores)
	s.P25 = Percentile(scores, 25)
	s.P50 = Percentile(scores, 50)
	s.P75 = Percentile(scores, 75)
	return s
}

// LevelStats holds counts and proportions for levels 1..5.
type LevelStats struct {
	Counts [6]int // index 0 unused; 1..5
	Total  int
	Props  [6]float64
}

// ComputeLevelStats returns LevelStats for the given levels (1..5).
func ComputeLevelStats(levels []int) LevelStats {
	var ls LevelStats
	for _, L := range levels {
		if L >= 1 && L <= 5 {
			ls.Counts[L]++
			ls.Total++
		}
	}
	if ls.Total > 0 {
		for i := 1; i <= 5; i++ {
			ls.Props[i] = float64(ls.Counts[i]) / float64(ls.Total)
		}
	}
	return ls
}

// CorrelationPearson returns Pearson correlation between x and y. Both must
// have same length and n>=2. Returns 0 if invalid.
func CorrelationPearson(x, y []float64) float64 {
	if len(x) != len(y) || len(x) < 2 {
		return 0
	}
	mx, my := Mean(x), Mean(y)
	var sumXY, sumX2, sumY2 float64
	for i := range x {
		dx, dy := x[i]-mx, y[i]-my
		sumXY += dx * dy
		sumX2 += dx * dx
		sumY2 += dy * dy
	}
	if sumX2 == 0 || sumY2 == 0 {
		return 0
	}
	return sumXY / (math.Sqrt(sumX2) * math.Sqrt(sumY2))
}

// RMSE returns root mean squared error between predicted and reference values.
func RMSE(pred, ref []float64) float64 {
	if len(pred) != len(ref) || len(pred) == 0 {
		return 0
	}
	var sum float64
	for i := range pred {
		d := pred[i] - ref[i]
		sum += d * d
	}
	return math.Sqrt(sum / float64(len(pred)))
}

// MAE returns mean absolute error between pred and ref.
func MAE(pred, ref []float64) float64 {
	if len(pred) != len(ref) || len(pred) == 0 {
		return 0
	}
	var sum float64
	for i := range pred {
		sum += math.Abs(pred[i] - ref[i])
	}
	return sum / float64(len(pred))
}

// WithinTolerance returns the proportion of pairs (pred[i], ref[i]) for
// which |pred[i]-ref[i]| <= tol.
func WithinTolerance(pred, ref []float64, tol float64) float64 {
	if len(pred) != len(ref) || len(pred) == 0 {
		return 0
	}
	var c int
	for i := range pred {
		if math.Abs(pred[i]-ref[i]) <= tol {
			c++
		}
	}
	return float64(c) / float64(len(pred))
}

// ExactAgreement returns the proportion of pairs where pred[i]==ref[i] (for int slices).
func ExactAgreement(pred, ref []int) float64 {
	if len(pred) != len(ref) || len(pred) == 0 {
		return 0
	}
	var c int
	for i := range pred {
		if pred[i] == ref[i] {
			c++
		}
	}
	return float64(c) / float64(len(pred))
}

// WithinLevel returns the proportion of pairs where |pred[i]-ref[i]| <= 1.
func WithinLevel(pred, ref []int) float64 {
	if len(pred) != len(ref) || len(pred) == 0 {
		return 0
	}
	var c int
	for i := range pred {
		if absInt(pred[i]-ref[i]) <= 1 {
			c++
		}
	}
	return float64(c) / float64(len(pred))
}

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
