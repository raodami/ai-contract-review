package nlp

// ContractType represents the type of contract being analyzed
type ContractType string

const (
	ContractTypeEmployment  ContractType = "employment"
	ContractTypeService     ContractType = "service"
	ContractTypeNDA         ContractType = "nda"
	ContractTypeSales       ContractType = "sales"
	ContractTypeLicense     ContractType = "license"
	ContractTypeLease       ContractType = "lease"
	ContractTypePartnership ContractType = "partnership"
	ContractTypeOther       ContractType = "other"
)

// Issue represents a specific issue in contract type
type Issue struct {
	Title     string `json:"title"`
	Severity  string `json:"severity"` // high, medium, low
	Found     bool   `json:"found"`
	Recommend string `json:"recommend"`
}

// Requirement represents a required clause for contract type
type Requirement struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Found    bool   `json:"found"`
}

// ScenarioAnalysis contains scenario-specific analysis
type ScenarioAnalysis struct {
	Type        ContractType   `json:"type"`
	TypeName    string         `json:"type_name"`
	Issues      []Issue        `json:"issues"`
	Requirements []Requirement `json:"requirements"`
	ScoreAdjust int            `json:"score_adjust"`
}

// TypeDetector detects contract type from text
type TypeDetector struct{}

func NewTypeDetector() *TypeDetector {
	return &TypeDetector{}
}

// ScenarioAnalyzer analyzes contracts based on their type
type ScenarioAnalyzer struct {
	detector *TypeDetector
}

func NewScenarioAnalyzer() *ScenarioAnalyzer {
	return &ScenarioAnalyzer{
		detector: NewTypeDetector(),
	}
}
