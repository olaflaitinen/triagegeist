package score

import (
	"testing"
)

func TestVitalComponent(t *testing.T) {
	v := Vitals{HR: 80, RR: 16, SBP: 120}
	w := VitalWeights
	c := VitalComponent(v, w)
	if c < 0 || c > 1 {
		t.Errorf("VitalComponent(normal vitals) = %f, want in [0,1]", c)
	}

	v2 := Vitals{HR: 160, RR: 32, SpO2: 80}
	c2 := VitalComponent(v2, w)
	if c2 <= c {
		t.Errorf("higher-deviation vitals should give higher component: got %f <= %f", c2, c)
	}
}

func TestAcuity(t *testing.T) {
	v := Vitals{HR: 120, RR: 24, SBP: 90, SpO2: 92}
	s := Acuity(v, 3, 6, VitalWeights, 0.25)
	if s < 0 || s > 1 {
		t.Errorf("Acuity = %f, want in [0,1]", s)
	}
}

func TestNormalize(t *testing.T) {
	if n := Normalize(0.5, 1.0); n != 0.5 {
		t.Errorf("Normalize(0.5, 1) = %f", n)
	}
	if n := Normalize(1.5, 1.0); n != 1.0 {
		t.Errorf("Normalize(1.5, 1) = %f, want 1", n)
	}
	if n := Normalize(-0.1, 1.0); n != 0 {
		t.Errorf("Normalize(-0.1, 1) = %f, want 0", n)
	}
}

var benchVitals = Vitals{HR: 120, RR: 24, SBP: 90, SpO2: 92}

func BenchmarkScore_Acuity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Acuity(benchVitals, 3, 6, VitalWeights, 0.25)
	}
}
