package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"ai-contract-review/internal/store"
	"ai-contract-review/internal/auth"
)

func TestRouterSetup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s, _ := store.NewDB(":memory:")
	defer s.Close()

	r := gin.New()
	SetupRoutes(r, s)

	// Create a test user
	userID := "test-user-123"
	s.CreateUser(userID, "test@example.com", "hashedpassword")

	// Generate a token for the test user
	token, _ := auth.GenerateToken(userID)

	// Test GET /api/user/usage with token
	req, _ := http.NewRequest("GET", "/api/user/usage", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["quota_limit"] != float64(30) {
		t.Errorf("Expected quota_limit 30, got %v", resp["quota_limit"])
	}
	if resp["allowed"] != true {
		t.Errorf("Expected allowed true, got %v", resp["allowed"])
	}

	// Test GET /api/payment/plans
	req2, _ := http.NewRequest("GET", "/api/payment/plans", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("Expected 200 for plans, got %d", w2.Code)
	}
}
