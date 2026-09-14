package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	pass := "test-password-123"
	hash, err := HashPassword(pass)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	if hash == pass {
		t.Error("Hash should not equal plaintext password")
	}
	if !CheckPassword(pass, hash) {
		t.Error("CheckPassword should return true for correct password")
	}
	if CheckPassword("wrong-password", hash) {
		t.Error("CheckPassword should return false for wrong password")
	}
}
