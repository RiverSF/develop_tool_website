package json2struct_test

import (
	"strings"
	"testing"

	"develop_tools/pkg/json2struct"
)

func TestJson2struct(t *testing.T) {
	out := json2struct.Json2struct(`{"name":"demo","count":1}`)
	if out == "" {
		t.Fatal("expected struct output")
	}
	if !strings.Contains(out, "type Request struct") || !strings.Contains(out, "Name string") || !strings.Contains(out, "Count int") {
		t.Fatalf("unexpected output: %s", out)
	}
}
