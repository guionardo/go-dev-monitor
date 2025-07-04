package detectors

import (
	"os"
	"path"
	"path/filepath"
)

func searchFiles(folder string, files ...string) string {
	for _, file := range files {
		if stat, err := os.Stat(path.Join(folder, file)); err == nil && !stat.IsDir() {
			return path.Join(folder, file)
		}
		if matches, err := filepath.Glob(path.Join(folder, file)); err == nil && len(matches) > 0 {
			return matches[0]
		}
	}
	return ""
}
