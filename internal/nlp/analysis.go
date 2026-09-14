package nlp

import (
	"fmt"
	"strings"
)

// RiskLevel represents the severity of a risk clause
type RiskLevel int

const (
	RiskLow RiskLevel = iota
	RiskMedium
	RiskHigh
	RiskCritical
)

func (r RiskLevel) String() string {
	switch r {
	case RiskLow:
		return "low"
	case RiskMedium:
		return "medium"
	case RiskHigh:
		return "high"
	case RiskCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// RiskClause represents a detected risk in a contract
type RiskClause struct {
	Term       string     `json:"term"`
	Risk       RiskLevel  `json:"risk"`
	Clause     string     `json:"clause"`
	Suggestion string     `json:"suggestion,omitempty"`
}

// ContractAnalysis contains the full analysis result
type ContractAnalysis struct {
	Summary          string       `json:"summary"`
	DangerousClauses []RiskClause `json:"dangerous_clauses"`
	Score            int          `json:"score"` // 0-100, higher is safer
}

// patternInfo holds pattern matching info
type patternInfo struct {
	clue       string
	risk       string
	suggestion string
}

// AnalyzeKeywords performs keyword-based risk detection
func AnalyzeKeywords(text string) []RiskClause {
	var clauses []RiskClause

	criticalPatterns := map[string]patternInfo{
		"indemnification":             {"Indemnification Clause", "high", "Consider limiting indemnification scope"},
		"automatic renewal":           {"Auto-Renewal Clause", "high", "Add opt-out provision"},
		"sole discretion":             {"Sole Discretion", "critical", "Add mutual agreement requirement"},
		"non-negotiable":              {"Non-Negotiable Terms", "critical", "Negotiate flexibility"},
		"force majeure":               {"Force Majeure", "medium", "Review liability limitations"},
		"unlimited liability":         {"Unlimited Liability", "critical", "Cap liability at contract value"},
		"permanent confidentiality":   {"Perpetual Confidentiality", "high", "Add time limit to NDA"},
		"exclusive license":           {"Exclusive License", "high", "Consider non-exclusive option"},
		"liquidated damages":          {"Liquidated Damages", "high", "Review reasonableness"},
		"jurisdiction":                {"Jurisdiction Clause", "medium", "Verify favorable jurisdiction"},
	}

	textLower := strings.ToLower(text)
	for pattern, info := range criticalPatterns {
		if strings.Contains(textLower, pattern) {
			clauses = append(clauses, RiskClause{
				Term:       info.clue,
				Risk:       parseRiskLevel(info.risk),
				Clause:     fmt.Sprintf("Contains: %s", pattern),
				Suggestion: info.suggestion,
			})
		}
	}

	// Scan lines for actual clause text
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 20 {
			continue
		}
		lowerLine := strings.ToLower(line)
		for pattern, info := range criticalPatterns {
			if strings.Contains(lowerLine, pattern) {
				clauses = append(clauses, RiskClause{
					Term:       info.clue,
					Risk:       parseRiskLevel(info.risk),
					Clause:     line,
					Suggestion: info.suggestion,
				})
				break
			}
		}
	}

	return clauses
}

func parseRiskLevel(level string) RiskLevel {
	switch level {
	case "critical":
		return RiskCritical
	case "high":
		return RiskHigh
	case "medium":
		return RiskMedium
	default:
		return RiskLow
	}
}

// ExtractKeyTerms extracts key terms from contract text
func ExtractKeyTerms(text string) map[string]string {
	terms := make(map[string]string)
	patterns := map[string]string{
		"contract period|term of agreement|duration": "Contract Term",
		"payment terms|net [0-9]+|payment due":      "Payment Terms",
		"termination|cancel|end date":                "Termination",
		"confidential|nda|non-disclosure":            "Confidentiality",
		"intellectual property|ip rights|copyright":  "IP Rights",
		"liability|damages|indemnif":                 "Liability",
		"jurisdiction|governing law|venue":           "Governing Law",
	}

	textLower := strings.ToLower(text)
	for pattern, label := range patterns {
		if match, _ := regexpMatch(pattern, textLower); match {
			terms[label] = extractContext(textLower, pattern, 100)
		}
	}

	return terms
}

func regexpMatch(pattern, text string) (bool, string) {
	// For MVP, split by | and check each alternative
	parts := strings.Split(pattern, "|")
	for _, part := range parts {
		if strings.Contains(text, part) {
			return true, part
		}
	}
	return false, ""
}

func extractContext(text, pattern string, contextLen int) string {
	idx := strings.Index(text, pattern)
	if idx == -1 {
		return ""
	}
	start := 0
	if idx > contextLen {
		start = idx - contextLen
	}
	end := idx + len(pattern) + contextLen
	if end > len(text) {
		end = len(text)
	}
	return strings.TrimSpace(text[start:end])
}
