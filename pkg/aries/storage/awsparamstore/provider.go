package awsparamstore

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/hyperledger/aries-framework-go/spi/storage"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) OpenStore(name string) (storage.Store, error) {
	store, err := newStore(name)
	if err != nil {
		return nil, fmt.Errorf("failed to open awsparamstore: %w", err)
	}
	return store, nil
}

func (p *Provider) SetStoreConfig(name string, config storage.StoreConfiguration) error {
	return errors.New("not supported")
}

func (p *Provider) GetStoreConfig(name string) (storage.StoreConfiguration, error) {
	return storage.StoreConfiguration{}, errors.New("not supported")
}

func (p *Provider) GetOpenStores() []storage.Store {
	log.Panic("not implemented")
	return nil
}

func (p *Provider) Close() error {
	return nil
}
