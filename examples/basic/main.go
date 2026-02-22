// Copyright (c) triagegeist authors: Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov.
// Licensed under the EUPL.
//
// Basic example: single evaluation with default parameters, validation, and export.
//
// Run from repository root:
//
//	go run ./examples/basic
package main

import (
	"fmt"
	"log"

	"github.com/olaflaitinen/triagegeist"
	"github.com/olaflaitinen/triagegeist/export"
	"github.com/olaflaitinen/triagegeist/score"
	"github.com/olaflaitinen/triagegeist/validate"
)

func main() {
	// 1. Build parameters and engine
	p := triagegeist.DefaultParams()
	if !p.Validate() {
		log.Fatal("default params should be valid")
	}
	eng := triagegeist.NewEngine(p)

	// 2. Prepare vitals (0 = missing)
	v := score.Vitals{
		HR:   120,
		RR:   24,
		SBP:  90,
		DBP:  60,
		SpO2: 92,
	}
	resourceCount := 3

	// 3. Validate inputs (optional but recommended)
	report := validate.Vitals(v)
	if !report.Valid {
		v = validate.ClampVitals(v)
	}
	resourceCount = validate.ResourceCount(resourceCount, p.MaxResources)

	// 4. Evaluate
	acuity, level := eng.ScoreAndLevel(v, resourceCount)

	// 5. Output
	fmt.Printf("Acuity: %.4f\n", acuity)
	fmt.Printf("Level:  %d (%s)\n", level.Int(), level.String())
	fmt.Printf("Wait:   %d min (guidance)\n", level.WaitTimeMinutes())

	// 6. Export to struct for JSON/CSV
	res := export.FromVitalsScoreLevel(v, resourceCount, acuity, level.Int(), level.String())
	_ = res
}
