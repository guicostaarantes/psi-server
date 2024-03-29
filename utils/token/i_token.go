package token

import "time"

// ITokenUtil is an abstraction for a utility that creates a string that can
// later be verified to have been created by this utility
type ITokenUtil interface {
	GenerateToken(payload string, secondsToExpire time.Duration) (string, error)
}
