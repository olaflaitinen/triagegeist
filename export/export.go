// Copyright (c) triagegeist authors: Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov.
// Licensed under the EUPL.
//
// Package export provides structures and helpers for serialising triage
// results to JSON, CSV, or other formats for logging, auditing, and
// research. All types use exported fields for encoding/json compatibility.
//
// # Output formats
//
//	| Format | Use case                    |
//	|--------|-----------------------------|
//	| JSON   | APIs, logs, single record   |
//	| CSV    | Batch export, spreadsheets  |
//	| Row    | Tabular in-memory          |
package export

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/olaflaitinen/triagegeist/score"
)

// Result holds one triage evaluation for export: vitals, resource count,
// acuity score, level, and optional metadata.
type Result struct {
	// Vitals at time of evaluation
	HR   int     `json:"hr"`
	RR   int     `json:"rr"`
	SBP  int     `json:"sbp"`
	DBP  int     `json:"dbp"`
	Temp float64 `json:"temp"`
	SpO2 int     `json:"spo2"`
	GCS  int     `json:"gcs"`
	// ResourceCount is the expected number of resources
	ResourceCount int     `json:"resource_count"`
	Acuity        float64 `json:"acuity"`
	Level         int     `json:"level"`
	LevelLabel    string  `json:"level_label"`
	// Timestamp is optional; zero value means not set
	Timestamp time.Time `json:"timestamp,omitempty"`
	// ID is optional (e.g. encounter or record ID)
	ID string `json:"id,omitempty"`
}

// FromVitalsScoreLevel builds a Result from score.Vitals, acuity, level (1..5), and label.
func FromVitalsScoreLevel(v score.Vitals, resourceCount int, acuity float64, level int, levelLabel string) Result {
	return Result{
		HR:             v.HR,
		RR:             v.RR,
		SBP:            v.SBP,
		DBP:            v.DBP,
		Temp:           v.Temp,
		SpO2:           v.SpO2,
		GCS:            v.GCS,
		ResourceCount:  resourceCount,
		Acuity:         acuity,
		Level:          level,
		LevelLabel:     levelLabel,
	}
}

// ToJSON writes r as a single JSON object to w (no newline array).
func (r Result) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(r)
}

// CSVHeader returns the header row for CSV export of Result.
func CSVHeader() []string {
	return []string{
		"hr", "rr", "sbp", "dbp", "temp", "spo2", "gcs",
		"resource_count", "acuity", "level", "level_label",
		"timestamp", "id",
	}
}

// ToCSVRow returns a slice of strings for one Result (same order as CSVHeader).
func (r Result) ToCSVRow() []string {
	ts := ""
	if !r.Timestamp.IsZero() {
		ts = r.Timestamp.Format(time.RFC3339)
	}
	return []string{
		strconv.Itoa(r.HR),
		strconv.Itoa(r.RR),
		strconv.Itoa(r.SBP),
		strconv.Itoa(r.DBP),
		strconv.FormatFloat(r.Temp, 'f', -1, 64),
		strconv.Itoa(r.SpO2),
		strconv.Itoa(r.GCS),
		strconv.Itoa(r.ResourceCount),
		strconv.FormatFloat(r.Acuity, 'f', -1, 64),
		strconv.Itoa(r.Level),
		r.LevelLabel,
		ts,
		r.ID,
	}
}

// WriteCSV writes the header and all results to w using encoding/csv.
func WriteCSV(w io.Writer, results []Result) error {
	cw := csv.NewWriter(w)
	if err := cw.Write(CSVHeader()); err != nil {
		return err
	}
	for _, r := range results {
		if err := cw.Write(r.ToCSVRow()); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

// Batch holds multiple Result and optional metadata for batch export.
type Batch struct {
	Results   []Result  `json:"results"`
	Generated time.Time `json:"generated"`
	Source    string    `json:"source,omitempty"`
}

// ToJSON writes the batch as one JSON object to w.
func (b Batch) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(b)
}

// ReportRow is a flattened row for tabular reporting (e.g. summary stats).
type ReportRow struct {
	Level      int
	LevelLabel string
	Count      int
	Pct        float64
	MeanAcuity float64
	MinAcuity  float64
	MaxAcuity  float64
}

// LevelReport builds one ReportRow per level 1..5 from results.
func LevelReport(results []Result) []ReportRow {
	var counts [6]int
	var sumAcuity [6]float64
	var minAcuity [6]float64
	var maxAcuity [6]float64
	for i := 1; i <= 5; i++ {
		minAcuity[i] = 1
		maxAcuity[i] = 0
	}
	for _, r := range results {
		if r.Level < 1 || r.Level > 5 {
			continue
		}
		counts[r.Level]++
		sumAcuity[r.Level] += r.Acuity
		if r.Acuity < minAcuity[r.Level] {
			minAcuity[r.Level] = r.Acuity
		}
		if r.Acuity > maxAcuity[r.Level] {
			maxAcuity[r.Level] = r.Acuity
		}
	}
	var total int
	for i := 1; i <= 5; i++ {
		total += counts[i]
	}
	labels := []string{"", "Resuscitation", "Emergent", "Urgent", "Less urgent", "Non-urgent"}
	var out []ReportRow
	for i := 1; i <= 5; i++ {
		pct := 0.0
		if total > 0 {
			pct = float64(counts[i]) / float64(total) * 100
		}
		mean := 0.0
		if counts[i] > 0 {
			mean = sumAcuity[i] / float64(counts[i])
		}
		min, max := minAcuity[i], maxAcuity[i]
		if counts[i] == 0 {
			min, max = 0, 0
		}
		out = append(out, ReportRow{
			Level:      i,
			LevelLabel: labels[i],
			Count:      counts[i],
			Pct:        pct,
			MeanAcuity: mean,
			MinAcuity:  min,
			MaxAcuity:  max,
		})
	}
	return out
}

// ReportRowHeader returns CSV header for ReportRow.
func ReportRowHeader() []string {
	return []string{"level", "level_label", "count", "pct", "mean_acuity", "min_acuity", "max_acuity"}
}

// ReportRowToCSV returns a string slice for one ReportRow.
func (r ReportRow) ReportRowToCSV() []string {
	return []string{
		strconv.Itoa(r.Level),
		r.LevelLabel,
		strconv.Itoa(r.Count),
		strconv.FormatFloat(r.Pct, 'f', 2, 64),
		strconv.FormatFloat(r.MeanAcuity, 'f', 4, 64),
		strconv.FormatFloat(r.MinAcuity, 'f', 4, 64),
		strconv.FormatFloat(r.MaxAcuity, 'f', 4, 64),
	}
}

// WriteLevelReportCSV writes LevelReport(results) as CSV to w.
func WriteLevelReportCSV(w io.Writer, results []Result) error {
	rows := LevelReport(results)
	cw := csv.NewWriter(w)
	if err := cw.Write(ReportRowHeader()); err != nil {
		return err
	}
	for _, row := range rows {
		if err := cw.Write(row.ReportRowToCSV()); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

// ReadResultJSON decodes a single Result from r.
func ReadResultJSON(r io.Reader) (Result, error) {
	var res Result
	dec := json.NewDecoder(r)
	err := dec.Decode(&res)
	return res, err
}

// ReadBatchJSON decodes a Batch from r.
func ReadBatchJSON(r io.Reader) (Batch, error) {
	var b Batch
	dec := json.NewDecoder(r)
	err := dec.Decode(&b)
	return b, err
}

// ResultToVitals converts a Result back to score.Vitals (for re-scoring or validation).
func ResultToVitals(r Result) score.Vitals {
	return score.Vitals{
		HR:   r.HR,
		RR:   r.RR,
		SBP:  r.SBP,
		DBP:  r.DBP,
		Temp: r.Temp,
		SpO2: r.SpO2,
		GCS:  r.GCS,
	}
}

// Summary holds aggregate stats over a slice of Result.
type Summary struct {
	N          int     `json:"n"`
	MeanAcuity float64 `json:"mean_acuity"`
	MinAcuity  float64 `json:"min_acuity"`
	MaxAcuity  float64 `json:"max_acuity"`
	LevelDist  [6]int  `json:"level_dist"` // index 0 unused; 1..5
}

// ComputeSummary returns Summary from results.
func ComputeSummary(results []Result) Summary {
	var s Summary
	if len(results) == 0 {
		return s
	}
	s.N = len(results)
	s.MinAcuity = 1
	for _, r := range results {
		s.MeanAcuity += r.Acuity
		if r.Acuity < s.MinAcuity {
			s.MinAcuity = r.Acuity
		}
		if r.Acuity > s.MaxAcuity {
			s.MaxAcuity = r.Acuity
		}
		if r.Level >= 1 && r.Level <= 5 {
			s.LevelDist[r.Level]++
		}
	}
	s.MeanAcuity /= float64(s.N)
	return s
}
