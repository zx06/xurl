package history

import (
	"path/filepath"
	"testing"
)

func TestAppendAndList(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(filepath.Join(dir, "history.json"))
	if xe := store.Append(Entry{Path: "/tmp/a.http", Profile: "dev", OK: true}); xe != nil {
		t.Fatal(xe)
	}
	if xe := store.Append(Entry{Path: "/tmp/b.http", Profile: "prod", OK: false}); xe != nil {
		t.Fatal(xe)
	}
	list, xe := store.ListNewestFirst()
	if xe != nil {
		t.Fatal(xe)
	}
	if len(list) != 2 {
		t.Fatalf("len %d", len(list))
	}
	if list[0].Path != "/tmp/b.http" {
		t.Fatalf("newest first: %+v", list)
	}
	e, xe := store.GetByIndex(2)
	if xe != nil {
		t.Fatal(xe)
	}
	if e.Path != "/tmp/a.http" {
		t.Fatalf("index 2: %+v", e)
	}
}

func TestCapEntries(t *testing.T) {
	dir := t.TempDir()
	store := Store{Path: filepath.Join(dir, "h.json"), MaxEntries: 3}
	for i := 0; i < 5; i++ {
		p := filepath.Join("/tmp", "f"+string(rune('0'+i))+".http")
		if xe := store.Append(Entry{Path: p, Profile: "p"}); xe != nil {
			t.Fatal(xe)
		}
	}
	list, xe := store.Load()
	if xe != nil {
		t.Fatal(xe)
	}
	if len(list) != 3 {
		t.Fatalf("expected cap 3, got %d", len(list))
	}
}
