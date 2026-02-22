// Copyright (c) triagegeist authors: Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov.
// Licensed under the EUPL.

package triagegeist

// Level is the discrete triage level (1 = highest acuity, 5 = lowest).
// Five-level systems are common in emergency departments (e.g. ESI, CTAS, MTS).
// This type does not implement any proprietary algorithm; it represents the
// outcome of thresholding a continuous acuity score as defined in this package.
//
// Assignment is by thresholding the normalised acuity score s:
//
//   - Level 1: s >= T1
//
//   - Level 2: T2 <= s < T1
//
//   - Level 3: T3 <= s < T2
//
//   - Level 4: T4 <= s < T3
//
//   - Level 5: s < T4
//
//     | Level | Constant            | Typical meaning        | Wait (min) |
//     |-------|---------------------|-------------------------|------------|
//     | 1     | Level1Resuscitation | Immediate life-saving   | 0          |
//     | 2     | Level2Emergent      | Emergent, high risk     | 15         |
//     | 3     | Level3Urgent        | Urgent but stable       | 60         |
//     | 4     | Level4LessUrgent    | Less urgent             | 120        |
//     | 5     | Level5NonUrgent     | Non-urgent              | 240        |
type Level int

const (
	Level1Resuscitation Level = 1
	Level2Emergent      Level = 2
	Level3Urgent        Level = 3
	Level4LessUrgent    Level = 4
	Level5NonUrgent     Level = 5
)

// String returns a short label for the level.
func (l Level) String() string {
	switch l {
	case Level1Resuscitation:
		return "Resuscitation"
	case Level2Emergent:
		return "Emergent"
	case Level3Urgent:
		return "Urgent"
	case Level4LessUrgent:
		return "Less urgent"
	case Level5NonUrgent:
		return "Non-urgent"
	default:
		return "Unknown"
	}
}

// FromScore maps a normalized acuity score in [0, 1] to a Level using the given thresholds.
func FromScore(score float64, p Params) Level {
	if score >= p.T1 {
		return Level1Resuscitation
	}
	if score >= p.T2 {
		return Level2Emergent
	}
	if score >= p.T3 {
		return Level3Urgent
	}
	if score >= p.T4 {
		return Level4LessUrgent
	}
	return Level5NonUrgent
}

// WaitTimeMinutes returns a suggested maximum wait time in minutes for the level.
// These are guidance only; institutional protocols override.
func (l Level) WaitTimeMinutes() int {
	switch l {
	case Level1Resuscitation:
		return 0
	case Level2Emergent:
		return 15
	case Level3Urgent:
		return 60
	case Level4LessUrgent:
		return 120
	case Level5NonUrgent:
		return 240
	default:
		return 240
	}
}

// IsHighAcuity returns true for levels 1 and 2.
func (l Level) IsHighAcuity() bool {
	return l == Level1Resuscitation || l == Level2Emergent
}

// IsLowAcuity returns true for levels 4 and 5.
func (l Level) IsLowAcuity() bool {
	return l == Level4LessUrgent || l == Level5NonUrgent
}

// Int returns the level as int (1..5). For unknown, returns 0.
func (l Level) Int() int {
	if l >= 1 && l <= 5 {
		return int(l)
	}
	return 0
}

// Valid returns true if l is 1..5.
func (l Level) Valid() bool {
	return l >= 1 && l <= 5
}

// Description returns a longer description for the level.
func (l Level) Description() string {
	switch l {
	case Level1Resuscitation:
		return "Requires immediate life-saving intervention; do not delay."
	case Level2Emergent:
		return "High risk; should be seen within 15 minutes."
	case Level3Urgent:
		return "Urgent but stable; target within 60 minutes."
	case Level4LessUrgent:
		return "Less urgent; target within 120 minutes."
	case Level5NonUrgent:
		return "Non-urgent; target within 240 minutes."
	default:
		return "Unknown level."
	}
}

// AllLevels returns a slice of all five levels in order (1..5).
func AllLevels() []Level {
	return []Level{
		Level1Resuscitation, Level2Emergent, Level3Urgent, Level4LessUrgent, Level5NonUrgent,
	}
}

// LevelFromInt converts i (1..5) to Level. Returns 0 for invalid i.
func LevelFromInt(i int) Level {
	if i >= 1 && i <= 5 {
		return Level(i)
	}
	return 0
}

// LessAcuteThan returns true if l has lower acuity than other (e.g. 4 < 2).
func (l Level) LessAcuteThan(other Level) bool {
	return l.Int() > other.Int() && other.Valid() && l.Valid()
}

// MoreAcuteThan returns true if l has higher acuity than other (e.g. 1 > 4).
func (l Level) MoreAcuteThan(other Level) bool {
	return l.Int() < other.Int() && l.Valid() && other.Valid()
}

// Distance returns |l - other| as integers (1..5). For use in weighted metrics.
func (l Level) Distance(other Level) int {
	a, b := l.Int(), other.Int()
	if a == 0 || b == 0 {
		return 4
	}
	d := a - b
	if d < 0 {
		return -d
	}
	return d
}

// LevelStrings returns a map from level (1..5) to label.
func LevelStrings() map[Level]string {
	return map[Level]string{
		Level1Resuscitation: "Resuscitation",
		Level2Emergent:      "Emergent",
		Level3Urgent:        "Urgent",
		Level4LessUrgent:    "Less urgent",
		Level5NonUrgent:     "Non-urgent",
	}
}

// LevelsFromInts converts a slice of int (1..5) to []Level. Invalid values become 0.
func LevelsFromInts(x []int) []Level {
	out := make([]Level, len(x))
	for i, v := range x {
		out[i] = LevelFromInt(v)
	}
	return out
}

// IntsFromLevels converts []Level to []int (1..5).
func IntsFromLevels(lvls []Level) []int {
	out := make([]int, len(lvls))
	for i, l := range lvls {
		out[i] = l.Int()
	}
	return out
}

// Compare returns -1 if l is more acute than other, 0 if equal, 1 if less acute.
func (l Level) Compare(other Level) int {
	a, b := l.Int(), other.Int()
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ShortCode returns a one-letter or short code for the level (e.g. "R", "E", "U", "L", "N").
func (l Level) ShortCode() string {
	switch l {
	case Level1Resuscitation:
		return "R"
	case Level2Emergent:
		return "E"
	case Level3Urgent:
		return "U"
	case Level4LessUrgent:
		return "L"
	case Level5NonUrgent:
		return "N"
	default:
		return "?"
	}
}

// ParseLevel parses a string label (case-insensitive) to Level. Returns 0 if unknown.
func ParseLevel(s string) Level {
	switch s {
	case "1", "resuscitation", "Resuscitation", "R":
		return Level1Resuscitation
	case "2", "emergent", "Emergent", "E":
		return Level2Emergent
	case "3", "urgent", "Urgent", "U":
		return Level3Urgent
	case "4", "less urgent", "Less urgent", "L":
		return Level4LessUrgent
	case "5", "non-urgent", "Non-urgent", "N":
		return Level5NonUrgent
	default:
		return 0
	}
}

// RecommendedActions returns a short list of recommended actions for the level.
// For display or decision support only; not a substitute for protocol.
func (l Level) RecommendedActions() []string {
	switch l {
	case Level1Resuscitation:
		return []string{"Immediate assessment", "Life-saving interventions as indicated", "Continuous monitoring"}
	case Level2Emergent:
		return []string{"Rapid assessment", "Stabilisation", "Re-evaluate within 15 min"}
	case Level3Urgent:
		return []string{"Assessment within 60 min", "Routine monitoring", "Re-evaluate as needed"}
	case Level4LessUrgent:
		return []string{"Assessment within 120 min", "Routine care", "Re-evaluate if condition changes"}
	case Level5NonUrgent:
		return []string{"Assessment within 240 min", "Routine care", "May use fast-track if available"}
	default:
		return nil
	}
}

// LevelCounts returns counts per level (index 1..5) for the given slice. Index 0 is unused.
func LevelCounts(lvls []Level) [6]int {
	var c [6]int
	for _, l := range lvls {
		if l >= 1 && l <= 5 {
			c[l.Int()]++
		}
	}
	return c
}

// LevelProportions returns proportions (0..1) per level for the given slice. Index 0 is 0.
func LevelProportions(lvls []Level) [6]float64 {
	var p [6]float64
	c := LevelCounts(lvls)
	var total int
	for i := 1; i <= 5; i++ {
		total += c[i]
	}
	if total > 0 {
		for i := 1; i <= 5; i++ {
			p[i] = float64(c[i]) / float64(total)
		}
	}
	return p
}

// CountHighAcuity returns the number of levels that are 1 or 2 in lvls.
func CountHighAcuity(lvls []Level) int {
	var n int
	for _, l := range lvls {
		if l.IsHighAcuity() {
			n++
		}
	}
	return n
}

// CountLowAcuity returns the number of levels that are 4 or 5 in lvls.
func CountLowAcuity(lvls []Level) int {
	var n int
	for _, l := range lvls {
		if l.IsLowAcuity() {
			n++
		}
	}
	return n
}
