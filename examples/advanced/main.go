// Copyright (c) triagegeist authors: Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov.
// Licensed under the EUPL.
//
// Advanced example: batch evaluation, metrics, statistics, and export.
// Demonstrates the full pipeline for research or auditing.
//
// Run from repository root: go run ./examples/advanced
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/olaflaitinen/triagegeist"
	"github.com/olaflaitinen/triagegeist/export"
	"github.com/olaflaitinen/triagegeist/metrics"
	"github.com/olaflaitinen/triagegeist/score"
	"github.com/olaflaitinen/triagegeist/stats"
	"github.com/olaflaitinen/triagegeist/validate"
)

func main() {
	eng := triagegeist.NewDefaultEngine()
	p := eng.Params()

	// Synthetic vitals and resource counts (e.g. from a study or audit)
	vitals, resourceCounts := sampleData()

	// Validate and clamp where needed
	for i := range vitals {
		if !validate.VitalsValid(vitals[i]) {
			vitals[i] = validate.ClampVitals(vitals[i])
		}
		resourceCounts[i] = validate.ResourceCount(resourceCounts[i], p.MaxResources)
	}

	// Batch evaluation
	acuities, levels := eng.BatchScoreAndLevel(vitals, resourceCounts)
	if acuities == nil {
		log.Fatal("batch length mismatch")
	}

	// Build results for export and metrics
	results := make([]export.Result, len(acuities))
	for i := range acuities {
		results[i] = export.FromVitalsScoreLevel(
			vitals[i], resourceCounts[i],
			acuities[i], levels[i].Int(), levels[i].String(),
		)
	}

	// Descriptive statistics
	scoreStats := stats.ComputeScoreStats(acuities)
	fmt.Println("--- Acuity score statistics ---")
	fmt.Printf("N: %d, Mean: %.4f, StdDev: %.4f\n", scoreStats.N, scoreStats.Mean, scoreStats.StdDev)
	fmt.Printf("95%% CI: [%.4f, %.4f]\n", scoreStats.CI95Lo, scoreStats.CI95Hi)
	fmt.Printf("Min: %.4f, Max: %.4f, P25: %.4f, P50: %.4f, P75: %.4f\n",
		scoreStats.Min, scoreStats.Max, scoreStats.P25, scoreStats.P50, scoreStats.P75)

	levelInts := triagegeist.IntsFromLevels(levels)
	levelStats := stats.ComputeLevelStats(levelInts)
	fmt.Println("--- Level distribution ---")
	for i := 1; i <= 5; i++ {
		fmt.Printf("Level %d: count=%d, pct=%.1f\n", i, levelStats.Counts[i], levelStats.Props[i]*100)
	}

	// If we had reference (ground truth) levels, we could compute metrics
	// For demo, use levels as both predicted and reference (perfect agreement)
	cm := metrics.NewConfusionMatrix(levelInts, levelInts)
	fmt.Println("--- Agreement (predicted vs reference, demo) ---")
	fmt.Printf("Overall accuracy: %.4f\n", cm.OverallAccuracy())
	fmt.Printf("Cohen's kappa: %.4f\n", cm.CohenKappa())
	fmt.Printf("Macro sensitivity: %.4f, macro specificity: %.4f\n", cm.MacroSensitivity(), cm.MacroSpecificity())

	// Binary classification: high acuity (1-2) vs low (3-5)
	binary := metrics.NewBinaryCM(levelInts, levelInts, []int{1, 2})
	fmt.Println("--- Binary (high vs low acuity) ---")
	fmt.Printf("Sensitivity: %.4f, Specificity: %.4f, PPV: %.4f, NPV: %.4f, F1: %.4f\n",
		binary.Sensitivity(), binary.Specificity(), binary.PPV(), binary.NPV(), binary.F1())

	// Weighted kappa (adjacent level agreement)
	wk := metrics.WeightedKappa(levelInts, levelInts)
	fmt.Printf("Weighted kappa (adjacent): %.4f\n", wk)

	// Exact and within-one-level agreement (if we had a reference)
	exact := stats.ExactAgreement(levelInts, levelInts)
	within1 := stats.WithinLevel(levelInts, levelInts)
	fmt.Printf("Exact agreement: %.4f, within 1 level: %.4f\n", exact, within1)

	// Summary struct and optional file export
	sum := export.ComputeSummary(results)
	fmt.Println("--- Export summary ---")
	fmt.Printf("N=%d, mean_acuity=%.4f, min=%.4f, max=%.4f\n", sum.N, sum.MeanAcuity, sum.MinAcuity, sum.MaxAcuity)
	for i := 1; i <= 5; i++ {
		fmt.Printf("Level %d count: %d\n", i, sum.LevelDist[i])
	}

	// Level report (counts and acuity stats per level)
	reportRows := export.LevelReport(results)
	fmt.Println("--- Level report ---")
	for _, row := range reportRows {
		fmt.Printf("L%d %s: n=%d, mean_acuity=%.4f\n", row.Level, row.LevelLabel, row.Count, row.MeanAcuity)
	}

	// Write CSV to stdout or file
	if err := export.WriteCSV(os.Stdout, results); err != nil {
		log.Printf("WriteCSV: %v", err)
	}
}

func sampleData() ([]score.Vitals, []int) {
	vitals := []score.Vitals{
		{HR: 120, RR: 24, SBP: 90, SpO2: 92},
		{HR: 80, RR: 16, SBP: 120, DBP: 80, SpO2: 98},
		{HR: 140, RR: 28, SBP: 85, SpO2: 88},
		{HR: 70, RR: 14, SBP: 130, Temp: 36.8, SpO2: 99},
		{HR: 100, RR: 20, SBP: 100, SpO2: 94},
		{HR: 90, RR: 18, SBP: 115, DBP: 75, Temp: 37.0, SpO2: 97},
		{HR: 130, RR: 26, SBP: 88, SpO2: 91},
		{HR: 75, RR: 12, SBP: 125, SpO2: 98},
		{HR: 110, RR: 22, SBP: 95, SpO2: 93},
		{HR: 85, RR: 16, SBP: 118, DBP: 78, Temp: 36.9, SpO2: 98},
		{HR: 95, RR: 18, SBP: 112, DBP: 72, SpO2: 96},
		{HR: 125, RR: 24, SBP: 92, SpO2: 91},
		{HR: 78, RR: 15, SBP: 122, Temp: 37.1, SpO2: 98},
		{HR: 105, RR: 20, SBP: 98, SpO2: 95},
		{HR: 88, RR: 17, SBP: 116, DBP: 76, SpO2: 97},
		{HR: 135, RR: 27, SBP: 86, SpO2: 89},
		{HR: 72, RR: 13, SBP: 128, SpO2: 99},
		{HR: 115, RR: 22, SBP: 96, SpO2: 94},
		{HR: 82, RR: 15, SBP: 119, DBP: 77, Temp: 36.7, SpO2: 98},
		{HR: 128, RR: 25, SBP: 89, SpO2: 90},
	}
	resourceCounts := []int{3, 1, 4, 0, 2, 1, 4, 0, 2, 1, 1, 3, 0, 2, 2, 4, 0, 3, 1, 3}
	for len(resourceCounts) < len(vitals) {
		resourceCounts = append(resourceCounts, 0)
	}
	resourceCounts = resourceCounts[:len(vitals)]
	return vitals, resourceCounts
}

