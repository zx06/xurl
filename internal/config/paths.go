package config

import (
	"path/filepath"

	"github.com/zx06/xurl/internal/paths"
)

func resolveLibraryDir(cfg File, homeDir, workDir string) string {
	raw := cfg.LibraryDir
	if raw == "" {
		return paths.DefaultLibraryDir(homeDir)
	}
	raw = paths.ExpandHome(raw)
	if !filepath.IsAbs(raw) && workDir != "" {
		raw = filepath.Join(workDir, raw)
	}
	abs, err := filepath.Abs(raw)
	if err != nil {
		return raw
	}
	return abs
}

func resolveHistoryPath(homeDir string) string {
	return paths.DefaultHistoryPath(homeDir)
}

// ProfileByName returns a profile from resolved map or error.
func ProfileByName(all map[string]Profile, name string) (Profile, bool) {
	p, ok := all[name]
	return p, ok
}
