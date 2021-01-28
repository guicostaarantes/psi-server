package token

import (
	"crypto/rand"
)

type rngToken struct {
	runes string
	size  int
}

func (r rngToken) GenerateToken(payload string, secondsToExpire int64) (string, error) {
	bytes := make([]byte, r.size)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = r.runes[b%byte(len(r.runes))]
	}

	return string(bytes), nil
}

// RngTokenUtil is an implementation of ITokenUtil that uses crypto/rand
var RngTokenUtil = rngToken{
	runes: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
	size:  64,
}
