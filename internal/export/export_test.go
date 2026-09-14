package export

import (
	"strings"
	"testing"
)

func TestGenerateTextReport(t *testing.T) {
	data := &ReportData{
		Filename:   "test_contract.docx",
		FileSize:   10240,
		AnalyzedAt: 1700000000,
		Score:      65,
		Summary:    "Sample contract with medium risk.",
		Risks: []map[string]any{
			{"type": "high", "risk_level": "high", "clause": "Non-compete clause", "explanation": "Overly broad restriction"},
		},
		Suggestions: []string{"Review non-compete terms", "Add geographic limitations"},
	}

	report := GenerateTextReport(data)
	if len(report) == 0 {
		t.Error("Expected non-empty report")
	}
	if !strings.Contains(report, "AI CONTRACT REVIEW") {
		t.Error("Report missing header")
	}
	if !strings.Contains(report, "test_contract.docx") {
		t.Error("Report missing filename")
	}
	if !strings.Contains(report, "65/100") {
		t.Error("Report missing score")
	}
}

func TestGenerateHTMLReport(t *testing.T) {
	data := &ReportData{
		Filename: "test.pdf",
		Score:    90,
		Summary:  "Low risk contract",
	}

	html := GenerateHTMLReport(data)
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("Missing HTML doctype")
	}
	if !strings.Contains(html, "AI Contract Review Report") {
		t.Error("Missing title")
	}
	if !strings.Contains(html, "90") {
		t.Error("Missing score")
	}
}

func TestGenerateJSONReport(t *testing.T) {
	data := &ReportData{
		Filename: "test.txt",
		Score:    75,
	}

	jsonStr, err := GenerateJSONReport(data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !strings.Contains(jsonStr, "test.txt") {
		t.Error("Missing filename in JSON")
	}
	if !strings.Contains(jsonStr, `"score"`) {
		t.Error("Missing score in JSON")
	}
}

func TestGenerateCSVReport(t *testing.T) {
	data := &ReportData{
		Risks: []map[string]any{
			{"type": "high", "risk_level": "high", "clause": "Clause A", "explanation": "Explanation A"},
		},
		Clauses: []map[string]any{
			{"type": "termination", "risk_level": "medium", "content": "Either party may terminate"},
		},
	}

	csv := GenerateCSVReport(data)
	if !strings.Contains(csv, "TYPE") {
		t.Error("Missing CSV header TYPE")
	}
	if !strings.Contains(csv, "Clause A") {
		t.Error("Missing risk data")
	}
}

func TestGenerateExport(t *testing.T) {
	data := &ReportData{Filename: "test.docx", Score: 80}

	// Test TXT export
	filename, content := GenerateExport(data, ExportOptions{Format: "text"})
	if filename == "" {
		t.Error("Empty filename")
	}
	if !strings.Contains(string(content), "AI CONTRACT REVIEW") {
		t.Error("Missing report content")
	}

	// Test HTML export
	_, htmlContent := GenerateExport(data, ExportOptions{Format: "html"})
	if !strings.Contains(string(htmlContent), "<!DOCTYPE html>") {
		t.Error("Missing HTML doctype")
	}

	// Test DOCX export
	_, docxContent := GenerateExport(data, ExportOptions{Format: "docx"})
	if len(docxContent) == 0 {
		t.Error("Empty DOCX content")
	}
}

func TestNewReportData(t *testing.T) {
	result := map[string]any{
		"score":       float64(72),
		"summary":     "Test summary",
		"risks":       []any{map[string]any{"type": "high", "risk_level": "high", "clause": "C1", "explanation": "E1"}},
		"key_terms":   map[string]any{"salary": "75000", "term": "2 years"},
		"clauses":     []any{map[string]any{"type": "payment", "content": "Pay monthly"}},
		"compliance":  []any{map[string]any{"item": "Termination", "status": "present"}},
		"suggestions": []any{"Add clause", "Review terms"},
		"scenario":    map[string]any{"type_name": "Employment Contract"},
	}

	data := NewReportData("contract.docx", 5000, result)
	if data.Score != 72 {
		t.Errorf("Expected score 72, got %d", data.Score)
	}
	if len(data.Risks) != 1 {
		t.Errorf("Expected 1 risk, got %d", len(data.Risks))
	}
	if len(data.Suggestions) != 2 {
		t.Errorf("Expected 2 suggestions, got %d", len(data.Suggestions))
	}
	if data.Scenario == nil {
		t.Error("Expected scenario")
	}
}
