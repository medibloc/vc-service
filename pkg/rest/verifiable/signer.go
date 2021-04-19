package verifiable

import (
	"strings"

	ariescrypto "github.com/hyperledger/aries-framework-go/pkg/crypto"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
)

type kmsSigner struct {
	keyHandle interface{}
	crypto    ariescrypto.Crypto
	bbs       bool
}

func getKID(opts *ProofOptions) string {
	idSplit := strings.Split(opts.VerificationMethod, "#")
	if len(idSplit) == creatorParts {
		return idSplit[1]
	}

	return ""
}

func newKMSSigner(keyManager kms.KeyManager, c ariescrypto.Crypto, kid string) (*kmsSigner, error) {
	keyHandler, err := keyManager.Get(kid)
	if err != nil {
		return nil, err
	}

	return &kmsSigner{keyHandle: keyHandler, crypto: c}, nil
}

func (s *kmsSigner) textToLines(txt string) [][]byte {
	lines := strings.Split(txt, "\n")
	linesBytes := make([][]byte, 0, len(lines))

	for i := range lines {
		if strings.TrimSpace(lines[i]) != "" {
			linesBytes = append(linesBytes, []byte(lines[i]))
		}
	}

	return linesBytes
}

func (s *kmsSigner) Sign(data []byte) ([]byte, error) {
	if s.bbs {
		return s.crypto.SignMulti(s.textToLines(string(data)), s.keyHandle)
	}

	v, err := s.crypto.Sign(data, s.keyHandle)
	if err != nil {
		return nil, err
	}

	return v, nil
}
