package panacea

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	vdr, err := New("https://testnet-api.gopanacea.org")
	require.NoError(t, err)
	require.NotNil(t, vdr)

	docResolution, err := vdr.Read("did:panacea:E8wdMGreu9PPtaXJAZLYsFgGWzmH1yHvy4eh2nfSxzsE")
	require.NoError(t, err)
	t.Log(docResolution.DIDDocument)
}
