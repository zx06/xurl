package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/zx06/xurl/internal/errors"
	"github.com/zx06/xurl/internal/paths"
)

const defaultMaxEntries = 100

type Entry struct {
	ID          string    `json:"id"`
	At          time.Time `json:"at"`
	Path        string    `json:"path"`
	RequestName string    `json:"request_name,omitempty"`
	Profile     string    `json:"profile"`
	OK          bool      `json:"ok"`
}

type Store struct {
	Path       string
	MaxEntries int
}

func NewStore(path string) Store {
	max := defaultMaxEntries
	if max <= 0 {
		max = defaultMaxEntries
	}
	return Store{Path: path, MaxEntries: max}
}

type fileData struct {
	Entries []Entry `json:"entries"`
}

func (s Store) max() int {
	if s.MaxEntries > 0 {
		return s.MaxEntries
	}
	return defaultMaxEntries
}

func (s Store) EnsureDir() *errors.XError {
	dir := filepath.Dir(s.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return errors.Wrap(errors.CodeCfgInvalid, "failed to create history directory", map[string]any{"path": dir}, err)
	}
	return nil
}

func (s Store) Load() ([]Entry, *errors.XError) {
	if s.Path == "" {
		return nil, errors.New(errors.CodeCfgInvalid, "history path not configured", nil)
	}
	b, err := os.ReadFile(s.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, errors.Wrap(errors.CodeCfgInvalid, "failed to read history", map[string]any{"path": s.Path}, err)
	}
	var data fileData
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, errors.Wrap(errors.CodeCfgInvalid, "invalid history file", map[string]any{"path": s.Path}, err)
	}
	if data.Entries == nil {
		return []Entry{}, nil
	}
	return data.Entries, nil
}

func (s Store) save(entries []Entry) *errors.XError {
	if xe := s.EnsureDir(); xe != nil {
		return xe
	}
	b, err := json.MarshalIndent(fileData{Entries: entries}, "", "  ")
	if err != nil {
		return errors.Wrap(errors.CodeCfgInvalid, "failed to encode history", nil, err)
	}
	b = append(b, '\n')
	if err := os.WriteFile(s.Path, b, 0o600); err != nil {
		return errors.Wrap(errors.CodeCfgInvalid, "failed to write history", map[string]any{"path": s.Path}, err)
	}
	return nil
}

func (s Store) Append(e Entry) *errors.XError {
	entries, xe := s.Load()
	if xe != nil {
		return xe
	}
	if e.ID == "" {
		e.ID = time.Now().UTC().Format("20060102T150405.000")
	}
	if e.At.IsZero() {
		e.At = time.Now().UTC()
	}
	if e.Path != "" {
		if abs, err := filepath.Abs(e.Path); err == nil {
			e.Path = abs
		}
	}
	entries = append([]Entry{e}, entries...)
	if len(entries) > s.max() {
		entries = entries[:s.max()]
	}
	return s.save(entries)
}

// ListNewestFirst returns entries sorted newest first.
func (s Store) ListNewestFirst() ([]Entry, *errors.XError) {
	entries, xe := s.Load()
	if xe != nil {
		return nil, xe
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].At.After(entries[j].At)
	})
	return entries, nil
}

// GetByIndex returns 1-based index (1 = most recent).
func (s Store) GetByIndex(index int) (Entry, *errors.XError) {
	list, xe := s.ListNewestFirst()
	if xe != nil {
		return Entry{}, xe
	}
	if index < 1 || index > len(list) {
		return Entry{}, errors.New(errors.CodeCfgInvalid, "history index out of range",
			map[string]any{"index": index, "count": len(list)})
	}
	return list[index-1], nil
}

func (s Store) Latest() (Entry, *errors.XError) {
	return s.GetByIndex(1)
}

func DefaultStore(homeDir string) Store {
	return NewStore(paths.DefaultHistoryPath(homeDir))
}
