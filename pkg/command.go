package pkg

import (
	"encoding/base64"
	"fmt"

	"github.com/medibloc/vc-sdk/pkg/vc"
)

func issueCredential(req *issueCredentialRequest) ([]byte, error) {
	privKey, err := base64.StdEncoding.DecodeString(req.Options.PrivateKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode a base64 private key: %w", err)
	}

	cred, err := vc.SignCredential(req.Credential, privKey, req.Options.vcProofOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to sign credential: %w", err)
	}
	return cred, nil
}
