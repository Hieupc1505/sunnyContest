package token

import "time"

// Maker is an interface manager token
type Maker interface {
	// CreateToken generates a JWT token with the given user ID, role, and duration.
	CreateToken(userId int64, role int32, duration time.Duration) (string, *Payload, error)

	// VerifyToken verifies the token and returns the payload if valid, otherwise returns an error.
	VerifyToken(token string) (*Payload, error)
}
