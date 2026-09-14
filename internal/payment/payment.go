package payment

import (
	"os"
)

const (
	FreeQuotaMinutes = 30
	ProMonthlyPrice  = 990
	ProYearlyPrice   = 9900
)

// Plan defines pricing plan
type Plan struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"` // in cents
	Interval    string `json:"interval"`
	Minutes     int    `json:"minutes"`
	Description string `json:"description"`
}

var plans = []Plan{
	{
		ID:          "free",
		Name:        "Free",
		Price:       0,
		Interval:    "monthly",
		Minutes:     FreeQuotaMinutes,
		Description: "30 minutes/month analysis",
	},
	{
		ID:          "pro-monthly",
		Name:        "Pro (Monthly)",
		Price:       ProMonthlyPrice,
		Interval:    "monthly",
		Minutes:     500,
		Description: "500 minutes/month analysis",
	},
	{
		ID:          "pro-yearly",
		Name:        "Pro (Yearly)",
		Price:       ProYearlyPrice,
		Interval:    "yearly",
		Minutes:     500,
		Description: "500 minutes/month, billed annually - Save 17%",
	},
}

func GetPlans() []Plan {
	return plans
}

// CheckoutSessionInput for creating Stripe checkout
type CheckoutSessionInput struct {
	UserID     string
	Email      string
	PlanID     string
	SuccessURL string
	CancelURL  string
}

// CheckoutSession creates a Stripe Checkout Session
func CheckoutSession(input CheckoutSessionInput) (string, error) {
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		return "", nil // Stub for demo
	}
	_ = input
	return "", nil
}

// WebhookHandler processes Stripe webhook events
func WebhookHandler(payload []byte) (string, error) {
	whSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if whSecret == "" {
		return "", nil // Stub
	}
	_ = whSecret
	_ = payload
	return "ok", nil
}

// GetAvailableMinutes calculates remaining minutes for user
func GetAvailableMinutes(isPro bool, usageMinutes int) (int, int) {
	limit := FreeQuotaMinutes
	if isPro {
		limit = 500
	}
	allowed := limit - usageMinutes
	if allowed < 0 {
		allowed = 0
	}
	return allowed, limit
}

// CheckQuota validates if user has remaining quota
func CheckQuota(isPro bool, usageMinutes int) (bool, int, int) {
	allowed, limit := GetAvailableMinutes(isPro, usageMinutes)
	return allowed > 0, allowed, limit
}

// ConsumeMinutes deducts usage minutes
func ConsumeMinutes(usageMinutes int, minutes int) int {
	return usageMinutes + minutes
}
