package path_test

import (
	"path/filepath"
	"testing"

	"develop_tools/pkg/path"
)

func TestJoinUsesRoot(t *testing.T) {
	t.Setenv("APP_ROOT", filepath.FromSlash("D:/demo"))
	got := path.Join("web", "templates", "a.html")
	want := filepath.Join("D:/demo", "web", "templates", "a.html")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
