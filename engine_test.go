package triagegeist

import (
	"testing"

	"github.com/olaflaitinen/triagegeist/score"
)

func TestEngine_AcuityAndLevel(t *testing.T) {
	p := DefaultParams()
	eng := NewEngine(p)

	v := score.Vitals{HR: 120, RR: 24, SBP: 90, SpO2: 92}
	acuity := eng.Acuity(v, 3)
	if acuity < 0 || acuity > 1 {
		t.Errorf("acuity %f not in [0,1]", acuity)
	}

	level := eng.Level(v, 3)
	if level < 1 || level > 5 {
		t.Errorf("level %d not in 1..5", level)
	}

	acuity2, level2 := eng.ScoreAndLevel(v, 3)
	if acuity2 != acuity || level2 != level {
		t.Errorf("ScoreAndLevel mismatch: got %.3f/%d, want %.3f/%d", acuity2, level2, acuity, level)
	}
}

func TestFromScore(t *testing.T) {
	p := DefaultParams()
	if s := FromScore(0.90, p); s != Level1Resuscitation {
		t.Errorf("FromScore(0.90) = %v, want Level1", s)
	}
	if s := FromScore(0.70, p); s != Level2Emergent {
		t.Errorf("FromScore(0.70) = %v, want Level2", s)
	}
	if s := FromScore(0.50, p); s != Level3Urgent {
		t.Errorf("FromScore(0.50) = %v, want Level3", s)
	}
	if s := FromScore(0.20, p); s != Level4LessUrgent {
		t.Errorf("FromScore(0.20) = %v, want Level4", s)
	}
	if s := FromScore(0.10, p); s != Level5NonUrgent {
		t.Errorf("FromScore(0.10) = %v, want Level5", s)
	}
}

func TestParams_Validate(t *testing.T) {
	p := DefaultParams()
	if !p.Validate() {
		t.Error("DefaultParams() should validate")
	}
	p.MaxResources = -1
	if p.Validate() {
		t.Error("MaxResources -1 should invalidate")
	}
}

var benchVitals = score.Vitals{HR: 120, RR: 24, SBP: 90, DBP: 60, SpO2: 92}
const benchResources = 3

func BenchmarkEngine_ScoreAndLevel(b *testing.B) {
	p := DefaultParams()
	eng := NewEngine(p)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = eng.ScoreAndLevel(benchVitals, benchResources)
	}
}

func BenchmarkEngine_Acuity(b *testing.B) {
	p := DefaultParams()
	eng := NewEngine(p)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = eng.Acuity(benchVitals, benchResources)
	}
}
