package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// SHA256 hashes bytes data using SHA256
// return the digest string of base64 encoded hashes
func SHA256(bytes []byte) string {
	h := sha256.New()
	_, err := h.Write(bytes)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// GenerateRandomBytes generates random bytes
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString generates random string
func GenerateRandomString(len int) string {
	b, err := GenerateRandomBytes(len)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)[0:len]
}

// GeneratePassword returns the hash of the password
func GeneratePassword(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

// ComparePassword compares a hashed password with its possible
// plaintext equivalent. Returns true on success, or an false on failure.
func ComparePassword(hashed string, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}
