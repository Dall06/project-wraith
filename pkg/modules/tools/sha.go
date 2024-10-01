package tools

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

func Sha512(secret, input string) string {
	key := []byte(secret)
	h := hmac.New(sha512.New, key)
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}
