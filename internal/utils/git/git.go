package git

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// GetGitRoot searches for the root of a Git repository starting from the specified folder.
func GetGitRoot(folderName string) (root string, err error) {
	var folder string
	if folder, err = filepath.Abs(folderName); err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	for {
		if stat, err := os.Stat(path.Join(folder, ".git")); err == nil && stat.IsDir() {
			return folder, nil
		}
		folder = filepath.Dir(folder)
		if filepath.Clean(folder) == string(filepath.Separator) {
			break // Reached the root directory without finding a .git folder
		}

	}
	return "", fmt.Errorf("no git repository found in %s or its parent directories", folderName)
}
