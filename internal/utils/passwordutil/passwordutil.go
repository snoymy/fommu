package passwordutil

import (
	"crypto/sha256"
	"encoding/base64"
)

func HashPassword(password string) string {
    h := sha256.New()
    h.Write([]byte(password))
    passwordHash := string(base64.StdEncoding.EncodeToString(h.Sum([]byte(password))))

    return passwordHash
}
