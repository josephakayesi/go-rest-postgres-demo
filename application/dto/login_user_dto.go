package dto

import (
	"strings"

	vo "github.com/josephakayesi/go-cerbos-abac/domain/core/valueobject"
)

type LoginUserDto struct {
	UsernameOrEmail string      `json:"username_or_email" validate:"required,ascii,min=2,max=30"`
	Password        vo.Password `json:"password"          validate:"required,min=8"`
}

func (l *LoginUserDto) Trim() {
	l.UsernameOrEmail = strings.TrimSpace(l.UsernameOrEmail)
	l.Password = vo.Password(strings.TrimSpace(l.Password.String()))
}

type UserLoggedInEvent struct {
	EventId   string `json:"event_id"`
	UserId    string `json:"user_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
}
