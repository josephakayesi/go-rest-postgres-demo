package dto

type PublicKeyDto struct {
	PublicKey string `json:"public_key" validate:"required,ascii"`
}
