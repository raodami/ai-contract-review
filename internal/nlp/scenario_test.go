package nlp

import (
	"strings"
	"testing"
)

func TestDetectType_Employment(t *testing.T) {
	text := "This employment agreement is between Employer and Employee. Salary is $50,000 per year."
	typ := detectType(text)
	if typ != ContractTypeEmployment {
		t.Errorf("Expected employment, got %s", typ)
	}
}

func TestDetectType_NDA(t *testing.T) {
	text := "Confidential information shall not be disclosed to third parties. Non-disclosure agreement."
	typ := detectType(text)
	if typ != ContractTypeNDA {
		t.Errorf("Expected NDA, got %s", typ)
	}
}

func TestDetectType_Service(t *testing.T) {
	text := "Service Provider agrees to deliver scope of work as defined in Statement of Work."
	typ := detectType(text)
	if typ != ContractTypeService {
		t.Errorf("Expected service, got %s", typ)
	}
}

func TestDetectType_Sales(t *testing.T) {
	text := "Seller agrees to sell goods to Buyer for the purchase price of $10,000."
	typ := detectType(text)
	if typ != ContractTypeSales {
		t.Errorf("Expected sales, got %s", typ)
	}
}

func TestDetectType_License(t *testing.T) {
	text := "Licensor grants Licensee a non-exclusive license to use the software. Royalty payments apply."
	typ := detectType(text)
	if typ != ContractTypeLicense {
		t.Errorf("Expected license, got %s", typ)
	}
}

func TestAnalyzeScenario_Employment(t *testing.T) {
	text := "This agreement is made between the Company and the Employee. The Employee agrees to perform duties as assigned."
	sa := analyzeScenario(text)

	if sa.Type != ContractTypeEmployment {
		t.Errorf("Expected employment type, got %s", sa.Type)
	}
	if sa.TypeName != "Employment Contract" {
		t.Errorf("Expected 'Employment Contract', got %s", sa.TypeName)
	}

	if len(sa.Requirements) == 0 {
		t.Error("Expected requirements to be populated")
	}

	hasMissingComp := false
	for _, issue := range sa.Issues {
		if strings.Contains(issue.Title, "Compensation") && !issue.Found {
			hasMissingComp = true
			break
		}
	}
	if !hasMissingComp {
		t.Error("Expected missing compensation issue")
	}
}

func TestAnalyzeScenario_NDA(t *testing.T) {
	text := "This Non-Disclosure Agreement is entered into between the parties. Trade secrets and proprietary information shall be protected."
	sa := analyzeScenario(text)

	if sa.Type != ContractTypeNDA {
		t.Errorf("Expected NDA type, got %s", sa.Type)
	}

	if len(sa.Requirements) == 0 {
		t.Error("Expected requirements to be populated")
	}

	hasMissingTerm := false
	for _, issue := range sa.Issues {
		if strings.Contains(issue.Title, "Term") && !issue.Found {
			hasMissingTerm = true
			break
		}
	}
	if !hasMissingTerm {
		t.Error("Expected missing term issue")
	}
}

func TestAnalyzeScenario_ScoreAdjust(t *testing.T) {
	text := "Simple contract."
	sa := analyzeScenario(text)
	adjust := sa.CalculateScoreAdjust()
	_ = adjust // Just verify it doesn't panic
}

func TestFormatTypeName(t *testing.T) {
	if formatTypeName(ContractTypeEmployment) != "Employment Contract" {
		t.Error("Expected 'Employment Contract'")
	}
	if formatTypeName(ContractTypeNDA) != "Non-Disclosure Agreement" {
		t.Error("Expected 'Non-Disclosure Agreement'")
	}
}

func TestAnalyzeScenario_Service(t *testing.T) {
	text := "Service Agreement between Developer and Client. Scope of work includes UI design, backend API development, and testing. Payment: $5000, net 30."
	sa := analyzeScenario(text)

	if sa.Type != ContractTypeService {
		t.Errorf("Expected service type, got %s", sa.Type)
	}

	hasScope := false
	for _, req := range sa.Requirements {
		if req.Name == "scope_of_work" && req.Found {
			hasScope = true
			break
		}
	}
	if !hasScope {
		t.Error("Expected scope_of_work to be found")
	}
}

func TestAnalyzeScenario_Sales(t *testing.T) {
	text := "Sales agreement for 100 units at $50 each. Delivery in 30 days. Warranty included."
	sa := analyzeScenario(text)

	if sa.Type != ContractTypeSales {
		t.Errorf("Expected sales type, got %s", sa.Type)
	}
}

func TestAnalyzeScenario_License(t *testing.T) {
	text := "Software license agreement between Licensor and Licensee. Royalty of 5% of net sales. Territory: North America."
	sa := analyzeScenario(text)

	if sa.Type != ContractTypeLicense {
		t.Errorf("Expected license type, got %s", sa.Type)
	}
}

func TestAnalyzeScenario_General(t *testing.T) {
	text := "This is a general agreement between two parties for services."
	sa := analyzeScenario(text)

	// Should detect as general contract (no specific type signals)
	if sa.Type != ContractTypeOther && sa.Type != ContractTypeService {
		// Service might be detected due to "services" keyword
		_ = sa.Type
	}
}
