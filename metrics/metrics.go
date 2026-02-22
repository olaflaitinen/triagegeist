// Copyright (c) triagegeist authors: Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov.
// Licensed under the EUPL.
//
// Package metrics provides accuracy and performance metrics for triage
// level predictions versus reference (ground truth) levels. Use for
// validation, calibration studies, and model comparison.
//
// # Metrics overview
//
//	| Metric        | Formula                    | Use case                |
//	|---------------|----------------------------|-------------------------|
//	| Sensitivity   | TP / (TP + FN)              | Detection of high acuity |
//	| Specificity   | TN / (TN + FP)             | Avoiding over-triage     |
//	| PPV           | TP / (TP + FP)             | Precision (positive)    |
//	| NPV           | TN / (TN + FN)             | Precision (negative)    |
//	| F1            | 2*PPV*Sens / (PPV+Sens)    | Balance                  |
//	| Accuracy      | (TP+TN) / N                | Overall agreement        |
//	| Cohen's Kappa  | (p_o - p_e) / (1 - p_e)    | Agreement vs chance      |
//
// All metrics return values in [0, 1] where applicable; callers must
// provide counts or slices of equal length (predicted, reference).
package metrics

import "math"

// ConfusionMatrix holds counts for a binary or multi-class classification.
// Rows = reference (true) class, Cols = predicted class. Level 1..5 map to
// indices 0..4. For binary (e.g. high acuity vs low), use BinaryCM.
type ConfusionMatrix struct {
	// N [i][j] = count where true level was i+1 and predicted was j+1
	N [5][5]int
	// Total number of samples
	Total int
}

// NewConfusionMatrix builds a 5x5 matrix from paired predicted and reference
// levels. Levels must be 1..5; others are skipped (not counted in Total).
func NewConfusionMatrix(predicted, reference []int) ConfusionMatrix {
	var cm ConfusionMatrix
	if len(predicted) != len(reference) {
		return cm
	}
	for k := range predicted {
		p, r := predicted[k], reference[k]
		if p < 1 || p > 5 || r < 1 || r > 5 {
			continue
		}
		cm.N[r-1][p-1]++
		cm.Total++
	}
	return cm
}

// TP returns true positives for the given class (1..5) when that class
// is treated as positive and the rest as negative.
func (cm ConfusionMatrix) TP(class int) int {
	if class < 1 || class > 5 {
		return 0
	}
	i := class - 1
	return cm.N[i][i]
}

// FP returns false positives for the given class (predicted=class, ref!=class).
func (cm ConfusionMatrix) FP(class int) int {
	if class < 1 || class > 5 {
		return 0
	}
	i := class - 1
	var fp int
	for r := 0; r < 5; r++ {
		if r != i {
			fp += cm.N[r][i]
		}
	}
	return fp
}

// FN returns false negatives for the given class (ref=class, pred!=class).
func (cm ConfusionMatrix) FN(class int) int {
	if class < 1 || class > 5 {
		return 0
	}
	i := class - 1
	var fn int
	for c := 0; c < 5; c++ {
		if c != i {
			fn += cm.N[i][c]
		}
	}
	return fn
}

// TN returns true negatives for the given class (ref!=class, pred!=class).
func (cm ConfusionMatrix) TN(class int) int {
	if class < 1 || class > 5 {
		return 0
	}
	i := class - 1
	var tn int
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if r != i && c != i {
				tn += cm.N[r][c]
			}
		}
	}
	return tn
}

// Sensitivity returns TP / (TP + FN) for the given class. Returns 0 if denominator 0.
func (cm ConfusionMatrix) Sensitivity(class int) float64 {
	tp := cm.TP(class)
	fn := cm.FN(class)
	if tp+fn == 0 {
		return 0
	}
	return float64(tp) / float64(tp+fn)
}

// Specificity returns TN / (TN + FP) for the given class. Returns 0 if denominator 0.
func (cm ConfusionMatrix) Specificity(class int) float64 {
	tn := cm.TN(class)
	fp := cm.FP(class)
	if tn+fp == 0 {
		return 0
	}
	return float64(tn) / float64(tn+fp)
}

// PPV returns positive predictive value TP / (TP + FP). Returns 0 if denominator 0.
func (cm ConfusionMatrix) PPV(class int) float64 {
	tp := cm.TP(class)
	fp := cm.FP(class)
	if tp+fp == 0 {
		return 0
	}
	return float64(tp) / float64(tp+fp)
}

// NPV returns negative predictive value TN / (TN + FN). Returns 0 if denominator 0.
func (cm ConfusionMatrix) NPV(class int) float64 {
	tn := cm.TN(class)
	fn := cm.FN(class)
	if tn+fn == 0 {
		return 0
	}
	return float64(tn) / float64(tn+fn)
}

// F1 returns 2 * PPV * Sensitivity / (PPV + Sensitivity) for the class. 0 if denominator 0.
func (cm ConfusionMatrix) F1(class int) float64 {
	ppv := cm.PPV(class)
	sens := cm.Sensitivity(class)
	if ppv+sens == 0 {
		return 0
	}
	return 2 * ppv * sens / (ppv + sens)
}

// Accuracy returns (TP+TN) / Total when treating the given class as positive.
func (cm ConfusionMatrix) Accuracy(class int) float64 {
	if cm.Total == 0 {
		return 0
	}
	return float64(cm.TP(class)+cm.TN(class)) / float64(cm.Total)
}

// MacroSensitivity returns the mean of Sensitivity(1)..Sensitivity(5).
func (cm ConfusionMatrix) MacroSensitivity() float64 {
	var sum float64
	for c := 1; c <= 5; c++ {
		sum += cm.Sensitivity(c)
	}
	return sum / 5
}

// MacroSpecificity returns the mean of Specificity(1)..Specificity(5).
func (cm ConfusionMatrix) MacroSpecificity() float64 {
	var sum float64
	for c := 1; c <= 5; c++ {
		sum += cm.Specificity(c)
	}
	return sum / 5
}

// OverallAccuracy returns the fraction of correct predictions (diagonal / Total).
func (cm ConfusionMatrix) OverallAccuracy() float64 {
	if cm.Total == 0 {
		return 0
	}
	var diag int
	for i := 0; i < 5; i++ {
		diag += cm.N[i][i]
	}
	return float64(diag) / float64(cm.Total)
}

// CohenKappa returns Cohen's kappa for agreement between predicted and reference
// levels (1..5). Returns 0 if Total is 0 or agreement is undefined.
func (cm ConfusionMatrix) CohenKappa() float64 {
	if cm.Total == 0 {
		return 0
	}
	var pObs float64
	for i := 0; i < 5; i++ {
		pObs += float64(cm.N[i][i])
	}
	pObs /= float64(cm.Total)
	var sumPred, sumRef [5]float64
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			sumPred[j] += float64(cm.N[i][j])
			sumRef[i] += float64(cm.N[i][j])
		}
	}
	t := float64(cm.Total)
	var pExp float64
	for i := 0; i < 5; i++ {
		pExp += sumPred[i] * sumRef[i] / (t * t)
	}
	if pExp >= 1 {
		return 0
	}
	return (pObs - pExp) / (1 - pExp)
}

// BinaryCM is a 2x2 confusion matrix for binary classification (e.g. high acuity 1-2 vs low 3-5).
type BinaryCM struct {
	TP, FP, FN, TN int
}

// NewBinaryCM builds a binary matrix by treating classes in positive as positive,
// e.g. positive = [1, 2] for "high acuity".
func NewBinaryCM(predicted, reference []int, positive []int) BinaryCM {
	if len(predicted) != len(reference) {
		return BinaryCM{}
	}
	posSet := make(map[int]bool)
	for _, p := range positive {
		posSet[p] = true
	}
	var b BinaryCM
	for k := range predicted {
		p, r := predicted[k], reference[k]
		if p < 1 || p > 5 || r < 1 || r > 5 {
			continue
		}
		pPos := posSet[p]
		rPos := posSet[r]
		if rPos && pPos {
			b.TP++
		} else if !rPos && pPos {
			b.FP++
		} else if rPos && !pPos {
			b.FN++
		} else {
			b.TN++
		}
	}
	return b
}

// Sensitivity returns TP / (TP + FN).
func (b BinaryCM) Sensitivity() float64 {
	d := b.TP + b.FN
	if d == 0 {
		return 0
	}
	return float64(b.TP) / float64(d)
}

// Specificity returns TN / (TN + FP).
func (b BinaryCM) Specificity() float64 {
	d := b.TN + b.FP
	if d == 0 {
		return 0
	}
	return float64(b.TN) / float64(d)
}

// PPV returns TP / (TP + FP).
func (b BinaryCM) PPV() float64 {
	d := b.TP + b.FP
	if d == 0 {
		return 0
	}
	return float64(b.TP) / float64(d)
}

// NPV returns TN / (TN + FN).
func (b BinaryCM) NPV() float64 {
	d := b.TN + b.FN
	if d == 0 {
		return 0
	}
	return float64(b.TN) / float64(d)
}

// F1 returns 2*PPV*Sensitivity / (PPV + Sensitivity).
func (b BinaryCM) F1() float64 {
	s, p := b.Sensitivity(), b.PPV()
	if s+p == 0 {
		return 0
	}
	return 2 * s * p / (s + p)
}

// Accuracy returns (TP+TN) / (TP+TN+FP+FN).
func (b BinaryCM) Accuracy() float64 {
	total := b.TP + b.TN + b.FP + b.FN
	if total == 0 {
		return 0
	}
	return float64(b.TP+b.TN) / float64(total)
}

// AUC trapezoidal from sorted (score, binary outcome) pairs.
// scores and outcomes must have same length; outcomes are 0 or 1.
// Higher score should correspond to positive (1). Returns value in [0, 1].
func AUC(scores []float64, outcomes []int) float64 {
	if len(scores) != len(outcomes) || len(scores) == 0 {
		return 0
	}
	// Sort by score ascending and count positives
	type pair struct {
		s float64
		o int
	}
	pairs := make([]pair, len(scores))
	for i := range scores {
		pairs[i] = pair{scores[i], outcomes[i]}
	}
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].s < pairs[i].s {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}
	var pos int
	for _, p := range pairs {
		if p.o == 1 {
			pos++
		}
	}
	neg := len(pairs) - pos
	if pos == 0 || neg == 0 {
		return 0.5
	}
	var sum float64
	var cumPos int
	for i, p := range pairs {
		if p.o == 1 {
			cumPos++
			sum += float64(i - cumPos + 1)
		}
	}
	return sum / float64(pos*neg)
}

// CalibrationError returns mean absolute error between predicted scores and
// observed binary outcomes. scores and outcomes same length; outcomes 0 or 1.
func CalibrationError(scores []float64, outcomes []int) float64 {
	if len(scores) != len(outcomes) || len(scores) == 0 {
		return 0
	}
	var sum float64
	for i := range scores {
		out := 0.0
		if outcomes[i] == 1 {
			out = 1
		}
		s := scores[i]
		if s < 0 {
			s = 0
		}
		if s > 1 {
			s = 1
		}
		sum += math.Abs(s - out)
	}
	return sum / float64(len(scores))
}

// kappaWeight computes linear weight 1 - |p-r|/4
func kappaWeight(p, r int) float64 {
	w := 1 - math.Abs(float64(p-r))/4
	if w < 0 {
		return 0
	}
	return w
}

// clampLevel ensures level is 1..5, defaulting to 3
func clampLevel(L int) int {
	if L < 1 || L > 5 {
		return 3
	}
	return L
}

// WeightedKappa returns linear weighted kappa with unit weights for adjacent
// level difference. pred and ref are level 1..5; equal length.
func WeightedKappa(pred, ref []int) float64 {
	if len(pred) != len(ref) || len(pred) == 0 {
		return 0
	}
	n := float64(len(pred))
	var obsWeight, expWeight float64
	for i := range pred {
		p, r := clampLevel(pred[i]), clampLevel(ref[i])
		obsWeight += kappaWeight(p, r)
	}
	obsWeight /= n
	countP, countR := [6]float64{}, [6]float64{}
	for i := range pred {
		p, r := pred[i], ref[i]
		if p >= 1 && p <= 5 {
			countP[p]++
		}
		if r >= 1 && r <= 5 {
			countR[r]++
		}
	}
	for i := 1; i <= 5; i++ {
		for j := 1; j <= 5; j++ {
			expWeight += (countP[i] / n) * (countR[j] / n) * kappaWeight(i, j)
		}
	}
	if expWeight >= 1 {
		return 0
	}
	return (obsWeight - expWeight) / (1 - expWeight)
}
