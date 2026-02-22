package metrics

import "testing"

func TestNewConfusionMatrix(t *testing.T) {
	pred := []int{1, 2, 3, 1, 2}
	ref := []int{1, 2, 3, 2, 2}
	cm := NewConfusionMatrix(pred, ref)
	if cm.Total != 5 {
		t.Errorf("Total = %d, want 5", cm.Total)
	}
	if cm.TP(1) != 1 {
		t.Errorf("TP(1) = %d, want 1", cm.TP(1))
	}
	if cm.TP(2) != 2 {
		t.Errorf("TP(2) = %d, want 2", cm.TP(2))
	}
}

func TestConfusionMatrix_Perfect(t *testing.T) {
	x := []int{1, 2, 3, 4, 5}
	cm := NewConfusionMatrix(x, x)
	if cm.OverallAccuracy() != 1.0 {
		t.Errorf("perfect agreement: accuracy = %v", cm.OverallAccuracy())
	}
	if cm.CohenKappa() != 1.0 {
		t.Errorf("perfect agreement: kappa = %v", cm.CohenKappa())
	}
}

func TestConfusionMatrix_SensitivitySpecificity(t *testing.T) {
	pred := []int{1, 1, 2, 2}
	ref := []int{1, 2, 1, 2}
	cm := NewConfusionMatrix(pred, ref)
	sens := cm.Sensitivity(1)
	if sens < 0 || sens > 1 {
		t.Errorf("Sensitivity(1) = %v", sens)
	}
	spec := cm.Specificity(1)
	if spec < 0 || spec > 1 {
		t.Errorf("Specificity(1) = %v", spec)
	}
}

func TestBinaryCM(t *testing.T) {
	pred := []int{1, 2, 3, 4, 5}
	ref := []int{1, 2, 3, 4, 5}
	b := NewBinaryCM(pred, ref, []int{1, 2})
	if b.TP+b.TN+b.FP+b.FN != 5 {
		t.Errorf("binary total should be 5, got %d", b.TP+b.TN+b.FP+b.FN)
	}
	acc := b.Accuracy()
	if acc < 0 || acc > 1 {
		t.Errorf("Accuracy = %v", acc)
	}
}

func TestAUC(t *testing.T) {
	scores := []float64{0.1, 0.3, 0.5, 0.7, 0.9}
	outcomes := []int{0, 0, 1, 1, 1}
	a := AUC(scores, outcomes)
	if a < 0 || a > 1 {
		t.Errorf("AUC = %v", a)
	}
}

func TestCalibrationError(t *testing.T) {
	scores := []float64{0.5, 0.5}
	outcomes := []int{0, 1}
	e := CalibrationError(scores, outcomes)
	if e < 0 || e > 1 {
		t.Errorf("CalibrationError = %v", e)
	}
}

func TestWeightedKappa(t *testing.T) {
	pred := []int{1, 2, 3, 4, 5}
	ref := []int{1, 2, 3, 4, 5}
	k := WeightedKappa(pred, ref)
	if k != 1.0 {
		t.Errorf("perfect agreement: weighted kappa = %v", k)
	}
}
