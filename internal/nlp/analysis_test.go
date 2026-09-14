package nlp

import (
	"testing"
)

func TestAnalyzeKeywords_CriticalTerms(t *testing.T) {
	text := "This agreement includes automatic renewal and sole discretion clauses. The indemnification clause is unlimited."
	clauses := AnalyzeKeywords(text)
	
	found := make(map[string]bool)
	for _, c := range clauses {
		found[c.Term] = true
	}
	
	expected := []string{"Indemnification Clause", "Sole Discretion", "Auto-Renewal Clause"}
	for _, exp := range expected {
		if !found[exp] {
			t.Errorf("Expected clause %s not found", exp)
		}
	}
}

func TestAnalyzeKeywords_Empty(t *testing.T) {
	clauses := AnalyzeKeywords("This is a simple agreement with no special clauses.")
	if len(clauses) > 0 {
		t.Error("Expected no critical clauses")
	}
}

func TestExtractKeyTerms(t *testing.T) {
	text := "The contract period is 12 months. Payment terms require net 30 payment. Termination requires 30 days notice. Confidential information must be protected indefinitely. Intellectual property rights belong to the creator. Unlimited liability is accepted. Governing law jurisdiction is New York."
	terms := ExtractKeyTerms(text)
	
	expected := []string{"Contract Term", "Payment Terms", "Termination", "Confidentiality"}
	for _, exp := range expected {
		if _, ok := terms[exp]; !ok {
			t.Errorf("Expected key term %s not found", exp)
		}
	}
}

func TestRiskLevelString(t *testing.T) {
	if RiskHigh.String() != "high" {
		t.Errorf("Expected 'high', got %s", RiskHigh.String())
	}
}

func TestAnalyzeKeywords_MixedContent(t *testing.T) {
	text := `
INDEMNIFICATION. Party A agrees to indemnify Party B against all claims.
TERM. This agreement begins on January 1, 2024 and continues for 12 months.
TERMINATION. Either party may terminate with 30 days notice.
`
	clauses := AnalyzeKeywords(text)
	found := false
	for _, c := range clauses {
		if c.Term == "Indemnification Clause" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find Indemnification Clause")
	}
}

func TestExtractKeyTerms_Empty(t *testing.T) {
	terms := ExtractKeyTerms("")
	if len(terms) != 0 {
		t.Error("Expected empty terms for empty input")
	}
}
