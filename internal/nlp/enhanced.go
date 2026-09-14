package nlp

import (
	"fmt"
	"strings"
)

// Enhanced contract analysis with clause classification and compliance checks
func AnalyzeContractEnhanced(text string) *EnhancedAnalysis {
	analysis := &EnhancedAnalysis{}
	
	// Risk analysis
	clauses := AnalyzeKeywords(text)
	analysis.Risks = convertToRiskItems(clauses)
	
	// Extract key terms
	terms := ExtractKeyTerms(text)
	analysis.KeyTerms = terms
	
	// Clause classification
	analysis.Clauses = classifyClauses(text)
	
	// Liability analysis
	liability := analyzeLiability(text)
	analysis.Liability = liability
	
	// Compliance checks
	compliance := checkCompliance(text)
	analysis.Compliance = compliance
	
	// Calculate score
	analysis.Score = calculateScore(clauses, compliance)
	
	// Generate suggestions
	analysis.Suggestions = generateSuggestions(clauses, compliance)
	
	// Generate summary
	analysis.Summary = generateSummary(analysis)
	
	return analysis
}

// EnhancedAnalysis contains full analysis result
type EnhancedAnalysis struct {
	Summary      string           `json:"summary"`
	Risks        []RiskItem       `json:"risks"`
	KeyTerms     map[string]string `json:"key_terms"`
	Clauses      []ClauseInfo     `json:"clauses"`
	Liability    LiabilityAnalysis `json:"liability"`
	Compliance   []ComplianceCheck `json:"compliance"`
	Score        int              `json:"score"`
	Suggestions  []string         `json:"suggestions"`
}

// RiskItem is a risk found in the contract
type RiskItem struct {
	Clause      string `json:"clause"`
	RiskLevel   string `json:"risk_level"`
	Explanation string `json:"explanation"`
}

// ClauseInfo represents a classified contract clause
type ClauseInfo struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	RiskLevel string `json:"risk_level,omitempty"`
}

// LiabilityAnalysis analyzes liability provisions
type LiabilityAnalysis struct {
	HasLimitation bool   `json:"has_limitation"`
	LimitType     string `json:"limit_type"`
	LimitAmount   string `json:"limit_amount"`
	Exclusions    []string `json:"exclusions"`
	Risk          string `json:"risk"`
}

// ComplianceCheck represents a compliance requirement
type ComplianceCheck struct {
	Item      string `json:"item"`
	Status    string `json:"status"` // "present", "missing", "needs_review"
	Issue     string `json:"issue,omitempty"`
	Severity  string `json:"severity,omitempty"` // "high", "medium", "low"
}

// Clause types for classification
var clausePatterns = map[string][]string{
	"termination": {"termination", "cancel", "end of term", "renewal"},
	"payment": {"payment", "fee", "compensation", "invoice", "net 30", "net 60"},
	"confidentiality": {"confidential", "nda", "non-disclosure", "proprietary"},
	"intellectual_property": {"intellectual property", "ip rights", "copyright", "patent", "trademark"},
	"liability": {"liability", "damages", "indemnif", "warranty", "disclaimer"},
	"governing_law": {"governing law", "jurisdiction", "venue", "arbitration"},
	"non_compete": {"non-compete", "non-solicitation", "restrictive covenant"},
	"force_majeure": {"force majeure", "act of god", "unforeseeable"},
	"assignment": {"assignment", "transfer", "successors"},
	"amendment": {"amendment", "modification", "waiver"},
}

func classifyClauses(text string) []ClauseInfo {
	var clauses []ClauseInfo
	lines := strings.Split(text, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 10 || len(line) > 500 {
			continue
		}
		
		lowerLine := strings.ToLower(line)
		
		for clauseType, patterns := range clausePatterns {
			for _, pattern := range patterns {
				if strings.Contains(lowerLine, pattern) {
					clauses = append(clauses, ClauseInfo{
						Type:    clauseType,
						Title:   formatClauseTitle(clauseType),
						Content: line,
					})
					break
				}
			}
		}
	}
	
	return clauses
}

func formatClauseTitle(t string) string {
	switch t {
	case "termination":
		return "Termination Clause"
	case "payment":
		return "Payment Terms"
	case "confidentiality":
		return "Confidentiality Clause"
	case "intellectual_property":
		return "Intellectual Property"
	case "liability":
		return "Liability Clause"
	case "governing_law":
		return "Governing Law"
	case "non_compete":
		return "Non-Compete Clause"
	case "force_majeure":
		return "Force Majeure"
	case "assignment":
		return "Assignment Clause"
	case "amendment":
		return "Amendment Clause"
	default:
		return strings.Title(t) + " Clause"
	}
}

func analyzeLiability(text string) LiabilityAnalysis {
	la := LiabilityAnalysis{}
	lower := strings.ToLower(text)
	
	// Check for limitation of liability
	if strings.Contains(lower, "limitation of liability") || strings.Contains(lower, "liability shall not exceed") {
		la.HasLimitation = true
		la.LimitType = "capped"
		la.Risk = "low"
	} else if strings.Contains(lower, "unlimited liability") || strings.Contains(lower, "full liability") {
		la.HasLimitation = false
		la.LimitType = "uncapped"
		la.Risk = "high"
	} else {
		la.HasLimitation = false
		la.LimitType = "none"
		la.Risk = "medium"
	}
	
	// Extract exclusions
	if strings.Contains(lower, "except") || strings.Contains(lower, "excluding") {
		la.Exclusions = extractExclusions(text)
	}
	
	// Check for liquidated damages
	if strings.Contains(lower, "liquidated damages") {
		la.Risk = "high"
	}
	
	return la
}

func extractExclusions(text string) []string {
	var exclusions []string
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		lower := strings.ToLower(line)
		if (strings.Contains(lower, "except") || strings.Contains(lower, "excluding")) && len(line) > 20 {
			exclusions = append(exclusions, line)
		}
	}
	return exclusions
}

func checkCompliance(text string) []ComplianceCheck {
	var checks []ComplianceCheck
	lower := strings.ToLower(text)
	
	// Essential clauses that should be present
	essentialClauses := map[string]struct {
		keywords []string
		issue    string
		severity string
	}{
		"parties_identification": {
			[]string{"party a", "party b", "between", "hereby agree"},
			"No clear party identification",
			"high",
		},
		"effective_date": {
			[]string{"effective date", "commencement date", "date of agreement"},
			"Missing effective date",
			"high",
		},
		"term_duration": {
			[]string{"term of", "duration", "period of", "shall remain in effect"},
			"Missing contract term/duration",
			"medium",
		},
		"governing_law": {
			[]string{"governing law", "jurisdiction", "governed by"},
			"Missing governing law clause",
			"medium",
		},
		"dispute_resolution": {
			[]string{"dispute", "arbitration", "litigation", "court"},
			"Missing dispute resolution clause",
			"medium",
		},
		"signatures": {
			[]string{"signature", "signed", "executed"},
			"Missing signature block",
			"low",
		},
	}
	
	for name, check := range essentialClauses {
		found := false
		for _, kw := range check.keywords {
			if strings.Contains(lower, kw) {
				found = true
				break
			}
		}
		
		status := "missing"
		if found {
			status = "present"
		}
		
		checks = append(checks, ComplianceCheck{
			Item:     name,
			Status:   status,
			Issue:    check.issue,
			Severity: check.severity,
		})
	}
	
	// Check for risk indicators
	riskChecks := map[string]struct {
		keywords []string
		issue    string
		severity string
	}{
		"auto_renewal": {
			[]string{"automatic renewal", "renews automatically", "deemed renewed"},
			"Auto-renewal clause may trap parties",
			"high",
		},
		"sole_discretion": {
			[]string{"sole discretion", "sole and exclusive"},
			"One-sided discretion may be unfair",
			"high",
		},
		"perpetual_obligation": {
			[]string{"perpetual", "in perpetuity", "survive termination"},
			"Perpetual obligations may be burdensome",
			"medium",
		},
	}
	
	for name, check := range riskChecks {
		for _, kw := range check.keywords {
			if strings.Contains(lower, kw) {
				checks = append(checks, ComplianceCheck{
					Item:     name,
					Status:   "present",
					Issue:    check.issue,
					Severity: check.severity,
				})
				break
			}
		}
	}
	
	return checks
}

func calculateScore(clauses []RiskClause, compliance []ComplianceCheck) int {
	score := 100
	
	// Deduct for risks
	for _, clause := range clauses {
		switch clause.Risk {
		case RiskCritical:
			score -= 20
		case RiskHigh:
			score -= 12
		case RiskMedium:
			score -= 5
		}
	}
	
	// Bonus for compliance
	for _, check := range compliance {
		if check.Status == "present" && check.Severity == "high" {
			score += 3
		}
	}
	
	// Penalty for missing essential items
	for _, check := range compliance {
		if check.Status == "missing" {
			switch check.Severity {
			case "high":
				score -= 10
			case "medium":
				score -= 5
			case "low":
				score -= 2
			}
		}
	}
	
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	
	return score
}

func generateSuggestions(clauses []RiskClause, compliance []ComplianceCheck) []string {
	var suggestions []string
	
	// Suggestions based on risks
	for _, clause := range clauses {
		switch clause.Risk {
		case RiskCritical:
			suggestions = append(suggestions, "CRITICAL: "+clause.Suggestion+" - Recommend immediate review")
		case RiskHigh:
			suggestions = append(suggestions, "HIGH RISK: "+clause.Suggestion)
		case RiskMedium:
			suggestions = append(suggestions, "Consider: "+clause.Suggestion)
		}
	}
	
	// Suggestions based on compliance
	for _, check := range compliance {
		if check.Status == "missing" {
			suggestions = append(suggestions, "Add clause: "+check.Issue)
		}
	}
	
	return suggestions
}

func generateSummary(a *EnhancedAnalysis) string {
	var sb strings.Builder
	sb.WriteString("Contract Analysis Summary\n\n")
	
	sb.WriteString(fmt.Sprintf("Risk Score: %d/100\n", a.Score))
	
	if len(a.Risks) > 0 {
		sb.WriteString(fmt.Sprintf("Detected %d risk(s):\n", len(a.Risks)))
		for i, risk := range a.Risks {
			if i < 5 {
				sb.WriteString(fmt.Sprintf("- [%s] %s\n", strings.ToUpper(risk.RiskLevel), risk.Clause))
			}
		}
		sb.WriteString("\n")
	}
	
	if len(a.Compliance) > 0 {
		missing := 0
		for _, c := range a.Compliance {
			if c.Status == "missing" {
				missing++
			}
		}
		if missing > 0 {
			sb.WriteString(fmt.Sprintf("Missing %d essential clause(s)\n", missing))
		}
	}
	
	if len(a.Suggestions) > 0 {
		sb.WriteString("\nKey Suggestions:\n")
		for _, s := range a.Suggestions[:min(len(a.Suggestions), 5)] {
			sb.WriteString("- " + s + "\n")
		}
	}
	
	return strings.TrimSpace(sb.String())
}

func convertToRiskItems(clauses []RiskClause) []RiskItem {
	var items []RiskItem
	for _, c := range clauses {
		items = append(items, RiskItem{
			Clause:      c.Clause,
			RiskLevel:   c.Risk.String(),
			Explanation: c.Suggestion,
		})
	}
	return items
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
