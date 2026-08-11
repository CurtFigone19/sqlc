package migrations

import (
	"testing"
	"testing/fstest"
)

func TestRead(t *testing.T) {
	fsys := fstest.MapFS{
		"migrations/10_posts.sql":             {Data: []byte("")},
		"migrations/2_users.sql":              {Data: []byte("")},
		"migrations/1_init.sql":               {Data: []byte("")},
		"migrations/20231024000000_init.sql": {Data: []byte("")},
		"migrations/a_init.sql":               {Data: []byte("")},
	}

	files, err := Read(fsys, "migrations")
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"migrations/1_init.sql",
		"migrations/2_users.sql",
		"migrations/10_posts.sql",
		"migrations/20231024000000_init.sql",
		"migrations/a_init.sql",
	}

	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for i, f := range files {
		if f != expected[i] {
			t.Errorf("at index %d: expected %s, got %s", i, expected[i], f)
		}
	}
}
