package validation

import (
	"errors"

	"github.com/josephakayesi/go-cerbos-abac/application/dto"
)

func RegisterUserDTOValidation(r *dto.RegisterUserDTO) error {

	r.Trim()

	// FirstName Validation
	err := Validator.Var(r.FirstName, "required")
	if err != nil {
		return errors.New("first name is required")
	}

	err = Validator.Var(r.FirstName, "ascii")
	if err != nil {
		return errors.New("first name must be valid ascii characters")
	}

	err = Validator.Var(r.FirstName, "min=2")
	if err != nil {
		return errors.New("first name must be at least 2 characters")
	}

	err = Validator.Var(r.FirstName, "max=30")
	if err != nil {
		return errors.New("first name must be at most 30 characters")
	}

	// LastName Validaion
	err = Validator.Var(r.LastName, "required")
	if err != nil {
		return errors.New("last name is required")
	}

	err = Validator.Var(r.LastName, "ascii")
	if err != nil {
		return errors.New("last name must be valid ascii characters")
	}

	err = Validator.Var(r.LastName, "min=2")
	if err != nil {
		return errors.New("last name must be at least 2 characters")
	}

	err = Validator.Var(r.LastName, "max=30")
	if err != nil {
		return errors.New("last name must be at most 30 characters")
	}

	// Username Validation
	err = Validator.Var(r.Username, "required")
	if err != nil {
		return errors.New("username is required")
	}

	err = Validator.Var(r.Username, "username")
	if err != nil {
		return errors.New("username must contain valid characters. allowed characters: [a-z, A-Z, 0-9, _, .]. cannot contain special characters consecutively")
	}

	err = Validator.Var(r.Username, "min=2")
	if err != nil {
		return errors.New("username must be at least 2 characters")
	}

	err = Validator.Var(r.Username, "max=30")
	if err != nil {
		return errors.New("username must be at most 30 characters")
	}

	// Email Validation
	err = Validator.Var(r.Email, "required")
	if err != nil {
		return errors.New("email is required")
	}

	err = Validator.Var(r.Email, "email")
	if err != nil {
		return errors.New("email must be a valid email")
	}

	err = Validator.Var(r.Email, "max=128")
	if err != nil {
		return errors.New("email must be at most 128 characters")
	}

	err = Validator.Var(r.Password, "min=8")
	if err != nil {
		return errors.New("password must be a minimum of 8 characters")
	}

	return nil
}

func LoginUserDTOValidation(l *dto.LoginUserDto) error {

	l.Trim()

	err := Validator.Var(l.UsernameOrEmail, "required")
	if err != nil {
		return errors.New("username or email is required")
	}

	err = Validator.Var(l.Password, "min=8")
	if err != nil {
		return errors.New("login details are incorrect")
	}

	return nil
}
