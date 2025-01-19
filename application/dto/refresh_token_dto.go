package dto

import (
	"strings"
)

type RefreshTokenDto struct {
	RefreshToken string `json:"refresh_token" validate:"required,ascii"`
}

func (r *RefreshTokenDto) Trim() {
	r.RefreshToken = strings.TrimSpace(r.RefreshToken)
}

type UserRefreshedTokenEvent struct {
	EventId       string `json:"event_id"`
	UserId        string `json:"user_id"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	UserSessionId string `json:"user_session_id"`
	ClientIP      string `json:"client_ip"`
	UserAgent     string `json:"user_agent"`
}
