package adx_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"develop_tools/internal/config"
	"develop_tools/internal/handler/adx"
)

func setupUserSyncConfig(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	iniPath := filepath.Join(dir, "app.ini")
	content := strings.TrimSpace(`
RUN_MODE = test

[server]
HTTP_PORT = 9080
READ_TIMEOUT = 60
WRITE_TIMEOUT = 60

[host]
HOST_TEST = http://localhost:9080

[mysql]
MYSQL_HOST = 127.0.0.1
MYSQL_PORT = 3306
MYSQL_USER = root
MYSQL_PASSWORD =
MYSQL_DB = test
`) + "\n"
	if err := os.WriteFile(iniPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CONFIG_PATH", iniPath)
	t.Setenv("RUN_MODE", "test")
	if err := config.Init(); err != nil {
		t.Fatal(err)
	}
}

func TestAdxUserSyncNoRedirect(t *testing.T) {
	setupUserSyncConfig(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/adx/userSync", nil)

	adx.AdxUserSync(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestAdxUserSyncBlocksDisallowedHost(t *testing.T) {
	setupUserSyncConfig(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	q := url.Values{}
	q.Set("my_redirect_url", "https://evil.example/phish")
	c.Request = httptest.NewRequest(http.MethodGet, "/adx/userSync?"+q.Encode(), nil)

	adx.AdxUserSync(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v body=%s", err, w.Body.String())
	}
	if body.Status != 400 {
		t.Fatalf("want status 400, got %+v body=%s", body, w.Body.String())
	}
}

func TestAdxUserSyncAllowsLocalhostRedirect(t *testing.T) {
	setupUserSyncConfig(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	q := url.Values{}
	q.Set("my_redirect_url", "http://localhost/callback")
	q.Set("my_dsp_uid", "dsp:test-uid")
	c.Request = httptest.NewRequest(http.MethodGet, "/adx/userSync?"+q.Encode(), nil)

	adx.AdxUserSync(c)

	if w.Code != http.StatusFound {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	loc := w.Header().Get("Location")
	if !strings.Contains(loc, "localhost") {
		t.Fatalf("unexpected Location %q", loc)
	}
}

func TestAdxUserSyncReplacesMacros(t *testing.T) {
	setupUserSyncConfig(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	q := url.Values{}
	q.Set("my_redirect_url", "http://localhost/cb?uid="+adx.AdxUserIdMacro+"&g="+adx.GdprMacro+"&c="+adx.GdprConsentMacro)
	q.Set("my_dsp_uid", "dsp:abc")
	q.Set("my_gdpr", "1")
	q.Set("my_gdpr_consent", "consent-token")
	c.Request = httptest.NewRequest(http.MethodGet, "/adx/userSync?"+q.Encode(), nil)

	adx.AdxUserSync(c)

	if w.Code != http.StatusFound {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	loc := w.Header().Get("Location")
	if strings.Contains(loc, adx.AdxUserIdMacro) || strings.Contains(loc, adx.GdprMacro) || strings.Contains(loc, adx.GdprConsentMacro) {
		t.Fatalf("macros not replaced: %q", loc)
	}
	if !strings.Contains(loc, "dsp:abc") || !strings.Contains(loc, "consent-token") {
		t.Fatalf("expected substituted values in %q", loc)
	}
}
