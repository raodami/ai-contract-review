package export

import (
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
	if !contains(report, "AI CONTRACT REVIEW") {
		t.Error("Report missing header")
	}
	if !contains(report, "test_contract.docx") {
		t.Error("Report missing filename")
	}
	if !contains(report, "65/100") {
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
	if !contains(html, "<!DOCTYPE html>") {
		t.Error("Missing HTML doctype")
	}
	if !contains(html, "AI Contract Review Report") {
		t.Error("Missing title")
	}
	if !contains(html, "90") {
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
	if !contains(jsonStr, "test.txt") {
		t.Error("Missing filename in JSON")
	}
	if !contains(jsonStr, `"score"`) {
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
	if !contains(csv, "TYPE") {
		t.Error("Missing CSV header TYPE")
	}
	if !contains(csv, "Clause A") {
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
	if !contains(content, "AI CONTRACT REVIEW") {
		t.Error("Missing report content")
	}

	// Test HTML export
	_, htmlContent := GenerateExport(data, ExportOptions{Format: "html"})
	if !contains(htmlContent, "<!DOCTYPE html>") {
		t.Error("Missing HTML doctype", htmlContent[:200])
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
	if data.Score != 72 {
		t.Error("Score mismatch")
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

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (len(s) >= len(substr)) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
