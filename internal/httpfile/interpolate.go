package httpfile

import (
	"os"
	"regexp"
	"strings"
)

var varPattern = regexp.MustCompile(`\{\{([^}]+)\}\}`)

func Interpolate(s string, vars map[string]string) string {
	return varPattern.ReplaceAllStringFunc(s, func(match string) string {
		key := strings.TrimSpace(varPattern.FindStringSubmatch(match)[1])
		if v, ok := vars[key]; ok {
			return v
		}
		if env := os.Getenv(key); env != "" {
			return env
		}
		return match
	})
}

func BuildVars(baseURL string, profileVars map[string]string, fileVars map[string]string, attrs map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range profileVars {
		out[k] = v
	}
	if baseURL != "" {
		out["base_url"] = baseURL
		out["host"] = baseURL
	}
	for k, v := range fileVars {
		out[k] = v
	}
	for k, v := range attrs {
		out[k] = v
	}
	return out
}
