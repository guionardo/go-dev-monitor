package finder

import (
	"iter"
	"os"
	"path"
)

func FindRepositories(folder string, maxFolderDept int) iter.Seq[string] {
	return func(yield func(string) bool) {
		for f := range findRepositories(folder, maxFolderDept, 0) {
			if !yield(f) {
				return // Stop if yield returns false
			}
		}
	}
}

func findRepositories(folder string, maxFolderDept int, level int) iter.Seq[string] {
	return func(yield func(string) bool) {
		if level > maxFolderDept {
			return // Stop if the maximum folder depth is reached
		}

		entries, err := os.ReadDir(folder)
		if err != nil {
			return //TODO: handle error
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue // Skip non-directory entries
			}
			currentFolder := path.Join(folder, entry.Name())
			gitFolder := path.Join(currentFolder, ".git")
			// Check if the .git folder exists
			if stat, err := os.Stat(gitFolder); err == nil && stat.IsDir() {
				if !yield(currentFolder) {
					return
				}
				continue
			}
			// Is another folder
			for sub := range findRepositories(currentFolder, maxFolderDept, level+1) {
				if !yield(sub) {
					return
				}
			}

		}

	}
}
