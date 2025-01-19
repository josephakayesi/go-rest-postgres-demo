package dto

import (
	"strings"
)

type UpdateOrderDto struct {
	Description string `json:"description" validate:"required,ascii,min=2"`
	Amount      int    `json:"amount"      validate:"required,number,min=1"`
}

func (u *UpdateOrderDto) Trim() {
	u.Description = strings.TrimSpace(u.Description)
}
