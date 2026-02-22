// Copyright (c) triagegeist authors: Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov.
// Licensed under the EUPL.
//
// Package validate provides input validation and sanitisation for vitals,
// parameters, and resource counts used by triagegeist. Use before passing
// data to the engine to avoid undefined behaviour or out-of-range results.
//
// # Validation rules (summary)
//
//	| Input      | Rule                          | Action on failure      |
//	|------------|--------------------------------|------------------------|
//	| HR         | 20 <= HR <= 300 or 0 (missing) | Clamp or mark invalid  |
//	| RR         | 0 <= RR <= 60 or 0 (missing)   | Clamp or mark invalid  |
//	| SBP        | 40 <= SBP <= 300 or 0          | Clamp or mark invalid  |
//	| DBP        | 20 <= DBP <= 200 or 0          | Clamp or mark invalid  |
//	| Temp       | 30 <= Temp <= 45 or 0          | Clamp or mark invalid  |
//	| SpO2       | 0 <= SpO2 <= 100 or 0         | Clamp or mark invalid  |
//	| GCS        | 3 <= GCS <= 15 or 0            | Clamp or mark invalid  |
//	| Resources  | 0 <= count <= max (e.g. 20)    | Clamp                  |
//	| Params     | T1>T2>T3>T4, weights in [0,1]   | Return error           |
package validate

import (
	"math"

	"github.com/olaflaitinen/triagegeist/score"
)

func finite(f float64) bool { return !math.IsNaN(f) && !math.IsInf(f, 0) }

// VitalsReport holds validation results for a single Vitals struct.
type VitalsReport struct {
	Valid   bool
	HR      string // "ok" | "clamped" | "invalid" | "missing"
	RR      string
	SBP     string
	DBP     string
	Temp    string
	SpO2    string
	GCS     string
	Clamped score.Vitals // If clamping was applied, the clamped values
}

const (
	StatusOK      = "ok"
	StatusClamped = "clamped"
	StatusInvalid = "invalid"
	StatusMissing = "missing"
)

// Bounds for each vital (min, max). 0 for a vital means "missing" and is allowed.
var (
	HRBounds  = [2]int{20, 300}
	RRBounds  = [2]int{0, 60}
	SBPBounds = [2]int{40, 300}
	DBPBounds = [2]int{20, 200}
	SpO2Bounds = [2]int{0, 100}
	GCSBounds = [2]int{3, 15}
)

var (
	TempBounds = [2]float64{30, 45}
)

// Vitals checks v against bounds and returns a report. It does not modify v.
func Vitals(v score.Vitals) VitalsReport {
	var r VitalsReport
	r.Valid = true
	r.Clamped = v

	if v.HR != 0 {
		if v.HR < HRBounds[0] || v.HR > HRBounds[1] {
			r.HR = StatusInvalid
			r.Valid = false
		} else {
			r.HR = StatusOK
		}
	} else {
		r.HR = StatusMissing
	}

	if v.RR != 0 {
		if v.RR < RRBounds[0] || v.RR > RRBounds[1] {
			r.RR = StatusInvalid
			r.Valid = false
		} else {
			r.RR = StatusOK
		}
	} else {
		r.RR = StatusMissing
	}

	if v.SBP != 0 {
		if v.SBP < SBPBounds[0] || v.SBP > SBPBounds[1] {
			r.SBP = StatusInvalid
			r.Valid = false
		} else {
			r.SBP = StatusOK
		}
	} else {
		r.SBP = StatusMissing
	}

	if v.DBP != 0 {
		if v.DBP < DBPBounds[0] || v.DBP > DBPBounds[1] {
			r.DBP = StatusInvalid
			r.Valid = false
		} else {
			r.DBP = StatusOK
		}
	} else {
		r.DBP = StatusMissing
	}

	if v.Temp != 0 {
		if v.Temp < TempBounds[0] || v.Temp > TempBounds[1] {
			r.Temp = StatusInvalid
			r.Valid = false
		} else {
			r.Temp = StatusOK
		}
	} else {
		r.Temp = StatusMissing
	}

	if v.SpO2 != 0 {
		if v.SpO2 < SpO2Bounds[0] || v.SpO2 > SpO2Bounds[1] {
			r.SpO2 = StatusInvalid
			r.Valid = false
		} else {
			r.SpO2 = StatusOK
		}
	} else {
		r.SpO2 = StatusMissing
	}

	if v.GCS != 0 {
		if v.GCS < GCSBounds[0] || v.GCS > GCSBounds[1] {
			r.GCS = StatusInvalid
			r.Valid = false
		} else {
			r.GCS = StatusOK
		}
	} else {
		r.GCS = StatusMissing
	}

	return r
}

// ClampVitals returns a copy of v with all present vitals clamped to bounds.
// Missing (0) values are left as 0.
func ClampVitals(v score.Vitals) score.Vitals {
	out := v
	if v.HR != 0 {
		if v.HR < HRBounds[0] {
			out.HR = HRBounds[0]
		} else if v.HR > HRBounds[1] {
			out.HR = HRBounds[1]
		}
	}
	if v.RR != 0 {
		if v.RR < RRBounds[0] {
			out.RR = RRBounds[0]
		} else if v.RR > RRBounds[1] {
			out.RR = RRBounds[1]
		}
	}
	if v.SBP != 0 {
		if v.SBP < SBPBounds[0] {
			out.SBP = SBPBounds[0]
		} else if v.SBP > SBPBounds[1] {
			out.SBP = SBPBounds[1]
		}
	}
	if v.DBP != 0 {
		if v.DBP < DBPBounds[0] {
			out.DBP = DBPBounds[0]
		} else if v.DBP > DBPBounds[1] {
			out.DBP = DBPBounds[1]
		}
	}
	if v.Temp != 0 {
		if v.Temp < TempBounds[0] {
			out.Temp = TempBounds[0]
		} else if v.Temp > TempBounds[1] {
			out.Temp = TempBounds[1]
		}
	}
	if v.SpO2 != 0 {
		if v.SpO2 < SpO2Bounds[0] {
			out.SpO2 = SpO2Bounds[0]
		} else if v.SpO2 > SpO2Bounds[1] {
			out.SpO2 = SpO2Bounds[1]
		}
	}
	if v.GCS != 0 {
		if v.GCS < GCSBounds[0] {
			out.GCS = GCSBounds[0]
		} else if v.GCS > GCSBounds[1] {
			out.GCS = GCSBounds[1]
		}
	}
	return out
}

// VitalsValid returns true if all present vitals are within bounds.
func VitalsValid(v score.Vitals) bool {
	return Vitals(v).Valid
}

// ResourceCount returns count clamped to [0, maxResources]. If maxResources <= 0, returns 0.
func ResourceCount(count, maxResources int) int {
	if maxResources <= 0 {
		return 0
	}
	if count < 0 {
		return 0
	}
	if count > maxResources {
		return maxResources
	}
	return count
}

// ParamsReport holds validation results for triagegeist.Params.
type ParamsReport struct {
	Valid       bool
	WeightsOK   bool
	ThresholdsOK bool
	MaxResOK    bool
	ResourceWOK bool
}

// Params validates p (weights in [0,1], T1>T2>T3>T4, MaxResources>=0, ResourceWeight>=0).
// It does not import triagegeist to avoid cycle; the caller passes the struct.
type ParamsLike struct {
	VitalWeights   [7]float64
	MaxResources   int
	ResourceWeight float64
	T1, T2, T3, T4 float64
}

// Params validates a parameter set and returns a report.
func Params(p ParamsLike) ParamsReport {
	var r ParamsReport
	r.Valid = true

	for _, w := range p.VitalWeights {
		if w < 0 || w > 1 || !finite(w) {
			r.WeightsOK = false
			r.Valid = false
			break
		}
	}
	if r.Valid {
		r.WeightsOK = true
	}

	if p.MaxResources < 0 {
		r.MaxResOK = false
		r.Valid = false
	} else {
		r.MaxResOK = true
	}

	if p.ResourceWeight < 0 || !finite(p.ResourceWeight) {
		r.ResourceWOK = false
		r.Valid = false
	} else {
		r.ResourceWOK = true
	}

	if !(p.T1 > p.T2 && p.T2 > p.T3 && p.T3 > p.T4 && p.T4 > 0 && p.T1 <= 1) {
		r.ThresholdsOK = false
		r.Valid = false
	} else {
		r.ThresholdsOK = true
	}

	return r
}

// ParamsValid returns true if p is valid.
func ParamsValid(p ParamsLike) bool {
	return Params(p).Valid
}

// AtLeastOneVital returns true if at least one of HR, RR, SBP, DBP, Temp, SpO2, GCS is present (non-zero).
func AtLeastOneVital(v score.Vitals) bool {
	return v.HR > 0 || v.RR > 0 || v.SBP > 0 || v.DBP > 0 || v.Temp != 0 || v.SpO2 > 0 || v.GCS > 0
}

// VitalsAndResources returns a combined check: vitals valid and resourceCount in [0, maxResources].
func VitalsAndResources(v score.Vitals, resourceCount, maxResources int) bool {
	if !VitalsValid(v) {
		return false
	}
	if maxResources <= 0 {
		return resourceCount == 0
	}
	return resourceCount >= 0 && resourceCount <= maxResources
}

// SanitizeVitals returns clamped vitals and true if at least one vital was clamped; otherwise (v, false).
func SanitizeVitals(v score.Vitals) (score.Vitals, bool) {
	report := Vitals(v)
	if report.Valid {
		return v, false
	}
	clamped := ClampVitals(v)
	return clamped, true
}
