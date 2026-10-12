package runner_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zx06/xurl/internal/config"
	"github.com/zx06/xurl/internal/httpfile"
	"github.com/zx06/xurl/internal/runner"
)

func TestRunIntegrationMockServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test") != "yes" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"hello":"world"}`)
	}))
	defer srv.Close()

	reqs, xe := httpfile.Parse([]byte("GET " + srv.URL + "/api\nX-Test: yes\n"))
	if xe != nil {
		t.Fatal(xe)
	}

	result := runner.Run(reqs[0], runner.Options{
		Profile: config.Profile{AllowPlaintext: true},
	})
	if !result.OK {
		t.Fatalf("expected ok: %+v", result)
	}
	if result.Status != 200 {
		t.Fatalf("status %d", result.Status)
	}
	if !strings.Contains(result.Body, "world") {
		t.Fatalf("body %q", result.Body)
	}
}

func findModuleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}

func TestCLIRunE2E(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "pong")
	}))
	defer srv.Close()

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "xurl.yaml")
	cfg := `
profiles:
  local:
    base_url: "` + srv.URL + `"
    allow_plaintext: true
`
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o600); err != nil {
		t.Fatal(err)
	}
	httpFile := filepath.Join(dir, "example.http")
	content := `### ping
GET /ping
`
	if err := os.WriteFile(httpFile, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	bin := filepath.Join(t.TempDir(), "xurl")
	root := findModuleRoot(t)
	build := exec.Command("go", "build", "-o", bin, "./cmd/xurl")
	build.Dir = root
	out, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}

	cmd := exec.Command(bin, "run", httpFile, "ping", "-p", "local", "--config", cfgPath, "-f", "json")
	raw, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, raw)
	}

	var result map[string]any
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("json: %v\n%s", err, raw)
	}
	if ok, _ := result["ok"].(bool); !ok {
		t.Fatalf("not ok: %s", raw)
	}
}
