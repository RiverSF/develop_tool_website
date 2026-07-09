package common_test

import (
	"strings"
	"testing"

	"develop_tools/pkg/common"
)

func TestUniqueCountLines(t *testing.T) {
	out := common.UniqueCountLines("a\na\nb\n")
	if out == "" {
		t.Fatal("expected output")
	}
	if !strings.Contains(out, "2 a") || !strings.Contains(out, "1 b") {
		t.Fatalf("unexpected output: %q", out)
	}
}
