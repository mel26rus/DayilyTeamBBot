package filesystem

import (
	"os"
	"path/filepath"
	"strings"

	"expense-bot/internal/platform"
)

func BaseDir(configFile string) string {
	isService, _ := platform.IsService()

	if wd, err := os.Getwd(); err == nil {
		if _, err = os.Stat(filepath.Join(wd, configFile)); err == nil && !isService {
			return wd
		}

		// go run builds the executable under the temporary go-build directory.
		// On the first launch there is no config yet, so use the project directory.
		if !isService {
			if exe, exeErr := os.Executable(); exeErr == nil {
				if rel, relErr := filepath.Rel(os.TempDir(), exe); relErr == nil &&
					!strings.HasPrefix(rel, "..") && strings.HasPrefix(filepath.Clean(rel), "go-build") {
					return wd
				}
			}
		}
	}

	exe, err := os.Executable()
	if err != nil {
		return "."
	}

	return filepath.Dir(exe)
}
