package payment

import (
	"testing"
)

func TestGetPlans(t *testing.T) {
	plans := GetPlans()
	if len(plans) != 3 {
		t.Errorf("Expected 3 plans, got %d", len(plans))
	}
	if plans[0].ID != "free" {
		t.Errorf("Expected first plan ID 'free', got %s", plans[0].ID)
	}
	if plans[1].Minutes != 500 {
		t.Errorf("Expected Pro plan 500 minutes, got %d", plans[1].Minutes)
	}
}

func TestCheckQuota_Free(t *testing.T) {
	allowed, remaining, limit := CheckQuota(false, 10)
	if !allowed {
		t.Error("Expected allowed for free user with usage < limit")
	}
	if remaining != 20 {
		t.Errorf("Expected 20 remaining, got %d", remaining)
	}
	if limit != 30 {
		t.Errorf("Expected limit 30, got %d", limit)
	}
}

func TestCheckQuota_FreeExceeded(t *testing.T) {
	allowed, remaining, _ := CheckQuota(false, 40)
	if allowed {
		t.Error("Expected not allowed for free user with exceeded quota")
	}
	if remaining != 0 {
		t.Errorf("Expected 0 remaining, got %d", remaining)
	}
}

func TestCheckQuota_Pro(t *testing.T) {
	allowed, remaining, limit := CheckQuota(true, 100)
	if !allowed {
		t.Error("Expected allowed for pro user")
	}
	if remaining != 400 {
		t.Errorf("Expected 400 remaining, got %d", remaining)
	}
	if limit != 500 {
		t.Errorf("Expected limit 500, got %d", limit)
	}
}

func TestConsumeMinutes(t *testing.T) {
	newUsage := ConsumeMinutes(10, 5)
	if newUsage != 15 {
		t.Errorf("Expected 15, got %d", newUsage)
	}
}

func TestGetAvailableMinutes_Pro(t *testing.T) {
	allowed, limit := GetAvailableMinutes(true, 0)
	if allowed != 500 {
		t.Errorf("Expected 500 allowed, got %d", allowed)
	}
	if limit != 500 {
		t.Errorf("Expected limit 500, got %d", limit)
	}
}
