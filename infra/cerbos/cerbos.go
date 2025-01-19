package infra

import (
	"github.com/cerbos/cerbos-sdk-go/cerbos"
	"github.com/josephakayesi/go-cerbos-abac/infra/config"
)

func NewCerbos(config *config.Config) (*cerbos.GRPCClient, error) {
	c, err := cerbos.New(config.CERBOS_URL, cerbos.WithPlaintext())
	if err != nil {
		return nil, err
	}

	return c, nil
}
