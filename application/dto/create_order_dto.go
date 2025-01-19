package dto

import (
	"strings"
)

type CreateOrderDto struct {
	Description string `json:"description" validate:"required,ascii,min=2"`
	Amount      int    `json:"amount"      validate:"required,number,min=1"`
	UserID      string
}

func (r *CreateOrderDto) Trim() {
	r.Description = strings.TrimSpace(r.Description)
}
