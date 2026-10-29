package main

import (
	"os"

	"github.com/zx06/xurl/internal/config"
	"github.com/zx06/xurl/internal/errors"
	"github.com/zx06/xurl/internal/history"
	"github.com/zx06/xurl/internal/httpfile"
	"github.com/zx06/xurl/internal/library"
	"github.com/zx06/xurl/internal/output"
	"github.com/zx06/xurl/internal/runner"
)

type runParams struct {
	Target      string
	RequestName string
	ProfileName string
	Profile     config.Profile
	RecordHist  bool
}

func executeRun(w output.Writer, format output.Format, p runParams) error {
	workDir, _ := os.Getwd()
	libDir := global.Resolved.LibraryDir
	if xe := library.EnsureDir(libDir); xe != nil {
		return xe
	}

	path, xe := library.ResolveTarget(workDir, libDir, p.Target)
	if xe != nil {
		return xe
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrap(errors.CodeHTTPFile, "failed to read HTTP file", map[string]any{"path": path}, err)
	}

	requests, xe := httpfile.Parse(b)
	if xe != nil {
		return xe
	}

	req, ok := httpfile.FindByName(requests, p.RequestName)
	if !ok {
		if p.RequestName == "" {
			return errors.New(errors.CodeHTTPFile, "file contains multiple requests; specify a request name", map[string]any{"count": len(requests)})
		}
		return errors.New(errors.CodeHTTPFile, "request not found", map[string]any{"name": p.RequestName})
	}

	profile := p.Profile
	profileName := p.ProfileName
	if profileName == "" {
		profileName = global.ProfileStr
	}

	result := runner.Run(req, runner.Options{
		Profile: profile,
		Attrs:   global.Attrs,
	})

	if p.RecordHist {
		store := history.NewStore(global.Resolved.HistoryPath)
		_ = store.Append(history.Entry{
			Path:        path,
			RequestName: p.RequestName,
			Profile:     profileName,
			OK:          result.OK && result.Error == nil,
		})
	}

	return w.WriteRunResult(format, result)
}

func profileForRun(name string) (config.Profile, *errors.XError) {
	if name == "" {
		return config.Profile{}, errors.New(errors.CodeCfgInvalid, "profile is required; use -p/--profile or XURL_PROFILE", nil)
	}
	if global.Resolved.ProfileName == name && global.Resolved.ProfileName != "" {
		return global.Resolved.Profile, nil
	}
	p, ok := global.Resolved.AllProfiles[name]
	if !ok {
		return config.Profile{}, errors.New(errors.CodeCfgInvalid, "profile not found", map[string]any{"profile": name})
	}
	return p, nil
}
