package paths

import (
	"os"
	"path/filepath"
	"strings"
)

func HomeDir() string {
	if h, err := os.UserHomeDir(); err == nil {
		return h
	}
	return ""
}

func ExpandHome(path string) string {
	if path == "" {
		return path
	}
	home := HomeDir()
	if home == "" {
		return path
	}
	if path == "~" {
		return home
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, strings.TrimPrefix(path, "~/"))
	}
	return path
}

func DefaultLibraryDir(homeDir string) string {
	if homeDir == "" {
		homeDir = HomeDir()
	}
	return filepath.Join(homeDir, ".config", "xurl", "library")
}

func DefaultHistoryPath(homeDir string) string {
	if homeDir == "" {
		homeDir = HomeDir()
	}
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = filepath.Join(homeDir, ".local", "share")
	}
	return filepath.Join(dataHome, "xurl", "history.json")
}

func AbsPath(workDir, p string) (string, error) {
	p = ExpandHome(p)
	if !filepath.IsAbs(p) && workDir != "" {
		p = filepath.Join(workDir, p)
	}
	return filepath.Abs(p)
}
