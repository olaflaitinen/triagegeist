package validate

import (
	"testing"

	"github.com/olaflaitinen/triagegeist/score"
)

func TestVitals_Valid(t *testing.T) {
	v := score.Vitals{HR: 80, RR: 16, SBP: 120, SpO2: 98}
	report := Vitals(v)
	if !report.Valid {
		t.Error("normal vitals should be valid")
	}
}

func TestVitals_Invalid(t *testing.T) {
	v := score.Vitals{HR: 500}
	report := Vitals(v)
	if report.Valid {
		t.Error("HR 500 should be invalid")
	}
}

func TestClampVitals(t *testing.T) {
	v := score.Vitals{HR: 500, RR: 16}
	clamped := ClampVitals(v)
	if clamped.HR != 300 {
		t.Errorf("HR 500 should clamp to 300, got %d", clamped.HR)
	}
	if clamped.RR != 16 {
		t.Errorf("RR 16 should stay 16, got %d", clamped.RR)
	}
}

func TestResourceCount(t *testing.T) {
	if c := ResourceCount(5, 3); c != 3 {
		t.Errorf("ResourceCount(5,3) = %d, want 3", c)
	}
	if c := ResourceCount(-1, 5); c != 0 {
		t.Errorf("ResourceCount(-1,5) = %d, want 0", c)
	}
	if c := ResourceCount(2, 5); c != 2 {
		t.Errorf("ResourceCount(2,5) = %d, want 2", c)
	}
}

func TestParams(t *testing.T) {
	pl := ParamsLike{
		VitalWeights:   [7]float64{0.2, 0.2, 0.2, 0.1, 0.1, 0.1, 0.1},
		MaxResources:   6,
		ResourceWeight: 0.25,
		T1:             0.85, T2: 0.6, T3: 0.35, T4: 0.15,
	}
	report := Params(pl)
	if !report.Valid {
		t.Errorf("Params report: %+v", report)
	}
	pl.MaxResources = -1
	report = Params(pl)
	if report.Valid {
		t.Error("MaxResources -1 should be invalid")
	}
}

func TestAtLeastOneVital(t *testing.T) {
	if AtLeastOneVital(score.Vitals{}) {
		t.Error("empty vitals should have none present")
	}
	if !AtLeastOneVital(score.Vitals{HR: 80}) {
		t.Error("HR present should be at least one")
	}
}
