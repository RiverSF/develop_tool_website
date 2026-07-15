package adx_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/adx"
	"develop_tools/pkg/logger"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	_ = os.Setenv("APP_ROOT", root)
	if err := logger.Init(); err != nil {
		panic(err)
	}
	code := m.Run()
	_ = logger.Close()
	os.Exit(code)
}

func TestAdxPagesRender(t *testing.T) {
	cases := []struct {
		name string
		fn   gin.HandlerFunc
	}{
		{"index", adx.AdxIndex},
		{"cn", adx.AdxCn},
		{"dsp", adx.AdxDSP},
		{"adm", adx.AdxAdm},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/adx", nil)
			tc.fn(c)
			if w.Code != http.StatusOK {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if w.Body.Len() == 0 {
				t.Fatal("empty body")
			}
		})
	}
}
