package path

import (
	"os"
	"path/filepath"
)

// Root returns application root directory.
// Override with APP_ROOT when running binary outside project root.
func Root() string {
	if v := os.Getenv("APP_ROOT"); v != "" {
		return v
	}
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

// Join joins paths under application root.
func Join(elem ...string) string {
	return filepath.Join(append([]string{Root()}, elem...)...)
}
