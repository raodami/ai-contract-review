package nlp

import (
	"fmt"
	"strings"
)

func analyzeScenario(text string) *ScenarioAnalysis {
	sa := &ScenarioAnalysis{}
	sa.Type = detectType(text)
	sa.TypeName = formatTypeName(sa.Type)

	switch sa.Type {
	case ContractTypeEmployment:
		analyzeEmployment(sa, text)
	case ContractTypeService:
		analyzeService(sa, text)
	case ContractTypeNDA:
		analyzeNDA(sa, text)
	case ContractTypeSales:
		analyzeSales(sa, text)
	case ContractTypeLicense:
		analyzeLicense(sa, text)
	default:
		analyzeGeneral(sa, text)
	}

	return sa
}

func detectType(text string) ContractType {
	lower := strings.ToLower(text)

	employmentSignals := []string{
		"employee", "employer", "wage", "salary", "benefits", "probation",
		"termination of employment", "severance", "working hours",
		"job description", "duties and responsibilities",
	}
	employmentScore := 0
	for _, s := range employmentSignals {
		if strings.Contains(lower, s) {
			employmentScore++
		}
	}

	serviceSignals := []string{
		"service provider", "client", "scope of services",
		"deliverables", "service fees", "statement of work",
	}
	serviceScore := 0
	for _, s := range serviceSignals {
		if strings.Contains(lower, s) {
			serviceScore++
		}
	}

	ndaSignals := []string{
		"confidential information", "non-disclosure", "proprietary information",
		"trade secrets", "disclosure obligations",
	}
	ndaScore := 0
	for _, s := range ndaSignals {
		if strings.Contains(lower, s) {
			ndaScore++
		}
	}

	salesSignals := []string{
		"purchase order", "buyer", "seller", "goods", "products",
		"price", "payment terms", "delivery", "warranty",
	}
	salesScore := 0
	for _, s := range salesSignals {
		if strings.Contains(lower, s) {
			salesScore++
		}
	}

	licenseSignals := []string{
		"license", "licensor", "licensee", "intellectual property",
		"royalty", "patent", "copyright", "trademark",
	}
	licenseScore := 0
	for _, s := range licenseSignals {
		if strings.Contains(lower, s) {
			licenseScore++
		}
	}

	maxScore := employmentScore
	maxType := ContractTypeEmployment
	if serviceScore > maxScore {
		maxScore = serviceScore
		maxType = ContractTypeService
	}
	if ndaScore > maxScore {
		maxScore = ndaScore
		maxType = ContractTypeNDA
	}
	if salesScore > maxScore {
		maxScore = salesScore
		maxType = ContractTypeSales
	}
	if licenseScore > maxScore {
		maxScore = licenseScore
		maxType = ContractTypeLicense
	}

	if maxScore == 0 {
		return ContractTypeOther
	}
	return maxType
}

func formatTypeName(t ContractType) string {
	switch t {
	case ContractTypeEmployment:
		return "Employment Contract"
	case ContractTypeService:
		return "Service Agreement"
	case ContractTypeNDA:
		return "Non-Disclosure Agreement"
	case ContractTypeSales:
		return "Sales/Purchase Agreement"
	case ContractTypeLicense:
		return "License Agreement"
	case ContractTypeLease:
		return "Lease Agreement"
	case ContractTypePartnership:
		return "Partnership Agreement"
	default:
		return "General Contract"
	}
}

func analyzeEmployment(sa *ScenarioAnalysis, text string) {
	lower := strings.ToLower(text)

	addRequirement(sa, "compensation_terms", strings.Contains(lower, "salary") || strings.Contains(lower, "wage"))
	addRequirement(sa, "benefits", strings.Contains(lower, "benefit") || strings.Contains(lower, "insurance"))
	addRequirement(sa, "working_hours", strings.Contains(lower, "working hour") || strings.Contains(lower, "full-time") || strings.Contains(lower, "part-time"))
	addRequirement(sa, "probation_period", strings.Contains(lower, "probation"))
	addRequirement(sa, "termination_notice", strings.Contains(lower, "notice period") || strings.Contains(lower, "termination notice"))
	addRequirement(sa, "non_compete", strings.Contains(lower, "non-compete") || strings.Contains(lower, "noncompetition"))
	addRequirement(sa, "confidentiality", strings.Contains(lower, "confidential") || strings.Contains(lower, "proprietary"))
	addRequirement(sa, "ip_assignment", strings.Contains(lower, "intellectual property") || strings.Contains(lower, "work made for hire"))

	if !strings.Contains(lower, "salary") && !strings.Contains(lower, "wage") && !strings.Contains(lower, "compensation") {
		addIssue(sa, "Missing Compensation Terms", "high", false, "Add clear salary/wage provisions including payment schedule")
	}
	if !strings.Contains(lower, "probation") {
		addIssue(sa, "Missing Probation Period", "medium", false, "Define probation period and conditions for confirmation")
	}
	if strings.Contains(lower, "non-compete") && !strings.Contains(lower, "geographic") && !strings.Contains(lower, "time limit") {
		addIssue(sa, "Unreasonable Non-Compete", "high", false, "Add geographic and temporal limitations to non-compete clause")
	}
	if !strings.Contains(lower, "working hour") && !strings.Contains(lower, "full-time") {
		addIssue(sa, "Missing Working Hours Definition", "medium", false, "Define working hours and overtime policies")
	}
}

func analyzeService(sa *ScenarioAnalysis, text string) {
	lower := strings.ToLower(text)

	addRequirement(sa, "scope_of_work", strings.Contains(lower, "scope") || strings.Contains(lower, "deliverable"))
	addRequirement(sa, "payment_terms", strings.Contains(lower, "payment") || strings.Contains(lower, "fee"))
	addRequirement(sa, "term_duration", strings.Contains(lower, "term") || strings.Contains(lower, "duration"))
	addRequirement(sa, "termination", strings.Contains(lower, "termination") || strings.Contains(lower, "cancel"))
	addRequirement(sa, "warranty", strings.Contains(lower, "warranty") || strings.Contains(lower, "guarantee"))
	addRequirement(sa, "sla", strings.Contains(lower, "sla") || strings.Contains(lower, "service level"))
	addRequirement(sa, "change_order", strings.Contains(lower, "change order") || strings.Contains(lower, "amendment"))

	if !strings.Contains(lower, "scope") && !strings.Contains(lower, "deliverable") {
		addIssue(sa, "Missing Scope of Work", "high", false, "Define clear scope of services and deliverables")
	}
	if !strings.Contains(lower, "payment") && !strings.Contains(lower, "fee") {
		addIssue(sa, "Missing Payment Terms", "high", false, "Specify payment amount, schedule, and method")
	}
	if !strings.Contains(lower, "sla") && !strings.Contains(lower, "service level") {
		addIssue(sa, "Missing SLA", "medium", false, "Add service level agreements with metrics")
	}
}

func analyzeNDA(sa *ScenarioAnalysis, text string) {
	lower := strings.ToLower(text)

	addRequirement(sa, "definition_confidential", strings.Contains(lower, "confidential information") || strings.Contains(lower, "proprietary"))
	addRequirement(sa, "exclusions", strings.Contains(lower, "except") || strings.Contains(lower, "excluding") || strings.Contains(lower, "not confidential"))
	addRequirement(sa, "obligations", strings.Contains(lower, "shall not disclose") || strings.Contains(lower, "obligation"))
	addRequirement(sa, "term_duration", strings.Contains(lower, "term") || strings.Contains(lower, "period"))
	addRequirement(sa, "return_of_info", strings.Contains(lower, "return") || strings.Contains(lower, "destroy"))
	addRequirement(sa, "remedies", strings.Contains(lower, "injunction") || strings.Contains(lower, "remedy"))
	addRequirement(sa, "permitted_disclosure", strings.Contains(lower, "required by law") || strings.Contains(lower, "court order"))

	if !strings.Contains(lower, "except") && !strings.Contains(lower, "excluding") {
		addIssue(sa, "Missing Confidentiality Exclusions", "high", false, "Define what information is NOT confidential (public domain, already known, etc.)")
	}
	if !strings.Contains(lower, "term") && !strings.Contains(lower, "period") {
		addIssue(sa, "Missing NDA Term", "high", false, "Specify the duration of confidentiality obligations")
	}
	if !strings.Contains(lower, "return") && !strings.Contains(lower, "destroy") {
		addIssue(sa, "Missing Return of Information", "medium", false, "Add provision for return/destruction of confidential info upon termination")
	}
	if strings.Contains(lower, "perpetual") || strings.Contains(lower, "in perpetuity") {
		addIssue(sa, "Perpetual Confidentiality", "medium", true, "Consider adding a time limit rather than perpetual obligation")
	}
}

func analyzeSales(sa *ScenarioAnalysis, text string) {
	lower := strings.ToLower(text)

	addRequirement(sa, "price", strings.Contains(lower, "price") || strings.Contains(lower, "purchase price"))
	addRequirement(sa, "payment_terms", strings.Contains(lower, "payment") || strings.Contains(lower, "net"))
	addRequirement(sa, "delivery", strings.Contains(lower, "delivery") || strings.Contains(lower, "shipping"))
	addRequirement(sa, "warranty", strings.Contains(lower, "warranty") || strings.Contains(lower, "guarantee"))
	addRequirement(sa, "inspection", strings.Contains(lower, "inspection") || strings.Contains(lower, "acceptance"))
	addRequirement(sa, "risk_of_loss", strings.Contains(lower, "risk of loss") || strings.Contains(lower, "title"))
	addRequirement(sa, "force_majeure", strings.Contains(lower, "force majeure"))
	addRequirement(sa, "governing_law", strings.Contains(lower, "governing law") || strings.Contains(lower, "jurisdiction"))

	if !strings.Contains(lower, "price") && !strings.Contains(lower, "purchase price") {
		addIssue(sa, "Missing Price Terms", "high", false, "Specify exact pricing and currency")
	}
	if !strings.Contains(lower, "delivery") && !strings.Contains(lower, "shipping") {
		addIssue(sa, "Missing Delivery Terms", "high", false, "Define delivery method, timeline, and incoterms")
	}
	if !strings.Contains(lower, "inspection") && !strings.Contains(lower, "acceptance") {
		addIssue(sa, "Missing Inspection/Acceptance", "medium", false, "Add inspection period and acceptance criteria")
	}
}

func analyzeLicense(sa *ScenarioAnalysis, text string) {
	lower := strings.ToLower(text)

	addRequirement(sa, "grant_of_license", strings.Contains(lower, "grant") || strings.Contains(lower, "licensee"))
	addRequirement(sa, "scope", strings.Contains(lower, "scope") || strings.Contains(lower, "restricted") || strings.Contains(lower, "exclusive"))
	addRequirement(sa, "royalty", strings.Contains(lower, "royalty") || strings.Contains(lower, "license fee"))
	addRequirement(sa, "territory", strings.Contains(lower, "territory") || strings.Contains(lower, "geographic"))
	addRequirement(sa, "term", strings.Contains(lower, "term") || strings.Contains(lower, "duration"))
	addRequirement(sa, "sublicense", strings.Contains(lower, "sublicense") || strings.Contains(lower, "assign"))
	addRequirement(sa, "quality_control", strings.Contains(lower, "quality") || strings.Contains(lower, "approve"))
	addRequirement(sa, "improvements", strings.Contains(lower, "improvement") || strings.Contains(lower, "derivative"))

	if !strings.Contains(lower, "exclusive") && !strings.Contains(lower, "non-exclusive") {
		addIssue(sa, "Unclear License Scope", "high", false, "Specify whether license is exclusive or non-exclusive")
	}
	if !strings.Contains(lower, "territory") && !strings.Contains(lower, "geographic") {
		addIssue(sa, "Missing Territory", "medium", false, "Define licensed territory/geographic scope")
	}
	if !strings.Contains(lower, "royalty") && !strings.Contains(lower, "license fee") {
		addIssue(sa, "Missing Royalty Terms", "high", false, "Specify royalty rate or license fee structure")
	}
}

func analyzeGeneral(sa *ScenarioAnalysis, text string) {
	lower := strings.ToLower(text)

	addRequirement(sa, "parties", strings.Contains(lower, "party") || strings.Contains(lower, "between"))
	addRequirement(sa, "consideration", strings.Contains(lower, "consideration") || strings.Contains(lower, "value"))
	addRequirement(sa, "signatures", strings.Contains(lower, "signature") || strings.Contains(lower, "signed"))

	if !strings.Contains(lower, "party") && !strings.Contains(lower, "between") {
		addIssue(sa, "Missing Party Identification", "high", false, "Clearly identify all contracting parties")
	}
}

func addRequirement(sa *ScenarioAnalysis, name string, found bool) {
	if sa.Requirements == nil {
		sa.Requirements = make([]Requirement, 0)
	}
	sa.Requirements = append(sa.Requirements, Requirement{
		Name:     name,
		Found:    found,
		Required: true,
	})
}

func addIssue(sa *ScenarioAnalysis, title string, severity string, found bool, recommend string) {
	if sa.Issues == nil {
		sa.Issues = make([]Issue, 0)
	}
	sa.Issues = append(sa.Issues, Issue{
		Title:     title,
		Severity:  severity,
		Found:     found,
		Recommend: recommend,
	})
}

func (sa *ScenarioAnalysis) CalculateScoreAdjust() int {
	adjust := 0
	for _, issue := range sa.Issues {
		if !issue.Found {
			switch issue.Severity {
			case "high":
				adjust -= 15
			case "medium":
				adjust -= 8
			case "low":
				adjust -= 3
			}
		}
	}
	for _, req := range sa.Requirements {
		if req.Required && req.Found {
			adjust += 3
		}
	}
	return adjust
}

func formatIssueList(sa *ScenarioAnalysis) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## %s Analysis\n\n", sa.TypeName))
	if len(sa.Issues) == 0 {
		sb.WriteString("No critical issues found.\n")
		return sb.String()
	}
	for _, issue := range sa.Issues {
		status := "✓ Present"
		if !issue.Found {
			status = "✗ Missing"
		}
		sb.WriteString(fmt.Sprintf("[%s] %s - %s\n", strings.ToUpper(issue.Severity), issue.Title, status))
		if issue.Recommend != "" && !issue.Found {
			sb.WriteString(fmt.Sprintf("  → %s\n", issue.Recommend))
		}
	}
	return sb.String()
}

func formatRequirementList(sa *ScenarioAnalysis) string {
	var sb strings.Builder
	sb.WriteString("## Required Clauses\n\n")
	for _, req := range sa.Requirements {
		status := "○ Missing"
		if req.Found {
			status = "✓ Present"
		}
		sb.WriteString(fmt.Sprintf("- %s: %s\n", req.Name, status))
	}
	return sb.String()
}
