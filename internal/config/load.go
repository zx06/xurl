package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/zx06/xurl/internal/errors"
)

func defaultConfigPaths(workDir, homeDir string) []string {
	paths := make([]string, 0, 2)
	if workDir != "" {
		paths = append(paths, filepath.Join(workDir, "xurl.yaml"))
	}
	if homeDir != "" {
		paths = append(paths, filepath.Join(homeDir, ".config", "xurl", "xurl.yaml"))
	}
	return paths
}

func readFile(path string) (File, *errors.XError) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return File{}, errors.New(errors.CodeCfgNotFound, "config file not found", map[string]any{"path": path})
		}
		return File{}, errors.Wrap(errors.CodeCfgInvalid, "failed to read config file", map[string]any{"path": path}, err)
	}
	var f File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return File{}, errors.Wrap(errors.CodeCfgInvalid, "invalid config file", map[string]any{"path": path}, err)
	}
	if f.Profiles == nil {
		f.Profiles = map[string]Profile{}
	}
	return f, nil
}

func FindConfigPath(opts Options) string {
	workDir := opts.WorkDir
	if workDir == "" {
		wd, _ := os.Getwd()
		workDir = wd
	}
	homeDir := opts.HomeDir
	if homeDir == "" {
		if hd, err := os.UserHomeDir(); err == nil {
			homeDir = hd
		}
	}
	if opts.ConfigPath != "" {
		abs := opts.ConfigPath
		if !filepath.IsAbs(abs) {
			abs = filepath.Join(workDir, abs)
		}
		return abs
	}
	for _, p := range defaultConfigPaths(workDir, homeDir) {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func LoadConfig(opts Options) (File, string, *errors.XError) {
	workDir := opts.WorkDir
	if workDir == "" {
		wd, _ := os.Getwd()
		workDir = wd
	}
	if opts.HomeDir == "" {
		if hd, err := os.UserHomeDir(); err == nil {
			opts.HomeDir = hd
		}
	}

	if opts.ConfigPath != "" {
		abs := opts.ConfigPath
		if !filepath.IsAbs(abs) {
			abs = filepath.Join(workDir, abs)
		}
		f, xe := readFile(abs)
		if xe != nil {
			return File{}, "", xe
		}
		return f, abs, nil
	}

	for _, p := range defaultConfigPaths(workDir, opts.HomeDir) {
		f, xe := readFile(p)
		if xe != nil {
			if xe.Code == errors.CodeCfgNotFound {
				continue
			}
			return File{}, "", xe
		}
		return f, p, nil
	}

	return File{Profiles: map[string]Profile{}}, "", nil
}
