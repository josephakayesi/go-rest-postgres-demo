package dto

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	entity "github.com/josephakayesi/go-cerbos-abac/domain/core/entity"
	vo "github.com/josephakayesi/go-cerbos-abac/domain/core/valueobject"
	"github.com/josephakayesi/go-cerbos-abac/internal"
)

type RegisterUserDTO struct {
	FirstName string      `json:"first_name" validate:"required,ascii,min=2,max=30"`
	LastName  string      `json:"last_name"  validate:"required,ascii,min=2,max=30"`
	Username  string      `json:"username"   validate:"required,username,min=2,max=30"`
	Email     string      `json:"email"      validate:"required,email,max=128"`
	Password  vo.Password `json:"password"   validate:"required,min=8"`
}

func (r *RegisterUserDTO) Trim() {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.TrimSpace(r.Email)
	r.Password = vo.Password(strings.TrimSpace(r.Password.String()))
}

type UserCreatedEvent struct {
	EventId   string `json:"event_id"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
}

// NewUserCreatedEvent creates a new UserCreatedEvent.
func NewUserCreatedEvent(u *entity.User, c context.Context) *UserCreatedEvent {
	key := internal.GetBrowserFingerPrintKey()

	return &UserCreatedEvent{
		EventId:   uuid.New().String(),
		FirstName: u.FirstName,
		Email:     u.Email,
		ClientIP:  c.Value(key).(internal.BrowserFingerprint).ClientIP,
		UserAgent: c.Value(key).(internal.BrowserFingerprint).UserAgent,
	}
}

func (uc *UserCreatedEvent) ToJSON() (*[]byte, error) {
	js, err := json.Marshal(uc)

	if err != nil {
		return nil, err
	}

	return &js, nil
}

func (uc *UserCreatedEvent) FromJSON(jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), &uc)

	if err != nil {
		return err
	}

	return nil
}
