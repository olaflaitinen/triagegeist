// Copyright (c) triagegeist authors: Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov.
// Licensed under the EUPL.

package triagegeist

import "math"

// Params holds all tunable parameters for acuity scoring and level assignment.
// Defaults are chosen for general emergency department use; override for
// site-specific or research calibration.
//
//	| Field           | Type      | Valid range / note                          |
//	|-----------------|-----------|---------------------------------------------|
//	| VitalWeights    | [7]float64| Each in [0, 1]; order HR, RR, SBP, DBP, Temp, SpO2, GCS |
//	| MaxResources    | int       | >= 0                                        |
//	| ResourceWeight  | float64   | >= 0                                        |
//	| T1, T2, T3, T4  | float64   | T1 > T2 > T3 > T4, all in (0, 1]           |
type Params struct {
	VitalWeights   [7]float64
	MaxResources   int
	ResourceWeight float64
	T1, T2, T3, T4 float64
}

// DefaultParams returns parameters tuned for a typical five-level ED triage.
// Thresholds follow a geometric spacing in (0, 1).
func DefaultParams() Params {
	return Params{
		VitalWeights: [7]float64{
			0.18, 0.22, 0.16, 0.10, 0.08, 0.16, 0.10,
		},
		MaxResources:   6,
		ResourceWeight: 0.25,
		T1:             0.85,
		T2:             0.60,
		T3:             0.35,
		T4:             0.15,
	}
}

// PresetStrict uses higher thresholds so that more patients are classified
// as higher acuity (conservative). Use when under-triage must be minimised.
func PresetStrict() Params {
	p := DefaultParams()
	p.T1, p.T2, p.T3, p.T4 = 0.80, 0.55, 0.30, 0.12
	return p
}

// PresetLenient uses lower thresholds so that fewer patients are classified
// as highest acuity. Use when over-triage is a concern.
func PresetLenient() Params {
	p := DefaultParams()
	p.T1, p.T2, p.T3, p.T4 = 0.90, 0.68, 0.42, 0.18
	return p
}

// PresetResearch returns parameters with equal level widths (0.2 each) for
// balanced research cohorts.
func PresetResearch() Params {
	return Params{
		VitalWeights: [7]float64{
			0.18, 0.22, 0.16, 0.10, 0.08, 0.16, 0.10,
		},
		MaxResources:   6,
		ResourceWeight: 0.25,
		T1:             0.80,
		T2:             0.60,
		T3:             0.40,
		T4:             0.20,
	}
}

// Validate returns true if all fields are within admissible ranges.
func (p Params) Validate() bool {
	if p.MaxResources < 0 || p.ResourceWeight < 0 {
		return false
	}
	for _, w := range p.VitalWeights {
		if w < 0 || w > 1 {
			return false
		}
	}
	return p.T1 > p.T2 && p.T2 > p.T3 && p.T3 > p.T4 && p.T4 > 0 && p.T1 <= 1
}

// WeightSum returns the sum of VitalWeights (for normalisation divisor).
func (p Params) WeightSum() float64 {
	var s float64
	for _, w := range p.VitalWeights {
		s += w
	}
	return s
}

// Divisor returns WeightSum() + ResourceWeight (denominator for score normalisation).
func (p Params) Divisor() float64 {
	return p.WeightSum() + p.ResourceWeight
}

// Clone returns a copy of p.
func (p Params) Clone() Params {
	q := p
	q.VitalWeights = [7]float64{}
	copy(q.VitalWeights[:], p.VitalWeights[:])
	return q
}

// SetThresholds sets T1..T4. No validation; use Validate() after.
func (p *Params) SetThresholds(t1, t2, t3, t4 float64) {
	p.T1, p.T2, p.T3, p.T4 = t1, t2, t3, t4
}

// SetVitalWeight sets VitalWeights[i] (i 0..6). No validation.
func (p *Params) SetVitalWeight(i int, w float64) {
	if i >= 0 && i < 7 {
		p.VitalWeights[i] = w
	}
}

// ThresholdForLevel returns the lower bound threshold for level L (1..5).
// Level 1 has no upper bound (use 1.0); Level 5 has no lower bound (use 0.0).
func (p Params) ThresholdForLevel(L int) (low, high float64) {
	switch L {
	case 1:
		return p.T1, 1.0
	case 2:
		return p.T2, p.T1
	case 3:
		return p.T3, p.T2
	case 4:
		return p.T4, p.T3
	case 5:
		return 0.0, p.T4
	default:
		return 0, 0
	}
}

// ScoreToLevelContinuous returns a continuous "level" in [1, 5] by linear
// interpolation between thresholds. For display or smoothing only; discrete
// level should use FromScore.
func (p Params) ScoreToLevelContinuous(s float64) float64 {
	if s >= p.T1 {
		return 1.0 + (1.0-s)/(1.0-p.T1)*0.5
	}
	if s >= p.T2 {
		return 1.5 + (p.T1-s)/(p.T1-p.T2)*0.5
	}
	if s >= p.T3 {
		return 2.0 + (p.T2-s)/(p.T2-p.T3)*0.5
	}
	if s >= p.T4 {
		return 2.5 + (p.T3-s)/(p.T3-p.T4)*0.5
	}
	return 3.0 + (p.T4-s)/p.T4*2.0
}

// Equal returns true if p and q have the same field values.
func (p Params) Equal(q Params) bool {
	if p.MaxResources != q.MaxResources || p.ResourceWeight != q.ResourceWeight {
		return false
	}
	if p.T1 != q.T1 || p.T2 != q.T2 || p.T3 != q.T3 || p.T4 != q.T4 {
		return false
	}
	for i := range p.VitalWeights {
		if p.VitalWeights[i] != q.VitalWeights[i] {
			return false
		}
	}
	return true
}

// ScaleWeights multiplies all VitalWeights by factor and re-normalises so
// that the max weight is 1.0 (if factor > 0). Use to emphasise or de-emphasise
// all vitals proportionally.
func (p *Params) ScaleWeights(factor float64) {
	if factor <= 0 {
		return
	}
	var max float64
	for i := range p.VitalWeights {
		p.VitalWeights[i] *= factor
		if p.VitalWeights[i] > max {
			max = p.VitalWeights[i]
		}
	}
	if max > 0 {
		for i := range p.VitalWeights {
			p.VitalWeights[i] /= max
		}
	}
}

// NormalizeWeights scales VitalWeights so they sum to 1.0. If sum is 0, no-op.
func (p *Params) NormalizeWeights() {
	sum := p.WeightSum()
	if sum <= 0 {
		return
	}
	for i := range p.VitalWeights {
		p.VitalWeights[i] /= sum
	}
}

// EntropyWeights sets VitalWeights to uniform (1/7 each). Useful as baseline.
func (p *Params) EntropyWeights() {
	for i := range p.VitalWeights {
		p.VitalWeights[i] = 1.0 / 7.0
	}
}

// GeometricThresholds sets T1..T4 so that they are evenly spaced in log-space
// between min and max (e.g. 0.1 and 0.9). Useful for consistent spacing.
func (p *Params) GeometricThresholds(min, max float64) {
	if min <= 0 || max <= min || max > 1 {
		return
	}
	logMin := math.Log(min)
	logMax := math.Log(max)
	step := (logMax - logMin) / 5
	p.T4 = math.Exp(logMin + step)
	p.T3 = math.Exp(logMin + 2*step)
	p.T2 = math.Exp(logMin + 3*step)
	p.T1 = math.Exp(logMin + 4*step)
	if p.T1 > 1 {
		p.T1 = 1
	}
}

// IsStricterThan returns true if p classifies more patients as higher acuity than q
// (i.e. p's thresholds are lower so more scores fall into levels 1-2).
func (p Params) IsStricterThan(q Params) bool {
	return p.T1 < q.T1 && p.T2 < q.T2 && p.T3 < q.T3 && p.T4 < q.T4
}

// CopyWeightsFrom copies VitalWeights from q into p.
func (p *Params) CopyWeightsFrom(q Params) {
	copy(p.VitalWeights[:], q.VitalWeights[:])
}

// SetAllThresholds sets T1, T2, T3, T4 from a slice (len 4). No-op if len != 4.
func (p *Params) SetAllThresholds(t []float64) {
	if len(t) != 4 {
		return
	}
	p.T1, p.T2, p.T3, p.T4 = t[0], t[1], t[2], t[3]
}

// Thresholds returns [T1, T2, T3, T4].
func (p Params) Thresholds() [4]float64 {
	return [4]float64{p.T1, p.T2, p.T3, p.T4}
}
