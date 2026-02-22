package export

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/olaflaitinen/triagegeist/score"
)

func TestFromVitalsScoreLevel(t *testing.T) {
	v := score.Vitals{HR: 120, RR: 24, SBP: 90, SpO2: 92}
	r := FromVitalsScoreLevel(v, 3, 0.75, 2, "Emergent")
	if r.HR != 120 || r.Acuity != 0.75 || r.Level != 2 || r.LevelLabel != "Emergent" {
		t.Errorf("FromVitalsScoreLevel: %+v", r)
	}
}

func TestToCSVRow(t *testing.T) {
	r := Result{HR: 80, Acuity: 0.5, Level: 2, LevelLabel: "Emergent"}
	row := r.ToCSVRow()
	if len(row) != len(CSVHeader()) {
		t.Errorf("ToCSVRow length = %d, header = %d", len(row), len(CSVHeader()))
	}
}

func TestToJSON(t *testing.T) {
	r := Result{HR: 80, Acuity: 0.5, Level: 2}
	var buf bytes.Buffer
	if err := r.ToJSON(&buf); err != nil {
		t.Fatal(err)
	}
	var dec Result
	if err := json.NewDecoder(&buf).Decode(&dec); err != nil {
		t.Fatal(err)
	}
	if dec.HR != r.HR || dec.Level != r.Level {
		t.Errorf("roundtrip: got %+v", dec)
	}
}

func TestLevelReport(t *testing.T) {
	results := []Result{
		{Level: 1, Acuity: 0.9},
		{Level: 1, Acuity: 0.85},
		{Level: 2, Acuity: 0.6},
	}
	rows := LevelReport(results)
	if len(rows) != 5 {
		t.Errorf("LevelReport should return 5 rows, got %d", len(rows))
	}
	if rows[0].Count != 2 || rows[1].Count != 1 {
		t.Errorf("LevelReport counts: %+v", rows)
	}
}

func TestComputeSummary(t *testing.T) {
	results := []Result{
		{Acuity: 0.2, Level: 5},
		{Acuity: 0.8, Level: 1},
	}
	s := ComputeSummary(results)
	if s.N != 2 || s.MeanAcuity != 0.5 {
		t.Errorf("ComputeSummary: N=%d Mean=%v", s.N, s.MeanAcuity)
	}
}

func TestResultToVitals(t *testing.T) {
	r := Result{HR: 100, RR: 20, SBP: 110}
	v := ResultToVitals(r)
	if v.HR != 100 || v.RR != 20 || v.SBP != 110 {
		t.Errorf("ResultToVitals: %+v", v)
	}
}

func TestWriteCSV(t *testing.T) {
	results := []Result{
		FromVitalsScoreLevel(score.Vitals{HR: 80}, 0, 0.3, 5, "Non-urgent"),
	}
	var buf bytes.Buffer
	if err := WriteCSV(&buf, results); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Error("WriteCSV produced no output")
	}
}
