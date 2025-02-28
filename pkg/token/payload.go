package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	app "go-rest-api-boilerplate/internal"
	"time"
)

var (
	ErrExpiredToken = errors.New("token is expired")
	ErrInvalidToken = errors.New("invalid token")
)

// Payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	Role      int32     `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

//func (p *Payload) GetExpirationTime() (*time.Time, error) {
//	if time.Now().After(p.ExpiredAt) {
//		return nil, ErrExpiredToken
//	}
//	return &p.ExpiredAt, nil
//}

// NewPayload creates a new payload
func NewPayload(userId int64, role int32, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		UserID:    userId,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (m *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}
func (m *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	if time.Now().After(m.ExpiredAt) {
		return nil, app.ErrTokenExpired
	}
	return &jwt.NumericDate{Time: m.ExpiredAt}, nil
}
func (m *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: m.IssuedAt}, nil
}
func (m *Payload) GetIssuer() (string, error) {
	return "", nil
}
func (m *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: m.IssuedAt}, nil
}
func (m *Payload) GetSubject() (string, error) {
	return m.ID.String(), nil
}
