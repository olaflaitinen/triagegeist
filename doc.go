// Copyright (c) triagegeist authors: Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov.
// Licensed under the EUPL.
//
// Package triagegeist provides a high-performance, parametric toolkit for
// AI-assisted emergency medicine triage and acuity scoring.
//
// The library is designed for low latency and minimal allocations: all
// core scoring is formula-based and deterministic. Optional model interfaces
// allow integration with external predictors (e.g. ONNX, TensorFlow Lite)
// without pulling heavy ML dependencies.
//
// # Subpackages
//
//	| Package   | Purpose                                                                 |
//	|-----------|-------------------------------------------------------------------------|
//	| score     | Acuity formula, Vitals struct, VitalComponent, ResourceComponent, Acuity, AcuityWithNorms, default norms and weights. |
//	| norm      | Reference ranges (Ranges), Deviation, NormalizeLinear, ClampToRange, CriticalBounds, WeightedDeviationSum, DefaultRanges, PediatricRanges. |
//	| metrics   | ConfusionMatrix, TP/FP/FN/TN, Sensitivity, Specificity, PPV, NPV, F1, CohenKappa, BinaryCM, AUC, CalibrationError, WeightedKappa. |
//	| stats     | Mean, Variance, StdDev, SE, CI95, Median, Percentile, LevelDistribution, ScoreStats, LevelStats, RMSE, MAE, ExactAgreement, WithinLevel. |
//	| validate  | Vitals validation (Vitals, ClampVitals, VitalsValid), ResourceCount, Params validation (ParamsLike, Params, ParamsValid), AtLeastOneVital. |
//	| export    | Result struct, FromVitalsScoreLevel, ToJSON, CSVHeader, ToCSVRow, WriteCSV, Batch, LevelReport, ComputeSummary, ReadResultJSON, ResultToVitals. |
//
// # Acuity score
//
// The parametric acuity score combines vital-sign deviations and expected
// resource consumption into a single continuous index. For level assignment,
// thresholds are applied as in the following table:
//
//	| Level | Label         | Score range (s)     | Typical wait |
//	|-------|---------------|---------------------|---------------|
//	| 1     | Resuscitation | s >= 0.85           | Immediate     |
//	| 2     | Emergent      | 0.60 <= s < 0.85    | < 15 min      |
//	| 3     | Urgent        | 0.35 <= s < 0.60    | < 60 min      |
//	| 4     | Less urgent   | 0.15 <= s < 0.35    | < 120 min     |
//	| 5     | Non-urgent    | s < 0.15            | < 240 min     |
//
// The normalized score s is computed from vitals and resource count; see
// [score.Acuity] and [Params] for the exact formula and parameters.
//
// # Formula (summary)
//
// Vital deviation: d_i = min(1, |x_i - mu_i| / sigma_i).
// Vital component: V = (sum w_i d_i) / (sum w_i) over present vitals.
// Resource component: R = alpha * min(1, resourceCount / maxResources).
// Raw = V + R; s = Raw / (sum w_i + alpha), clamped to [0, 1].
//
// # Default weights (VitalWeights)
//
//	| Index | Vital | Default weight |
//	|-------|-------|----------------|
//	| 0     | HR    | 0.18           |
//	| 1     | RR    | 0.22           |
//	| 2     | SBP   | 0.16           |
//	| 3     | DBP   | 0.10           |
//	| 4     | Temp  | 0.08           |
//	| 5     | SpO2  | 0.16           |
//	| 6     | GCS   | 0.10           |
//
// # Default thresholds
//
//	| Param | Value | Level boundary        |
//	|-------|-------|------------------------|
//	| T1    | 0.85  | s >= T1 => Level 1     |
//	| T2    | 0.60  | T2 <= s < T1 => Level 2 |
//	| T3    | 0.35  | T3 <= s < T2 => Level 3 |
//	| T4    | 0.15  | T4 <= s < T3 => Level 4 |
//	|       |       | s < T4 => Level 5       |
//
// # Example
//
//	p := triagegeist.DefaultParams()
//	eng := triagegeist.NewEngine(p)
//	v := score.Vitals{HR: 120, RR: 24, SBP: 90, SpO2: 92}
//	acuity, level := eng.ScoreAndLevel(v, 3)
//
// # Validation
//
// Use the validate package to check vitals and params before calling the engine.
// Use the metrics package to compute sensitivity, specificity, and other
// accuracy metrics when reference (ground truth) levels are available.
//
// # Licence and authors
//
// triagegeist is licensed under the European Union Public Licence v. 1.2 (EUPL-1.2).
// Authors: Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov.
package triagegeist
