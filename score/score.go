// Copyright (c) triagegeist authors: Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov.
// Licensed under the EUPL.
//
// Package score provides the parametric acuity formula used by triagegeist.
//
// The raw score is a weighted combination of vital-sign deviations (from
// normal ranges) and expected resource consumption, then normalized to [0, 1].
package score

import "math"

// Vitals holds one set of vital signs. Units: HR (bpm), RR (per min),
// SBP/DBP (mmHg), Temp (Celsius), SpO2 (%), GCS (3-15). Use 0 for unknown;
// unknown values are excluded from the weighted sum.
type Vitals struct {
	HR   int     // Heart rate, beats per minute
	RR   int     // Respiratory rate, per minute
	SBP  int     // Systolic blood pressure, mmHg
	DBP  int     // Diastolic blood pressure, mmHg
	Temp float64 // Temperature, Celsius
	SpO2 int     // Oxygen saturation, percent
	GCS  int     // Glasgow Coma Scale, 3-15
}

// VitalWeights is the default weight vector (HR, RR, SBP, DBP, Temp, SpO2, GCS).
var VitalWeights = [7]float64{0.18, 0.22, 0.16, 0.10, 0.08, 0.16, 0.10}

// Normal ranges (mid and half-width) for deviation calculation.
// Deviation = |value - mid| / halfWidth, capped at 1.
var (
	HRNorm   = [2]float64{80, 40}
	RRNorm   = [2]float64{16, 10}
	SBPNorm  = [2]float64{120, 40}
	DBPNorm  = [2]float64{80, 30}
	TempNorm = [2]float64{37.0, 2.0}
	SpO2Norm = [2]float64{98, 8}
	GCSNorm  = [2]float64{15, 6}
)

// deviation returns |v - mid| / hw capped to 1. If hw <= 0 or v is "unknown", returns 0.
func deviation(v float64, mid, hw float64) float64 {
	if hw <= 0 {
		return 0
	}
	if v <= 0 && mid > 0 {
		return 0
	}
	d := math.Abs(v - mid) / hw
	if d > 1 {
		return 1
	}
	return d
}

// VitalComponent returns the weighted sum of vital deviations in [0, 1].
// Uses the package-level VitalWeights; pass a custom slice if needed via AcuityRaw.
func VitalComponent(v Vitals, weights [7]float64) float64 {
	var sum, wSum float64
	if v.HR > 0 {
		d := deviation(float64(v.HR), HRNorm[0], HRNorm[1])
		sum += weights[0] * d
		wSum += weights[0]
	}
	if v.RR > 0 {
		d := deviation(float64(v.RR), RRNorm[0], RRNorm[1])
		sum += weights[1] * d
		wSum += weights[1]
	}
	if v.SBP > 0 {
		d := deviation(float64(v.SBP), SBPNorm[0], SBPNorm[1])
		sum += weights[2] * d
		wSum += weights[2]
	}
	if v.DBP > 0 {
		d := deviation(float64(v.DBP), DBPNorm[0], DBPNorm[1])
		sum += weights[3] * d
		wSum += weights[3]
	}
	if v.Temp > 0 {
		d := deviation(v.Temp, TempNorm[0], TempNorm[1])
		sum += weights[4] * d
		wSum += weights[4]
	}
	if v.SpO2 > 0 {
		d := deviation(float64(v.SpO2), SpO2Norm[0], SpO2Norm[1])
		sum += weights[5] * d
		wSum += weights[5]
	}
	if v.GCS > 0 {
		d := deviation(float64(v.GCS), GCSNorm[0], GCSNorm[1])
		sum += weights[6] * d
		wSum += weights[6]
	}
	if wSum <= 0 {
		return 0
	}
	raw := sum / wSum
	if raw > 1 {
		return 1
	}
	return raw
}

// ResourceComponent returns the resource contribution in [0, 1] for
// resourceCount with given maxResources and weight.
func ResourceComponent(resourceCount, maxResources int, weight float64) float64 {
	if maxResources <= 0 || weight <= 0 {
		return 0
	}
	if resourceCount <= 0 {
		return 0
	}
	r := float64(resourceCount) / float64(maxResources)
	if r > 1 {
		r = 1
	}
	return weight * r
}

// AcuityRaw returns the unnormalized acuity: vitalComponent + resourceComponent.
// vitalComponent and resourceComponent should be in [0, 1]; weight sum can exceed 1.
func AcuityRaw(vitalComponent, resourceComponent float64) float64 {
	return vitalComponent + resourceComponent
}

// Normalize maps raw score to [0, 1] using the divisor (vitalWeightSum + resourceWeight).
// If divisor <= 0, returns 0.
func Normalize(raw, divisor float64) float64 {
	if divisor <= 0 {
		return 0
	}
	s := raw / divisor
	if s > 1 {
		return 1
	}
	if s < 0 {
		return 0
	}
	return s
}

// Acuity returns the normalized acuity score in [0, 1] for the given vitals,
// resource count, and weights. It uses VitalWeights and the provided
// maxResources and resourceWeight to compute the divisor for normalization.
func Acuity(v Vitals, resourceCount, maxResources int, vitalWeights [7]float64, resourceWeight float64) float64 {
	vSum := VitalComponent(v, vitalWeights)
	var wSum float64
	for _, w := range vitalWeights {
		wSum += w
	}
	rComp := ResourceComponent(resourceCount, maxResources, resourceWeight)
	raw := AcuityRaw(vSum, rComp)
	div := wSum + resourceWeight
	return Normalize(raw, div)
}

// VitalsToValues returns [7]float64 with HR, RR, SBP, DBP, Temp, SpO2, GCS in order.
// Used when interfacing with packages that expect a fixed array (e.g. norm.WeightedDeviationSum).
func VitalsToValues(v Vitals) [7]float64 {
	return [7]float64{
		float64(v.HR), float64(v.RR), float64(v.SBP), float64(v.DBP),
		v.Temp, float64(v.SpO2), float64(v.GCS),
	}
}

// PresentCount returns the number of vitals that are present (non-zero).
// Temp is present if != 0.
func PresentCount(v Vitals) int {
	var n int
	if v.HR > 0 {
		n++
	}
	if v.RR > 0 {
		n++
	}
	if v.SBP > 0 {
		n++
	}
	if v.DBP > 0 {
		n++
	}
	if v.Temp != 0 {
		n++
	}
	if v.SpO2 > 0 {
		n++
	}
	if v.GCS > 0 {
		n++
	}
	return n
}

// VitalComponentWithNorms computes the vital component using custom norms (mid, halfWidth) per vital.
// norms[i] = [mid, halfWidth] for vital i (0..6). If norms[i][1] <= 0, that vital is skipped.
func VitalComponentWithNorms(v Vitals, weights [7]float64, norms [7][2]float64) float64 {
	var sum, wSum float64
	if v.HR > 0 && norms[0][1] > 0 {
		d := deviation(float64(v.HR), norms[0][0], norms[0][1])
		sum += weights[0] * d
		wSum += weights[0]
	}
	if v.RR > 0 && norms[1][1] > 0 {
		d := deviation(float64(v.RR), norms[1][0], norms[1][1])
		sum += weights[1] * d
		wSum += weights[1]
	}
	if v.SBP > 0 && norms[2][1] > 0 {
		d := deviation(float64(v.SBP), norms[2][0], norms[2][1])
		sum += weights[2] * d
		wSum += weights[2]
	}
	if v.DBP > 0 && norms[3][1] > 0 {
		d := deviation(float64(v.DBP), norms[3][0], norms[3][1])
		sum += weights[3] * d
		wSum += weights[3]
	}
	if v.Temp != 0 && norms[4][1] > 0 {
		d := deviation(v.Temp, norms[4][0], norms[4][1])
		sum += weights[4] * d
		wSum += weights[4]
	}
	if v.SpO2 > 0 && norms[5][1] > 0 {
		d := deviation(float64(v.SpO2), norms[5][0], norms[5][1])
		sum += weights[5] * d
		wSum += weights[5]
	}
	if v.GCS > 0 && norms[6][1] > 0 {
		d := deviation(float64(v.GCS), norms[6][0], norms[6][1])
		sum += weights[6] * d
		wSum += weights[6]
	}
	if wSum <= 0 {
		return 0
	}
	raw := sum / wSum
	if raw > 1 {
		return 1
	}
	return raw
}

// DefaultNorms returns the package-level norm arrays as [7][2]float64 for use with VitalComponentWithNorms.
func DefaultNorms() [7][2]float64 {
	return [7][2]float64{
		HRNorm, RRNorm, SBPNorm, DBPNorm, TempNorm, SpO2Norm, GCSNorm,
	}
}

// AcuityWithNorms is like Acuity but uses VitalComponentWithNorms with the given norms.
func AcuityWithNorms(v Vitals, resourceCount, maxResources int, vitalWeights [7]float64, resourceWeight float64, norms [7][2]float64) float64 {
	vSum := VitalComponentWithNorms(v, vitalWeights, norms)
	var wSum float64
	for _, w := range vitalWeights {
		wSum += w
	}
	rComp := ResourceComponent(resourceCount, maxResources, resourceWeight)
	raw := AcuityRaw(vSum, rComp)
	div := wSum + resourceWeight
	return Normalize(raw, div)
}

// WeightSum returns the sum of the given weight vector.
func WeightSum(w [7]float64) float64 {
	var s float64
	for _, x := range w {
		s += x
	}
	return s
}

// CloneVitals returns a copy of v.
func CloneVitals(v Vitals) Vitals {
	return v
}

// ZeroVitals returns a Vitals struct with all fields zero.
func ZeroVitals() Vitals {
	return Vitals{}
}
