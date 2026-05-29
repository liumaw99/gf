package password

import (
	"testing"
)

func TestHashAndVerify(t *testing.T) {
	pwd := "my-secret-password-123"

	hash, err := Hash(pwd)
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}
	if hash == "" {
		t.Fatal("hash is empty")
	}
	if hash == pwd {
		t.Fatal("hash should not equal plaintext")
	}

	// verify correct password
	if err := Verify(pwd, hash); err != nil {
		t.Fatalf("verify correct password failed: %v", err)
	}

	// verify wrong password
	if err := Verify("wrong-password", hash); err == nil {
		t.Fatal("verify wrong password should fail")
	}
}

func TestHashDifferent(t *testing.T) {
	pwd := "same-password"

	h1, _ := Hash(pwd)
	h2, _ := Hash(pwd)

	if h1 == h2 {
		t.Fatal("same password should produce different hashes")
	}
}
