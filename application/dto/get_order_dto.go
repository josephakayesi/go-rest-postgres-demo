package dto

import "time"

type GetOrderDto struct {
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Amount      int             `json:"amount"`
	Status      string          `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	User        GetOrderUserDto `json:"user"`
}

type GetOrderUserDto struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
}
