package library

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveTargetFilesystem(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "demo.http")
	if err := os.WriteFile(f, []byte("GET https://example.com\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, xe := ResolveTarget(dir, filepath.Join(dir, "lib"), "demo.http")
	if xe != nil {
		t.Fatal(xe)
	}
	if got != f {
		t.Fatalf("got %q want %q", got, f)
	}
}

func TestResolveTargetLibrary(t *testing.T) {
	work := t.TempDir()
	lib := filepath.Join(work, "library")
	if err := os.MkdirAll(lib, 0o755); err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(lib, "ping.http")
	if err := os.WriteFile(f, []byte("GET https://example.com\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, xe := ResolveTarget(work, lib, "ping")
	if xe != nil {
		t.Fatal(xe)
	}
	if got != f {
		t.Fatalf("got %q want %q", got, f)
	}
}

func TestResolveTargetSuffix(t *testing.T) {
	work := t.TempDir()
	lib := filepath.Join(work, "library")
	f := filepath.Join(lib, "api.http")
	if err := os.MkdirAll(lib, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f, []byte("GET https://example.com\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, xe := ResolveTarget(work, lib, "api")
	if xe != nil {
		t.Fatal(xe)
	}
	if got != f {
		t.Fatalf("got %q want %q", got, f)
	}
}
