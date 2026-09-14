package parser

import (
	"bytes"
	"io"
	"regexp"
	"strings"
)

// ParsePDF extracts text from a PDF file
func ParsePDF(data []byte) (string, error) {
	// PDF parsing: extract text from content streams
	text := extractTextFromPDF(data)
	return strings.TrimSpace(text), nil
}

// extractTextFromPDF parses PDF structure and extracts text content
func extractTextFromPDF(data []byte) string {
	var result strings.Builder
	
	// Find all BT/ET blocks (text operators)
	btPattern := regexp.MustCompile(`BT[\s\S]*?ET`)
	matches := btPattern.FindAll(data, -1)
	
	for _, match := range matches {
		text := extractTextOperators(match)
		if text != "" {
			result.WriteString(text)
			result.WriteString("\n")
		}
	}
	
	// Also try direct text extraction for simple PDFs
	if result.Len() == 0 {
		text := extractRawText(data)
		return text
	}
	
	return result.String()
}

// extractTextOperators extracts text from PDF text operators
func extractTextOperators(data []byte) string {
	var result strings.Builder
	
	// Match (text) Tj operator
	tjPattern := regexp.MustCompile(`\(([^)]*)\) *Tj`)
	for _, match := range tjPattern.FindAllSubmatch(data, -1) {
		if len(match) > 1 {
			result.WriteString(string(match[1]))
			result.WriteString("\n")
		}
	}
	
	// Match (text) TJ array operator
	tjArrayPattern := regexp.MustCompile(`\[([^\]]*)\] *TJ`)
	for _, match := range tjArrayPattern.FindAllSubmatch(data, -1) {
		if len(match) > 1 {
			text := extractStringFromArray(string(match[1]))
			result.WriteString(text)
			result.WriteString("\n")
		}
	}
	
	// Match ' operator (continued text)
	singleQuotePattern := regexp.MustCompile(`\(([^)]*)\) *'`)
	for _, match := range singleQuotePattern.FindAllSubmatch(data, -1) {
		if len(match) > 1 {
			result.WriteString(string(match[1]))
		}
	}
	
	return result.String()
}

// extractStringFromArray extracts text from PDF array notation
func extractStringFromArray(arr string) string {
	var result strings.Builder
	pattern := regexp.MustCompile(`\(([^)]*)\)`)
	for _, match := range pattern.FindAllStringSubmatch(arr, -1) {
		if len(match) > 1 {
			result.WriteString(match[1])
		}
	}
	return result.String()
}

// extractRawText is a fallback for simple PDFs
func extractRawText(data []byte) string {
	var result strings.Builder
	
	// Look for stream...endstream blocks
	streamPattern := regexp.MustCompile(`stream\s*\n([\s\S]*?)\s*endstream`)
	for _, match := range streamPattern.FindAllSubmatch(data, -1) {
		streamData := match[1]
		text := extractTextFromStream(streamData)
		if text != "" {
			result.WriteString(text)
			result.WriteString("\n")
		}
	}
	
	return result.String()
}

// extractTextFromStream extracts readable text from PDF stream
func extractTextFromStream(data []byte) string {
	var result strings.Builder
	
	// Extract printable characters
	for _, b := range data {
		if (b >= 32 && b <= 126) || b == 10 || b == 13 || b == 9 {
			result.WriteByte(b)
		} else if b < 32 {
			// Replace control characters with space
			result.WriteByte(' ')
		}
	}
	
	text := result.String()
	
	// Clean up excessive whitespace
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	
	// Remove PDF operators
	text = regexp.MustCompile(`\b(BT|ET|TD|Tm|Ts|Tf|Tj|TJ|TJ\)|T*)\b`).ReplaceAllString(text, "")
	
	return strings.TrimSpace(text)
}

// IsPDF checks if data is a PDF file
func IsPDF(data []byte) bool {
	return bytes.HasPrefix(data, []byte("%PDF"))
}

// ReadAll reads all data from reader
func ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
