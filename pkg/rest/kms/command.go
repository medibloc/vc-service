package kms

import (
	"encoding/base64"
	"fmt"

	"github.com/hyperledger/aries-framework-go/pkg/kms"
)

type command struct {
	ctx provider
}

// provider contains dependencies for the service command and is typically created by using aries.Context().
type provider interface {
	KMS() kms.KeyManager
}

func newCommand(p provider) *command {
	return &command{ctx: p}
}

func (c *command) CreateKeySet(request *CreateKeySetRequest) (*CreateKeySetResponse, error) {
	if request.KeyType == "" {
		return nil, fmt.Errorf("invalid key type: %v", request.KeyType)
	}

	keyID, pubKeyBytes, err := c.ctx.KMS().CreateAndExportPubKeyBytes(kms.KeyType(request.KeyType))
	if err != nil {
		return nil, fmt.Errorf("failed to create/export public key bytes: %w", err)
	}

	return &CreateKeySetResponse{
		KeyID:     keyID,
		PublicKey: base64.RawURLEncoding.EncodeToString(pubKeyBytes),
	}, nil
}
