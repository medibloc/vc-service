package verifiable

import (
	"encoding/json"
	"time"

	docverifiable "github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
)

type IssueCredentialRequest struct {
	Credential json.RawMessage `json:"credential"`
	Options    *ProofOptions   `json:"options"`
}

type VerifyCredentialRequest struct {
	VerifiableCredential json.RawMessage `json:"verifiableCredential"`
}

type DeriveCredentialRequest struct {
	VerifiableCredential json.RawMessage `json:"verifiableCredential"`
	Frame                json.RawMessage `json:"frame"`
	Options              *DeriveOptions  `json:"options,omitempty"`
}

type ProvePresentationRequest struct {
	Presentation json.RawMessage `json:"presentation"`
	Options      *ProofOptions   `json:"options,omitempty"`
}

type VerifyPresentationRequest struct {
	VerifiablePresentation json.RawMessage `json:"verifiablePresentation"`
}

// ProofOptions is model to allow the dynamic proofing options by the user.
type ProofOptions struct {
	KID string `json:"kid,omitempty"`
	// VerificationMethod is the URI of the verificationMethod used for the proof.
	VerificationMethod      string                                 `json:"verificationMethod,omitempty"`
	SignatureRepresentation *docverifiable.SignatureRepresentation `json:"signatureRepresentation,omitempty"`
	// Created date of the proof. If omitted current system time will be used.
	Created *time.Time `json:"created,omitempty"`
	// Domain is operational domain of a digital proof.
	Domain string `json:"domain,omitempty"`
	// Challenge is a random or pseudo-random value option authentication
	Challenge string `json:"challenge,omitempty"`
	// SignatureType signature type used for signing
	SignatureType string `json:"signatureType,omitempty"`
	// proofPurpose is purpose of the proof.
	proofPurpose string
}

type DeriveOptions struct {
	Nonce string `json:"nonce"`
}
