package infra

import (
	"crypto/ed25519"
	"encoding/hex"
	"log"
	"sync"

	"github.com/josephakayesi/go-cerbos-abac/infra/config"
	"github.com/o1egl/paseto"

	vo "github.com/josephakayesi/go-cerbos-abac/domain/core/valueobject"
)

var c = config.GetConfig()

type Paseto struct {
	v2 *paseto.V2
}

// GetPasetoKeys : Generate and get paseto prvate key and public keys from secrets respectively
func GetPasetoKeys() (ed25519.PrivateKey, ed25519.PublicKey) {
	b, err := hex.DecodeString(c.PASETO_PRIVATE_KEY_SECRET)
	privateKey := ed25519.PrivateKey(b)

	if err != nil {
		log.Fatalln(err)
	}

	b, err = hex.DecodeString(c.PASETO_PUBLIC_KEY_SECRET)
	publicKey := ed25519.PublicKey(b)

	if err != nil {
		log.Fatalln(err)
	}

	return privateKey, publicKey
}

// NewPaseto : Creates a new instance of Paseto
func LoadPaseto() *Paseto {
	return &Paseto{
		v2: paseto.NewV2(),
	}
}

var NewPaseto = sync.OnceValue(LoadPaseto)

// Sign : Signs a payload using a private key
func (p *Paseto) Sign(payload interface{}) (string, error) {
	priv, _ := GetPasetoKeys()

	token, err := p.v2.Sign(priv, &payload, nil)

	if err != nil {
		return "", err
	}

	return token, nil
}

// VerifyRefreshToken : Verifies a paseto token and returns the decoded payload
func (p *Paseto) VerifyRefreshToken(token string) (*vo.RefreshTokenPayload, error) {
	payload := &vo.RefreshTokenPayload{}

	_, pub := GetPasetoKeys()

	err := p.v2.Verify(token, pub, payload, nil)

	if err != nil {
		return nil, err
	}

	return payload, nil
}

// VerifyAccessToken : Verifies a paseto token and returns the decoded payload
func (p *Paseto) VerifyAccessToken(token string) (*vo.AccessTokenPayload, error) {
	payload := &vo.AccessTokenPayload{}

	_, pub := GetPasetoKeys()

	err := p.v2.Verify(token, pub, payload, nil)

	if err != nil {
		return nil, err
	}

	return payload, nil
}
