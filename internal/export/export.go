package export

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ReportData holds the data for export
type ReportData struct {
	Filename    string                 `json:"filename"`
	FileSize    int64                  `json:"file_size"`
	AnalyzedAt  int64                  `json:"analyzed_at"`
	Score       int                    `json:"score"`
	Summary     string                 `json:"summary"`
	Risks       []map[string]any       `json:"risks"`
	KeyTerms    map[string]string      `json:"key_terms"`
	Clauses     []map[string]any       `json:"clauses"`
	Compliance  []map[string]any       `json:"compliance"`
	Suggestions []string               `json:"suggestions"`
	Scenario    map[string]any         `json:"scenario,omitempty"`
}

// GenerateTextReport creates a plain text report
func GenerateTextReport(data *ReportData) string {
	var b strings.Builder

	b.WriteString("═══════════════════════════════════════════════\n")
	b.WriteString("        AI CONTRACT REVIEW - ANALYSIS REPORT    \n")
	b.WriteString("═══════════════════════════════════════════════\n\n")
	b.WriteString(fmt.Sprintf("Document: %s\n", data.Filename))
	b.WriteString(fmt.Sprintf("File Size: %d bytes\n", data.FileSize))
	b.WriteString(fmt.Sprintf("Analyzed At: %s\n\n", time.Unix(data.AnalyzedAt, 0).Format("2006-01-02 15:04:05")))

	// Risk Score
	b.WriteString("═══════════════════════════════════════════════\n")
	b.WriteString("                    RISK SCORE                    \n")
	b.WriteString("═══════════════════════════════════════════════\n\n")
	b.WriteString(fmt.Sprintf("  Overall Score: %d/100\n", data.Score))
	if data.Score >= 80 {
		b.WriteString("  Risk Level: LOW (Good)\n")
	} else if data.Score >= 50 {
		b.WriteString("  Risk Level: MEDIUM (Review Recommended)\n")
	} else {
		b.WriteString("  Risk Level: HIGH (Critical Issues Found)\n")
	}
	b.WriteString("\n")

	// Summary
	if data.Summary != "" {
		b.WriteString("═══════════════════════════════════════════════\n")
		b.WriteString("                   EXECUTIVE SUMMARY              \n")
		b.WriteString("═══════════════════════════════════════════════\n\n")
		b.WriteString(data.Summary)
		b.WriteString("\n\n")
	}

	// Risks
	if len(data.Risks) > 0 {
		b.WriteString("═══════════════════════════════════════════════\n")
		b.WriteString(fmt.Sprintf("             DETECTED RISKS (%d)                \n", len(data.Risks)))
		b.WriteString("═══════════════════════════════════════════════\n\n")
		for i, risk := range data.Risks {
			level := risk["risk_level"]
			clause := risk["clause"]
			explanation := risk["explanation"]
			b.WriteString(fmt.Sprintf("  [%d] %s - %s\n", i+1, level, clause))
			b.WriteString(fmt.Sprintf("      %s\n\n", explanation))
		}
	}

	// Clauses
	if len(data.Clauses) > 0 {
		b.WriteString("═══════════════════════════════════════════════\n")
		b.WriteString(fmt.Sprintf("          CLASSIFIED CLAUSES (%d)              \n", len(data.Clauses)))
		b.WriteString("═══════════════════════════════════════════════\n\n")
		for i, clause := range data.Clauses {
			typ := clause["type"]
			content := clause["content"]
			risk := clause["risk_level"]
			b.WriteString(fmt.Sprintf("  [%d] %s", i+1, typ))
			if risk != nil {
				b.WriteString(fmt.Sprintf(" [%s]", risk))
			}
			b.WriteString("\n")
			b.WriteString(fmt.Sprintf("      %s\n\n", content))
		}
	}

	// Compliance
	if len(data.Compliance) > 0 {
		b.WriteString("═══════════════════════════════════════════════\n")
		b.WriteString("             COMPLIANCE CHECKS                  \n")
		b.WriteString("═══════════════════════════════════════════════\n\n")
		present := 0
		missing := 0
		for _, check := range data.Compliance {
			if check["status"] == "present" {
				present++
			} else if check["status"] == "missing" {
				missing++
			}
		}
		b.WriteString(fmt.Sprintf("  Present: %d | Missing: %d | Total: %d\n\n", present, missing, len(data.Compliance)))
		for _, check := range data.Compliance {
			item := check["item"]
			status := check["status"]
			issue := check["issue"]
			statusStr := strings.ToUpper(status.(string))
			b.WriteString(fmt.Sprintf("  • %s: %s", item, statusStr))
			if issue != nil && issue != "" {
				b.WriteString(fmt.Sprintf(" - %s", issue))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Scenario Analysis
	if data.Scenario != nil {
		b.WriteString("═══════════════════════════════════════════════\n")
		b.WriteString("           SCENARIO ANALYSIS                    \n")
		b.WriteString("═══════════════════════════════════════════════\n\n")
		typeName := data.Scenario["type_name"]
		b.WriteString(fmt.Sprintf("  Contract Type: %s\n\n", typeName))

		issues := data.Scenario["issues"]
		if issuesList, ok := issues.([]any); ok {
			missingIssues := 0
			for _, issue := range issuesList {
				if issueMap, ok := issue.(map[string]any); ok {
					if found, ok := issueMap["found"].(bool); ok && !found {
						missingIssues++
					}
				}
			}
			if missingIssues > 0 {
				b.WriteString(fmt.Sprintf("  ⚠ Issues Found: %d missing requirements\n\n", missingIssues))
				for _, issue := range issuesList {
					if issueMap, ok := issue.(map[string]any); ok {
						title := issueMap["title"]
						found := issueMap["found"]
						recommend := issueMap["recommend"]
						if foundBool, ok := found.(bool); !ok || !foundBool {
							b.WriteString(fmt.Sprintf("  • %s\n", title))
							if recommendStr, ok := recommend.(string); ok && recommendStr != "" {
								b.WriteString(fmt.Sprintf("    Recommendation: %s\n", recommendStr))
							}
							b.WriteString("\n")
						}
					}
				}
			}
		}
		b.WriteString("\n")
	}

	// Suggestions
	if len(data.Suggestions) > 0 {
		b.WriteString("═══════════════════════════════════════════════\n")
		b.WriteString("              RECOMMENDATIONS                   \n")
		b.WriteString("═══════════════════════════════════════════════\n\n")
		for i, sug := range data.Suggestions {
			b.WriteString(fmt.Sprintf("  %d. %s\n\n", i+1, sug))
		}
	}

	// Key Terms
	if len(data.KeyTerms) > 0 {
		b.WriteString("═══════════════════════════════════════════════\n")
		b.WriteString("              KEY TERMS                         \n")
		b.WriteString("═══════════════════════════════════════════════\n\n")
		for term, value := range data.KeyTerms {
			b.WriteString(fmt.Sprintf("  • %s: %s\n", term, value))
		}
		b.WriteString("\n")
	}

	b.WriteString("═══════════════════════════════════════════════\n")
	b.WriteString("              END OF REPORT                     \n")
	b.WriteString("═══════════════════════════════════════════════\n")

	return b.String()
}

// GenerateHTMLReport creates an HTML report
func GenerateHTMLReport(data *ReportData) string {
	var b strings.Builder

	b.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Contract Analysis Report</title>
<style>
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #0f172a; color: #e2e8f0; padding: 40px; max-width: 900px; margin: 0 auto; }
h1 { color: #533afd; text-align: center; margin-bottom: 8px; }
h2 { color: #f8fafc; font-size: 18px; margin-top: 32px; border-bottom: 1px solid rgba(255,255,255,0.1); padding-bottom: 8px; }
.subtitle { text-align: center; color: #8899a6; margin-bottom: 40px; }
.score-box { text-align: center; padding: 32px; background: rgba(255,255,255,0.03); border-radius: 16px; margin: 24px 0; }
.score-number { font-size: 72px; font-weight: 700; color: #533afd; }
.score-label { color: #8899a6; font-size: 14px; margin-top: 8px; }
.risk-card { padding: 16px; margin: 12px 0; border-radius: 8px; border-left: 4px solid; }
.risk-high { background: rgba(239,68,68,0.1); border-color: #ef4444; }
.risk-medium { background: rgba(245,158,11,0.1); border-color: #f59e0b; }
.risk-low { background: rgba(34,197,94,0.1); border-color: #22c55e; }
.clause-card { padding: 12px 16px; margin: 8px 0; background: rgba(255,255,255,0.03); border-radius: 8px; border: 1px solid rgba(255,255,255,0.05); }
.clause-type { color: #533afd; font-weight: 600; font-size: 12px; }
.clause-content { color: #e2e8f0; margin-top: 4px; }
.meta { color: #8899a6; font-size: 14px; text-align: center; margin-bottom: 24px; }
.compliance-item { display: flex; justify-content: space-between; padding: 12px; margin: 6px 0; background: rgba(255,255,255,0.03); border-radius: 6px; }
.status-present { color: #22c55e; }
.status-missing { color: #ef4444; }
.suggestion { padding: 12px; margin: 8px 0; background: rgba(83,58,253,0.1); border-radius: 8px; border-left: 3px solid #533afd; }
.term-list { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 12px; }
.term-tag { padding: 4px 12px; background: rgba(255,255,255,0.05); border-radius: 20px; font-size: 13px; }
</style>
</head>
<body>`)

	b.WriteString(fmt.Sprintf(`<h1>AI Contract Review Report</h1>`))
	b.WriteString(fmt.Sprintf(`<div class="subtitle">Document: %s</div>`, data.Filename))
	b.WriteString(fmt.Sprintf(`<div class="meta">Generated: %s</div>`, time.Unix(data.AnalyzedAt, 0).Format("2006-01-02 15:04:05")))

	// Score
	b.WriteString(`<div class="score-box">`)
	scoreColor := "#22c55e"
	if data.Score < 50 {
		scoreColor = "#ef4444"
	} else if data.Score < 80 {
		scoreColor = "#f59e0b"
	}
	b.WriteString(fmt.Sprintf(`<div class="score-number" style="color: %s">%d</div>`, scoreColor, data.Score))
	b.WriteString(`<div class="score-label">Overall Risk Score</div>`)
	b.WriteString(fmt.Sprintf(`<div style="color: %s; margin-top: 8px;">%s</div>`, scoreColor, getRiskLabel(data.Score)))
	b.WriteString(`</div>`)

	// Summary
	if data.Summary != "" {
		b.WriteString(`<h2>Executive Summary</h2>`)
		b.WriteString(fmt.Sprintf(`<p style="line-height: 1.8; color: #cbd5e1;">%s</p>`, strings.ReplaceAll(data.Summary, "\n", "<br>")))
	}

	// Risks
	if len(data.Risks) > 0 {
		b.WriteString(`<h2>Detected Risks (<span>`)
		b.WriteString(fmt.Sprintf("%d", len(data.Risks)))
		b.WriteString(`)</span>)</h2>`)
		for _, risk := range data.Risks {
			level := risk["risk_level"].(string)
			cardClass := "risk-low"
			if level == "high" {
				cardClass = "risk-high"
			} else if level == "medium" {
				cardClass = "risk-medium"
			}
			b.WriteString(fmt.Sprintf(`<div class="risk-card %s">`, cardClass))
			b.WriteString(fmt.Sprintf(`<strong>%s - %s</strong><br><span style="color:#8899a6;font-size:13px">%s</span>`, risk["risk_level"], risk["clause"], risk["explanation"]))
			b.WriteString(`</div>`)
		}
	}

	// Clauses
	if len(data.Clauses) > 0 {
		b.WriteString(`<h2>Classified Clauses (<span>`)
		b.WriteString(fmt.Sprintf("%d", len(data.Clauses)))
		b.WriteString(`)</span>)</h2>`)
		for _, clause := range data.Clauses {
			typ := clause["type"].(string)
			content := clause["content"].(string)
			b.WriteString(fmt.Sprintf(`<div class="clause-card"><div class="clause-type">%s</div><div class="clause-content">%s</div></div>`, typ, strings.ReplaceAll(content, "\n", " ")))
		}
	}

	// Compliance
	if len(data.Compliance) > 0 {
		b.WriteString(`<h2>Compliance Checks</h2>`)
		present, missing := 0, 0
		for _, c := range data.Compliance {
			if c["status"] == "present" {
				present++
			} else {
				missing++
			}
		}
		b.WriteString(fmt.Sprintf(`<div style="margin-bottom:16px;color:#8899a6">Present: %d | Missing: %d</div>`, present, missing))
		for _, c := range data.Compliance {
			statusClass := "status-present"
			if c["status"] == "missing" {
				statusClass = "status-missing"
			}
			b.WriteString(fmt.Sprintf(`<div class="compliance-item"><span>%s</span><span class="%s">%s</span></div>`, c["item"], statusClass, c["status"]))
		}
	}

	// Scenario
	if data.Scenario != nil {
		b.WriteString(`<h2>Scenario Analysis</h2>`)
		b.WriteString(fmt.Sprintf(`<div style="padding:16px;background:rgba(83,58,253,0.1);border-radius:8px;margin:12px 0"><strong>%s</strong></div>`, data.Scenario["type_name"]))
		issues := data.Scenario["issues"]
		if issuesList, ok := issues.([]any); ok {
			for _, issue := range issuesList {
				if issueMap, ok := issue.(map[string]any); ok {
					found, _ := issueMap["found"].(bool)
					title := issueMap["title"].(string)
					rec, _ := issueMap["recommend"].(string)
					if found {
						b.WriteString(fmt.Sprintf(`<div class="compliance-item"><span>%s</span><span class="status-present">✓ Present</span></div>`, title))
					} else {
						b.WriteString(fmt.Sprintf(`<div class="compliance-item"><span>%s</span><span class="status-missing">✗ Missing</span></div>`, title))
						if rec != "" {
							b.WriteString(fmt.Sprintf(`<div style="color:#8899a6;font-size:12px;padding-left:12px;margin:4px 0">→ %s</div>`, rec))
						}
					}
				}
			}
		}
	}

	// Suggestions
	if len(data.Suggestions) > 0 {
		b.WriteString(`<h2>Recommendations</h2>`)
		for _, s := range data.Suggestions {
			b.WriteString(fmt.Sprintf(`<div class="suggestion">%s</div>`, s))
		}
	}

	// Key Terms
	if len(data.KeyTerms) > 0 {
		b.WriteString(`<h2>Key Terms</h2>`)
		b.WriteString(`<div class="term-list">`)
		for term, value := range data.KeyTerms {
			b.WriteString(fmt.Sprintf(`<span class="term-tag"><strong>%s:</strong> %s</span>`, term, value))
		}
		b.WriteString(`</div>`)
	}

	b.WriteString(`</body></html>`)
	return b.String()
}

func getRiskLabel(score int) string {
	if score >= 80 {
		return "LOW RISK - Contract is well-structured"
	} else if score >= 50 {
		return "MEDIUM RISK - Review recommended"
	}
	return "HIGH RISK - Critical issues detected"
}

// GenerateJSONReport returns JSON representation
func GenerateJSONReport(data *ReportData) (string, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// GenerateCSVReport creates CSV of risks and clauses
func GenerateCSVReport(data *ReportData) string {
	var b strings.Builder

	// Risks CSV
	b.WriteString("TYPE,LEVEL,CLAUSE,EXPLANATION\n")
	for _, risk := range data.Risks {
		b.WriteString(fmt.Sprintf("%s,%s,%s,%s\n",
			risk["type"], risk["risk_level"], risk["clause"], risk["explanation"]))
	}
	b.WriteString("\n")

	// Clauses CSV
	b.WriteString("CLAUSE_TYPE,RISK_LEVEL,CONTENT\n")
	for _, clause := range data.Clauses {
		b.WriteString(fmt.Sprintf("%s,%s,%s\n",
			clause["type"], clause["risk_level"], clause["content"]))
	}

	return b.String()
}

// ExportOptions for controlling export behavior
type ExportOptions struct {
	Format string // "text", "html", "json", "csv"
}

// GenerateExport creates report in requested format
func GenerateExport(data *ReportData, opts ExportOptions) (string, string) {
	ext := ".txt"
	content := ""

	switch opts.Format {
	case "html":
		ext = ".html"
		content = GenerateHTMLReport(data)
	case "json":
		ext = ".json"
		content, _ = GenerateJSONReport(data)
	case "csv":
		ext = ".csv"
		content = GenerateCSVReport(data)
	default:
		ext = ".txt"
		content = GenerateTextReport(data)
	}

	filename := fmt.Sprintf("contract-analysis-%s%s", time.Now().Format("20060102-150405"), ext)
	return filename, content
}

// NewReportData creates ReportData from analysis result
func NewReportData(filename string, fileSize int64, result map[string]any) *ReportData {
	data := &ReportData{
		Filename:   filename,
		FileSize:   fileSize,
		AnalyzedAt: time.Now().Unix(),
	}

	if score, ok := result["score"].(float64); ok {
		data.Score = int(score)
	}
	if summary, ok := result["summary"].(string); ok {
		data.Summary = summary
	}
	if risks, ok := result["risks"].([]any); ok {
		for _, r := range risks {
			if rm, ok := r.(map[string]any); ok {
				data.Risks = append(data.Risks, rm)
			}
		}
	}
	if terms, ok := result["key_terms"].(map[string]any); ok {
		data.KeyTerms = make(map[string]string)
		for k, v := range terms {
			data.KeyTerms[k] = fmt.Sprintf("%v", v)
		}
	}
	if clauses, ok := result["clauses"].([]any); ok {
		for _, c := range clauses {
			if cm, ok := c.(map[string]any); ok {
				data.Clauses = append(data.Clauses, cm)
			}
		}
	}
	if compliance, ok := result["compliance"].([]any); ok {
		for _, c := range compliance {
			if cm, ok := c.(map[string]any); ok {
				data.Compliance = append(data.Compliance, cm)
			}
		}
	}
	if suggestions, ok := result["suggestions"].([]any); ok {
		for _, s := range suggestions {
			if ss, ok := s.(string); ok {
				data.Suggestions = append(data.Suggestions, ss)
			}
		}
	}
	if scenario, ok := result["scenario"].(map[string]any); ok {
		data.Scenario = scenario
	}

	return data
}
