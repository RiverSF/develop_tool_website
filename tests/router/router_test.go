package router_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"develop_tools/internal/router"
	"develop_tools/pkg/logger"
	"develop_tools/pkg/path"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	_ = logger.Init()
	m.Run()
	_ = logger.Close()
}

func TestNewRegistersCoreAndAdxRoutes(t *testing.T) {
	r := router.New()

	want := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/ping"},
		{http.MethodGet, "/"},
		{http.MethodGet, "/data-calculator"},
		{http.MethodGet, "/xml"},
		{http.MethodGet, "/adx"},
		{http.MethodGet, "/adx/cn"},
		{http.MethodGet, "/adx/dsp"},
		{http.MethodGet, "/adx/adm"},
		{http.MethodPost, "/adx/adxGetDspList"},
		{http.MethodGet, "/adx/adxGetDspAdm"},
		{http.MethodGet, "/adx/adxGetDspResponse"},
		{http.MethodPost, "/adx/adxDspSave"},
		{http.MethodGet, "/adx/adxGetDspNotice"},
		{http.MethodPost, "/adx/cn/:uniqueKey"},
		{http.MethodGet, "/adx/:uniqueKey/:noticeType"},
		{http.MethodPost, "/adx/:uniqueKey"},
		{http.MethodPost, "/adx/dsp/:uniqueKey"},
		{http.MethodGet, "/adx/userSync"},
		{http.MethodPost, "/adx/bundle/extract"},
	}

	registered := map[string]bool{}
	for _, rt := range r.Routes() {
		registered[rt.Method+" "+rt.Path] = true
	}

	for _, w := range want {
		key := w.method + " " + w.path
		if !registered[key] {
			t.Errorf("missing route %s", key)
		}
	}
}

func TestPing(t *testing.T) {
	r := router.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if got := w.Body.String(); got == "" || got == "null" {
		t.Fatalf("unexpected body %q", got)
	}
}

func TestStaticAssetsMounted(t *testing.T) {
	r := router.New()
	for _, rt := range r.Routes() {
		if strings.HasPrefix(rt.Path, "/assets") {
			return
		}
	}
	t.Fatalf("assets static route not registered; root=%s routes=%v", path.Root(), routePaths(r))
}

func routePaths(r *gin.Engine) []string {
	out := make([]string, 0, len(r.Routes()))
	for _, rt := range r.Routes() {
		out = append(out, rt.Method+" "+rt.Path)
	}
	return out
}
