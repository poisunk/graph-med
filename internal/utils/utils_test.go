package utils

import "testing"

func TestHashPassword(t *testing.T) {
	hashedPassword, err := HashPassword("123456")
	if err != nil {
		t.Errorf("HashPassword error: %v", err)
	}
	t.Logf("Hashed password: %s", hashedPassword)

	valid := ComparePasswords(hashedPassword, "123456")
	if !valid {
		t.Errorf("ComparePasswords error")
	}

	valid = ComparePasswords(hashedPassword, "1234567")
	if valid {
		t.Errorf("ComparePasswords error")
	}
}
