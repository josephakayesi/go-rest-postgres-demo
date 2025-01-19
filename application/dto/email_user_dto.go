package dto

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/josephakayesi/go-cerbos-abac/internal"
)

type EmailUserDTO struct {
	FirstName string `json:"first_name" validate:"required,ascii,min=2,max=30"`
	LastName  string `json:"last_name"  validate:"required,ascii,min=2,max=30"`
	Email     string `json:"email"      validate:"required,email,max=128"`
}

func (r *EmailUserDTO) Trim() {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Email = strings.TrimSpace(r.Email)
}

type UserEmailedEvent struct {
	EventId   string `json:"event_id"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
}

// NewUserCreatedEvent creates a new UserCreatedEvent.
func NewUserEmailedEvent(e EmailUserDTO, c context.Context) *UserCreatedEvent {
	key := internal.GetBrowserFingerPrintKey()

	return &UserCreatedEvent{
		EventId:   uuid.New().String(),
		FirstName: e.FirstName,
		Email:     e.Email,
		ClientIP:  c.Value(key).(internal.BrowserFingerprint).ClientIP,
		UserAgent: c.Value(key).(internal.BrowserFingerprint).UserAgent,
	}
}

func (uc *UserEmailedEvent) ToJSON() (*[]byte, error) {
	js, err := json.Marshal(uc)

	if err != nil {
		return nil, err
	}

	return &js, nil
}

func (uc *UserEmailedEvent) FromJSON(jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), &uc)

	if err != nil {
		return err
	}

	return nil
}
