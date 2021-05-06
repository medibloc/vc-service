package pkg

import (
	"encoding/json"

	"github.com/medibloc/vc-sdk/pkg/vc"
)

type issueCredentialRequest struct {
	Credential json.RawMessage `json:"credential"`
	Options    signOptions     `json:"options"`
}

type signOptions struct {
	PrivateKeyBase64   string `json:"privateKeyBase64"`
	SignatureType      string `json:"signatureType"`
	VerificationMethod string `json:"verificationMethod"`
	ProofPurpose       string `json:"proofPurpose,omitempty"`
	Domain             string `json:"domain,omitempty"`
	Challenge          string `json:"challenge,omitempty"`
	CredentialStatus   string `json:"credentialStatus,omitempty"`
}

func (o *signOptions) vcProofOptions() *vc.ProofOptions {
	return &vc.ProofOptions{
		VerificationMethod: o.VerificationMethod,
		SignatureType:      o.SignatureType,
		Domain:             o.Domain,
		Challenge:          o.Challenge,
		// TODO: handle proofPurpose, credentialStatus
	}
}
