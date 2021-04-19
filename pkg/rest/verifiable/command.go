package verifiable

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/piprate/json-gold/ld"

	ariescrypto "github.com/hyperledger/aries-framework-go/pkg/crypto"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/hyperledger/aries-framework-go/pkg/doc/presexch"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/jsonld"
	verifiablesigner "github.com/hyperledger/aries-framework-go/pkg/doc/signature/signer"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/bbsblssignature2020"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2018"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/jsonwebsignature2020"
	ariesverifiable "github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries/api/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	"github.com/hyperledger/aries-framework-go/spi/storage"
	log "github.com/sirupsen/logrus"
)

const (
	creatorParts = 2

	// Ed25519Signature2018 ed25519 signature suite.
	Ed25519Signature2018 = "Ed25519Signature2018"
	// JSONWebSignature2020 json web signature suite.
	JSONWebSignature2020 = "JsonWebSignature2020"
	// BbsBlsSignature2020 BBS signature suite.
	BbsBlsSignature2020 = "BbsBlsSignature2020"

	// Ed25519VerificationKey ED25519 verification key type.
	Ed25519VerificationKey = "Ed25519VerificationKey"
)

const bbsContext = "https://w3id.org/security/bbs/v1"

// command contains command operations provided by service credential controller.
type command struct {
	kResolver keyResolver
	ctx       provider
}

type provable interface {
	AddLinkedDataProof(context *ariesverifiable.LinkedDataProofContext, jsonldOpts ...jsonld.ProcessorOpts) error
}

type keyResolver interface {
	PublicKeyFetcher() ariesverifiable.PublicKeyFetcher
}

// provider contains dependencies for the service command and is typically created by using aries.Context().
type provider interface {
	StorageProvider() storage.Provider
	VDRegistry() vdr.Registry
	KMS() kms.KeyManager
	Crypto() ariescrypto.Crypto
}

// newCommand returns newCommand service credential controller command instance.
func newCommand(p provider) *command {
	return &command{
		kResolver: ariesverifiable.NewDIDKeyResolver(p.VDRegistry()),
		ctx:       p,
	}
}

func (o *command) IssueCredential(request *IssueCredentialRequest) (*ariesverifiable.Credential, error) {
	credential, err := ariesverifiable.ParseCredential(request.Credential, ariesverifiable.WithDisabledProofCheck())
	if err != nil {
		return nil, fmt.Errorf("failed to parse credential: %w", err)
	}

	var didDoc *did.Doc
	doc, err := o.ctx.VDRegistry().Resolve(credential.Issuer.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve DID document of %v: %w", credential.Issuer.ID, err)
	}
	didDoc = doc.DIDDocument

	if err := o.addProof(credential, didDoc, request.Options); err != nil {
		return nil, fmt.Errorf("failed to add credential proof: %w", err)
	}

	return credential, nil
}

func (o *command) VerifyCredential(request *VerifyCredentialRequest) error {
	_, err := ariesverifiable.ParseCredential(
		request.VerifiableCredential,
		ariesverifiable.WithPublicKeyFetcher(ariesverifiable.NewDIDKeyResolver(o.ctx.VDRegistry()).PublicKeyFetcher()),
	)
	if err != nil {
		return fmt.Errorf("failed to parse(verify) credential: %w", err)
	}
	return nil
}

func (o *command) DeriveCredential(request *DeriveCredentialRequest) (*ariesverifiable.Credential, error) {
	vc, err := ariesverifiable.ParseCredential(request.VerifiableCredential, ariesverifiable.WithDisabledProofCheck())
	if err != nil {
		return nil, fmt.Errorf("failed to parse service credential: %w", err)
	}

	revealDoc, err := toMap(request.Frame)
	if err != nil {
		return nil, fmt.Errorf("failed to build a map from frame: %w", err)
	}

	revealVC, err := vc.GenerateBBSSelectiveDisclosure(
		revealDoc,
		[]byte(request.Options.Nonce),
		ariesverifiable.WithPublicKeyFetcher(ariesverifiable.NewDIDKeyResolver(o.ctx.VDRegistry()).PublicKeyFetcher()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate BBS selective disclosure: %w", err)
	}

	return revealVC, nil
}

func (o *command) ProvePresentation(request *ProvePresentationRequest) (*ariesverifiable.Presentation, error) {
	presentation, err := ariesverifiable.ParsePresentation(request.Presentation, ariesverifiable.WithPresDisabledProofCheck())
	if err != nil {
		return nil, fmt.Errorf("failed to parse presentation: %w", err)
	}

	var didDoc *did.Doc
	doc, err := o.ctx.VDRegistry().Resolve(presentation.Holder)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve DID document of %v: %w", presentation.Holder, err)
	}
	didDoc = doc.DIDDocument

	if err := o.addProof(presentation, didDoc, request.Options); err != nil {
		return nil, fmt.Errorf("failed to add presentation proof: %w", err)
	}

	return presentation, nil
}

func (o *command) VerifyPresentation(request *VerifyPresentationRequest) error {
	_, err := ariesverifiable.ParsePresentation(
		request.VerifiablePresentation,
		ariesverifiable.WithPresPublicKeyFetcher(ariesverifiable.NewDIDKeyResolver(o.ctx.VDRegistry()).PublicKeyFetcher()),
	)
	if err != nil {
		return fmt.Errorf("failed to parse(verify) presentation: %w", err)
	}
	return nil
}

func (o *command) addLinkedDataProof(p provable, opts *ProofOptions) error {
	s, err := newKMSSigner(o.ctx.KMS(), o.ctx.Crypto(), getKID(opts))
	if err != nil {
		return err
	}

	var signatureSuite verifiablesigner.SignatureSuite

	switch opts.SignatureType {
	case Ed25519Signature2018:
		signatureSuite = ed25519signature2018.New(suite.WithSigner(s))
	case JSONWebSignature2020:
		signatureSuite = jsonwebsignature2020.New(suite.WithSigner(s))
	case BbsBlsSignature2020:
		s.bbs = true
		signatureSuite = bbsblssignature2020.New(suite.WithSigner(s))
	default:
		return fmt.Errorf("signature type unsupported %s", opts.SignatureType)
	}

	signatureRepresentation := ariesverifiable.SignatureProofValue

	if opts.SignatureRepresentation == nil {
		opts.SignatureRepresentation = &signatureRepresentation
	}

	signingCtx := &ariesverifiable.LinkedDataProofContext{
		VerificationMethod:      opts.VerificationMethod,
		SignatureRepresentation: *opts.SignatureRepresentation,
		SignatureType:           opts.SignatureType,
		Suite:                   signatureSuite,
		Created:                 opts.Created,
		Domain:                  opts.Domain,
		Challenge:               opts.Challenge,
		Purpose:                 opts.proofPurpose,
	}

	bbsLoader, err := bbsJSONLDDocumentLoader()
	if err != nil {
		return err
	}

	err = p.AddLinkedDataProof(signingCtx, jsonld.WithDocumentLoader(bbsLoader))
	if err != nil {
		return fmt.Errorf("failed to add linked data proof: %w", err)
	}

	return nil
}

func (o *command) getCredentialOpts(disableProofCheck bool) []ariesverifiable.CredentialOpt {
	if disableProofCheck {
		return []ariesverifiable.CredentialOpt{ariesverifiable.WithDisabledProofCheck()}
	}

	return []ariesverifiable.CredentialOpt{ariesverifiable.WithPublicKeyFetcher(
		ariesverifiable.NewDIDKeyResolver(o.ctx.VDRegistry()).PublicKeyFetcher(),
	)}
}

func prepareOpts(opts *ProofOptions, didDoc *did.Doc, method did.VerificationRelationship) (*ProofOptions, error) {
	if opts == nil {
		opts = &ProofOptions{}
	}

	var err error

	opts.proofPurpose, err = getProofPurpose(method)
	if err != nil {
		return nil, err
	}

	vMs := didDoc.VerificationMethods(method)[method]

	vmMatched := opts.VerificationMethod == ""

	for _, vm := range vMs {
		if opts.VerificationMethod != "" {
			// if verification method is provided as an option, then validate if it belongs to given method
			if opts.VerificationMethod == vm.VerificationMethod.ID {
				vmMatched = true

				break
			}

			continue
		} else {
			// by default first authentication public key
			opts.VerificationMethod = vm.VerificationMethod.ID

			break
		}
	}

	if !vmMatched {
		return nil, fmt.Errorf("unable to find matching '%s' key IDs for given verification method", opts.proofPurpose)
	}

	// this is the fallback logic kept for DIDs not having authentication method
	if opts.VerificationMethod == "" {
		log.Warnf("Could not find matching verification method for '%s' proof purpose", opts.proofPurpose)

		defaultVM, err := getDefaultVerificationMethod(didDoc)
		if err != nil {
			return nil, fmt.Errorf("failed to get default verification method: %w", err)
		}

		opts.VerificationMethod = defaultVM
	}

	return opts, nil
}

func getDefaultVerificationMethod(didDoc *did.Doc) (string, error) {
	switch {
	case len(didDoc.VerificationMethod) > 0:
		var publicKeyID string

		for _, k := range didDoc.VerificationMethod {
			if strings.HasPrefix(k.Type, Ed25519VerificationKey) {
				publicKeyID = k.ID

				break
			}
		}

		// if there isn't any ed25519 key then pick first one
		if publicKeyID == "" {
			publicKeyID = didDoc.VerificationMethod[0].ID
		}

		if !isDID(publicKeyID) {
			return didDoc.ID + publicKeyID, nil
		}

		return publicKeyID, nil
	case len(didDoc.Authentication) > 0:
		return didDoc.Authentication[0].VerificationMethod.ID, nil
	default:
		return "", errors.New("public key not found in DID Document")
	}
}

func (o *command) addProof(p provable, didDoc *did.Doc, opts *ProofOptions) error {
	var err error

	opts, err = prepareOpts(opts, didDoc, did.AssertionMethod)
	if err != nil {
		return err
	}

	return o.addLinkedDataProof(p, opts)
}

func isDID(str string) bool {
	return strings.HasPrefix(str, "did:")
}

func getProofPurpose(method did.VerificationRelationship) (string, error) {
	if method != did.Authentication && method != did.AssertionMethod {
		return "", fmt.Errorf("unsupported proof purpose, only authentication or assertionMethod are supported")
	}

	if method == did.Authentication {
		return "authentication", nil
	}

	return "assertionMethod", nil
}

func bbsJSONLDDocumentLoader() (*ld.CachingDocumentLoader, error) {
	loader := presexch.CachingJSONLDLoader()

	reader, err := ld.DocumentFromReader(strings.NewReader(contextBBSContent))
	if err != nil {
		return nil, err
	}

	loader.AddDocument(bbsContext, reader)

	return loader, nil
}

const contextBBSContent = `{
  "@context": {
    "@version": 1.1,
    "id": "@id",
    "type": "@type",
    "ldssk": "https://w3id.org/security#",
    "BbsBlsSignature2020": {
      "@id": "https://w3id.org/security#BbsBlsSignature2020",
      "@context": {
        "@version": 1.1,
        "@protected": true,
        "id": "@id",
        "type": "@type",
        "sec": "https://w3id.org/security#",
        "xsd": "http://www.w3.org/2001/XMLSchema#",
        "challenge": "sec:challenge",
        "created": {
          "@id": "http://purl.org/dc/terms/created",
          "@type": "xsd:dateTime"
        },
        "domain": "sec:domain",
        "proofValue": "sec:proofValue",
        "nonce": "sec:nonce",
        "proofPurpose": {
          "@id": "sec:proofPurpose",
          "@type": "@vocab",
          "@context": {
            "@version": 1.1,
            "@protected": true,
            "id": "@id",
            "type": "@type",
            "sec": "https://w3id.org/security#",
            "assertionMethod": {
              "@id": "sec:assertionMethod",
              "@type": "@id",
              "@container": "@set"
            },
            "authentication": {
              "@id": "sec:authenticationMethod",
              "@type": "@id",
              "@container": "@set"
            }
          }
        },
        "verificationMethod": {
          "@id": "sec:verificationMethod",
          "@type": "@id"
        }
      }
    },
    "BbsBlsSignatureProof2020": {
      "@id": "https://w3id.org/security#BbsBlsSignatureProof2020",
      "@context": {
        "@version": 1.1,
        "@protected": true,
        "id": "@id",
        "type": "@type",
        "sec": "https://w3id.org/security#",
        "xsd": "http://www.w3.org/2001/XMLSchema#",
        "challenge": "sec:challenge",
        "created": {
          "@id": "http://purl.org/dc/terms/created",
          "@type": "xsd:dateTime"
        },
        "domain": "sec:domain",
        "nonce": "sec:nonce",
        "proofPurpose": {
          "@id": "sec:proofPurpose",
          "@type": "@vocab",
          "@context": {
            "@version": 1.1,
            "@protected": true,
            "id": "@id",
            "type": "@type",
            "sec": "https://w3id.org/security#",
            "assertionMethod": {
              "@id": "sec:assertionMethod",
              "@type": "@id",
              "@container": "@set"
            },
            "authentication": {
              "@id": "sec:authenticationMethod",
              "@type": "@id",
              "@container": "@set"
            }
          }
        },
        "proofValue": "sec:proofValue",
        "verificationMethod": {
          "@id": "sec:verificationMethod",
          "@type": "@id"
        }
      }
    },
    "Bls12381G2Key2020": "ldssk:Bls12381G2Key2020"
  }
}`

func toMap(v interface{}) (map[string]interface{}, error) {
	var (
		b   []byte
		err error
	)

	switch cv := v.(type) {
	case []byte:
		b = cv
	case string:
		b = []byte(cv)
	default:
		b, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
	}

	var m map[string]interface{}

	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
