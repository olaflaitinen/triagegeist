package triagegeist_test

import (
	"fmt"

	"github.com/olaflaitinen/triagegeist"
	"github.com/olaflaitinen/triagegeist/score"
)

func ExampleEngine_ScoreAndLevel() {
	p := triagegeist.DefaultParams()
	eng := triagegeist.NewEngine(p)

	v := score.Vitals{
		HR:   120,
		RR:   24,
		SBP:  90,
		SpO2: 92,
	}
	resourceCount := 3

	acuity, level := eng.ScoreAndLevel(v, resourceCount)
	fmt.Printf("acuity: %.3f, level: %d (%s)\n", acuity, level, level.String())
}

func ExampleFromScore() {
	p := triagegeist.DefaultParams()
	level := triagegeist.FromScore(0.72, p)
	fmt.Println(level.String())
	// Output: Emergent
}
