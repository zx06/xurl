package config

import (
	"strings"

	"github.com/zx06/xurl/internal/secret"
)

func MaskProfile(name string, p Profile) ProfileDetail {
	d := ProfileDetail{
		Name:           name,
		Description:    p.Description,
		BaseURL:        p.BaseURL,
		Vars:           copyMap(p.Vars),
		Headers:        maskMap(p.Headers),
		AllowPlaintext: p.AllowPlaintext,
	}
	if p.Auth.Bearer != "" {
		d.Auth.Bearer = maskSecret(p.Auth.Bearer)
	}
	if p.Auth.Basic != nil {
		d.Auth.Basic = &struct {
			Username string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
		}{
			Username: p.Auth.Basic.Username,
			Password: maskSecret(p.Auth.Basic.Password),
		}
	}
	if p.Auth.Header != nil {
		d.Auth.Header = &struct {
			Name  string `json:"name,omitempty"`
			Value string `json:"value,omitempty"`
		}{
			Name:  p.Auth.Header.Name,
			Value: maskSecret(p.Auth.Header.Value),
		}
	}
	return d
}

func maskMap(m map[string]string) map[string]string {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = maskSecret(v)
	}
	return out
}

func maskSecret(s string) string {
	if s == "" {
		return ""
	}
	if secret.IsKeyringRef(s) {
		return s
	}
	if strings.HasPrefix(s, "{{") {
		return s
	}
	return "***"
}

func copyMap(m map[string]string) map[string]string {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func ListProfiles(profiles map[string]Profile) []ProfileInfo {
	names := make([]string, 0, len(profiles))
	for n := range profiles {
		names = append(names, n)
	}
	// simple sort
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			if names[j] < names[i] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}
	out := make([]ProfileInfo, 0, len(names))
	for _, n := range names {
		p := profiles[n]
		out = append(out, ProfileInfo{Name: n, Description: p.Description, BaseURL: p.BaseURL})
	}
	return out
}
