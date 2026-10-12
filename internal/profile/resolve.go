package profile

import (
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/zx06/xurl/internal/config"
	"github.com/zx06/xurl/internal/errors"
	"github.com/zx06/xurl/internal/secret"
)

type Resolved struct {
	BaseURL string
	Headers map[string]string
	Auth    config.Auth
}

func ResolveSecrets(p config.Profile, kr secret.KeyringAPI) (Resolved, *errors.XError) {
	opts := secret.Options{AllowPlaintext: p.AllowPlaintext, Keyring: kr}
	headers := map[string]string{}
	for k, v := range p.Headers {
		resolved, xe := resolveMaybeSecret(v, opts, false)
		if xe != nil {
			return Resolved{}, xe
		}
		headers[k] = resolved
	}

	auth := config.Auth{}
	if p.Auth.Bearer != "" {
		tok, xe := resolveMaybeSecret(p.Auth.Bearer, opts, true)
		if xe != nil {
			return Resolved{}, xe
		}
		auth.Bearer = tok
	}
	if p.Auth.Basic != nil {
		user := p.Auth.Basic.Username
		pass, xe := resolveMaybeSecret(p.Auth.Basic.Password, opts, true)
		if xe != nil {
			return Resolved{}, xe
		}
		auth.Basic = &config.BasicAuth{Username: user, Password: pass}
	}
	if p.Auth.Header != nil {
		val, xe := resolveMaybeSecret(p.Auth.Header.Value, opts, true)
		if xe != nil {
			return Resolved{}, xe
		}
		auth.Header = &config.HeaderKey{Name: p.Auth.Header.Name, Value: val}
	}

	return Resolved{
		BaseURL: p.BaseURL,
		Headers: headers,
		Auth:    auth,
	}, nil
}

func resolveMaybeSecret(raw string, opts secret.Options, isSecretField bool) (string, *errors.XError) {
	if raw == "" {
		return "", nil
	}
	if secret.IsKeyringRef(raw) {
		return secret.Resolve(raw, opts)
	}
	if strings.Contains(raw, "{{") {
		// interpolation happens later; not a stored secret
		return raw, nil
	}
	if isSecretField && !opts.AllowPlaintext && secret.LooksLikeSecret(raw) {
		return "", errors.New(errors.CodeCfgInvalid,
			"plaintext secret not allowed; use keyring:service/user reference or set allow_plaintext: true on the profile",
			map[string]any{"value_hint": "***"})
	}
	if isSecretField {
		return secret.Resolve(raw, opts)
	}
	return raw, nil
}

func ResolveURL(baseURL, reqURL string) (string, *errors.XError) {
	reqURL = strings.TrimSpace(reqURL)
	if reqURL == "" {
		return "", errors.New(errors.CodeHTTPFile, "empty URL", nil)
	}
	if strings.HasPrefix(reqURL, "http://") || strings.HasPrefix(reqURL, "https://") {
		return reqURL, nil
	}
	if baseURL == "" {
		return reqURL, nil
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", errors.Wrap(errors.CodeCfgInvalid, "invalid base_url", map[string]any{"base_url": baseURL}, err)
	}
	ref, err := url.Parse(reqURL)
	if err != nil {
		return "", errors.Wrap(errors.CodeHTTPFile, "invalid request URL", map[string]any{"url": reqURL}, err)
	}
	return base.ResolveReference(ref).String(), nil
}

func ApplyAuth(headers map[string]string, auth config.Auth) {
	if auth.Bearer != "" {
		headers["Authorization"] = "Bearer " + auth.Bearer
	}
	if auth.Basic != nil {
		cred := auth.Basic.Username + ":" + auth.Basic.Password
		headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(cred))
	}
	if auth.Header != nil && auth.Header.Name != "" {
		headers[auth.Header.Name] = auth.Header.Value
	}
}
