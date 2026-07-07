package auth

import (
	"testing"
)

func TestHasher_HashAndVerify(t *testing.T) {
	h := NewHasher()

	password := "MyP@ssw0rd!"
	hash, err := h.Hash(password)
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	if hash == "" {
		t.Fatal("Hash() returned empty string")
	}
	if !h.Verify(password, hash) {
		t.Fatal("Verify() returned false for correct password")
	}
	if h.Verify("wrong-password", hash) {
		t.Fatal("Verify() returned true for wrong password")
	}
}

func TestHasher_EmptyString(t *testing.T) {
	h := NewHasher()
	hash, err := h.Hash("")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	if !h.Verify("", hash) {
		t.Fatal("Verify() returned false for empty string")
	}
}

func TestTokenGenerator_GenerateAndValidate(t *testing.T) {
	tg := NewTokenGenerator("test-secret", 24)

	tokenStr, err := tg.Generate(123, "customer")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if tokenStr == "" {
		t.Fatal("Generate() returned empty token")
	}

	claims, err := tg.Validate(tokenStr)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if claims.UserID != 123 {
		t.Errorf("Validate() UserID = %d, want %d", claims.UserID, 123)
	}
	if claims.Role != "customer" {
		t.Errorf("Validate() Role = %s, want %s", claims.Role, "customer")
	}
}

func TestTokenGenerator_GenerateProvider(t *testing.T) {
	tg := NewTokenGenerator("test-secret", 24)

	tokenStr, err := tg.Generate(456, "provider")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	claims, err := tg.Validate(tokenStr)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if claims.UserID != 456 {
		t.Errorf("Validate() UserID = %d, want %d", claims.UserID, 456)
	}
	if claims.Role != "provider" {
		t.Errorf("Validate() Role = %s, want %s", claims.Role, "provider")
	}
}

func TestTokenGenerator_InvalidToken(t *testing.T) {
	tg := NewTokenGenerator("test-secret", 24)

	_, err := tg.Validate("invalid-token-string")
	if err == nil {
		t.Fatal("Validate() expected error for invalid token")
	}
}

func TestTokenGenerator_TamperedToken(t *testing.T) {
	tg1 := NewTokenGenerator("secret-1", 24)
	tg2 := NewTokenGenerator("secret-2", 24)

	tokenStr, err := tg1.Generate(1, "customer")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	_, err = tg2.Validate(tokenStr)
	if err == nil {
		t.Fatal("Validate() expected error for token signed with different secret")
	}
}

func TestTokenGenerator_ExpiredToken(t *testing.T) {
	tg := NewTokenGenerator("test-secret", 0)

	tokenStr, err := tg.Generate(1, "customer")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	_, err = tg.Validate(tokenStr)
	if err == nil {
		t.Fatal("Validate() expected error for expired token")
	}
}
