package norm

import "testing"

func TestDefaultRanges(t *testing.T) {
	r := DefaultRanges()
	if !r.Valid() {
		t.Error("DefaultRanges should be valid")
	}
	if r.HR[0] != 80 || r.HR[1] != 40 {
		t.Errorf("HR: got [%v, %v]", r.HR[0], r.HR[1])
	}
}

func TestDeviation(t *testing.T) {
	if d := Deviation(80, 80, 40); d != 0 {
		t.Errorf("Deviation(80,80,40) = %v, want 0", d)
	}
	if d := Deviation(120, 80, 40); d != 1 {
		t.Errorf("Deviation(120,80,40) = %v, want 1", d)
	}
	if d := Deviation(100, 80, 40); d != 0.5 {
		t.Errorf("Deviation(100,80,40) = %v, want 0.5", d)
	}
	if d := Deviation(50, 80, 0); d != 0 {
		t.Errorf("Deviation with halfWidth 0 should be 0, got %v", d)
	}
}

func TestNormalizeLinear(t *testing.T) {
	if n := NormalizeLinear(50, 0, 100); n != 0.5 {
		t.Errorf("NormalizeLinear(50,0,100) = %v", n)
	}
	if n := NormalizeLinear(-1, 0, 100); n != 0 {
		t.Errorf("below range should clamp to 0, got %v", n)
	}
	if n := NormalizeLinear(150, 0, 100); n != 1 {
		t.Errorf("above range should clamp to 1, got %v", n)
	}
}

func TestClampToRange(t *testing.T) {
	if c := ClampToRange(50, 0, 100); c != 50 {
		t.Errorf("ClampToRange(50,0,100) = %v", c)
	}
	if c := ClampToRange(-1, 0, 100); c != 0 {
		t.Errorf("ClampToRange(-1,0,100) = %v", c)
	}
	if c := ClampToRange(101, 0, 100); c != 100 {
		t.Errorf("ClampToRange(101,0,100) = %v", c)
	}
}

func TestRanges_At_Set(t *testing.T) {
	var r Ranges
	r.Set(0, 90, 45)
	mid, hw := r.At(0)
	if mid != 90 || hw != 45 {
		t.Errorf("At(0) = %v, %v", mid, hw)
	}
	_, _ = r.At(7)
	_, _ = r.At(-1)
}

func TestCriticalBounds(t *testing.T) {
	lo, hi := CriticalBounds(VitalHR)
	if lo != 20 || hi != 300 {
		t.Errorf("HR bounds: got %v, %v", lo, hi)
	}
	lo, hi = CriticalBounds(10)
	if lo != 0 || hi != 0 {
		t.Errorf("invalid index should return 0,0: got %v, %v", lo, hi)
	}
}

func TestRanges_WeightedDeviationSum(t *testing.T) {
	r := DefaultRanges()
	values := [7]float64{80, 16, 120, 80, 37, 98, 15}
	weights := [7]float64{0.2, 0.2, 0.2, 0.2, 0.1, 0.1, 0}
	sum, wSum := r.WeightedDeviationSum(values, weights)
	if sum != 0 {
		t.Errorf("normal values should give 0 deviation sum, got %v", sum)
	}
	if wSum <= 0 {
		t.Errorf("weight sum should be positive, got %v", wSum)
	}
}
