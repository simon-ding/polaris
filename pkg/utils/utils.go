package utils

import (
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func isASCII(s string) bool {
    for _, c := range s {
        if c > unicode.MaxASCII {
            return false
        }
    }
    return true
}

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func IsChineseChar(str string) bool {

    for _, r := range str {
        if unicode.Is(unicode.Han, r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
            return true
        }
    }
    return false
}
