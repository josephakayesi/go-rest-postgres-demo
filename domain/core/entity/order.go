package domain

import (
	"time"

	vo "github.com/josephakayesi/go-cerbos-abac/domain/core/valueobject"
)

type Order struct {
	Model
	Description string         `json:"description" gorm:"not null"`
	Amount      int            `json:"amount"      gorm:"not null"`
	UserID      string         `json:"user_id"     gorm:"not null"`
	User        User           `json:"user"`
	Status      vo.OrderStatus `json:"status"      gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"  gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at"  gorm:"autoUpdateTime:milli"`
}
