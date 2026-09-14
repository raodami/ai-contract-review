package parser

import (
	"archive/zip"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// ParseDOCX extracts text content from a .docx file
func ParseDOCX(data []byte) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("failed to open docx: %w", err)
	}

	var docFile *zip.File
	for _, f := range reader.File {
		if strings.Contains(f.Name, "word/document.xml") {
			docFile = f
			break
		}
	}

	if docFile == nil {
		return "", fmt.Errorf("document.xml not found")
	}

	content, err := docFile.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open document.xml: %w", err)
	}
	defer content.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(content)
	xmlContent := buf.String()

	// Extract text from <w:t>...</w:t> tags
	re := regexp.MustCompile(`<w:t[^>]*>([^<]*)</w:t>`)
	matches := re.FindAllStringSubmatch(xmlContent, -1)

	var sb strings.Builder
	for _, m := range matches {
		if len(m) > 1 && strings.TrimSpace(m[1]) != "" {
			sb.WriteString(m[1])
			sb.WriteString(" ")
		}
	}

	result := strings.TrimSpace(sb.String())
	if result == "" {
		return "", fmt.Errorf("no text found in document")
	}
	return result, nil
}
