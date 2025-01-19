package domain

import (
	"time"

	vo "github.com/josephakayesi/go-cerbos-abac/domain/core/valueobject"
	"github.com/josephakayesi/go-cerbos-abac/infra/config"
	"github.com/josephakayesi/go-cerbos-abac/internal"
)

type User struct {
	Model
	FirstName string              `json:"first_name" gorm:"not null"`
	LastName  string              `json:"last_name"  gorm:"not null"`
	Username  string              `json:"username"   gorm:"not null"`
	Email     string              `json:"email"      gorm:"not null"`
	Password  vo.Password         `json:"password"   gorm:"not null"`
	Role      vo.Role             `json:"roles"      gorm:"not null; default:'user'"`
	Status    internal.UserStatus `json:"status"     gorm:"not null"`
	CreatedAt time.Time           `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time           `json:"updated_at" gorm:"autoUpdateTime:milli"`
	Orders    []Order             `json:"orders"`
}

func (u *User) CreateAccessToken() (string, *vo.AccessTokenPayload, error) {
	c := config.GetConfig()

	payload := &vo.AccessTokenPayload{
		ID:         u.ID,
		Email:      u.Email,
		Username:   u.Username,
		Subject:    u.Email,
		Audience:   "go cerbos abac users",
		Issuer:     "go cerbos abac org",
		TokenID:    "example",
		KeyID:      "example",
		Role:       u.Role,
		NotBefore:  time.Now(),
		IssuedAt:   time.Now(),
		Expiration: time.Now().Add(time.Duration(c.PASETO_ACCESS_TOKEN_TTL) * time.Minute),
	}

	paseto := internal.NewPaseto()

	token, err := paseto.Sign(payload)

	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

func (u *User) CreateRefreshToken() (string, *vo.RefreshTokenPayload, error) {
	c := config.GetConfig()

	payload := &vo.RefreshTokenPayload{
		ID:            u.ID,
		UserSessionID: internal.GenerateUserSessionId(),
		Email:         u.Email,
		Username:      u.Username,
		Subject:       u.Email,
		Audience:      "kale capital users",
		Issuer:        "kale capital org",
		TokenID:       "example",
		KeyID:         "example",
		NotBefore:     time.Now(),
		IssuedAt:      time.Now(),
		Expiration:    time.Now().Add(time.Duration(c.PASETO_REFRESH_TOKEN_TTL) * time.Minute),
	}

	paseto := internal.NewPaseto()

	token, err := paseto.Sign(payload)

	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}
