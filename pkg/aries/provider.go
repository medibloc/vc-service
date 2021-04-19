package aries

import (
	"github.com/hyperledger/aries-framework-go/component/storage/leveldb"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/hyperledger/aries-framework-go/pkg/framework/context"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	"github.com/hyperledger/aries-framework-go/pkg/kms/localkms"
	"github.com/hyperledger/aries-framework-go/pkg/secretlock"
	"github.com/hyperledger/aries-framework-go/pkg/secretlock/noop"
	"github.com/hyperledger/aries-framework-go/spi/storage"
	"github.com/medibloc/verifiable/pkg/aries/storage/awsparamstore"
	"github.com/medibloc/verifiable/pkg/aries/vdr/panacea"
)

func NewProvider() (*context.Provider, error) {
	panaceaVDR, err := panacea.New() // TODO: use a DID universal resolver (using httpbinding.VDR of Aries)
	if err != nil {
		return nil, err
	}

	//TODO: Use AWS Parameter Store as a KMS
	framework, err := aries.New(aries.WithStoreProvider(getStorageProvider()), aries.WithVDR(panaceaVDR), aries.WithKMS(getKMSCreator()))
	if err != nil {
		return nil, err
	}

	ctx, err := framework.Context()
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func getStorageProvider() storage.Provider {
	return leveldb.NewProvider("") //TODO: path (prefix)
}

type kmsProvider struct {
	storageProvider storage.Provider
	secretLock      secretlock.Service
}

func (k kmsProvider) StorageProvider() storage.Provider {
	return k.storageProvider
}

func (k kmsProvider) SecretLock() secretlock.Service {
	return k.secretLock
}

func getKMSCreator() func(provider kms.Provider) (kms.KeyManager, error) {
	return func(provider kms.Provider) (kms.KeyManager, error) {
		return localkms.New("local-lock://custom/master/key/", &kmsProvider{awsparamstore.NewProvider(), &noop.NoLock{}})
	}
}
