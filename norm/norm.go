// Copyright (c) triagegeist authors: Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov.
// Licensed under the EUPL.
//
// Package norm provides reference ranges (normal values) and normalisation
// utilities for vital signs used in triagegeist acuity scoring.
//
// Reference ranges are given as (midpoint, half-width) pairs. The midpoint
// is the central value of the normal range; the half-width is the allowed
// deviation before the normalised deviation reaches 1. For a value x:
//
//	deviation = min(1, |x - mid| / halfWidth)
//
// This package does not depend on the score or triagegeist root package;
// it can be used standalone for clamping, validation, or custom formulae.
//
// # Index of reference ranges (adult, general ED)
//
//	| Vital | Mid | Half-width | Unit   | Source note        |
//	|-------|-----|------------|--------|--------------------|
//	| HR    | 80  | 40         | bpm    | Adult resting      |
//	| RR    | 16  | 10         | /min   | Adult at rest       |
//	| SBP   | 120 | 40         | mmHg   | Adult               |
//	| DBP   | 80  | 30         | mmHg   | Adult               |
//	| Temp  | 37.0| 2.0        | Celsius| Core                |
//	| SpO2  | 98  | 8          | %      | Room air baseline   |
//	| GCS   | 15  | 6          | 3-15   | Full scale         |
//
// # Usage
//
// Use DefaultRanges() for the standard triagegeist ranges, or build a
// custom Ranges for site-specific or population-specific calibration.
// Deviation and NormalizeLinear are pure functions and safe for concurrent use.
package norm

import "math"

func finite(f float64) bool { return !math.IsNaN(f) && !math.IsInf(f, 0) }

// Ranges holds (midpoint, half-width) for each of the seven vitals.
// Order: HR, RR, SBP, DBP, Temp, SpO2, GCS.
// Zero half-width means that vital is not used in deviation calculation.
type Ranges struct {
	HR   [2]float64 // [mid, halfWidth]
	RR   [2]float64
	SBP  [2]float64
	DBP  [2]float64
	Temp [2]float64
	SpO2 [2]float64
	GCS  [2]float64
}

// DefaultRanges returns the standard adult ED reference ranges used by
// triagegeist for deviation-based acuity. All half-widths are positive.
func DefaultRanges() Ranges {
	return Ranges{
		HR:   [2]float64{80, 40},
		RR:   [2]float64{16, 10},
		SBP:  [2]float64{120, 40},
		DBP:  [2]float64{80, 30},
		Temp: [2]float64{37.0, 2.0},
		SpO2: [2]float64{98, 8},
		GCS:  [2]float64{15, 6},
	}
}

// PediatricRanges returns example ranges for a paediatric population.
// These are illustrative only; calibrate to your own protocol.
func PediatricRanges() Ranges {
	return Ranges{
		HR:   [2]float64{100, 50},
		RR:   [2]float64{24, 14},
		SBP:  [2]float64{90, 30},
		DBP:  [2]float64{60, 25},
		Temp: [2]float64{37.0, 2.0},
		SpO2: [2]float64{98, 8},
		GCS:  [2]float64{15, 6},
	}
}

// Deviation returns the normalised deviation of value from the reference:
//
//	d = min(1, |value - mid| / halfWidth)
//
// If halfWidth <= 0, returns 0. If value is considered "missing" (e.g. 0 when
// mid is positive), the caller should not call Deviation or pass a sentinel;
// this function does not treat 0 specially. Result is in [0, 1].
func Deviation(value, mid, halfWidth float64) float64 {
	if halfWidth <= 0 {
		return 0
	}
	d := math.Abs(value - mid) / halfWidth
	if d > 1 {
		return 1
	}
	return d
}

// NormalizeLinear maps x from [low, high] to [0, 1]. Values outside the
// interval are clamped. If low >= high, returns 0.
func NormalizeLinear(x, low, high float64) float64 {
	if low >= high {
		return 0
	}
	if x <= low {
		return 0
	}
	if x >= high {
		return 1
	}
	return (x - low) / (high - low)
}

// ClampToRange returns x clamped to [low, high]. If low > high, returns low.
func ClampToRange(x, low, high float64) float64 {
	if low > high {
		return low
	}
	if x < low {
		return low
	}
	if x > high {
		return high
	}
	return x
}

// InRange returns true if low <= x <= high. If low > high, always false.
func InRange(x, low, high float64) bool {
	if low > high {
		return false
	}
	return x >= low && x <= high
}

// VitalIndex constants for indexing into Ranges or weight vectors.
const (
	VitalHR   = 0
	VitalRR   = 1
	VitalSBP  = 2
	VitalDBP  = 3
	VitalTemp = 4
	VitalSpO2 = 5
	VitalGCS  = 6
	NumVitals = 7
)

// At returns the [mid, halfWidth] for vital index i (0..6). If i is out of
// range, returns zero values.
func (r Ranges) At(i int) (mid, halfWidth float64) {
	switch i {
	case 0:
		return r.HR[0], r.HR[1]
	case 1:
		return r.RR[0], r.RR[1]
	case 2:
		return r.SBP[0], r.SBP[1]
	case 3:
		return r.DBP[0], r.DBP[1]
	case 4:
		return r.Temp[0], r.Temp[1]
	case 5:
		return r.SpO2[0], r.SpO2[1]
	case 6:
		return r.GCS[0], r.GCS[1]
	default:
		return 0, 0
	}
}

// Set sets the [mid, halfWidth] for vital index i. No-op if i out of range.
func (r *Ranges) Set(i int, mid, halfWidth float64) {
	switch i {
	case 0:
		r.HR[0], r.HR[1] = mid, halfWidth
	case 1:
		r.RR[0], r.RR[1] = mid, halfWidth
	case 2:
		r.SBP[0], r.SBP[1] = mid, halfWidth
	case 3:
		r.DBP[0], r.DBP[1] = mid, halfWidth
	case 4:
		r.Temp[0], r.Temp[1] = mid, halfWidth
	case 5:
		r.SpO2[0], r.SpO2[1] = mid, halfWidth
	case 6:
		r.GCS[0], r.GCS[1] = mid, halfWidth
	}
}

// Valid returns true if all half-widths are non-negative and finite.
func (r Ranges) Valid() bool {
	for i := 0; i < NumVitals; i++ {
		m, hw := r.At(i)
		if hw < 0 || !finite(m) || !finite(hw) {
			return false
		}
	}
	return true
}

// DeviationHR returns Deviation(v, r.HR[0], r.HR[1]).
func (r Ranges) DeviationHR(v float64) float64 { return Deviation(v, r.HR[0], r.HR[1]) }

// DeviationRR returns Deviation(v, r.RR[0], r.RR[1]).
func (r Ranges) DeviationRR(v float64) float64 { return Deviation(v, r.RR[0], r.RR[1]) }

// DeviationSBP returns Deviation(v, r.SBP[0], r.SBP[1]).
func (r Ranges) DeviationSBP(v float64) float64 { return Deviation(v, r.SBP[0], r.SBP[1]) }

// DeviationDBP returns Deviation(v, r.DBP[0], r.DBP[1]).
func (r Ranges) DeviationDBP(v float64) float64 { return Deviation(v, r.DBP[0], r.DBP[1]) }

// DeviationTemp returns Deviation(v, r.Temp[0], r.Temp[1]).
func (r Ranges) DeviationTemp(v float64) float64 { return Deviation(v, r.Temp[0], r.Temp[1]) }

// DeviationSpO2 returns Deviation(v, r.SpO2[0], r.SpO2[1]).
func (r Ranges) DeviationSpO2(v float64) float64 { return Deviation(v, r.SpO2[0], r.SpO2[1]) }

// DeviationGCS returns Deviation(v, r.GCS[0], r.GCS[1]).
func (r Ranges) DeviationGCS(v float64) float64 { return Deviation(v, r.GCS[0], r.GCS[1]) }

// Copy returns a copy of r.
func (r Ranges) Copy() Ranges {
	return Ranges{
		HR:   r.HR,
		RR:   r.RR,
		SBP:  r.SBP,
		DBP:  r.DBP,
		Temp: r.Temp,
		SpO2: r.SpO2,
		GCS:  r.GCS,
	}
}

// MergeWith overwrites ranges in r with non-zero half-widths from other.
// Used to apply overrides (e.g. only change Temp) without touching the rest.
func (r *Ranges) MergeWith(other Ranges) {
	if other.HR[1] > 0 {
		r.HR = other.HR
	}
	if other.RR[1] > 0 {
		r.RR = other.RR
	}
	if other.SBP[1] > 0 {
		r.SBP = other.SBP
	}
	if other.DBP[1] > 0 {
		r.DBP = other.DBP
	}
	if other.Temp[1] > 0 {
		r.Temp = other.Temp
	}
	if other.SpO2[1] > 0 {
		r.SpO2 = other.SpO2
	}
	if other.GCS[1] > 0 {
		r.GCS = other.GCS
	}
}

// ScaleHalfWidths multiplies all half-widths by factor. Use to widen or
// narrow the "normal" band without changing midpoints. No-op if factor <= 0.
func (r *Ranges) ScaleHalfWidths(factor float64) {
	if factor <= 0 {
		return
	}
	r.HR[1] *= factor
	r.RR[1] *= factor
	r.SBP[1] *= factor
	r.DBP[1] *= factor
	r.Temp[1] *= factor
	r.SpO2[1] *= factor
	r.GCS[1] *= factor
}

// CriticalBounds returns recommended absolute bounds [min, max] for vital i
// for use in validation (e.g. flag out-of-range inputs). Values are
// conservative; i must be 0..6. Returns (0, 0) if i out of range.
//
//	| Vital | Min  | Max   |
//	| HR    | 20   | 300   |
//	| RR    | 0    | 60    |
//	| SBP   | 40   | 300   |
//	| DBP   | 20   | 200   |
//	| Temp  | 30   | 45    |
//	| SpO2  | 0    | 100   |
//	| GCS   | 3    | 15    |
func CriticalBounds(i int) (min, max float64) {
	switch i {
	case VitalHR:
		return 20, 300
	case VitalRR:
		return 0, 60
	case VitalSBP:
		return 40, 300
	case VitalDBP:
		return 20, 200
	case VitalTemp:
		return 30, 45
	case VitalSpO2:
		return 0, 100
	case VitalGCS:
		return 3, 15
	default:
		return 0, 0
	}
}

// IsWithinCriticalBounds returns true if value is within CriticalBounds(i).
func IsWithinCriticalBounds(i int, value float64) bool {
	lo, hi := CriticalBounds(i)
	return value >= lo && value <= hi
}

// WeightedDeviationSum computes sum over present vitals of weight[i] * deviation[i],
// and the sum of weights for present vitals. Present means value > 0 for
// integer vitals or value != 0 for Temp. Used by callers to build V without
// depending on score.Vitals.
func (r Ranges) WeightedDeviationSum(values [7]float64, weights [7]float64) (sum, weightSum float64) {
	for i := 0; i < NumVitals; i++ {
		if weights[i] <= 0 {
			continue
		}
		v := values[i]
		if i == VitalTemp {
			if v == 0 {
				continue
			}
		} else if v <= 0 {
			continue
		}
		mid, hw := r.At(i)
		if hw <= 0 {
			continue
		}
		d := Deviation(v, mid, hw)
		sum += weights[i] * d
		weightSum += weights[i]
	}
	return sum, weightSum
}
