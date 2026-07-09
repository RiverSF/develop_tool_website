package path

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	wdOnce sync.Once
	wdVal  string
)

// Root returns application root directory.
// Override with APP_ROOT when running binary outside project root.
func Root() string {
	if v := os.Getenv("APP_ROOT"); v != "" {
		return v
	}
	wdOnce.Do(func() {
		wd, err := os.Getwd()
		if err != nil {
			wdVal = "."
			return
		}
		wdVal = wd
	})
	return wdVal
}

// Join joins paths under application root.
func Join(elem ...string) string {
	if len(elem) == 0 {
		return Root()
	}
	parts := make([]string, 0, len(elem)+1)
	parts = append(parts, Root())
	parts = append(parts, elem...)
	return filepath.Join(parts...)
}
