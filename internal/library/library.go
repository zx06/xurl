package library

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zx06/xurl/internal/errors"
	"github.com/zx06/xurl/internal/paths"
)

func EnsureDir(dir string) *errors.XError {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return errors.Wrap(errors.CodeCfgInvalid, "failed to create library directory", map[string]any{"path": dir}, err)
	}
	return nil
}

// ResolveTarget finds an HTTP file on disk or under libraryDir.
func ResolveTarget(workDir, libraryDir, target string) (string, *errors.XError) {
	target = strings.TrimSpace(target)
	if target == "" {
		return "", errors.New(errors.CodeHTTPFile, "empty target", nil)
	}

	candidates := []string{target}
	if !strings.HasSuffix(strings.ToLower(target), ".http") {
		candidates = append(candidates, target+".http")
	}

	for _, c := range candidates {
		if abs, ok := existingFile(workDir, c); ok {
			return abs, nil
		}
	}

	for _, c := range candidates {
		libPath := filepath.Join(libraryDir, c)
		if abs, ok := existingFile("", libPath); ok {
			return abs, nil
		}
	}

	return "", errors.New(errors.CodeHTTPFile, "HTTP file not found",
		map[string]any{"target": target, "library_dir": libraryDir})
}

func existingFile(workDir, p string) (string, bool) {
	abs, err := paths.AbsPath(workDir, p)
	if err != nil {
		return "", false
	}
	st, err := os.Stat(abs)
	if err != nil || st.IsDir() {
		return "", false
	}
	return abs, true
}

type FileInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func List(dir string) ([]FileInfo, *errors.XError) {
	if xe := EnsureDir(dir); xe != nil {
		return nil, xe
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(errors.CodeCfgInvalid, "failed to read library", map[string]any{"path": dir}, err)
	}
	out := make([]FileInfo, 0)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(e.Name()), ".http") {
			continue
		}
		full := filepath.Join(dir, e.Name())
		out = append(out, FileInfo{Name: e.Name(), Path: full})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func Add(srcPath, libraryDir string, symlink bool) (string, *errors.XError) {
	if xe := EnsureDir(libraryDir); xe != nil {
		return "", xe
	}
	srcAbs, err := filepath.Abs(srcPath)
	if err != nil {
		return "", errors.Wrap(errors.CodeHTTPFile, "invalid source path", map[string]any{"path": srcPath}, err)
	}
	if _, err := os.Stat(srcAbs); err != nil {
		return "", errors.Wrap(errors.CodeHTTPFile, "source file not found", map[string]any{"path": srcAbs}, err)
	}
	dest := filepath.Join(libraryDir, filepath.Base(srcAbs))
	if _, err := os.Stat(dest); err == nil {
		return "", errors.New(errors.CodeCfgInvalid, "library file already exists", map[string]any{"path": dest})
	}
	if symlink {
		if err := os.Symlink(srcAbs, dest); err != nil {
			return "", errors.Wrap(errors.CodeCfgInvalid, "failed to create symlink", map[string]any{"dest": dest}, err)
		}
		return dest, nil
	}
	in, err := os.Open(srcAbs)
	if err != nil {
		return "", errors.Wrap(errors.CodeHTTPFile, "failed to open source", nil, err)
	}
	defer in.Close()
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o644)
	if err != nil {
		return "", errors.Wrap(errors.CodeCfgInvalid, "failed to create library file", map[string]any{"dest": dest}, err)
	}
	_, cpErr := io.Copy(out, in)
	closeErr := out.Close()
	if cpErr != nil {
		_ = os.Remove(dest)
		return "", errors.Wrap(errors.CodeCfgInvalid, "failed to copy file", nil, cpErr)
	}
	if closeErr != nil {
		return "", errors.Wrap(errors.CodeCfgInvalid, "failed to close library file", nil, closeErr)
	}
	return dest, nil
}
