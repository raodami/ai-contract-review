package parser

import (
	"strings"
	"testing"
)

func TestParsePDF_Invalid(t *testing.T) {
	_, err := ParsePDF([]byte("not a pdf"))
	if err != nil {
		t.Fatalf("Expected no error for invalid PDF, got: %v", err)
	}
}

func TestParsePDF_Empty(t *testing.T) {
	text, err := ParsePDF([]byte{})
	if err != nil {
		t.Fatalf("ParsePDF failed: %v", err)
	}
	if text != "" {
		t.Errorf("Expected empty text, got: %s", text)
	}
}

func TestIsPDF(t *testing.T) {
	if !IsPDF([]byte("%PDF-1.4")) {
		t.Error("Expected IsPDF to return true for PDF header")
	}
	if IsPDF([]byte("not a pdf")) {
		t.Error("Expected IsPDF to return false")
	}
}

func TestExtractTextOperators(t *testing.T) {
	data := []byte("BT /F1 12 Tf (Hello World) Tj ET")
	text := extractTextOperators(data)
	if strings.TrimSpace(text) != "Hello World" {
		t.Errorf("Expected 'Hello World', got: %s", strings.TrimSpace(text))
	}
}

func TestExtractTextFromArray(t *testing.T) {
	arr := "[ (Hello) (World) ]"
	text := extractStringFromArray(arr)
	if text != "HelloWorld" {
		t.Errorf("Expected 'HelloWorld', got: %s", text)
	}
}
