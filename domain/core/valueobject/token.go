package domain

import (
	"time"
)

// AccessTokenPayload: Describes structure of access token payload
type AccessTokenPayload struct {
	ID         string    `json:"_id"`
	UserID     string    `json:"user_id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Role       Role      `json:"role"`
	Subject    string    `json:"sub"`
	Audience   string    `json:"aud"`
	Issuer     string    `json:"iss"`
	TokenID    string    `json:"jti"`
	KeyID      string    `json:"kid"`
	NotBefore  time.Time `json:"nbf"`
	IssuedAt   time.Time `json:"iat"`
	Expiration time.Time `json:"exp"`
}

// RefreshTokenPayload: Describes structure of refresh token payload
type RefreshTokenPayload struct {
	ID            string    `json:"_id"`
	UserSessionID string    `json:"user_session_id"`
	UserID        string    `json:"user_id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	Subject       string    `json:"sub"`
	Audience      string    `json:"aud"`
	Issuer        string    `json:"iss"`
	TokenID       string    `json:"jti"`
	KeyID         string    `json:"kid"`
	NotBefore     time.Time `json:"nbf"`
	IssuedAt      time.Time `json:"iat"`
	Expiration    time.Time `json:"exp"`
}

// UserCredentials: Describes structure of user credentials
type UserCredentials struct {
	Email    string
	Username string
}
