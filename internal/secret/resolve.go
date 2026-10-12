package secret

import (
	"strings"

	"github.com/zx06/xurl/internal/errors"
)

const keyringPrefix = "keyring:"

const defaultService = "xurl"

type Options struct {
	AllowPlaintext bool
	Keyring        KeyringAPI
}

// parseKeyringRef parses keyring:service/user or keyring:account (default service xurl).
func parseKeyringRef(ref string) (service, account string, xe *errors.XError) {
	if ref == "" {
		return "", "", errors.New(errors.CodeCfgInvalid, "invalid keyring reference: empty path after keyring:", map[string]any{"ref": ref})
	}
	if i := strings.Index(ref, "/"); i >= 0 {
		service = ref[:i]
		account = ref[i+1:]
		if service == "" || account == "" {
			return "", "", errors.New(errors.CodeCfgInvalid, "invalid keyring reference: expected keyring:service/user", map[string]any{"ref": ref})
		}
		return service, account, nil
	}
	return defaultService, ref, nil
}

func Resolve(raw string, opts Options) (string, *errors.XError) {
	if strings.HasPrefix(raw, keyringPrefix) {
		ref := strings.TrimPrefix(raw, keyringPrefix)
		service, account, parseErr := parseKeyringRef(ref)
		if parseErr != nil {
			return "", parseErr
		}
		kr := opts.Keyring
		if kr == nil {
			kr = defaultKeyring()
		}
		val, err := kr.Get(service, account)
		if err != nil {
			return "", errors.Wrap(errors.CodeSecretNotFound, "failed to read secret from keyring",
				map[string]any{"service": service, "account": account}, err)
		}
		return val, nil
	}
	if opts.AllowPlaintext {
		return raw, nil
	}
	return "", errors.New(errors.CodeCfgInvalid,
		"plaintext secret not allowed; use keyring:service/user reference or set allow_plaintext: true on the profile",
		nil)
}

func IsKeyringRef(s string) bool {
	return strings.HasPrefix(s, keyringPrefix)
}

func LooksLikeSecret(s string) bool {
	if IsKeyringRef(s) {
		return true
	}
	// Heuristic: non-empty values in auth fields that aren't env refs
	return s != "" && !strings.HasPrefix(s, "{{")
}
