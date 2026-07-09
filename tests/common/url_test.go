package common_test

import (
	"testing"

	"develop_tools/pkg/common"
)

func TestValidateRedirectURL(t *testing.T) {
	allowed := []string{"river.site", "localhost"}
	ok, err := common.ValidateRedirectURL("https://river.site/callback", allowed...)
	if err != nil || ok == "" {
		t.Fatalf("expected allowed redirect, got %q err=%v", ok, err)
	}
	if _, err := common.ValidateRedirectURL("https://evil.example/phish", allowed...); err == nil {
		t.Fatal("expected blocked redirect")
	}
}
