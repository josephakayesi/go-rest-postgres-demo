package validation

import (
	"errors"

	"github.com/josephakayesi/go-cerbos-abac/application/dto"
)

func CreateOrderDtoValidation(r *dto.CreateOrderDto) error {

	r.Trim()

	// Description Validation
	err := Validator.Var(r.Description, "required")
	if err != nil {
		return errors.New("description is required")
	}

	err = Validator.Var(r.Description, "ascii")
	if err != nil {
		return errors.New("description must be valid ascii characters")
	}

	err = Validator.Var(r.Description, "min=2")
	if err != nil {
		return errors.New("description must be at least 2 characters")
	}

	// Amount Validaion
	err = Validator.Var(r.Amount, "required")
	if err != nil {
		return errors.New("amount is required")
	}

	err = Validator.Var(r.Amount, "number")
	if err != nil {
		return errors.New("amount must be valid ascii characters")
	}

	err = Validator.Var(r.Amount, "min=2")
	if err != nil {
		return errors.New("amount must be at least 2 characters")
	}

	return nil
}

func UpdateOrderDtoValidation(u *dto.UpdateOrderDto) error {

	u.Trim()

	// Description Validation
	err := Validator.Var(u.Description, "required")
	if err != nil {
		return errors.New("description is required")
	}

	err = Validator.Var(u.Description, "ascii")
	if err != nil {
		return errors.New("description must be valid ascii characters")
	}

	err = Validator.Var(u.Description, "min=2")
	if err != nil {
		return errors.New("description must be at least 2 characters")
	}

	// Amount Validaion
	err = Validator.Var(u.Amount, "required")
	if err != nil {
		return errors.New("amount is required")
	}

	err = Validator.Var(u.Amount, "number")
	if err != nil {
		return errors.New("amount must be valid ascii characters")
	}

	err = Validator.Var(u.Amount, "min=2")
	if err != nil {
		return errors.New("amount must be at least 2 characters")
	}

	return nil
}
