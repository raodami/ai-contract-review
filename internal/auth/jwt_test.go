package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndValidateToken(t *testing.T) {
	tokenStr, err := GenerateToken("test-user-123")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("Expected non-empty token")
	}

	userID, err := ParseToken(tokenStr)
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}
	if userID != "test-user-123" {
		t.Errorf("Expected 'test-user-123', got '%s'", userID)
	}
}

func TestParseToken_Invalid(t *testing.T) {
	_, err := ParseToken("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestParseToken_Expired(t *testing.T) {
	// Create expired token manually
	claims := jwt.MapClaims{
		"user_id": "test",
		"exp":     time.Now().Add(-1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(getTokenSecret()))
	
	_, err := ParseToken(tokenStr)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}
