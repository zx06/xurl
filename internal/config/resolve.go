package config

import (
	"os"
	"path/filepath"

	"github.com/zx06/xurl/internal/errors"
)

func Resolve(opts Options) (Resolved, *errors.XError) {
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

	var cfg File
	var cfgPath string
	if opts.ConfigPath != "" {
		abs := opts.ConfigPath
		if !filepath.IsAbs(abs) {
			abs = filepath.Join(workDir, abs)
		}
		f, xe := readFile(abs)
		if xe != nil {
			return Resolved{}, xe
		}
		cfg = f
		cfgPath = abs
	} else {
		for _, p := range defaultConfigPaths(workDir, opts.HomeDir) {
			f, xe := readFile(p)
			if xe != nil {
				if xe.Code == errors.CodeCfgNotFound {
					continue
				}
				return Resolved{}, xe
			}
			cfg = f
			cfgPath = p
			break
		}
	}

	profile := ""
	if opts.CLIProfileSet {
		profile = opts.CLIProfile
	} else if opts.EnvProfile != "" {
		profile = opts.EnvProfile
	} else if _, ok := cfg.Profiles["default"]; ok {
		profile = "default"
	}

	resolvedProfiles := make(map[string]Profile, len(cfg.Profiles))
	for name, p := range cfg.Profiles {
		if p.Headers == nil {
			p.Headers = map[string]string{}
		}
		if p.Vars == nil {
			p.Vars = map[string]string{}
		}
		resolvedProfiles[name] = p
	}

	var selected Profile
	if profile != "" {
		p, ok := resolvedProfiles[profile]
		if !ok {
			return Resolved{}, errors.New(errors.CodeCfgInvalid, "profile not found", map[string]any{"profile": profile})
		}
		selected = p
	}

	format := "json"
	if opts.CLIFormatSet {
		format = opts.CLIFormat
	} else if opts.EnvFormat != "" {
		format = opts.EnvFormat
	} else if selected.Format != "" {
		format = selected.Format
	} else if cfg.Format != "" {
		format = cfg.Format
	}

	libDir := resolveLibraryDir(cfg, opts.HomeDir, workDir)
	histPath := resolveHistoryPath(opts.HomeDir)

	return Resolved{
		ConfigPath:  cfgPath,
		ProfileName: profile,
		Format:      format,
		LibraryDir:  libDir,
		HistoryPath: histPath,
		Profile:     selected,
		AllProfiles: resolvedProfiles,
	}, nil
}
