package panacea

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	vdrapi "github.com/hyperledger/aries-framework-go/pkg/framework/aries/api/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
)

type VDR struct {
	endpoint string
}

func New(endpoint string) (*VDR, error) {
	return &VDR{endpoint: endpoint}, nil
}

func (v *VDR) Close() error {
	return nil
}

func (v *VDR) Accept(method string) bool {
	return method == "panacea"
}

func (v *VDR) Create(keyManager kms.KeyManager, didDoc *did.Doc, opts ...vdrapi.DIDMethodOption) (*did.DocResolution, error) {
	return nil, fmt.Errorf("not supported")
}

func (v *VDR) Update(didDoc *did.Doc, opts ...vdrapi.DIDMethodOption) error {
	return fmt.Errorf("not supported")
}

// Deactivate did doc.
func (v *VDR) Deactivate(didID string, opts ...vdrapi.DIDMethodOption) error {
	return fmt.Errorf("not supported")
}

type response struct {
	Height string          `json:"height"`
	Result *responseResult `json:"result"`
}

type responseResult struct {
	Document json.RawMessage `json:"document"`
	Sequence string          `json:"sequence"`
}

// Read implements didresolver.DidMethod.Read interface (https://w3c-ccg.github.io/did-resolution/#resolving-input)
func (v *VDR) Read(didID string, _ ...vdrapi.ResolveOption) (*did.DocResolution, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/did/%s", v.endpoint, didID))
	if err != nil {
		return nil, fmt.Errorf("failed to get did from panacea: %w", err)
	}
	defer resp.Body.Close()

	var res response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	didDoc, err := did.ParseDocument(res.Result.Document)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DID document: %w", err)
	}

	return &did.DocResolution{DIDDocument: didDoc}, nil
}

//
//var didDocMap = map[string]string{
//	"did:panacea:7Prd74ry1Uct87nZqL3ny7aR7Cg46JamVbJgk8azVgUm": issuerDIDDoc,
//	"did:panacea:aR7Cg46JamVbJgk8azVgUm7Prd74ry1Uct87nZqL3ny7": holderDIDDoc,
//}
//
//const issuerDIDDoc = `
//{
//    "@context": "https://www.w3.org/ns/did/v1",
//    "id": "did:panacea:7Prd74ry1Uct87nZqL3ny7aR7Cg46JamVbJgk8azVgUm",
//    "verificationMethod": [
//        {
//            "id": "did:panacea:7Prd74ry1Uct87nZqL3ny7aR7Cg46JamVbJgk8azVgUm#key1",
//            "type": "Bls12381G2Key2020",
//            "controller": "did:panacea:7Prd74ry1Uct87nZqL3ny7aR7Cg46JamVbJgk8azVgUm",
//            "publicKeyBase58": "22TdHn1mf2eR5CBab7ZjZbZfgVeoCyPZ8mzzzxgek19R5G5LLToDMYt3CAfBJGLgS1oNoGgzJB8DGhhViXux7fvdRGSvjapwaFQtkKQKCN26XtNJvSyYQ3vENYbU5bti23eF"
//        }
//    ],
//    "authentication": [
//        "did:panacea:7Prd74ry1Uct87nZqL3ny7aR7Cg46JamVbJgk8azVgUm#key1"
//    ],
//	"assertionMethod": [
//        "did:panacea:7Prd74ry1Uct87nZqL3ny7aR7Cg46JamVbJgk8azVgUm#key1"
//    ]
//}
//`
//
//const holderDIDDoc = `
//{
//    "@context": "https://www.w3.org/ns/did/v1",
//    "id": "did:panacea:aR7Cg46JamVbJgk8azVgUm7Prd74ry1Uct87nZqL3ny7",
//    "verificationMethod": [
//        {
//            "id": "did:panacea:aR7Cg46JamVbJgk8azVgUm7Prd74ry1Uct87nZqL3ny7#key1",
//            "type": "Bls12381G2Key2020",
//            "controller": "did:panacea:aR7Cg46JamVbJgk8azVgUm7Prd74ry1Uct87nZqL3ny7",
//            "publicKeyBase58": "23yHYqYHU7n9jY3Fsaf61yZp7czsQSJMyT3ZB89jMBvdF5hAvB33YPVeaddEzPwbySBmPEzqHqx6w5rS6GQAcLaVuhQXukavEWVvEpftLBpn72kapGqGXDEcEjqZZiMSHBxd"
//        }
//    ],
//    "authentication": [
//        "did:panacea:aR7Cg46JamVbJgk8azVgUm7Prd74ry1Uct87nZqL3ny7#key1"
//    ],
//	"assertionMethod": [
//        "did:panacea:aR7Cg46JamVbJgk8azVgUm7Prd74ry1Uct87nZqL3ny7#key1"
//    ]
//}
//`
