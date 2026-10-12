//go:build !windows

package secret

import (
	"github.com/zalando/go-keyring"
)

func (k *osKeyring) Get(service, account string) (string, error) {
	return keyring.Get(service, account)
}

func (k *osKeyring) Set(service, account, value string) error {
	return keyring.Set(service, account, value)
}

func (k *osKeyring) Delete(service, account string) error {
	return keyring.Delete(service, account)
}
