package httpfile

import (
	"testing"
)

func TestParseSingleRequest(t *testing.T) {
	content := []byte(`GET https://example.com/hello
Accept: application/json

`)
	reqs, xe := Parse(content)
	if xe != nil {
		t.Fatal(xe)
	}
	if len(reqs) != 1 {
		t.Fatalf("got %d requests", len(reqs))
	}
	if reqs[0].Method != "GET" || reqs[0].URL != "https://example.com/hello" {
		t.Fatalf("unexpected request: %+v", reqs[0])
	}
	if reqs[0].Headers["Accept"] != "application/json" {
		t.Fatalf("headers: %+v", reqs[0].Headers)
	}
}

func TestParseMultiWithNames(t *testing.T) {
	content := []byte(`### health
GET {{host}}/health

### create
POST {{host}}/items
Content-Type: application/json

{"name":"test"}
`)
	reqs, xe := Parse(content)
	if xe != nil {
		t.Fatal(xe)
	}
	if len(reqs) != 2 {
		t.Fatalf("got %d requests", len(reqs))
	}
	if reqs[0].Name != "health" || reqs[1].Name != "create" {
		t.Fatalf("names: %q, %q", reqs[0].Name, reqs[1].Name)
	}
	if reqs[1].Body != `{"name":"test"}` {
		t.Fatalf("body: %q", reqs[1].Body)
	}
}

func TestParseFileVars(t *testing.T) {
	content := []byte(`@host=https://api.test
GET {{host}}/v1
`)
	reqs, xe := Parse(content)
	if xe != nil {
		t.Fatal(xe)
	}
	if reqs[0].Vars["host"] != "https://api.test" {
		t.Fatalf("vars: %+v", reqs[0].Vars)
	}
}

func TestFindByName(t *testing.T) {
	reqs := []Request{{Name: "a"}, {Name: "b"}}
	_, ok := FindByName(reqs, "")
	if ok {
		t.Fatal("expected false for unnamed multi")
	}
	r, ok := FindByName(reqs, "b")
	if !ok || r.Name != "b" {
		t.Fatal("expected b")
	}
}

func TestInterpolate(t *testing.T) {
	vars := map[string]string{"host": "https://x.com", "id": "42"}
	got := Interpolate("GET {{host}}/items/{{id}}", vars)
	want := "GET https://x.com/items/42"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
