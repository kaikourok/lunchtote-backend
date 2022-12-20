package secure

import (
	"crypto/rand"
	"encoding/hex"
)

func generateSecureRandomBytes(bytes int) []byte {
	k := make([]byte, bytes)
	if _, err := rand.Read(k); err != nil {
		panic(err)
	}
	return k
}

func GenerateSecureRandomHex(bytes int) string {
	return hex.EncodeToString(generateSecureRandomBytes(bytes))
}
