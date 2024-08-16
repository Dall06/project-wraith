package apikey

import (
	"crypto/sha256"
	"encoding/hex"
)

func CrateApiKey(secret string) string {
	hash := sha256.New()
	hash.Write([]byte(secret))
	createdKey := hex.EncodeToString(hash.Sum(nil))

	return createdKey
}
