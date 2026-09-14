package export

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"
	"time"
)

// GenerateDOCXReport creates a DOCX report from analysis data
func GenerateDOCXReport(data *ReportData) (string, []byte) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	// [Content_Types].xml
	addFile(w, "[Content_Types].xml", contentTypesXML())

	// word/document.xml
	docXML := buildDocumentXML(data)
	addFile(w, "word/document.xml", docXML)

	// word/_rels/document.xml.rels
	addFile(w, "word/_rels/document.xml.rels", documentRelsXML())

	// _rels/.rels
	addFile(w, "_rels/.rels", rootRelsXML())

	w.Close()
	return fmt.Sprintf("contract-analysis-%s.docx", time.Now().Format("20060102-150405")), buf.Bytes()
}

func addFile(w *zip.Writer, name string, content []byte) {
	f, _ := w.Create(name)
	f.Write(content)
}

func contentTypesXML() []byte {
	return []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`)
}

func rootRelsXML() []byte {
	return []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`)
}

func documentRelsXML() []byte {
	return []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/2006/relationships">
</Relationships>`)
}

func buildDocumentXML(data *ReportData) []byte {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body>`)

	// Title
	b.WriteString(`<w:p><w:pPr><w:spacing w:after="200"/></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="36"/><w:color w:val="533afd"/></w:rPr><w:t>AI Contract Review Report</w:t></w:r></w:p>`)
	b.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:spacing w:after="100"/></w:pPr><w:r><w:rPr><w:sz w:val="20"/><w:color w:val="8899a6"/></w:rPr><w:t>%s</w:t></w:r></w:p>`, escapeXML(data.Filename)))
	b.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:spacing w:after="200"/></w:pPr><w:r><w:rPr><w:sz w:val="20"/><w:color w:val="8899a6"/></w:rPr><w:t>Generated: %s</w:t></w:r></w:p>`, time.Unix(data.AnalyzedAt, 0).Format("2006-01-02 15:04:05")))

	// Score section
	b.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:spacing w:before="200" w:after="100"/></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="24"/></w:rPr><w:t>Risk Score: %d/100</w:t></w:r></w:p>`, data.Score))

	// Summary
	if data.Summary != "" {
		b.WriteString(`<w:p><w:pPr><w:spacing w:before="200"/></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="24"/></w:rPr><w:t>Executive Summary</w:t></w:r></w:p>`)
		for _, line := range stringsToList(data.Summary) {
			b.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:spacing w:after="60"/></w:pPr><w:r><w:rPr><w:sz w:val="20"/></w:rPr><w:t>%s</w:t></w:r></w:p>`, escapeXML(line)))
		}
	}

	// Risks
	if len(data.Risks) > 0 {
		b.WriteString(`<w:p><w:pPr><w:spacing w:before="200"/></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="24"/></w:rPr><w:t>Detected Risks</w:t></w:r></w:p>`)
		for _, risk := range data.Risks {
			level := risk["risk_level"].(string)
			clause := risk["clause"].(string)
			explanation := risk["explanation"].(string)
			b.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:spacing w:after="60"/></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="20"/></w:rPr><w:t>[%.14s] %s</w:t></w:r></w:p>`, level, escapeXML(clause)))
			b.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:spacing w:after="60"/></w:pPr><w:r><w:rPr><w:sz w:val="20"/><w:color w:val="8899a6"/></w:rPr><w:t>%s</w:t></w:r></w:p>`, escapeXML(explanation)))
		}
	}

	// Clauses
	if len(data.Clauses) > 0 {
		b.WriteString(`<w:p><w:pPr><w:spacing w:before="200"/></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="24"/></w:rPr><w:t>Classified Clauses</w:t></w:r></w:p>`)
		for _, clause := range data.Clauses {
			typ := clause["type"].(string)
			content := clause["content"].(string)
			b.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:spacing w:after="60"/></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="20"/><w:color w:val="533afd"/></w:rPr><w:t>%s</w:t></w:r></w:p>`, escapeXML(typ)))
			b.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:spacing w:after="60"/></w:pPr><w:r><w:rPr><w:sz w:val="20"/></w:rPr><w:t>%s</w:t></w:r></w:p>`, escapeXML(content)))
		}
	}

	// Compliance
	if len(data.Compliance) > 0 {
		b.WriteString(`<w:p><w:pPr><w:spacing w:before="200"/></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="24"/></w:rPr><w:t>Compliance Check</w:t></w:r></w:p>`)
		present, missing := 0, 0
		for _, c := range data.Compliance {
			if c["status"] == "present" {
				present++
			} else {
				missing++
			}
		}
		b.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:spacing w:after="60"/></w:pPr><w:r><w:rPr><w:sz w:val="20"/></w:rPr><w:t>Present: %d | Missing: %d</w:t></w:r></w:p>`, present, missing))
		for _, c := range data.Compliance {
			b.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:spacing w:after="60"/></w:pPr><w:r><w:rPr><w:sz w:val="20"/></w:rPr><w:t>• %s: %s</w:t></w:r></w:p>`, c["item"], c["status"]))
		}
	}

	// Scenario
	if data.Scenario != nil {
		typeName, _ := data.Scenario["type_name"].(string)
		b.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:spacing w:before="200"/></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="24"/></w:rPr><w:t>Contract Type: %s</w:t></w:r></w:p>`, escapeXML(typeName)))

		issues := data.Scenario["issues"]
		if issuesList, ok := issues.([]any); ok {
			for _, issue := range issuesList {
				if issueMap, ok := issue.(map[string]any); ok {
					title, _ := issueMap["title"].(string)
					found, _ := issueMap["found"].(bool)
					if !found {
						b.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:spacing w:after="60"/></w:pPr><w:r><w:rPr><w:sz w:val="20"/><w:color w:val="ef4444"/></w:rPr><w:t>✗ %s</w:t></w:r></w:p>`, escapeXML(title)))
					} else {
						b.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:spacing w:after="60"/></w:pPr><w:r><w:rPr><w:sz w:val="20"/><w:color w:val="22c55e"/></w:rPr><w:t>✓ %s</w:t></w:r></w:p>`, escapeXML(title)))
					}
				}
			}
		}
	}

	// Suggestions
	if len(data.Suggestions) > 0 {
		b.WriteString(`<w:p><w:pPr><w:spacing w:before="200"/></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="24"/></w:rPr><w:t>Recommendations</w:t></w:r></w:p>`)
		for i, sug := range data.Suggestions {
			b.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:spacing w:after="60"/></w:pPr><w:r><w:rPr><w:sz w:val="20"/></w:rPr><w:t>%d. %s</w:t></w:r></w:p>`, i+1, escapeXML(sug)))
		}
	}

	// Key Terms
	if len(data.KeyTerms) > 0 {
		b.WriteString(`<w:p><w:pPr><w:spacing w:before="200"/></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="24"/></w:rPr><w:t>Key Terms</w:t></w:r></w:p>`)
		for term, value := range data.KeyTerms {
			b.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:spacing w:after="60"/></w:pPr><w:r><w:rPr><w:sz w:val="20"/></w:rPr><w:t>• %s: %s</w:t></w:r></w:p>`, escapeXML(term), escapeXML(value)))
		}
	}

	b.WriteString(`</w:body></w:document>`)
	return []byte(b.String())
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}

func stringsToList(s string) []string {
	lines := strings.Split(s, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}
