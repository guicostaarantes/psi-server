package token

import (
	"crypto/rand"
	"time"
)

type RngTokenUtil struct {
	Runes string
	Size  int
}

func (r RngTokenUtil) GenerateToken(payload string, secondsToExpire time.Duration) (string, error) {
	bytes := make([]byte, r.Size)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = r.Runes[b%byte(len(r.Runes))]
	}

	return string(bytes), nil
}
