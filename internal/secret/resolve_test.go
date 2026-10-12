package secret

import (
	"testing"

	"github.com/zx06/xurl/internal/errors"
)

type mockKeyring struct {
	secrets map[string]string
}

func (m *mockKeyring) key(service, account string) string {
	return service + "/" + account
}

func (m *mockKeyring) Get(service, account string) (string, error) {
	v, ok := m.secrets[m.key(service, account)]
	if !ok {
		return "", errors.New(errors.CodeSecretNotFound, "not found", nil)
	}
	return v, nil
}

func (m *mockKeyring) Set(service, account, value string) error {
	m.secrets[m.key(service, account)] = value
	return nil
}

func (m *mockKeyring) Delete(service, account string) error {
	delete(m.secrets, m.key(service, account))
	return nil
}

func TestResolveKeyringServiceUser(t *testing.T) {
	kr := &mockKeyring{secrets: map[string]string{"myapp/token": "secret-token"}}
	val, xe := Resolve("keyring:myapp/token", Options{Keyring: kr})
	if xe != nil {
		t.Fatal(xe)
	}
	if val != "secret-token" {
		t.Fatalf("got %q", val)
	}
}

func TestResolvePlaintextDenied(t *testing.T) {
	_, xe := Resolve("plain-secret", Options{AllowPlaintext: false})
	if xe == nil || xe.Code != errors.CodeCfgInvalid {
		t.Fatalf("expected CFG_INVALID, got %v", xe)
	}
}

func TestResolvePlaintextAllowed(t *testing.T) {
	val, xe := Resolve("plain-secret", Options{AllowPlaintext: true})
	if xe != nil {
		t.Fatal(xe)
	}
	if val != "plain-secret" {
		t.Fatalf("got %q", val)
	}
}
