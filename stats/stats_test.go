package stats

import "testing"

func TestMean(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	if m := Mean(x); m != 3 {
		t.Errorf("Mean = %v, want 3", m)
	}
	if Mean(nil) != 0 {
		t.Error("Mean(nil) should be 0")
	}
}

func TestVariance(t *testing.T) {
	x := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	v := Variance(x)
	if v < 0 {
		t.Errorf("Variance = %v", v)
	}
}

func TestStdDev(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	_ = StdDev(x)
}

func TestCI95(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	lo, hi := CI95(x)
	if lo >= hi {
		t.Errorf("CI95: lo=%v hi=%v", lo, hi)
	}
}

func TestMedian(t *testing.T) {
	x := []float64{1, 3, 5}
	if m := Median(x); m != 3 {
		t.Errorf("Median = %v, want 3", m)
	}
}

func TestPercentile(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	p50 := Percentile(x, 50)
	if p50 < 1 || p50 > 5 {
		t.Errorf("P50 = %v", p50)
	}
}

func TestLevelDistribution(t *testing.T) {
	levels := []int{1, 2, 2, 3, 5}
	dist := LevelDistribution(levels)
	if dist[1] != 1 || dist[2] != 2 {
		t.Errorf("LevelDistribution = %v", dist)
	}
}

func TestComputeScoreStats(t *testing.T) {
	scores := []float64{0.2, 0.5, 0.8}
	s := ComputeScoreStats(scores)
	if s.N != 3 {
		t.Errorf("ComputeScoreStats N = %d", s.N)
	}
}

func TestExactAgreement(t *testing.T) {
	pred := []int{1, 2, 3}
	ref := []int{1, 2, 3}
	if a := ExactAgreement(pred, ref); a != 1.0 {
		t.Errorf("ExactAgreement = %v", a)
	}
}

func TestRMSE(t *testing.T) {
	pred := []float64{1, 2, 3}
	ref := []float64{1, 2, 3}
	if r := RMSE(pred, ref); r != 0 {
		t.Errorf("RMSE = %v", r)
	}
}
