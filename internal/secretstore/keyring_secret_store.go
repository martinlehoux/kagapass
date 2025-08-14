package secretstore

import (
	"fmt"

	"github.com/99designs/keyring"
)

const (
	serviceName = "kagapass"
	keyPrefix   = "keepass_"
)

type keyringSecretStore struct {
	ring keyring.Keyring
}

var _ SecretStore = (*keyringSecretStore)(nil)

func NewKeyring() (keyringSecretStore, error) {
	config := keyring.Config{ //nolint:exhaustruct // Configuration is too wide
		ServiceName:              serviceName,
		AllowedBackends:          []keyring.BackendType{keyring.SecretServiceBackend},
		KeychainTrustApplication: true,
	}
	ring, err := keyring.Open(config)

	return keyringSecretStore{ring: ring}, err
}

func getKey(key string) string {
	return keyPrefix + key
}

func (k keyringSecretStore) Store(key string, secret []byte) error {
	item := keyring.Item{
		Key:                         getKey(key),
		Data:                        secret,
		Label:                       fmt.Sprintf("%s - %s", serviceName, key),
		Description:                 "",
		KeychainNotTrustApplication: false,
		KeychainNotSynchronizable:   false,
	}

	return k.ring.Set(item)
}

func (k keyringSecretStore) Get(key string) ([]byte, error) {
	item, err := k.ring.Get(getKey(key))
	if err != nil {
		return nil, err
	}

	return item.Data, nil
}

func (k keyringSecretStore) Remove(key string) error {
	return k.ring.Remove(getKey(key))
}
